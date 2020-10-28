use clarity::Address as EthAddress;
use deep_space::address::Address as CosmosAddress;
use num256::Uint256;
use web30::types::Log;

use crate::error::OrchestratorError;

/// A parsed struct representing the Ethereum event fired by the Peggy contract
/// when the validator set is updated.
pub struct ValsetUpdatedEvent {
    pub nonce: Uint256,
    // we currently don't parse members, but the data is there
    //members: Vec<ValsetMember>,
}

impl ValsetUpdatedEvent {
    pub fn from_log(input: &Log) -> Result<ValsetUpdatedEvent, OrchestratorError> {
        // we have one indexed event so we should fine two indexes, one the event itself
        // and one the indexed nonce
        if let Some(nonce_data) = input.topics.get(1) {
            let nonce = Uint256::from_bytes_be(nonce_data);
            Ok(ValsetUpdatedEvent { nonce })
        } else {
            Err(OrchestratorError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<ValsetUpdatedEvent>, OrchestratorError> {
        let mut res = Vec::new();
        for item in input {
            res.push(ValsetUpdatedEvent::from_log(item)?);
        }
        Ok(res)
    }
}

/// A parsed struct representing the Ethereum event fired by the Peggy contract when
/// a transaction batch is executed.
pub struct TransactionBatchExecutedEvent {
    pub nonce: Uint256,
    /// The ERC20 token contract address for the batch executed, since batches are uniform
    /// in token type there is only one
    pub erc20: EthAddress,
}

impl TransactionBatchExecutedEvent {
    pub fn from_log(input: &Log) -> Result<TransactionBatchExecutedEvent, OrchestratorError> {
        if let (Some(nonce_data), Some(erc20_data)) = (input.topics.get(1), input.topics.get(2)) {
            let nonce = Uint256::from_bytes_be(nonce_data);
            let erc20 = EthAddress::from_slice(&erc20_data)?;
            Ok(TransactionBatchExecutedEvent { nonce, erc20 })
        } else {
            Err(OrchestratorError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(
        input: &[Log],
    ) -> Result<Vec<TransactionBatchExecutedEvent>, OrchestratorError> {
        let mut res = Vec::new();
        for item in input {
            res.push(TransactionBatchExecutedEvent::from_log(item)?);
        }
        Ok(res)
    }
}

/// A parsed struct representing the Ethereum event fired when someone makes a deposit
/// on the Peggy contract
pub struct SendToCosmosEvent {
    /// The token contract address for the deposit
    pub erc20: EthAddress,
    /// The Ethereum Sender
    pub sender: EthAddress,
    /// The Cosmos destination
    pub destination: CosmosAddress,
    /// The amount of the erc20 token that is being sent
    pub amount: Uint256,
    /// The transaction's nonce, every event from the contract gets a unique nonce that forces
    /// the oracle stream to be in order and consistent
    pub nonce: Uint256,
}

impl SendToCosmosEvent {
    pub fn from_log(input: &Log) -> Result<SendToCosmosEvent, OrchestratorError> {
        let topics = (
            input.topics.get(1),
            input.topics.get(2),
            input.topics.get(3),
        );
        if let (Some(erc20_data), Some(sender_data), Some(destination_data)) = topics {
            let erc20 = EthAddress::from_slice(&erc20_data[12..32])?;
            let sender = EthAddress::from_slice(&sender_data[12..32])?;
            let mut c_address_bytes: [u8; 20] = [0; 20];
            // this is little endian encoded
            c_address_bytes.copy_from_slice(&destination_data[0..20]);
            let destination = CosmosAddress::from_bytes(c_address_bytes);
            let amount = Uint256::from_bytes_be(&input.data[..32]);
            let nonce = Uint256::from_bytes_be(&input.data[32..]);
            Ok(SendToCosmosEvent {
                erc20,
                sender,
                destination,
                amount,
                nonce,
            })
        } else {
            Err(OrchestratorError::InvalidEventLogError(
                "Too few topics".to_string(),
            ))
        }
    }
    pub fn from_logs(input: &[Log]) -> Result<Vec<SendToCosmosEvent>, OrchestratorError> {
        let mut res = Vec::new();
        for item in input {
            res.push(SendToCosmosEvent::from_log(item)?);
        }
        Ok(res)
    }
}
