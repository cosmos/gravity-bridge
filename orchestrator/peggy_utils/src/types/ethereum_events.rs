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
            error!(
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
            if event_nonce > u64::MAX.into() || batch_nonce > u64::MAX.into() {
                Err(PeggyError::InvalidEventLogError(
                    "Event nonce overflow, probably incorrect parsing".to_string(),
                ))
            } else {
                Ok(TransactionBatchExecutedEvent {
                    batch_nonce,
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
    /// The transaction's nonce, used to make sure there can be no accidntal duplication
    pub event_nonce: Uint256,
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
            if event_nonce > u64::MAX.into() {
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
