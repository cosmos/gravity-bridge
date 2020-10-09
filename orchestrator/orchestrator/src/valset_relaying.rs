//! This module contains code for the validator update lifecycle. Functioning as a way for this validator to observe
//! the state of both chains and perform the required operations.

use std::time::Duration;

use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::query::get_all_valset_confirms;
use cosmos_peggy::query::get_last_valset_requests;
use cosmos_peggy::query::get_oldest_unsigned_valset;
use cosmos_peggy::query::get_peggy_valset;
use cosmos_peggy::query::get_peggy_valset_request;
use cosmos_peggy::send::send_valset_confirm;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::utils::get_peggy_id;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::valset_update::send_eth_valset_update;
use web30::client::Web3;

/// This function makes all decisions about the validator set update lifecycle.
///
/// It goes roughly in this order
/// 1) Determine if we should request a validator set update
/// 2) See if we have any unsigned validator set updates, if so sign and submit them
/// 3) Check the last validator set on Ethereum, if it's lower than our latest validator
///    set then we should package and submit the update as an Ethereum transaction
pub async fn relay_valsets(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    contact: &Contact,
    contract_address: EthAddress,
    // the denom to pay fees in
    pay_fees_in: String,
    timeout: Duration,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    // TODO this should compare to the cosmos value and crash if incorrect
    let peggy_id = get_peggy_id(contract_address, our_ethereum_address, web3)
        .await
        .expect("Failed to get Peggy ID");

    // confirm arbitrarily many valsets
    let mut last_unsigned_valset = get_oldest_unsigned_valset(contact, our_cosmos_address).await;
    while let Ok(valset) = last_unsigned_valset {
        send_valset_confirm(
            contact,
            ethereum_key,
            Coin {
                denom: pay_fees_in.clone(),
                amount: 1u32.into(),
            },
            valset.result,
            cosmos_key,
            String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID"),
            None,
            None,
            None,
        )
        .await
        .unwrap();

        last_unsigned_valset = get_oldest_unsigned_valset(contact, our_cosmos_address).await;
    }

    // now that we have caught up on valset requests we should determine if we need to relay one
    // to Ethereum for that we will find the latest confirmed valset and compare it to the ethereum chain
    let latest_valsets = get_last_valset_requests(contact).await;
    if latest_valsets.is_err() {
        // there are no latest valsets to check, possible on a bootstrapping chain maybe handle better?
        return;
    }
    let latest_valsets = latest_valsets.unwrap();

    let mut latest_confirmed = None;
    let mut latest_valset = None;
    info!("Retrieving validator set signatures from the Cosmos chain");
    for set in latest_valsets.result {
        let confirms = get_all_valset_confirms(contact, set.nonce).await;
        if let Ok(confirms) = confirms {
            latest_confirmed = Some(confirms);
            latest_valset = Some(set);
            break;
        }
    }

    if latest_confirmed.is_none() {
        return;
    }
    let latest_confirmed = latest_confirmed.unwrap();
    let latest_valset = latest_valset.unwrap();

    let latest_ethereum_valset = get_valset_nonce(contract_address, our_ethereum_address, web3)
        .await
        .expect("Failed to get Ethereum valset");
    let latest_cosmos_valset: Uint256 = latest_confirmed.height.into();
    if latest_cosmos_valset > latest_ethereum_valset {
        info!(
            "We have detected latest valset {} but latest on Ethereum is {} sending an update!",
            latest_cosmos_valset, latest_ethereum_valset
        );

        let old_valset = if latest_ethereum_valset == 0u8.into() {
            // we need to have a special case for validator set zero, that valset was never stored on chain
            let mut valset = get_peggy_valset(contact)
                .await
                .expect("Failed to get old valset");
            // make sure the nonce here is zero
            valset.result.nonce = 0;
            valset
        } else {
            // get the old valset from the Cosmos chain
            get_peggy_valset_request(contact, latest_ethereum_valset)
                .await
                .expect("Failed to get old valset")
        };

        send_eth_valset_update(
            latest_valset,
            old_valset.result,
            &latest_confirmed.result,
            web3,
            timeout,
            contract_address,
            ethereum_key,
            cosmos_key,
        )
        .await
        .expect("Failed to update valset!");
    }
}
