//! This module contains code for the validator update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_peggy::query::get_all_valset_confirms;
use cosmos_peggy::query::get_latest_valsets;
use ethereum_peggy::valset_update::send_eth_valset_update;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use tonic::transport::Channel;
use web30::client::Web3;

use crate::find_latest_valset::find_latest_valset;

/// Check the last validator set on Ethereum, if it's lower than our latest validator
/// set then we should package and submit the update as an Ethereum transaction
pub async fn relay_valsets(
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    grpc_client: &mut PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    timeout: Duration,
) {
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    // we should determine if we need to relay one
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
            latest_confirmed = Some(confirms);
            latest_valset = Some(set);
            break;
        }
    }

    if latest_confirmed.is_none() {
        error!("We don't have a latest confirmed valset?");
        return;
    }
    let latest_cosmos_confirmed = latest_confirmed.unwrap();
    let latest_cosmos_valset = latest_valset.unwrap();

    let current_valset = find_latest_valset(
        grpc_client,
        our_ethereum_address,
        peggy_contract_address,
        web3,
    )
    .await;
    if current_valset.is_err() {
        error!("Could not get current valset!");
        return;
    }
    let current_valset = current_valset.unwrap();
    let latest_cosmos_valset_nonce = latest_cosmos_valset.nonce;
    if latest_cosmos_valset_nonce > current_valset.nonce {
        info!(
            "We have detected latest valset {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_valset.nonce, current_valset.nonce
        );

        // If the ENV var NO_GAS_OPT is not set at compile time then the resulting binary will not
        // have gas optimizations. In this case if we exit early if gas optimizations are enabled
        // (the default value)
        if option_env!("NO_GAS_OPT").is_none() {
            let diff = current_valset.power_diff(&latest_cosmos_valset);
            // if the power difference is less than one percent, skip updating
            // the validator set
            if diff < 0.01 {
                info!("Difference in power between valset {} and {} is less than 1% skipping update to save gas", current_valset.nonce, latest_cosmos_valset.nonce);
                return;
            }
        }

        let _res = send_eth_valset_update(
            latest_cosmos_valset,
            current_valset,
            &latest_cosmos_confirmed,
            web3,
            timeout,
            peggy_contract_address,
            ethereum_key,
        )
        .await;
    }
}
