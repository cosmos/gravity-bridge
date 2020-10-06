use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use std::thread;
use std::time::Duration;
use std::time::Instant;
use url::Url;
use web30::client::Web3;

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
pub async fn orchestrator_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    loop_speed: Duration,
) {
    let mut last_seen_block = web3.eth_get_latest_block();
    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_get_latest_block().await.unwrap();
        let latest_cosmos_block = contact.get_latest_block().await.unwrap();
        println!(
            "Latest Eth block {} Latest Cosmos block {}",
            latest_eth_block.number, latest_cosmos_block.block.header.version.block
        );

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < loop_speed {
            thread::sleep(loop_speed - elapsed)
        }
    }
}
