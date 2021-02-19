use crate::utils::{get_logic_call_nonce, GasCost};
use clarity::{abi::Token, utils::bytes_to_hex_str, PrivateKey as EthPrivateKey};
use clarity::{Address as EthAddress, Uint256};
use peggy_utils::error::PeggyError;
use peggy_utils::types::*;
use std::{cmp::min, time::Duration};
use web30::{client::Web3, types::TransactionRequest};

/// this function generates an appropriate Ethereum transaction
/// to submit the provided logic call
pub async fn send_eth_logic_call(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    web3: &Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<(), PeggyError> {
    let new_call_nonce = call.invalidation_nonce;
    let eth_address = our_eth_key.to_public_key().unwrap();
    info!(
        "Ordering signatures and submitting LogicCall {}:{} to Ethereum",
        bytes_to_hex_str(&call.invalidation_id),
        new_call_nonce
    );
    trace!("Call {:?}", call);

    let before_nonce = get_logic_call_nonce(
        peggy_contract_address,
        call.invalidation_id.clone(),
        eth_address,
        &web3,
    )
    .await?;
    let current_block_height = web3.eth_block_number().await?;
    if before_nonce >= new_call_nonce {
        info!(
            "Someone else updated the LogicCall to {}, exiting early",
            before_nonce
        );
        return Ok(());
    } else if current_block_height > call.timeout.into() {
        info!(
            "This LogicCall is timed out. timeout block: {} current block: {}, exiting early",
            current_block_height, call.timeout
        );
        return Ok(());
    }

    let payload = encode_logic_call_payload(current_valset, &call, confirms)?;

    let tx = web3
        .send_transaction(
            peggy_contract_address,
            payload,
            0u32.into(),
            eth_address,
            our_eth_key,
            vec![],
        )
        .await?;
    info!("Sent batch update with txid {:#066x}", tx);

    web3.wait_for_transaction(tx.clone(), timeout, None).await?;

    let last_nonce = get_logic_call_nonce(
        peggy_contract_address,
        call.invalidation_id,
        eth_address,
        &web3,
    )
    .await?;
    if last_nonce != new_call_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_call_nonce
        );
    } else {
        info!(
            "Successfully updated LogicCall with new Nonce {:?}",
            last_nonce
        );
    }
    Ok(())
}

/// Returns the cost in Eth of sending this batch
pub async fn estimate_logic_call_cost(
    current_valset: Valset,
    call: LogicCall,
    confirms: &[LogicCallConfirmResponse],
    web3: &Web3,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<GasCost, PeggyError> {
    let our_eth_address = our_eth_key.to_public_key().unwrap();
    let our_balance = web3.eth_get_balance(our_eth_address).await?;
    let our_nonce = web3.eth_get_transaction_count(our_eth_address).await?;
    let gas_limit = min((u64::MAX - 1).into(), our_balance.clone());
    let gas_price = web3.eth_gas_price().await?;
    let zero: Uint256 = 0u8.into();
    let val = web3
        .eth_estimate_gas(TransactionRequest {
            from: Some(our_eth_address),
            to: peggy_contract_address,
            nonce: Some(our_nonce.clone().into()),
            gas_price: Some(gas_price.clone().into()),
            gas: Some(gas_limit.into()),
            value: Some(zero.into()),
            data: Some(encode_logic_call_payload(current_valset, &call, confirms)?.into()),
        })
        .await?;

    Ok(GasCost {
        gas: val,
        gas_price,
    })
}

/// Encodes the logic call payload for both cost estimation and submission to EThereum
fn encode_logic_call_payload(
    current_valset: Valset,
    call: &LogicCall,
    confirms: &[LogicCallConfirmResponse],
) -> Result<Vec<u8>, PeggyError> {
    let (current_addresses, current_powers) = current_valset.filter_empty_addresses();
    let current_valset_nonce = current_valset.nonce;
    let sig_data = current_valset.order_sigs(confirms)?;
    let sig_arrays = to_arrays(sig_data);

    let mut transfer_amounts = Vec::new();
    let mut transfer_token_contracts = Vec::new();
    let mut fee_amounts = Vec::new();
    let mut fee_token_contracts = Vec::new();
    for item in call.transfers.iter() {
        transfer_amounts.push(Token::Uint(item.amount.clone()));
        transfer_token_contracts.push(item.token_contract_address);
    }
    for item in call.fees.iter() {
        fee_amounts.push(Token::Uint(item.amount.clone()));
        fee_token_contracts.push(item.token_contract_address);
    }

    // Solidity function signature
    // function submitBatch(
    // // The validators that approve the batch and new valset
    // address[] memory _currentValidators,
    // uint256[] memory _currentPowers,
    // uint256 _currentValsetNonce,
    // // These are arrays of the parts of the validators signatures
    // uint8[] memory _v,
    // bytes32[] memory _r,
    // bytes32[] memory _s,
    // // The LogicCall arguments, encoded as a struct (see the Ethereum ABI encoding documentation for the handling of structs as arguments)
    // uint256[] transferAmounts;
    // address[] transferTokenContracts;
    // // The fees (transferred to msg.sender)
    // uint256[] feeAmounts;
    // address[] feeTokenContracts;
    // // The arbitrary logic call
    // address logicContractAddress;
    // bytes payload;
    // // Invalidation metadata
    // uint256 timeOut;
    // bytes32 invalidationId;
    // uint256 invalidationNonce;
    let struct_tokens = &[
        Token::Dynamic(transfer_amounts),
        transfer_token_contracts.into(),
        Token::Dynamic(fee_amounts),
        fee_token_contracts.into(),
        call.logic_contract_address.into(),
        Token::UnboundedBytes(call.payload.clone()),
        call.timeout.into(),
        Token::Bytes(call.invalidation_id.clone()),
        call.invalidation_nonce.into(),
    ];
    let tokens = &[
        current_addresses.into(),
        current_powers.into(),
        current_valset_nonce.into(),
        sig_arrays.v,
        sig_arrays.r,
        sig_arrays.s,
        Token::Dynamic(struct_tokens.to_vec()),
    ];
    let payload = clarity::abi::encode_call(
        "submitLogicCall(address[],uint256[],uint256,uint8[],bytes32[],bytes32[],(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256))",
        tokens,
    )
    .unwrap();
    trace!("Tokens {:?}", tokens);

    Ok(payload)
}
