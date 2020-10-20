//! for things that don't belong in the cosmos or ethereum libraries but also don't belong
//! in a function specific library

use contact::jsonrpc::error::JsonRpcError;
use std::fmt::{self, Debug};
use web30::jsonrpc::error::Web3Error;

#[derive(Debug)]
pub enum OrchestratorError {
    CosmosRestErr(JsonRpcError),
    EthereumRestErr(Web3Error),
    InvalidBridgeStateError(String),
    FailedToUpdateValset,
}

impl fmt::Display for OrchestratorError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            OrchestratorError::CosmosRestErr(val) => write!(f, "Cosmos REST error {}", val),
            OrchestratorError::EthereumRestErr(val) => write!(f, "Ethereum REST error {}", val),
            OrchestratorError::InvalidBridgeStateError(val) => {
                write!(f, "Invalid bridge state! {}", val)
            }
            OrchestratorError::FailedToUpdateValset => write!(f, "ValidatorSetUpdate Failed!"),
        }
    }
}

impl std::error::Error for OrchestratorError {}

impl From<JsonRpcError> for OrchestratorError {
    fn from(error: JsonRpcError) -> Self {
        OrchestratorError::CosmosRestErr(error)
    }
}

impl From<Web3Error> for OrchestratorError {
    fn from(error: Web3Error) -> Self {
        OrchestratorError::EthereumRestErr(error)
    }
}
