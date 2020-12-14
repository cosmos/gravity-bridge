//! This module contains code for the batch update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use cosmos_peggy::query::get_transaction_batch_signatures;
use cosmos_peggy::query::get_valset;
use cosmos_peggy::query::{get_latest_transaction_batches, get_oldest_unsigned_transaction_batch};
use cosmos_peggy::send::send_batch_confirm;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::submit_batch::send_eth_transaction_batch;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::utils::{get_peggy_id, get_tx_batch_nonce};
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::{BatchConfirmResponse, TransactionBatch};
use std::time::Duration;
use tonic::transport::Channel;
use web30::client::Web3;

/// This function makes all decisions about the validator set update lifecycle.
///
/// It goes roughly in this order
/// 1) Determine if we should request a validator set update
/// 2) See if we have any unsigned validator set updates, if so sign and submit them
/// 3) Check the last validator set on Ethereum, if it's lower than our latest validator
///    set then we should package and submit the update as an Ethereum transaction
#[allow(clippy::too_many_arguments)]
pub async fn relay_batches(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    contact: &Contact,
    grpc_client: &mut PeggyQueryClient<Channel>,
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
    match get_oldest_unsigned_transaction_batch(grpc_client, our_cosmos_address).await {
        Ok(Some(last_unsigned_batch)) => {
            info!("Sending batch confirm for {}", last_unsigned_batch.nonce);
            let res = send_batch_confirm(
                contact,
                ethereum_key,
                fee.clone(),
                last_unsigned_batch,
                cosmos_key,
                peggy_id,
            )
            .await
            .unwrap();
            info!("Batch confirm result is {:?}", res);
        }
        Ok(None) => trace!("No unsigned batches! Everything good!"),
        Err(e) => trace!("Failed to get unsigned Batches with {:?}", e),
    }

    let latest_batches = get_latest_transaction_batches(grpc_client).await;
    trace!("Latest batches {:?}", latest_batches);
    if latest_batches.is_err() {
        return;
    }
    let latest_batches = latest_batches.unwrap();
    let mut oldest_signed_batch: Option<TransactionBatch> = None;
    let mut oldest_signatures: Option<Vec<BatchConfirmResponse>> = None;
    for batch in latest_batches {
        let sigs =
            get_transaction_batch_signatures(grpc_client, batch.nonce, batch.token_contract).await;
        trace!("Got sigs {:?}", sigs);
        if let Ok(sigs) = sigs {
            // todo check that enough people have signed
            oldest_signed_batch = Some(batch);
            oldest_signatures = Some(sigs);
        } else {
            error!(
                "could not get signatures for {}:{} with {:?}",
                batch.token_contract, batch.nonce, sigs
            );
        }
    }
    if oldest_signed_batch.is_none() {
        error!("Could not find batch with signatures! exiting");
        return;
    }
    let oldest_signed_batch = oldest_signed_batch.unwrap();
    let oldest_signatures = oldest_signatures.unwrap();
    let erc20_contract = oldest_signed_batch.token_contract;

    let latest_ethereum_batch = get_tx_batch_nonce(
        peggy_contract_address,
        erc20_contract,
        our_ethereum_address,
        web3,
    )
    .await
    .expect("Failed to get batch nonce from Ethereum");
    let latest_ethereum_valset =
        get_valset_nonce(peggy_contract_address, our_ethereum_address, web3)
            .await
            .expect("Failed to get Ethereum valset");
    let latest_cosmos_batch_nonce = oldest_signed_batch.clone().nonce;
    if latest_cosmos_batch_nonce > latest_ethereum_batch {
        info!(
            "We have detected latest batch {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_batch_nonce, latest_ethereum_batch
        );

        // get the current valset from the Cosmos chain
        let current_valset = get_valset(grpc_client, latest_ethereum_valset).await;
        if let Ok(Some(current_valset)) = current_valset {
            let _res = send_eth_transaction_batch(
                current_valset,
                oldest_signed_batch,
                &oldest_signatures,
                web3,
                timeout,
                peggy_contract_address,
                ethereum_key,
            )
            .await;
        } else {
            error!("Failed to get latest validator set!");
        }
    }
}
