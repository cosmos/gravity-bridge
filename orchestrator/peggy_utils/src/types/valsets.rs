use super::*;
use crate::error::PeggyError;
use clarity::Address as EthAddress;
use clarity::Signature as EthSignature;
use contact::{jsonrpc::error::JsonRpcError, types::parse_val};
use deep_space::address::Address as CosmosAddress;
use num256::Uint256;
use std::{cmp::Ordering, collections::HashMap, fmt};

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

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BatchConfirmResponse {
    #[serde(deserialize_with = "parse_val")]
    pub nonce: Uint256,
    #[serde(deserialize_with = "parse_val")]
    pub validator: CosmosAddress,
    #[serde(deserialize_with = "parse_val")]
    pub token_contract: EthAddress,
    #[serde(deserialize_with = "parse_val")]
    pub ethereum_signer: EthAddress,
    #[serde(deserialize_with = "parse_val", rename = "signature")]
    pub eth_signature: EthSignature,
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
    pub fn order_valset_sigs(
        &self,
        signatures: &[ValsetConfirmResponse],
    ) -> Result<Vec<PeggySignature>, JsonRpcError> {
        let mut out = Vec::new();
        let mut members = HashMap::new();
        for member in self.members.iter() {
            if let Some(address) = member.eth_address {
                members.insert(address, member);
            } else {
                return Err(JsonRpcError::BadInput(
                    "All Eth Addresses must be set".to_string(),
                ));
            }
        }
        for sig in signatures {
            if let Some(val) = members.get(&sig.eth_address) {
                out.push(PeggySignature {
                    power: val.power,
                    eth_address: sig.eth_address,
                    v: sig.eth_signature.v.clone(),
                    r: sig.eth_signature.r.clone(),
                    s: sig.eth_signature.s.clone(),
                })
            } else {
                return Err(JsonRpcError::BadInput(format!(
                    "No Match for sig! {} and {}",
                    sig.eth_address,
                    ValsetMember::display_vec(&self.members)
                )));
            }
        }
        // sort by power so that it is accepted by the contract
        out.sort();
        // go code sorts descending, rust sorts ascending, annoying
        out.reverse();

        Ok(out)
    }

    /// combines the provided signatures with the valset ensuring that ordering and signature data is correct
    pub fn order_batch_sigs(
        &self,
        signatures: &[BatchConfirmResponse],
    ) -> Result<Vec<PeggySignature>, JsonRpcError> {
        let mut out = Vec::new();
        let mut members = HashMap::new();
        for member in self.members.iter() {
            if let Some(address) = member.eth_address {
                members.insert(address, member);
            } else {
                return Err(JsonRpcError::BadInput(
                    "All Eth Addresses must be set".to_string(),
                ));
            }
        }
        for sig in signatures {
            if let Some(val) = members.get(&sig.ethereum_signer) {
                out.push(PeggySignature {
                    power: val.power,
                    eth_address: sig.ethereum_signer,
                    v: sig.eth_signature.v.clone(),
                    r: sig.eth_signature.r.clone(),
                    s: sig.eth_signature.s.clone(),
                })
            } else {
                return Err(JsonRpcError::BadInput(format!(
                    "No Match for sig! {} and {}",
                    sig.ethereum_signer,
                    ValsetMember::display_vec(&self.members)
                )));
            }
        }
        // sort by power so that it is accepted by the contract
        out.sort();
        // go code sorts descending, rust sorts ascending, annoying
        out.reverse();

        Ok(out)
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
