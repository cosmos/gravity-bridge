//! for things that don't belong in the cosmos or ethereum libraries but also don't belong
//! in a function specific library

use clarity::Error as ClarityError;
use contact::jsonrpc::error::JsonRpcError;
use deep_space::address::AddressError as CosmosAddressError;
use num_bigint::ParseBigIntError;
use std::fmt::{self, Debug};
use tokio::time::Elapsed;
use tonic::Status;
use web30::jsonrpc::error::Web3Error;

#[derive(Debug)]
pub enum PeggyError {
    InvalidBigInt(ParseBigIntError),
    CosmosRestError(JsonRpcError),
    CosmosAddressError(CosmosAddressError),
    EthereumRestError(Web3Error),
    InvalidBridgeStateError(String),
    FailedToUpdateValset,
    EthereumContractError(String),
    InvalidOptionsError(String),
    ClarityError(ClarityError),
    TimeoutError,
    InvalidEventLogError(String),
    CosmosgRPCError(Status),
    InsufficientVotingPowerToPass(String),
}

impl fmt::Display for PeggyError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            PeggyError::CosmosgRPCError(val) => write!(f, "Cosmos gRPC error {}", val),
            PeggyError::InvalidBigInt(val) => write!(f, "Got invalid BigInt from cosmos! {}", val),
            PeggyError::CosmosRestError(val) => write!(f, "Cosmos REST error {}", val),
            PeggyError::CosmosAddressError(val) => write!(f, "Cosmos Address error {}", val),
            PeggyError::EthereumRestError(val) => write!(f, "Ethereum REST error {}", val),
            PeggyError::InvalidOptionsError(val) => {
                write!(f, "Invalid TX options for this call {}", val)
            }
            PeggyError::InvalidBridgeStateError(val) => {
                write!(f, "Invalid bridge state! {}", val)
            }
            PeggyError::FailedToUpdateValset => write!(f, "ValidatorSetUpdate Failed!"),
            PeggyError::TimeoutError => write!(f, "Operation timed out!"),
            PeggyError::ClarityError(val) => write!(f, "Clarity Error {}", val),
            PeggyError::InvalidEventLogError(val) => write!(f, "InvalidEvent: {}", val),
            PeggyError::EthereumContractError(val) => {
                write!(f, "Contract operation failed: {}", val)
            }
            PeggyError::InsufficientVotingPowerToPass(val) => {
                write!(f, "{}", val)
            }
        }
    }
}

impl std::error::Error for PeggyError {}

impl From<JsonRpcError> for PeggyError {
    fn from(error: JsonRpcError) -> Self {
        PeggyError::CosmosRestError(error)
    }
}

impl From<Elapsed> for PeggyError {
    fn from(_error: Elapsed) -> Self {
        PeggyError::TimeoutError
    }
}

impl From<ClarityError> for PeggyError {
    fn from(error: ClarityError) -> Self {
        PeggyError::ClarityError(error)
    }
}

impl From<Web3Error> for PeggyError {
    fn from(error: Web3Error) -> Self {
        PeggyError::EthereumRestError(error)
    }
}
impl From<Status> for PeggyError {
    fn from(error: Status) -> Self {
        PeggyError::CosmosgRPCError(error)
    }
}
impl From<CosmosAddressError> for PeggyError {
    fn from(error: CosmosAddressError) -> Self {
        PeggyError::CosmosAddressError(error)
    }
}
impl From<ParseBigIntError> for PeggyError {
    fn from(error: ParseBigIntError) -> Self {
        PeggyError::InvalidBigInt(error)
    }
}
