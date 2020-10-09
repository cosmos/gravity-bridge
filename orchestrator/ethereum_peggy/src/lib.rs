//! This crate contains various components and utilities for interacting with the Peggy Ethereum contract.

#[macro_use]
extern crate log;

pub mod utils;
pub mod valset_update;

use clarity::abi::encode_call;
use clarity::abi::Token;
use clarity::Address as EthAddress;
use clarity::Uint256;
use cosmos_peggy::types::Valset;
use peggy_utils::error::OrchestratorError;
use web30::jsonrpc::error::Web3Error;
use web30::types::TransactionRequest;
use web30::{client::Web3, types::Data, types::UnpaddedHex};

#[derive(PartialEq, Eq, PartialOrd, Ord)]
pub struct ValidatorPower {
    power: u64,
    eth_address: EthAddress,
}

impl ValidatorPower {
    /// Gets a sorted list of validator powers properly associated with each other
    pub fn get_powers_list(input: Valset) -> Result<Vec<ValidatorPower>, OrchestratorError> {
        let mut out = Vec::new();
        for (power, address) in input.powers.iter().zip(input.eth_addresses.iter()) {
            if let Some(eth_address) = address {
                out.push(ValidatorPower {
                    power: *power,
                    eth_address: *eth_address,
                });
            } else {
                return Err(OrchestratorError::InvalidBridgeStateError(
                    "Can't update valset with unset key".to_string(),
                ));
            }
        }
        out.sort();
        Ok(out)
    }
}

/// A sortable struct of a validator and it's signatures
/// there's some black magic here TODO we should implement
/// ORD ourselves because the order of this structs members below
/// determines what is compared first to produce an order. In this case
/// it's powers, then eth addresses
#[derive(PartialEq, Eq, PartialOrd, Ord)]
struct ValsetSignature {
    power: u64,
    eth_address: EthAddress,
    v: Uint256,
    r: Uint256,
    s: Uint256,
}
