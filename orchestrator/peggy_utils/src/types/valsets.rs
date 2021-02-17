use super::*;
use crate::error::PeggyError;
use clarity::Address as EthAddress;
use clarity::Signature as EthSignature;
use contact::jsonrpc::error::JsonRpcError;
use deep_space::address::Address as CosmosAddress;
use std::{
    cmp::Ordering,
    collections::{HashMap, HashSet},
    fmt,
};

/// The total power in the Peggy bridge is normalized to u32 max every
/// time a validator set is created. This value of up to u32 max is then
/// stored in a u64 to prevent overflow during computation.
pub const TOTAL_PEGGY_POWER: u64 = u32::MAX as u64;

/// takes in an amount of power in the peggy bridge, returns a percentage of total
fn peggy_power_to_percent(input: u64) -> f32 {
    (input as f32 / TOTAL_PEGGY_POWER as f32) * 100f32
}

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetConfirmResponse {
    pub orchestrator: CosmosAddress,
    pub eth_address: EthAddress,
    pub nonce: u64,
    pub eth_signature: EthSignature,
}

impl ValsetConfirmResponse {
    pub fn from_proto(input: peggy_proto::peggy::MsgValsetConfirm) -> Result<Self, PeggyError> {
        Ok(ValsetConfirmResponse {
            orchestrator: input.orchestrator.parse()?,
            eth_address: input.eth_address.parse()?,
            nonce: input.nonce,
            eth_signature: input.signature.parse()?,
        })
    }
}

/// the response we get when querying for a batch confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BatchConfirmResponse {
    pub nonce: u64,
    pub orchestrator: CosmosAddress,
    pub token_contract: EthAddress,
    pub ethereum_signer: EthAddress,
    pub eth_signature: EthSignature,
}

impl BatchConfirmResponse {
    pub fn from_proto(input: peggy_proto::peggy::MsgConfirmBatch) -> Result<Self, PeggyError> {
        Ok(BatchConfirmResponse {
            nonce: input.nonce,
            orchestrator: input.orchestrator.parse()?,
            token_contract: input.token_contract.parse()?,
            ethereum_signer: input.eth_signer.parse()?,
            eth_signature: input.signature.parse()?,
        })
    }
}

fn valset_confirms_to_hashmap(
    input: &[ValsetConfirmResponse],
) -> HashMap<EthAddress, ValsetConfirmResponse> {
    let mut ret = HashMap::new();
    for item in input {
        ret.insert(item.eth_address, item.clone());
    }
    ret
}

fn batch_confirms_to_hashmap(
    input: &[BatchConfirmResponse],
) -> HashMap<EthAddress, BatchConfirmResponse> {
    let mut ret = HashMap::new();
    for item in input {
        ret.insert(item.ethereum_signer, item.clone());
    }
    ret
}

/// a list of validators, powers, and eth addresses at a given block height
#[derive(Serialize, Deserialize, Debug, Default, Clone, PartialEq, Eq)]
pub struct Valset {
    pub nonce: u64,
    pub members: Vec<ValsetMember>,
}

impl Valset {
    /// Takes an array of Option<EthAddress> and converts to EthAddress and replaces with zeros
    /// when none is found, Zeros are interpreted by the contract as 'no signature provided' and
    /// signature checks can pass with up to 33% of all voting power presented as zeroed addresses
    pub fn filter_empty_addresses(&self) -> (Vec<EthAddress>, Vec<u64>) {
        let mut addresses = Vec::new();
        let mut powers = Vec::new();
        for val in self.members.iter() {
            match val.eth_address {
                Some(a) => {
                    addresses.push(a);
                    powers.push(val.power);
                }
                None => {
                    addresses.push(EthAddress::default());
                    powers.push(val.power);
                }
            }
        }
        (addresses, powers)
    }

    pub fn get_power(&self, address: EthAddress) -> Result<u64, JsonRpcError> {
        for val in self.members.iter() {
            if val.eth_address == Some(address) {
                return Ok(val.power);
            }
        }
        Err(JsonRpcError::BadInput(
            "All Eth Addresses must be set".to_string(),
        ))
    }

