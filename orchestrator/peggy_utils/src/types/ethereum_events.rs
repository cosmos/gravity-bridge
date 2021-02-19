use std::unimplemented;

use super::ValsetMember;
use crate::error::PeggyError;
use clarity::Address as EthAddress;
use deep_space::address::Address as CosmosAddress;
use num256::Uint256;
use web30::types::Log;

/// A parsed struct representing the Ethereum event fired by the Peggy contract
/// when the validator set is updated.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ValsetUpdatedEvent {
    pub nonce: u64,
    pub members: Vec<ValsetMember>,
}

impl ValsetUpdatedEvent {
    /// This function is not an abi compatible bytes parser, but it's actually
    /// not hard at all to extract data like this by hand.
    pub fn from_log(input: &Log) -> Result<ValsetUpdatedEvent, PeggyError> {
        // we have one indexed event so we should fine two indexes, one the event itself
        // and one the indexed nonce
        if input.topics.get(1).is_none() {
            return Err(PeggyError::InvalidEventLogError(
                "Too few topics".to_string(),
            ));
        }
        let nonce_data = &input.topics[1];
        let nonce = Uint256::from_bytes_be(nonce_data);
        if nonce > u64::MAX.into() {
            return Err(PeggyError::InvalidEventLogError(
                "Nonce overflow, probably incorrect parsing".to_string(),
            ));
        }
        let nonce: u64 = nonce.to_string().parse().unwrap();
        // first two indexes contain event info we don't care about, third index is
        // the length of the eth addresses array
        let index_start = 2 * 32;
        let index_end = index_start + 32;
        let eth_addresses_offset = index_start + 32;
        let len_eth_addresses = Uint256::from_bytes_be(&input.data[index_start..index_end]);
        if len_eth_addresses > usize::MAX.into() {
            return Err(PeggyError::InvalidEventLogError(
                "Ethereum array len overflow, probably incorrect parsing".to_string(),
            ));
        }
        let len_eth_addresses: usize = len_eth_addresses.to_string().parse().unwrap();
        let index_start = (3 + len_eth_addresses) * 32;
        let index_end = index_start + 32;
        let powers_offset = index_start + 32;
        let len_powers = Uint256::from_bytes_be(&input.data[index_start..index_end]);
        if len_powers > usize::MAX.into() {
            return Err(PeggyError::InvalidEventLogError(
                "Powers array len overflow, probably incorrect parsing".to_string(),
            ));
        }
        let len_powers: usize = len_eth_addresses.to_string().parse().unwrap();
        if len_powers != len_eth_addresses {
            return Err(PeggyError::InvalidEventLogError(
                "Array len mismatch, probably incorrect parsing".to_string(),
            ));
        }

        let mut validators = Vec::new();
        for i in 0..len_eth_addresses {
            let power_start = (i * 32) + powers_offset;
            let power_end = power_start + 32;
            let address_start = (i * 32) + eth_addresses_offset;
            let address_end = address_start + 32;
            let power = Uint256::from_bytes_be(&input.data[power_start..power_end]);
            // an eth address at 20 bytes is 12 bytes shorter than the Uint256 it's stored in.
            let eth_address = EthAddress::from_slice(&input.data[address_start + 12..address_end]);
            if eth_address.is_err() {
                return Err(PeggyError::InvalidEventLogError(
                    "Ethereum Address parsing error, probably incorrect parsing".to_string(),
                ));
            }
            let eth_address = Some(eth_address.unwrap());
            if power > u64::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "Power greater than u64::MAX, probably incorrect parsing".to_string(),
                ));
            }
            let power: u64 = power.to_string().parse().unwrap();
            validators.push(ValsetMember { power, eth_address })
        }
        let mut check = validators.clone();
        check.sort();
        check.reverse();
        // if the validator set is not sorted we're in a bad spot
        if validators != check {
            trace!(
                "Someone submitted an unsorted validator set, this means all updates will fail until someone feeds in this unsorted value by hand {:?} instead of {:?}",
                validators, check
            );
        }

        Ok(ValsetUpdatedEvent {
            nonce,
            members: validators,
        })
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<ValsetUpdatedEvent>, PeggyError> {
        let mut res = Vec::new();
        for item in input {
            res.push(ValsetUpdatedEvent::from_log(item)?);
        }
        Ok(res)
    }
}

/// A parsed struct representing the Ethereum event fired by the Peggy contract when
/// a transaction batch is executed.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct TransactionBatchExecutedEvent {
    /// the nonce attached to the transaction batch that follows
    /// it throughout it's lifecycle
    pub batch_nonce: Uint256,
    /// The block height this event occurred at
    pub block_height: Uint256,
    /// The ERC20 token contract address for the batch executed, since batches are uniform
    /// in token type there is only one
    pub erc20: EthAddress,
    /// the event nonce representing a unique ordering of events coming out
    /// of the Peggy solidity contract. Ensuring that these events can only be played
    /// back in order
    pub event_nonce: Uint256,
}

