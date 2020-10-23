//! Helper functions for sending tokens to Cosmos

use std::time::Duration;

use clarity::abi::{encode_call, Token};
use clarity::PrivateKey as EthPrivateKey;
use clarity::{Address, Uint256};
use deep_space::address::Address as CosmosAddress;
use peggy_utils::error::OrchestratorError;
use web30::client::Web3;
use web30::types::SendTxOption;

const SEND_TO_COSMOS_GAS_LIMIT: u128 = 100_000;

pub async fn send_to_cosmos(
    erc20: Address,
    peggy_contract: Address,
    amount: Uint256,
    cosmos_destination: CosmosAddress,
    sender_secret: EthPrivateKey,
    wait_timeout: Option<Duration>,
    web3: &Web3,
    options: Vec<SendTxOption>,
) -> Result<Uint256, OrchestratorError> {
    let sender_address = sender_secret.to_public_key()?;
    let approved = web3
        .check_erc20_approved(erc20, sender_address, peggy_contract)
        .await?;
    if !approved {
        let txid = web3
            .approve_erc20_transfers(erc20, sender_secret, peggy_contract, None, options.clone())
            .await?;
        info!(
            "We are not approved for ERC20 transfers, approving txid: {:#066x}",
            txid
        );
        if let Some(timeout) = wait_timeout {
            web3.wait_for_transaction(txid, timeout, None).await?;
            info!("Approval finished!")
        }
    }

    // if the user sets a gas limit we should honor it, if they don't we
    // should add the default
    let mut has_gas_limit = false;
    let mut options = options;
    for option in options.iter() {
        if let SendTxOption::GasLimit(_) = option {
            has_gas_limit = true;
            break;
        }
    }
    if !has_gas_limit {
        options.push(SendTxOption::GasLimit(SEND_TO_COSMOS_GAS_LIMIT.into()));
    }

    // TODO figure out bigendian encoding here
    let encoded_destination_address = Token::Bytes(cosmos_destination.as_bytes().to_vec());

    let tx_hash = web3
        .send_transaction(
            peggy_contract,
            encode_call(
                "sendToCosmos(address,bytes32,uint256)",
                &[
                    erc20.into(),
                    encoded_destination_address,
                    amount.clone().into(),
                ],
            )?,
            0u32.into(),
            sender_address,
            sender_secret,
            options,
        )
        .await?;

    if let Some(timeout) = wait_timeout {
        web3.wait_for_transaction(tx_hash.clone(), timeout, None)
            .await?;
    }

    Ok(tx_hash)
}