    /// combines the provided signatures with the valset ensuring that ordering and signature data is correct
    /// Note the 'correct' ordering is the *same* ordering as the validator set members in 'self'. In some cases
    /// this will be sorted, in others it will be improperly sorted but must be maintained so that the signatures
    /// are accepted on the Ethereum chain, which requires the submitted addresses to match whatever the previously
    /// submitted ordering was and the signatures must be in parallel arrays to reduce shuffling.
    /// TODO give the signatures types a trait and de-duplicate
    pub fn order_valset_sigs(
        &self,
        signatures: &[ValsetConfirmResponse],
    ) -> Result<Vec<PeggySignature>, PeggyError> {
        if signatures.is_empty() {
            return Err(PeggyError::InsufficientVotingPowerToPass(
                "No signatures!".to_string(),
            ));
        }
        // if signatures is not empty this is safe.
        let new_valset_nonce = signatures[0].nonce;

        let mut out = Vec::new();
        let signatures_hashmap = valset_confirms_to_hashmap(signatures);
        let mut power_of_good_sigs = 0;
        let mut power_of_unset_keys = 0;
        let mut number_of_unset_key_validators = 0;
        let mut power_of_nonvoters = 0;
        let mut number_of_nonvoters = 0;
        for member in self.members.iter() {
            if let Some(eth_address) = member.eth_address {
                if let Some(sig) = signatures_hashmap.get(&eth_address) {
                    assert_eq!(sig.eth_address, eth_address);
                    out.push(PeggySignature {
                        power: member.power,
                        eth_address: sig.eth_address,
                        v: sig.eth_signature.v.clone(),
                        r: sig.eth_signature.r.clone(),
                        s: sig.eth_signature.s.clone(),
                    });
                    power_of_good_sigs += member.power;
                } else {
                    out.push(PeggySignature {
                        power: member.power,
                        eth_address,
                        v: 0u8.into(),
                        r: 0u8.into(),
                        s: 0u8.into(),
                    });
                    power_of_nonvoters += member.power;
                    number_of_nonvoters += 1;
                }
            } else {
                out.push(PeggySignature {
                    power: member.power,
                    eth_address: EthAddress::default(),
                    v: 0u8.into(),
                    r: 0u8.into(),
                    s: 0u8.into(),
                });
                power_of_unset_keys += member.power;
                number_of_unset_key_validators += 1;
            }
        }

        let num_validators = self.members.len();
        // now that we have collected the signatures we can determine if the measure has the votes to pass
        // and error early if it does not, otherwise the user will pay fees for a transaction that will
        // just throw
        if peggy_power_to_percent(power_of_good_sigs) < 66f32 {
            let message = format!(
                "
                Valset {} has {}/{} or {:.2}% power voting! Can not execute on Ethereum!
                {}/{} validators have unset Ethereum keys representing {}/{} or {:.2}% of the power required
                {}/{} validators have Ethereum keys set but have not voted representing {}/{} or {:.2}% of the power required
                This valset probably just needs to accumulate signatures for a moment.",
                new_valset_nonce,
                power_of_good_sigs,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_good_sigs),
                number_of_unset_key_validators,
                num_validators,
                power_of_unset_keys,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_unset_keys),
                number_of_nonvoters,
                num_validators,
                power_of_nonvoters,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_nonvoters),
            );
            Err(PeggyError::InsufficientVotingPowerToPass(message))
        } else {
            Ok(out)
        }
    }

    /// combines the provided signatures with the valset ensuring that ordering and signature data is correct
    /// TODO give the signatures types a trait and de-duplicate
    pub fn order_batch_sigs(
        &self,
        signatures: &[BatchConfirmResponse],
    ) -> Result<Vec<PeggySignature>, PeggyError> {
        if signatures.is_empty() {
            return Err(PeggyError::InsufficientVotingPowerToPass(
                "No signatures!".to_string(),
            ));
        }
        // if the signatures are not empty this is safe
        let batch_nonce = signatures[0].nonce;
        let batch_erc20 = signatures[0].token_contract;

        let mut out = Vec::new();
        let signatures_hashmap = batch_confirms_to_hashmap(signatures);
        let mut power_of_good_sigs = 0;
        let mut power_of_unset_keys = 0;
        let mut number_of_unset_key_validators = 0;
        let mut power_of_nonvoters = 0;
        let mut number_of_nonvoters = 0;
        for member in self.members.iter() {
            if let Some(eth_address) = member.eth_address {
                if let Some(sig) = signatures_hashmap.get(&eth_address) {
                    assert_eq!(sig.ethereum_signer, eth_address);
                    out.push(PeggySignature {
                        power: member.power,
                        eth_address: sig.ethereum_signer,
                        v: sig.eth_signature.v.clone(),
                        r: sig.eth_signature.r.clone(),
                        s: sig.eth_signature.s.clone(),
                    });
                    power_of_good_sigs += member.power;
                } else {
                    out.push(PeggySignature {
                        power: member.power,
                        eth_address,
                        v: 0u8.into(),
                        r: 0u8.into(),
                        s: 0u8.into(),
                    });
                    power_of_nonvoters += member.power;
                    number_of_nonvoters += 1;
                }
            } else {
                out.push(PeggySignature {
                    power: member.power,
                    eth_address: EthAddress::default(),
                    v: 0u8.into(),
                    r: 0u8.into(),
                    s: 0u8.into(),
                });
                power_of_unset_keys += member.power;
                number_of_unset_key_validators += 1;
            }
        }

        let num_validators = self.members.len();
        // now that we have collected the signatures we can determine if the measure has the votes to pass
        // and error early if it does not, otherwise the user will pay fees for a transaction that will
        // just throw
        if peggy_power_to_percent(power_of_good_sigs) < 66f32 {
            let message = format!(
                "
                Batch {} for token {} has {}/{} or {:.2}% power voting! Can not execute on Ethereum!
                {}/{} validators have unset Ethereum keys representing {}/{} or {:.2}% of the power required
                {}/{} validators have Ethereum keys set but have not voted representing {}/{} or {:.2}% of the power required
                This batch probably just needs to accumulate signatures, for a moment.",
                batch_nonce,
                batch_erc20,
                power_of_good_sigs,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_good_sigs),
                number_of_unset_key_validators,
                num_validators,
                power_of_unset_keys,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_unset_keys),
                number_of_nonvoters,
                num_validators,
                power_of_nonvoters,
                TOTAL_PEGGY_POWER,
                peggy_power_to_percent(power_of_nonvoters),
            );
            warn!("{}", message);
            Err(PeggyError::InsufficientVotingPowerToPass(message))
        } else {
            Ok(out)
        }
    }

    /// A utility function to provide a HashMap of members for easy lookups
    pub fn to_hashmap(&self) -> HashMap<EthAddress, u64> {
        let mut res = HashMap::new();
        for item in self.members.iter() {
            if let Some(address) = item.eth_address {
                res.insert(address, item.power);
            } else {
                error!("Validator in active set without Eth Address! This must be corrected immediately!")
            }
        }
        res
    }

    /// A utility function to provide a HashSet of members for union operations
    pub fn to_hashset(&self) -> HashSet<EthAddress> {
        let mut res = HashSet::new();
        for item in self.members.iter() {
            if let Some(address) = item.eth_address {
                res.insert(address);
            } else {
                error!("Validator in active set without Eth Address! This must be corrected immediately!")
            }
        }
        res
    }

    /// This function takes the current valset and compares it to a provided one
    /// returning a percentage difference in their power allocation. This is a very
    /// important function as it's used to decide when the validator sets are updated
    /// on the Ethereum chain and when new validator sets are requested on the Cosmos
    /// side. In theory an error here, if unnoticed for long enough, could allow funds
    /// to be stolen from the bridge without the validators in question still having stake
    /// to lose.
    /// Returned value must be less than or equal to two
    pub fn power_diff(&self, other: &Valset) -> f32 {
        let mut total_power_diff = 0u64;
        let a = self.to_hashmap();
        let b = other.to_hashmap();
        let a_map = self.to_hashset();
        let b_map = other.to_hashset();
        // items in A and B, we go through these and compute the absolute value of the
        // difference in power and sum it.
        let intersection = a_map.intersection(&b_map);
        // items in A but not in B or vice versa, since we're just trying to compute the difference
        // we can simply sum all of these up.
        let symmetric_difference = a_map.symmetric_difference(&b_map);
        for item in symmetric_difference {
            let mut power = None;
            if let Some(val) = a.get(item) {
                power = Some(val);
            } else if let Some(val) = b.get(item) {
                power = Some(val);
            }
            // impossible for this to panic without a failure in the logic
            // of the symmetric difference function
            let power = power.unwrap();
            total_power_diff += power;
        }
        for item in intersection {
            // can't panic since there must be an entry for both.
            let power_a = a[item];
            let power_b = b[item];
            if power_a > power_b {
                total_power_diff += power_a - power_b;
            } else {
                total_power_diff += power_b - power_a;
            }
        }

        (total_power_diff as f32) / (u32::MAX as f32)
    }
}

