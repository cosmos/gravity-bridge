use crate::{
    batch_relaying::relay_batches, find_latest_valset::find_latest_valset,
    logic_call_relaying::relay_logic_calls, valset_relaying::relay_valsets,
};
use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use ethereum_peggy::utils::get_peggy_id;
use peggy_proto::gravity::query_client::QueryClient as PeggyQueryClient;
use std::time::{Duration, Instant};
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub const LOOP_SPEED: Duration = Duration::from_secs(17);

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
pub async fn relayer_main_loop(
    ethereum_key: EthPrivateKey,
    web3: Web3,
    grpc_client: PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
) {
    let mut grpc_client = grpc_client;
    loop {
        let loop_start = Instant::now();

        let our_ethereum_address = ethereum_key.to_public_key().unwrap();
        let current_valset = find_latest_valset(
            &mut grpc_client,
            our_ethereum_address,
            peggy_contract_address,
            &web3,
        )
        .await;
        if current_valset.is_err() {
            error!("Could not get current valset! {:?}", current_valset);
            continue;
        }
        let current_valset = current_valset.unwrap();

        let peggy_id = get_peggy_id(peggy_contract_address, our_ethereum_address, &web3).await;
        if peggy_id.is_err() {
            error!("Failed to get PeggyID, check your Eth node");
            return;
        }
        let peggy_id = peggy_id.unwrap();
        let peggy_id = String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID");

        relay_valsets(
            current_valset.clone(),
            ethereum_key,
            &web3,
            &mut grpc_client,
            peggy_contract_address,
            peggy_id.clone(),
            LOOP_SPEED,
        )
        .await;

        relay_batches(
            current_valset.clone(),
            ethereum_key,
            &web3,
            &mut grpc_client,
            peggy_contract_address,
            peggy_id.clone(),
            LOOP_SPEED,
        )
        .await;

        relay_logic_calls(
            current_valset,
            ethereum_key,
            &web3,
            &mut grpc_client,
            peggy_contract_address,
            peggy_id.clone(),
            LOOP_SPEED,
        )
        .await;

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < LOOP_SPEED {
            delay_for(LOOP_SPEED - elapsed).await;
        }
    }
}
