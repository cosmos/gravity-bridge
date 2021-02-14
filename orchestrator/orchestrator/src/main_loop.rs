//! This file contains the main loops for two distinct functions that just happen to reside int his same binary for ease of use. The Ethereum Signer and the Ethereum Oracle are both roles in Peggy
//! that can only be run by a validator. This single binary the 'Orchestrator' runs not only these two rules but also the untrusted role of a relayer, that does not need any permissions and has it's
//! own crate and binary so that anyone may run it.

use crate::{ethereum_event_watcher::check_for_events, oracle_resync::get_last_checked_block};
use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::{
    query::{get_oldest_unsigned_transaction_batch, get_oldest_unsigned_valsets},
    send::{send_batch_confirm, send_valset_confirms},
};
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::utils::get_peggy_id;
use futures::future::join3;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use relayer::main_loop::relayer_main_loop;
use std::time::Duration;
use std::time::Instant;
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

/// The execution speed governing all loops in this file
/// which is to say all loops started by Orchestrator main
/// loop except the relayer loop
pub const ETH_SIGNER_LOOP_SPEED: Duration = Duration::from_secs(11);
pub const ETH_ORACLE_LOOP_SPEED: Duration = Duration::from_secs(13);

/// This loop combines the three major roles required to make
/// up the 'Orchestrator', all three of these are async loops
/// meaning they will occupy the same thread, but since they do
/// very little actual cpu bound work and spend the vast majority
/// of all execution time sleeping this shouldn't be an issue at all.
pub async fn orchestrator_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    pay_fees_in: String,
) {
    let fee = Coin {
        denom: pay_fees_in.clone(),
        amount: 1u32.into(),
    };

    let a = eth_oracle_main_loop(
        cosmos_key,
        web3.clone(),
        contact.clone(),
        grpc_client.clone(),
        peggy_contract_address,
        fee.clone(),
    );
    let b = eth_signer_main_loop(
        cosmos_key,
        ethereum_key,
        web3.clone(),
        contact.clone(),
        grpc_client.clone(),
        peggy_contract_address,
        fee.clone(),
    );
    let c = relayer_main_loop(
        ethereum_key,
        web3,
        grpc_client.clone(),
        peggy_contract_address,
    );
    join3(a, b, c).await;
}

/// This function is responsible for making sure that Ethereum events are retrieved from the Ethereum blockchain
/// and ferried over to Cosmos where they will be used to issue tokens or process batches.
pub async fn eth_oracle_main_loop(
    cosmos_key: CosmosPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    fee: Coin,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let long_timeout_web30 = Web3::new(&web3.get_url(), Duration::from_secs(120));
    let mut last_checked_block: Uint256 = get_last_checked_block(
        grpc_client.clone(),
        our_cosmos_address,
        peggy_contract_address,
        &long_timeout_web30,
    )
    .await;
    info!("Oracle resync complete, Oracle now operational");
    let mut grpc_client = grpc_client;

    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_block_number().await;
        let latest_cosmos_block = contact.get_latest_block_number().await;
        if let (Ok(latest_eth_block), Ok(latest_cosmos_block)) =
            (latest_eth_block, latest_cosmos_block)
        {
            trace!(
                "Latest Eth block {} Latest Cosmos block {}",
                latest_eth_block,
                latest_cosmos_block,
            );
        }

        // Relays events from Ethereum -> Cosmos
        match check_for_events(
            &web3,
            &contact,
            &mut grpc_client,
            peggy_contract_address,
            cosmos_key,
            fee.clone(),
            last_checked_block.clone(),
        )
        .await
        {
            Ok(new_block) => last_checked_block = new_block,
            Err(e) => error!(
                "Failed to get events for block range, Check your Eth node and Cosmos gRPC {:?}",
                e
            ),
        }

        // a bit of logic that tires to keep things running every LOOP_SPEED seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < ETH_ORACLE_LOOP_SPEED {
            delay_for(ETH_ORACLE_LOOP_SPEED - elapsed).await;
        }
    }
}

/// The eth_signer simply signs off on any batches or validator sets provided by the validator
/// since these are provided directly by a trusted Cosmsos node they can simply be assumed to be
/// valid and signed off on.
pub async fn eth_signer_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    fee: Coin,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();
    let mut grpc_client = grpc_client;
    let peggy_id = get_peggy_id(peggy_contract_address, our_ethereum_address, &web3).await;
    if peggy_id.is_err() {
        error!("Failed to get PeggyID, check your Eth node");
        return;
    }
    let peggy_id = peggy_id.unwrap();
    let peggy_id = String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID");

    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_block_number().await;
        let latest_cosmos_block = contact.get_latest_block_number().await;
        if let (Ok(latest_eth_block), Ok(latest_cosmos_block)) =
            (latest_eth_block, latest_cosmos_block)
        {
            trace!(
                "Latest Eth block {} Latest Cosmos block {}",
                latest_eth_block,
                latest_cosmos_block
            );
        }

        // sign the last unsigned valsets
        match get_oldest_unsigned_valsets(&mut grpc_client, our_cosmos_address).await {
            Ok(valsets) => {
                if valsets.is_empty() {
                    trace!("No validator sets to sign, node is caught up!")
                } else {
                    info!(
                        "Sending {} valset confirms starting with {}",
                        valsets.len(),
                        valsets[0].nonce
                    );
                    let res = send_valset_confirms(
                        &contact,
                        ethereum_key,
                        fee.clone(),
                        valsets,
                        cosmos_key,
                        peggy_id.clone(),
                    )
                    .await;
                    trace!("Valset confirm result is {:?}", res);
                }
            }
            Err(e) => trace!(
                "Failed to get unsigned valsets, check your Cosmos gRPC {:?}",
                e
            ),
        }

        // sign the last unsigned batch, TODO check if we already have signed this
        match get_oldest_unsigned_transaction_batch(&mut grpc_client, our_cosmos_address).await {
            Ok(Some(last_unsigned_batch)) => {
                info!(
                    "Sending batch confirm for {}:{} with {} in fees",
                    last_unsigned_batch.token_contract,
                    last_unsigned_batch.nonce,
                    last_unsigned_batch.total_fee.amount
                );
                let res = send_batch_confirm(
                    &contact,
                    ethereum_key,
                    fee.clone(),
                    last_unsigned_batch,
                    cosmos_key,
                    peggy_id.clone(),
                )
                .await;
                trace!("Batch confirm result is {:?}", res);
            }
            Ok(None) => trace!("No unsigned batches! Everything good!"),
            Err(e) => trace!(
                "Failed to get unsigned Batches, check your Cosmos gRPC {:?}",
                e
            ),
        }

        // a bit of logic that tires to keep things running every LOOP_SPEED seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < ETH_SIGNER_LOOP_SPEED {
            delay_for(ETH_SIGNER_LOOP_SPEED - elapsed).await;
        }
    }
}
