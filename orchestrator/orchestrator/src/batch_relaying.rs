//! This module contains code for the batch update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::query::get_oldest_unsigned_transaction_batch;
use cosmos_peggy::query::get_signed_transaction_batches;
use cosmos_peggy::query::get_valset;
use cosmos_peggy::send::send_batch_confirm;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::submit_batch::send_eth_transaction_batch;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::utils::{get_peggy_id, get_tx_batch_nonce};
use peggy_utils::types::SignedTransactionBatch;
use web30::client::Web3;

/// This function makes all decisions about the validator set update lifecycle.
///
/// It goes roughly in this order
/// 1) Determine if we should request a validator set update
/// 2) See if we have any unsigned validator set updates, if so sign and submit them
/// 3) Check the last validator set on Ethereum, if it's lower than our latest validator
///    set then we should package and submit the update as an Ethereum transaction
pub async fn relay_batches(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    contact: &Contact,
    peggy_contract_address: EthAddress,
    fee: Coin,
    timeout: Duration,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    // TODO this should compare to the cosmos value and crash if incorrect to do that
    // finish the bootstrapping message
    let peggy_id = get_peggy_id(peggy_contract_address, our_ethereum_address, web3).await;
    if peggy_id.is_err() {
        error!("Failed to get PeggyID");
        return;
    }
    let peggy_id = peggy_id.unwrap();
    let peggy_id = String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID");

    // sign the last unsigned batch, TODO check if we already have signed this
    match get_oldest_unsigned_transaction_batch(contact, our_cosmos_address).await {
        Ok(last_unsigned_batch) => {
            info!(
                "Sending batch confirm for {}",
                last_unsigned_batch.result.nonce
            );
            let res = send_batch_confirm(
                contact,
                ethereum_key,
                fee.clone(),
                last_unsigned_batch.result,
                cosmos_key,
                peggy_id,
            )
            .await
            .unwrap();
            info!("Batch confirm result is {:?}", res);
        }
        Err(e) => info!("Failed to get unsigned Batches with {:?}", e),
    }

    let latest_signed_batches = get_signed_transaction_batches(&contact).await;
    if latest_signed_batches.is_err() {
        error!("Failed to get latest signed batches");
        return;
    }
    let latest_signed_batches = latest_signed_batches.unwrap();
    let mut oldest_signed_batch: Option<SignedTransactionBatch> = None;
    for batch in latest_signed_batches.result {
        // todo support relaying valsets with less than 100% signing
        if batch.signatures.len() == batch.batch.valset.members.len() {
            oldest_signed_batch = Some(batch);
        }
    }
    if oldest_signed_batch.is_none() {
        return;
    }
    let oldest_signed_batch = oldest_signed_batch.unwrap();

    let latest_ethereum_batch =
        get_tx_batch_nonce(peggy_contract_address, our_ethereum_address, web3)
            .await
            .expect("Failed to get Ethereum valset");
    let latest_ethereum_valset =
        get_valset_nonce(peggy_contract_address, our_ethereum_address, web3)
            .await
            .expect("Failed to get Ethereum valset");
    let latest_cosmos_batch_nonce: Uint256 = oldest_signed_batch.clone().batch.nonce;
    if latest_cosmos_batch_nonce > latest_ethereum_batch {
        info!(
            "We have detected latest batch {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_batch_nonce, latest_ethereum_batch
        );

        let old_valset = if latest_ethereum_valset == 0u8.into() {
            // we need to have a special case for validator set zero, that valset was never stored on chain
            // right now we just use the current valset
            let mut latest_valset = oldest_signed_batch.batch.valset.clone();
            latest_valset.nonce = 0;
            latest_valset
        } else {
            // get the old valset from the Cosmos chain
            get_valset(contact, latest_ethereum_valset)
                .await
                .expect("Failed to get old valset")
                .result
        };

        let _res = send_eth_transaction_batch(
            old_valset,
            oldest_signed_batch,
            web3,
            timeout,
            peggy_contract_address,
            ethereum_key,
        )
        .await;
    }
}