impl TransactionBatchExecutedEvent {
    pub fn from_log(input: &Log) -> Result<TransactionBatchExecutedEvent, PeggyError> {
        if let (Some(batch_nonce_data), Some(erc20_data)) =
            (input.topics.get(1), input.topics.get(2))
        {
            let batch_nonce = Uint256::from_bytes_be(batch_nonce_data);
            let erc20 = EthAddress::from_slice(&erc20_data[12..32])?;
            let event_nonce = Uint256::from_bytes_be(&input.data);
            let block_height = if let Some(bn) = input.block_number.clone() {
                bn
            } else {
                return Err(PeggyError::InvalidEventLogError(
                    "Log does not have block number, we only search logs already in blocks?"
                        .to_string(),
                ));
            };
            if event_nonce > u64::MAX.into()
                || batch_nonce > u64::MAX.into()
                || block_height > u64::MAX.into()
            {
                Err(PeggyError::InvalidEventLogError(
                    "Event nonce overflow, probably incorrect parsing".to_string(),
                ))
            } else {
                Ok(TransactionBatchExecutedEvent {
                    batch_nonce,
                    block_height,
                    erc20,
                    event_nonce,
                })
            }
        } else {
            Err(PeggyError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<TransactionBatchExecutedEvent>, PeggyError> {
        let mut res = Vec::new();
        for item in input {
            res.push(TransactionBatchExecutedEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    pub fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce.into() {
                ret.push(item.clone())
            }
        }
        ret
    }
}

/// A parsed struct representing the Ethereum event fired when someone makes a deposit
/// on the Peggy contract
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct SendToCosmosEvent {
    /// The token contract address for the deposit
    pub erc20: EthAddress,
    /// The Ethereum Sender
    pub sender: EthAddress,
    /// The Cosmos destination
    pub destination: CosmosAddress,
    /// The amount of the erc20 token that is being sent
    pub amount: Uint256,
    /// The transaction's nonce, used to make sure there can be no accidental duplication
    pub event_nonce: Uint256,
    /// The block height this event occurred at
    pub block_height: Uint256,
}

impl SendToCosmosEvent {
    pub fn from_log(input: &Log) -> Result<SendToCosmosEvent, PeggyError> {
        let topics = (
            input.topics.get(1),
            input.topics.get(2),
            input.topics.get(3),
        );
        if let (Some(erc20_data), Some(sender_data), Some(destination_data)) = topics {
            let erc20 = EthAddress::from_slice(&erc20_data[12..32])?;
            let sender = EthAddress::from_slice(&sender_data[12..32])?;
            // this is required because deep_space requires a fixed length slice to
            // create an address from bytes.
            let mut c_address_bytes: [u8; 20] = [0; 20];
            c_address_bytes.copy_from_slice(&destination_data[12..32]);
            let destination = CosmosAddress::from_bytes(c_address_bytes);
            let amount = Uint256::from_bytes_be(&input.data[..32]);
            let event_nonce = Uint256::from_bytes_be(&input.data[32..]);
            let block_height = if let Some(bn) = input.block_number.clone() {
                bn
            } else {
                return Err(PeggyError::InvalidEventLogError(
                    "Log does not have block number, we only search logs already in blocks?"
                        .to_string(),
                ));
            };
            if event_nonce > u64::MAX.into() || block_height > u64::MAX.into() {
                Err(PeggyError::InvalidEventLogError(
                    "Event nonce overflow, probably incorrect parsing".to_string(),
                ))
            } else {
                Ok(SendToCosmosEvent {
                    erc20,
                    sender,
                    destination,
                    amount,
                    event_nonce,
                    block_height,
                })
            }
        } else {
            Err(PeggyError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<SendToCosmosEvent>, PeggyError> {
        let mut res = Vec::new();
        for item in input {
            res.push(SendToCosmosEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    pub fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce.into() {
                ret.push(item.clone())
            }
        }
        ret
    }
}

/// A parsed struct representing the Ethereum event fired when someone uses the Peggy
/// contract to deploy a new ERC20 contract representing a Cosmos asset
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20DeployedEvent {
    /// The denom on the Cosmos chain this contract is intended to represent
    pub cosmos_denom: String,
    /// The ERC20 address of the deployed contract, this may or may not be adopted
    /// by the Cosmos chain as the contract for this asset
    pub erc20_address: EthAddress,
    /// The name of the token in the ERC20 contract, should match the Cosmos denom
    /// but it is up to the Cosmos module to check that
    pub name: String,
    /// The symbol for the token in the ERC20 contract
    pub symbol: String,
    /// The number of decimals required to represent the smallest unit of this token
    pub decimals: u8,
    pub event_nonce: Uint256,
    pub block_height: Uint256,
}

impl ERC20DeployedEvent {
    pub fn from_log(input: &Log) -> Result<ERC20DeployedEvent, PeggyError> {
        let token_contract = input.topics.get(1);
        if let Some(new_token_contract_data) = token_contract {
            let erc20 = EthAddress::from_slice(&new_token_contract_data[12..32])?;
            let index_start = 3 * 32;
            let index_end = index_start + 32;
            let decimals = Uint256::from_bytes_be(&input.data[index_start..index_end]);
            if decimals > u8::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "Decimals overflow, probably incorrect parsing".to_string(),
                ));
            }
            let decimals: u8 = decimals.to_string().parse().unwrap();

            let index_start = 4 * 32;
            let index_end = index_start + 32;
            let nonce = Uint256::from_bytes_be(&input.data[index_start..index_end]);
            if nonce > u64::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "Nonce overflow, probably incorrect parsing".to_string(),
                ));
            }

            let index_start = 5 * 32;
            let index_end = index_start + 32;
            let denom_len = Uint256::from_bytes_be(&input.data[index_start..index_end]);
            // it's not probable that we have 4+ gigabytes of event data
            if denom_len > u32::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "denom length overflow, probably incorrect parsing".to_string(),
                ));
            }
            let denom_len: usize = denom_len.to_string().parse().unwrap();
            let index_start = 6 * 32;
            let index_end = index_start + denom_len;
            let denom = String::from_utf8(input.data[index_start..index_end].to_vec());
            trace!("Denom {:?}", denom);
            if denom.is_err() {
                return Err(PeggyError::InvalidEventLogError(format!(
                    "{:?} is not valid utf8, probably incorrect parsing",
                    denom
                )));
            }
            let denom = denom.unwrap();

            // beyond this point we are parsing strings placed
            // after a variable length string and we will need to compute offsets

            // this trick computes the next 32 byte (256 bit) word index, then multiplies by
            // 32 to get the bytes offset, this is required since we have dynamic length types but
            // the next entry always starts on a round 32 byte word.
            let index_start = ((index_end + 31) / 32) * 32;
            let index_end = index_start + 32;
            let erc20_name_len = Uint256::from_bytes_be(&input.data[index_start..index_end]);
            // it's not probable that we have 4+ gigabytes of event data
            if erc20_name_len > u32::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "ERC20 Name length overflow, probably incorrect parsing".to_string(),
                ));
            }
            let erc20_name_len: usize = erc20_name_len.to_string().parse().unwrap();
            let index_start = index_end;
            let index_end = index_start + erc20_name_len;
            let erc20_name = String::from_utf8(input.data[index_start..index_end].to_vec());
            if erc20_name.is_err() {
                return Err(PeggyError::InvalidEventLogError(format!(
                    "{:?} is not valid utf8, probably incorrect parsing",
                    erc20_name
                )));
            }
            trace!("ERC20 Name {:?}", erc20_name);
            let erc20_name = erc20_name.unwrap();

            let index_start = ((index_end + 31) / 32) * 32;
            let index_end = index_start + 32;
            let symbol_len = Uint256::from_bytes_be(&input.data[index_start..index_end]);
            // it's not probable that we have 4+ gigabytes of event data
            if symbol_len > u32::MAX.into() {
                return Err(PeggyError::InvalidEventLogError(
                    "Symbol length overflow, probably incorrect parsing".to_string(),
                ));
            }
            let symbol_len: usize = symbol_len.to_string().parse().unwrap();
            let index_start = index_end;
            let index_end = index_start + symbol_len;
            let symbol = String::from_utf8(input.data[index_start..index_end].to_vec());
            trace!("Symbol {:?}", symbol);
            if symbol.is_err() {
                return Err(PeggyError::InvalidEventLogError(format!(
                    "{:?} is not valid utf8, probably incorrect parsing",
                    symbol
                )));
            }
            let symbol = symbol.unwrap();

            let block_height = if let Some(bn) = input.block_number.clone() {
                bn
            } else {
                return Err(PeggyError::InvalidEventLogError(
                    "Log does not have block number, we only search logs already in blocks?"
                        .to_string(),
                ));
            };

            Ok(ERC20DeployedEvent {
                cosmos_denom: denom,
                name: erc20_name,
                decimals,
                event_nonce: nonce,
                erc20_address: erc20,
                symbol,
                block_height,
            })
        } else {
            Err(PeggyError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<ERC20DeployedEvent>, PeggyError> {
        let mut res = Vec::new();
        for item in input {
            res.push(ERC20DeployedEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    pub fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce.into() {
                ret.push(item.clone())
            }
        }
        ret
    }
}
/// A parsed struct representing the Ethereum event fired when someone uses the Peggy
/// contract to deploy a new ERC20 contract representing a Cosmos asset
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct LogicCallExecutedEvent {
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: Uint256,
    pub return_data: Vec<u8>,
    pub event_nonce: Uint256,
    pub block_height: Uint256,
}

impl LogicCallExecutedEvent {
    pub fn from_log(input: &Log) -> Result<LogicCallExecutedEvent, PeggyError> {
        unimplemented!()
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<LogicCallExecutedEvent>, PeggyError> {
        let mut res = Vec::new();
        for item in input {
            res.push(LogicCallExecutedEvent::from_log(item)?);
        }
        Ok(res)
    }
    /// returns all values in the array with event nonces greater
    /// than the provided value
    pub fn filter_by_event_nonce(event_nonce: u64, input: &[Self]) -> Vec<Self> {
        let mut ret = Vec::new();
        for item in input {
            if item.event_nonce > event_nonce.into() {
                ret.push(item.clone())
            }
        }
        ret
    }
}
