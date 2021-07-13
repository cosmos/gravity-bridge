use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, utils::bytes_to_hex_str};
use cosmos_gravity::query::{get_latest_logic_calls, get_logic_call_signatures};
use ethereum_gravity::one_eth;
use ethereum_gravity::{
    logic_call::send_eth_logic_call,
    utils::{downcast_to_u128, get_logic_call_nonce},
};
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::types::{LogicCallConfirmResponse, Valset};
use gravity_utils::{message_signatures::encode_logic_call_confirm_hashed, types::LogicCall};
use std::time::Duration;
use tonic::transport::Channel;
use web30::client::Web3;

pub async fn relay_logic_calls(
    // the validator set currently in the contract on Ethereum
    current_valset: Valset,
    ethereum_key: EthPrivateKey,
    web3: &Web3,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    gravity_id: String,
    timeout: Duration,
) {
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();

    let latest_calls = get_latest_logic_calls(grpc_client).await;
    trace!("Latest Logic calls {:?}", latest_calls);
    if latest_calls.is_err() {
        return;
    }
    let latest_calls = latest_calls.unwrap();
    let mut oldest_signed_call: Option<LogicCall> = None;
    let mut oldest_signatures: Option<Vec<LogicCallConfirmResponse>> = None;
    for call in latest_calls {
        let sigs = get_logic_call_signatures(
            grpc_client,
            call.invalidation_id.clone(),
            call.invalidation_nonce,
        )
        .await;
        trace!("Got sigs {:?}", sigs);
        if let Ok(sigs) = sigs {
            let hash = encode_logic_call_confirm_hashed(gravity_id.clone(), call.clone());
            // this checks that the signatures for the batch are actually possible to submit to the chain
            if current_valset.order_sigs(&hash, &sigs).is_ok() {
                oldest_signed_call = Some(call);
                oldest_signatures = Some(sigs);
            } else {
                warn!(
                    "LogicCall {}/{} can not be submitted yet, waiting for more signatures",
                    bytes_to_hex_str(&call.invalidation_id),
                    call.invalidation_nonce
                );
            }
        } else {
            error!(
                "could not get signatures for {}/{} with {:?}",
                bytes_to_hex_str(&call.invalidation_id),
                call.invalidation_nonce,
                sigs
            );
        }
    }
    if oldest_signed_call.is_none() {
        trace!("Could not find Call with signatures! exiting");
        return;
    }
    let oldest_signed_call = oldest_signed_call.unwrap();
    let oldest_signatures = oldest_signatures.unwrap();

    let latest_ethereum_call = get_logic_call_nonce(
        gravity_contract_address,
        oldest_signed_call.invalidation_id.clone(),
        our_ethereum_address,
        web3,
    )
    .await;
    if latest_ethereum_call.is_err() {
        error!(
            "Failed to get latest Ethereum LogicCall with {:?}",
            latest_ethereum_call
        );
        return;
    }
    let latest_ethereum_call = latest_ethereum_call.unwrap();
    let latest_cosmos_call_nonce = oldest_signed_call.clone().invalidation_nonce;
    if latest_cosmos_call_nonce > latest_ethereum_call {
        let cost = ethereum_gravity::logic_call::estimate_logic_call_cost(
            current_valset.clone(),
            oldest_signed_call.clone(),
            &oldest_signatures,
            web3,
            gravity_contract_address,
            gravity_id.clone(),
            ethereum_key,
        )
        .await;
        if cost.is_err() {
            error!("LogicCall cost estimate failed with {:?}", cost);
            return;
        }
        let cost = cost.unwrap();
        info!(
                "We have detected latest LogicCall {} but latest on Ethereum is {} This LogicCall is estimated to cost {} Gas / {:.4} ETH to submit",
                latest_cosmos_call_nonce,
                latest_ethereum_call,
                cost.gas_price.clone(),
                downcast_to_u128(cost.get_total()).unwrap() as f32
                    / downcast_to_u128(one_eth()).unwrap() as f32
            );

        let res = send_eth_logic_call(
            current_valset,
            oldest_signed_call,
            &oldest_signatures,
            web3,
            timeout,
            gravity_contract_address,
            gravity_id.clone(),
            ethereum_key,
        )
        .await;
        if res.is_err() {
            info!("LogicCall submission failed with {:?}", res);
        }
    }
}
