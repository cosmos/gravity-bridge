//! This module contains code for the batch update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use crate::find_latest_valset::find_latest_valset;
use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_peggy::query::get_latest_transaction_batches;
use cosmos_peggy::query::get_transaction_batch_signatures;
use ethereum_peggy::submit_batch::send_eth_transaction_batch;
use ethereum_peggy::utils::get_tx_batch_nonce;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::{BatchConfirmResponse, TransactionBatch};
use std::time::Duration;
use tonic::transport::Channel;
use web30::client::Web3;

/// Check the last validator set on Ethereum, if it's lower than our latest validator
/// set then we should package and submit the update as an Ethereum transaction
pub async fn relay_batches(
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    mut grpc_client: &mut PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    timeout: Duration,
) {
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

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
        trace!("Could not find batch with signatures! exiting");
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
    .await;
    if latest_ethereum_batch.is_err() {
        error!(
            "Failed to get latest Ethereum batch with {:?}",
            latest_ethereum_batch
        );
    }
    let latest_ethereum_batch = latest_ethereum_batch.unwrap();
    let latest_cosmos_batch_nonce = oldest_signed_batch.clone().nonce;
    if latest_cosmos_batch_nonce > latest_ethereum_batch {
        info!(
            "We have detected latest batch {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_batch_nonce, latest_ethereum_batch
        );
        let current_valset = find_latest_valset(
            &mut grpc_client,
            our_ethereum_address,
            peggy_contract_address,
            web3,
        )
        .await;
        if let Ok(current_valset) = current_valset {
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
            error!("Failed to find latest valset with {:?}", current_valset);
        }
    }
}
