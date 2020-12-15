use crate::{batch_relaying::relay_batches, valset_relaying::relay_valsets};
use clarity::address::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use std::time::{Duration, Instant};
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub const LOOP_SPEED: Duration = Duration::from_secs(10);

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
        relay_valsets(
            ethereum_key,
            &web3,
            &mut grpc_client,
            peggy_contract_address,
            LOOP_SPEED,
        )
        .await;

        relay_batches(
            ethereum_key,
            &web3,
            &mut grpc_client,
            peggy_contract_address,
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
