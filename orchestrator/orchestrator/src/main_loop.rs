use crate::ethereum_event_watcher::check_for_events;
use crate::valset_relaying::relay_valsets;
use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, Uint256};
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use std::time::Duration;
use std::time::Instant;
use tokio::time::delay_for;
use web30::client::Web3;

const BLOCK_DELAY: u128 = 50;

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
pub async fn orchestrator_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    contract_address: EthAddress,
    pay_fees_in: String,
    loop_speed: Duration,
) {
    let mut last_checked_block: Uint256 = 0u64.into();
    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_block_number().await.unwrap();
        let latest_cosmos_block = contact.get_latest_block().await.unwrap();
        info!(
            "Latest Eth block {} Latest Cosmos block {}",
            latest_eth_block, latest_cosmos_block.block.header.version.block
        );

        relay_valsets(
            cosmos_key,
            ethereum_key,
            &web3,
            &contact,
            contract_address,
            pay_fees_in.clone(),
            loop_speed,
        )
        .await;

        match check_for_events(&web3, contract_address, last_checked_block.clone()).await {
            Ok(new_block) => last_checked_block = new_block,
            Err(e) => error!("Failed to get events for block range {:?}", e),
        }

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < loop_speed {
            delay_for(loop_speed - elapsed).await;
        }
    }
}
