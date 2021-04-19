//! Basic utility functions to stubbornly get data
use clarity::Uint256;
use cosmos_gravity::query::get_last_event_nonce;
use deep_space::address::Address as CosmosAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use std::time::Duration;
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub const RETRY_TIME: Duration = Duration::from_secs(5);

/// gets the current block number, no matter how long it takes
pub async fn get_block_number_with_retry(web3: &Web3) -> Uint256 {
    let mut res = web3.eth_block_number().await;
    while res.is_err() {
        error!("Failed to get latest block! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = web3.eth_block_number().await;
    }
    res.unwrap()
}

/// gets the last event nonce, no matter how long it takes.
pub async fn get_last_event_nonce_with_retry(
    client: &mut GravityQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
) -> u64 {
    let mut res = get_last_event_nonce(client, our_cosmos_address).await;
    while res.is_err() {
        error!(
            "Failed to get last event nonce, is the Cosmos GRPC working? {:?}",
            res
        );
        delay_for(RETRY_TIME).await;
        res = get_last_event_nonce(client, our_cosmos_address).await;
    }
    res.unwrap()
}

/// gets the net version, no matter how long it takes
pub async fn get_net_version_with_retry(web3: &Web3) -> u64 {
    let mut res = web3.net_version().await;
    while res.is_err() {
        error!("Failed to get net version! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = web3.net_version().await;
    }
    res.unwrap()
}
