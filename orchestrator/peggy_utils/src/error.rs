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
pub enum GravityError {
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
    ParseBigIntError(ParseBigIntError),
}

impl fmt::Display for GravityError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            GravityError::CosmosgRPCError(val) => write!(f, "Cosmos gRPC error {}", val),
            GravityError::InvalidBigInt(val) => write!(f, "Got invalid BigInt from cosmos! {}", val),
            GravityError::CosmosRestError(val) => write!(f, "Cosmos REST error {}", val),
            GravityError::CosmosAddressError(val) => write!(f, "Cosmos Address error {}", val),
            GravityError::EthereumRestError(val) => write!(f, "Ethereum REST error {}", val),
            GravityError::InvalidOptionsError(val) => {
                write!(f, "Invalid TX options for this call {}", val)
            }
            GravityError::InvalidBridgeStateError(val) => {
                write!(f, "Invalid bridge state! {}", val)
            }
            GravityError::FailedToUpdateValset => write!(f, "ValidatorSetUpdate Failed!"),
            GravityError::TimeoutError => write!(f, "Operation timed out!"),
            GravityError::ClarityError(val) => write!(f, "Clarity Error {}", val),
            GravityError::InvalidEventLogError(val) => write!(f, "InvalidEvent: {}", val),
            GravityError::EthereumContractError(val) => {
                write!(f, "Contract operation failed: {}", val)
            }
            GravityError::InsufficientVotingPowerToPass(val) => {
                write!(f, "{}", val)
            }
            GravityError::ParseBigIntError(val) => write!(f, "Failed to parse big integer {}", val),
        }
    }
}

impl std::error::Error for GravityError {}

impl From<JsonRpcError> for GravityError {
    fn from(error: JsonRpcError) -> Self {
        GravityError::CosmosRestError(error)
    }
}

impl From<Elapsed> for GravityError {
    fn from(_error: Elapsed) -> Self {
        GravityError::TimeoutError
    }
}

impl From<ClarityError> for GravityError {
    fn from(error: ClarityError) -> Self {
        GravityError::ClarityError(error)
    }
}

impl From<Web3Error> for GravityError {
    fn from(error: Web3Error) -> Self {
        GravityError::EthereumRestError(error)
    }
}
impl From<Status> for GravityError {
    fn from(error: Status) -> Self {
        GravityError::CosmosgRPCError(error)
    }
}
impl From<CosmosAddressError> for GravityError {
    fn from(error: CosmosAddressError) -> Self {
        GravityError::CosmosAddressError(error)
    }
}
impl From<ParseBigIntError> for GravityError {
    fn from(error: ParseBigIntError) -> Self {
        GravityError::InvalidBigInt(error)
    }
}
