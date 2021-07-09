use clarity::Address as EthAddress;
use clarity::Uint256;
use std::{time::Duration, u128};
use web30::{client::Web3, jsonrpc::error::Web3Error};

pub const TIMEOUT: Duration = Duration::from_secs(60);

/// TODO revisit this for higher precision while
/// still representing the number to the user as a float
/// this takes a number like 0.37 eth and turns it into wei
/// or any erc20 with arbitrary decimals
pub fn fraction_to_exponent(num: f64, exponent: u8) -> Uint256 {
    let mut res = num;
    // in order to avoid floating point rounding issues we
    // multiply only by 10 each time. this reduces the rounding
    // errors enough to be ignored
    for _ in 0..exponent {
        res *= 10f64
    }
    (res as u128).into()
}

pub async fn get_erc20_decimals(
    web3: &Web3,
    erc20: EthAddress,
    caller_address: EthAddress,
) -> Result<Uint256, Web3Error> {
    let decimals = web3
        .contract_call(erc20, "decimals()", &[], caller_address, None)
        .await?;

    Ok(Uint256::from_bytes_be(match decimals.get(0..32) {
        Some(val) => val,
        None => {
            return Err(Web3Error::ContractCallError(
                "Bad response from ERC20 decimals".to_string(),
            ))
        }
    }))
}
