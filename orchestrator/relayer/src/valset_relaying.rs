//! This module contains code for the validator update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_peggy::query::get_latest_valsets;
use cosmos_peggy::query::{get_all_valset_confirms, get_valset};
use ethereum_peggy::{one_eth, utils::downcast_to_u128, valset_update::send_eth_valset_update};
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::Valset;
use tonic::transport::Channel;
use web30::client::Web3;

/// Check the last validator set on Ethereum, if it's lower than our latest validator
/// set then we should package and submit the update as an Ethereum transaction
pub async fn relay_valsets(
    // the validator set currently in the contract on Ethereum
    current_valset: Valset,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    grpc_client: &mut PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    timeout: Duration,
) {
    // we have to start with the current valset, we need to know what's currently
    // in the contract in order to determine if a new validator set is valid.
    // For example the contract has set A which contains validators x/y/z the
    // latest valset has set C which has validators z/e/f in order to have enough
    // power we actually need to submit validator set B with validators x/y/e in
    // order to know that we need a set from the history

    // we should determine if we need to relay one
    // to Ethereum for that we will find the latest confirmed valset and compare it to the ethereum chain
    let latest_valsets = get_latest_valsets(grpc_client).await;
    if latest_valsets.is_err() {
        trace!("Failed to get latest valsets!");
        // there are no latest valsets to check, possible on a bootstrapping chain maybe handle better?
        return;
    }
    let latest_valsets = latest_valsets.unwrap();
    if latest_valsets.is_empty() {
        return;
    }

    // we only use the latest valsets endpoint to get a starting point, from there we will iterate
    // backwards until we find the newest validator set that we can submit to the bridge. So if we
    // have sets A-Z and it's possible to submit only A, L, and Q before reaching Z this code will do
    // so.
    let mut latest_nonce = latest_valsets[0].nonce;
    let mut latest_confirmed = None;
    let mut latest_valset = None;
    // this is used to display the state of the last validator set to fail signature checks
    let mut last_error = None;
    while latest_nonce > 0 {
        let valset = get_valset(grpc_client, latest_nonce).await;
        if let Ok(Some(valset)) = valset {
            let confirms = get_all_valset_confirms(grpc_client, valset.nonce).await;
            if let Ok(confirms) = confirms {
                // order valset sigs prepares signatures for submission, notice we compare
                // them to the 'current' set in the bridge, this confirms for us that the validator set
                // we have here can be submitted to the bridge in it's current state
                let res = current_valset.order_sigs(&confirms);
                if res.is_ok() {
                    latest_confirmed = Some(confirms);
                    latest_valset = Some(valset);
                    // once we have the latest validator set we can submit exit
                    break;
                } else if let Err(e) = res {
                    last_error = Some(e);
                }
            }
        }

        latest_nonce -= 1
    }

    if latest_confirmed.is_none() {
        error!("We don't have a latest confirmed valset?");
        return;
    }
    // the latest cosmos validator set that it is possible to submit given the constraints
    // of the validator set currently in the bridge
    let latest_cosmos_valset = latest_valset.unwrap();
    // the signatures for the above
    let latest_cosmos_confirmed = latest_confirmed.unwrap();

    // this will print a message indicating the signing state of the latest validator
    // set if the latest available validator set is not the latest one that is possible
    // to submit. AKA if the bridge is behind where it should be
    if latest_nonce > latest_cosmos_valset.nonce && last_error.is_some() {
        warn!("{:?}", last_error)
    }

    let latest_cosmos_valset_nonce = latest_cosmos_valset.nonce;
    if latest_cosmos_valset_nonce > current_valset.nonce {
        let cost = ethereum_peggy::valset_update::estimate_valset_cost(
            &latest_cosmos_valset,
            &current_valset,
            &latest_cosmos_confirmed,
            web3,
            peggy_contract_address,
            ethereum_key,
        )
        .await;
        if cost.is_err() {
            error!("Valset cost estimate failed with {:?}", cost);
            return;
        }
        let cost = cost.unwrap();

        info!(
           "We have detected latest valset {} but latest on Ethereum is {} This valset is estimated to cost {} Gas / {:.4} ETH to submit",
            latest_cosmos_valset.nonce, current_valset.nonce,
            cost.gas_price.clone(),
            downcast_to_u128(cost.get_total()).unwrap() as f32
                / downcast_to_u128(one_eth()).unwrap() as f32
        );

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
