use super::*;
use crate::error::PeggyError;
use clarity::Address as EthAddress;
use clarity::Signature as EthSignature;
use contact::{jsonrpc::error::JsonRpcError, types::parse_val};
use deep_space::address::Address as CosmosAddress;
use std::{
    cmp::Ordering,
    collections::{HashMap, HashSet},
    fmt,
};

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetConfirmResponse {
    pub validator: CosmosAddress,
    pub eth_address: EthAddress,
    pub nonce: u64,
    pub eth_signature: EthSignature,
}

impl ValsetConfirmResponse {
    pub fn from_proto(input: peggy_proto::peggy::MsgValsetConfirm) -> Result<Self, PeggyError> {
        Ok(ValsetConfirmResponse {
            validator: input.validator.parse()?,
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
    pub validator: CosmosAddress,
    pub token_contract: EthAddress,
    pub ethereum_signer: EthAddress,
    pub eth_signature: EthSignature,
}

impl BatchConfirmResponse {
    pub fn from_proto(input: peggy_proto::peggy::MsgConfirmBatch) -> Result<Self, PeggyError> {
        Ok(BatchConfirmResponse {
            nonce: input.nonce,
            validator: input.validator.parse()?,
            token_contract: input.token_contract.parse()?,
            ethereum_signer: input.eth_signer.parse()?,
            eth_signature: input.signature.parse()?,
        })
    }
}

/// a list of validators, powers, and eth addresses at a given block height
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
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
    /// TODO give the signatures types a trait and de-duplicate
    pub fn order_valset_sigs(
        &self,
        signatures: &[ValsetConfirmResponse],
    ) -> Result<Vec<PeggySignature>, JsonRpcError> {
        let mut out = Vec::new();
        let mut valset_members = HashMap::new();
        // hashsets used for their efficient set theory types
        let mut valset_hashset = HashSet::new();
        let mut signatures_hashset = HashSet::new();
        for member in self.members.iter() {
            if let Some(address) = member.eth_address {
                valset_members.insert(address, member);
                valset_hashset.insert(address);
            } else {
                error!("Validator without set EthKey! Not included in Peggy Valset Update");
            }
        }
        for member in signatures {
            signatures_hashset.insert(member.eth_address);
        }
        // the validators who are in the valset, but not in the signatures list. Note this is order dependent
        // symmetric_difference or difference in a different order would include people who submitted signatures
        // but are not part of the validator set (we don't care about these people at all)
        // We need the validators who are in the set but have not submitted a signature in order to insert them
        // with zeroed out signatures, the contract will not accept a submission missing a signature otherwise.
        let validators_who_did_not_sign = valset_hashset.difference(&signatures_hashset);
        for sig in signatures {
            if let Some(val) = valset_members.get(&sig.eth_address) {
                out.push(PeggySignature {
                    power: val.power,
                    eth_address: sig.eth_address,
                    v: sig.eth_signature.v.clone(),
                    r: sig.eth_signature.r.clone(),
                    s: sig.eth_signature.s.clone(),
                })
            } else {
                // someone who is not a valset member submitted
                // a signature, this is fine to ignore
                info!(
                    "No Match for sig probably non-validator submitting a signature! {} and {}",
                    sig.eth_address,
                    ValsetMember::display_vec(&self.members)
                );
            }
        }
        for val in validators_who_did_not_sign {
            out.push(PeggySignature {
                // in order to be in the valset_hashset an address must both in valset_members and have a set
                // eth address, therefore we can disregard error handling here and do direct lookups
                power: valset_members[val].power,
                eth_address: valset_members[val].eth_address.unwrap(),
                v: 0u8.into(),
                r: 0u8.into(),
                s: 0u8.into(),
            });
        }
        // sort by power so that it is accepted by the contract
        out.sort();
        // go code sorts descending, rust sorts ascending, annoying
        out.reverse();

        Ok(out)
    }

    /// combines the provided signatures with the valset ensuring that ordering and signature data is correct
    /// TODO give the signatures types a trait and de-duplicate
    pub fn order_batch_sigs(
        &self,
        signatures: &[BatchConfirmResponse],
    ) -> Result<Vec<PeggySignature>, JsonRpcError> {
        let mut out = Vec::new();
        let mut valset_members = HashMap::new();
        // hashsets used for their efficient set theory types
        let mut valset_hashset = HashSet::new();
        let mut signatures_hashset = HashSet::new();
        for member in self.members.iter() {
            if let Some(address) = member.eth_address {
                valset_members.insert(address, member);
                valset_hashset.insert(address);
            } else {
                error!("Validator without set EthKey! Not included in Peggy Valset Update");
            }
        }
        for member in signatures {
            signatures_hashset.insert(member.ethereum_signer);
        }
        // the validators who are in the valset, but not in the signatures list. Note this is order dependent
        // symmetric_difference or difference in a different order would include people who submitted signatures
        // but are not part of the validator set (we don't care about these people at all)
        // We need the validators who are in the set but have not submitted a signature in order to insert them
        // with zeroed out signatures, the contract will not accept a submission missing a signature otherwise.
        let validators_who_did_not_sign = valset_hashset.difference(&signatures_hashset);
        for sig in signatures {
            if let Some(val) = valset_members.get(&sig.ethereum_signer) {
                out.push(PeggySignature {
                    power: val.power,
                    eth_address: sig.ethereum_signer,
                    v: sig.eth_signature.v.clone(),
                    r: sig.eth_signature.r.clone(),
                    s: sig.eth_signature.s.clone(),
                })
            } else {
                // someone who is not a valset member submitted
                // a signature, this is fine to ignore
                info!(
                    "No Match for sig probably non-validator submitting a signature! {} and {}",
                    sig.ethereum_signer,
                    ValsetMember::display_vec(&self.members)
                );
            }
        }
        for val in validators_who_did_not_sign {
            out.push(PeggySignature {
                // in order to be in the valset_hashset an address must both in valset_members and have a set
                // eth address, therefore we can disregard error handling here and do direct lookups
                power: valset_members[val].power,
                eth_address: valset_members[val].eth_address.unwrap(),
                v: 0u8.into(),
                r: 0u8.into(),
                s: 0u8.into(),
            });
        }
        // sort by power so that it is accepted by the contract
        out.sort();
        // go code sorts descending, rust sorts ascending, annoying
        out.reverse();

        Ok(out)
    }

    /// A utility function to provide a HashMap of members for easy lookups
    pub fn to_hashmap(&self) -> HashMap<EthAddress, u64> {
        let mut res = HashMap::new();
        for item in self.members.iter() {
            if let Some(address) = item.eth_address {
                res.insert(address, item.power);
            } else {
                panic!("Validator in active set without Eth Address! This must be corrected immediately!")
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
                panic!("Validator in active set without Eth Address! This must be corrected immediately!")
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
    /// Returned value must be less than one
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

        // if this is true then something has failed in the Cosmos module. Power is supposed to be allocated by dividing
        // between members at a resolution of u32 MAX anything greater than this value risks proposals passing with less
        // than the desired amount of power. For example if the Cosmos module switched to using u64 max as the cap but the
        // contract stayed the same (it always will without being redeployed the 'proposal pass' value is hardcoded on deploy)
        // then a vote could pass with less than 1% of all power.
        if total_power_diff > u32::MAX.into() {
            panic!("Power in bridge greater than u32 MAX! Bridge may be open to highjacking! Take action immediately!");
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
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq)]
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
    fn display_vec(input: &[ValsetMember]) -> String {
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

/// a list of validators, powers, and eth addresses at a given block height
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetMemberUnparsed {
    ethereum_address: String,
    #[serde(deserialize_with = "parse_val")]
    power: u64,
}

/// a list of validators, powers, and eth addresses at a given block height
/// this version is used by the endpoint to get the data and is then processed
/// by "convert" into ValsetResponse. Making this struct purely internal
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct ValsetUnparsed {
    #[serde(deserialize_with = "parse_val")]
    nonce: u64,
    members: Vec<ValsetMemberUnparsed>,
}

impl ValsetUnparsed {
    pub fn convert(self) -> Valset {
        let mut out = Vec::new();
        for member in self.members {
            if member.ethereum_address.is_empty() {
                out.push(ValsetMember {
                    power: member.power,
                    eth_address: None,
                });
            } else {
                match member.ethereum_address.parse() {
                    Ok(val) => out.push(ValsetMember {
                        power: member.power,
                        eth_address: Some(val),
                    }),
                    Err(_e) => out.push(ValsetMember {
                        power: member.power,
                        eth_address: None,
                    }),
                }
            }
        }
        Valset {
            nonce: self.nonce,
            members: out,
        }
    }
}

/// the query struct required to get the valset request sent by a specific
/// validator. This is required because the url encoded get methods don't
/// parse addresses well. So there's no way to get an individual validators
/// address without sending over a json body
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct QueryValsetConfirm {
    pub nonce: String,
    pub address: String,
}
