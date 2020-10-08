//! This crate contains various components and utilities for interacting with the Peggy Ethereum contract.

#[macro_use]
extern crate log;

pub mod utils;
pub mod valset_update;

use clarity::abi::encode_call;
use clarity::abi::Token;
use clarity::Address as EthAddress;
use clarity::Uint256;
use web30::jsonrpc::error::Web3Error;
use web30::types::TransactionRequest;
use web30::{client::Web3, types::Data, types::UnpaddedHex};

/// takes a valset and signatures as input and returns a sorted
/// array of each in order of descending validator power with the
/// appropriate signatures lined up correctly
fn prepare_sigs() {}

struct OrderedSignatures {
    addresses: Vec<EthAddress>,
    powers: Vec<u64>,
    v: Vec<Uint256>,
    r: Vec<Uint256>,
    s: Vec<Uint256>,
}

pub fn to_uint_vec(input: &[u64]) -> Vec<Uint256> {
    let mut new_vec = Vec::new();
    for value in input {
        let v: u64 = *value;
        new_vec.push(v.into())
    }
    new_vec
}