impl From<peggy_proto::peggy::Valset> for Valset {
    fn from(input: peggy_proto::peggy::Valset) -> Self {
        Valset {
            nonce: input.nonce,
            members: input.members.iter().map(|i| i.into()).collect(),
        }
    }
}

impl From<&peggy_proto::peggy::Valset> for Valset {
    fn from(input: &peggy_proto::peggy::Valset) -> Self {
        Valset {
            nonce: input.nonce,
            members: input.members.iter().map(|i| i.into()).collect(),
        }
    }
}

/// a list of validators, powers, and eth addresses at a given block height
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ValsetMember {
    // ord sorts on the first member first, so this produces the correct sorting
    pub power: u64,
    pub eth_address: Option<EthAddress>,
}

impl Ord for ValsetMember {
    // Alex wrote the Go sorting implementation for validator
    // sets as Greatest to Least, now this isn't the convention
    // for any standard sorting implementation and Rust doesn't
    // really like it when you implement sort yourself. It prefers
    // Ord. So here we implement Ord with the Eth address sorting
    // reversed, since they are also sorted greatest to least in
    // the Cosmos module. Then we can call .sort and .reverse and get
    // the same sorting as the Cosmos module.
    fn cmp(&self, other: &Self) -> Ordering {
        if self.power != other.power {
            self.power.cmp(&other.power)
        } else {
            self.eth_address.cmp(&other.eth_address).reverse()
        }
    }
}

impl PartialOrd for ValsetMember {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl ValsetMember {
    pub fn display_vec(input: &[ValsetMember]) -> String {
        let mut out = String::new();
        for val in input.iter() {
            out += &val.to_string()
        }
        out
    }
}

impl fmt::Display for ValsetMember {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self.eth_address {
            Some(a) => write!(f, "Address: {} Power: {}", a, self.power),
            None => write!(f, "Address: None Power: {}", self.power),
        }
    }
}

impl From<peggy_proto::peggy::BridgeValidator> for ValsetMember {
    fn from(input: peggy_proto::peggy::BridgeValidator) -> Self {
        let eth_address = match input.ethereum_address.parse() {
            Ok(e) => Some(e),
            Err(_) => None,
        };
        ValsetMember {
            power: input.power,
            eth_address,
        }
    }
}

impl From<&peggy_proto::peggy::BridgeValidator> for ValsetMember {
    fn from(input: &peggy_proto::peggy::BridgeValidator) -> Self {
        let eth_address = match input.ethereum_address.parse() {
            Ok(e) => Some(e),
            Err(_) => None,
        };
        ValsetMember {
            power: input.power,
            eth_address,
        }
    }
}
