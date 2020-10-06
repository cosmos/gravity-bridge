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

/// TODO modify code in web30 if this works for some reason the
/// geth node for the testnet seems convinced that we need to provide
/// gas
pub async fn contract_call(
    web30: &Web3,
    contract_address: EthAddress,
    sig: &str,
    tokens: &[Token],
    own_address: EthAddress,
) -> Result<Vec<u8>, Web3Error> {
    let gas_price = match web30.eth_gas_price().await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };

    let nonce = match web30.eth_get_transaction_count(own_address).await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };

    let payload = encode_call(sig, tokens).unwrap();

    let transaction = TransactionRequest {
        from: Some(own_address),
        to: contract_address,
        nonce: Some(UnpaddedHex(nonce)),
        gas: Some(UnpaddedHex(1_000_000u64.into())),
        gas_price: Some(UnpaddedHex(gas_price)),
        value: Some(UnpaddedHex(0u64.into())),
        data: Some(Data(payload)),
    };

    let bytes = match web30.eth_call(transaction).await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };
    Ok(bytes.0)
}

pub fn to_uint_vec(input: &[u64]) -> Vec<Uint256> {
    let mut new_vec = Vec::new();
    for value in input {
        let v: u64 = *value;
        new_vec.push(v.into())
    }
    new_vec
}
