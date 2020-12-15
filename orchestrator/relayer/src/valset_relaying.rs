//! This module contains code for the validator update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use cosmos_peggy::query::get_latest_valsets;
use cosmos_peggy::query::{get_all_valset_confirms, get_valset};
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::utils::get_peggy_id;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::valset_update::send_eth_valset_update;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
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
pub async fn relay_valsets(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    contact: &Contact,
    grpc_client: &mut PeggyQueryClient<Channel>,
    contract_address: EthAddress,
    fee: Coin,
    timeout: Duration,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    let peggy_id = get_peggy_id(contract_address, our_ethereum_address, web3).await;
    if peggy_id.is_err() {
        error!("Failed to get PeggyID!");
        return;
    }
    let peggy_id = peggy_id.unwrap();
    let peggy_id = String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID");

    // now that we have caught up on valset requests we should determine if we need to relay one
    // to Ethereum for that we will find the latest confirmed valset and compare it to the ethereum chain
    let latest_valsets = get_latest_valsets(grpc_client).await;
    if latest_valsets.is_err() {
        trace!("Failed to get latest valsets!");
        // there are no latest valsets to check, possible on a bootstrapping chain maybe handle better?
        return;
    }
    let latest_valsets = latest_valsets.unwrap();

    let mut latest_confirmed = None;
    let mut latest_valset = None;
    for set in latest_valsets {
        let confirms = get_all_valset_confirms(grpc_client, set.nonce).await;
        if let Ok(confirms) = confirms {
            // todo allow submission without signatures from all validators
            if confirms.len() == set.members.len() {
                latest_confirmed = Some(confirms);
                latest_valset = Some(set);
            } else {
                info!(
                    "Skipping incomplete valset {} we have {} confirms of {}",
                    set.nonce,
                    confirms.len(),
                    set.members.len()
                );
            }
            break;
        }
    }

    if latest_confirmed.is_none() {
        error!("We don't have a latest confirmed valset?");
        return;
    }
    let latest_cosmos_confirmed = latest_confirmed.unwrap();
    let latest_cosmos_valset = latest_valset.unwrap();

    let latest_ethereum_valset = get_valset_nonce(contract_address, our_ethereum_address, web3)
        .await
        .expect("Failed to get Ethereum valset");
    let latest_cosmos_valset_nonce = latest_cosmos_valset.nonce;
    if latest_cosmos_valset_nonce > latest_ethereum_valset {
        info!(
            "We have detected latest valset {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_valset.nonce, latest_ethereum_valset
        );

        let old_valset = if latest_ethereum_valset == 0 {
            info!("This is the first validator set update! Using the current set");
            // we need to have a special case for validator set zero, that valset was never stored on chain
            // right now we just use the current valset
            let mut latest_valset = latest_cosmos_valset.clone();
            latest_valset.nonce = 0;
            latest_valset
        } else {
            // get the old valset from the Cosmos chain
            if let Ok(Some(valset)) = get_valset(grpc_client, latest_ethereum_valset).await {
                valset
            } else {
                error!("Failed to get latest valset!");
                return;
            }
        };

        let _res = send_eth_valset_update(
            latest_cosmos_valset,
            old_valset,
            &latest_cosmos_confirmed,
            web3,
            timeout,
            contract_address,
            ethereum_key,
        )
        .await;
    }
}
