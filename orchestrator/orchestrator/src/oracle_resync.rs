use clarity::{Address, Uint256};
use cosmos_peggy::query::get_last_event_nonce;
use deep_space::address::Address as CosmosAddress;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::{SendToCosmosEvent, TransactionBatchExecutedEvent};
use tonic::transport::Channel;
use web30::client::Web3;

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
/// TODO this should simply be stored in the deposit or withdraw claim and we
/// ask the Cosmos chain, this searching is a total waste of work
pub async fn get_last_checked_block(
    grpc_client: PeggyQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    peggy_contract_address: Address,
    web3: &Web3,
) -> Uint256 {
    let mut grpc_client = grpc_client;
    let latest_block = web3.eth_block_number().await.unwrap();
    const BLOCKS_TO_SEARCH: u128 = 50_000u128;
    let mut current_block: Uint256 = latest_block.clone();

    let last_event_nonce: Uint256 = get_last_event_nonce(&mut grpc_client, our_cosmos_address)
        .await
        .unwrap()
        .into();

    if last_event_nonce == 0u8.into() {
        return web3.eth_block_number().await.unwrap();
    }

    while current_block.clone() > 0u8.into() {
        info!(
            "Oracle is resyncing, looking back into the history to find our last event nonce, on block {}",
            current_block
        );
        let end_search = if current_block.clone() < BLOCKS_TO_SEARCH.into() {
            0u8.into()
        } else {
            current_block.clone() - BLOCKS_TO_SEARCH.into()
        };
        let all_batch_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["TransactionBatchExecutedEvent(uint256,address,uint256)"],
            )
            .await
            .unwrap();
        let all_send_to_cosmos_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["SendToCosmosEvent(address,address,bytes32,uint256,uint256)"],
            )
            .await
            .unwrap();

        trace!(
            "Found events {:?} {:?}",
            all_batch_events,
            all_send_to_cosmos_events
        );
        for event in all_batch_events {
            match TransactionBatchExecutedEvent::from_log(&event) {
                Ok(batch) => {
                    if batch.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got batch event that we can't parse {}", e),
            }
        }
        for event in all_send_to_cosmos_events {
            match SendToCosmosEvent::from_log(&event) {
                Ok(send) => {
                    if send.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got valset event that we can't parse {}", e),
            }
        }
        current_block = end_search;
    }

    panic!("Could not find the last event relayed by {}, Last Event nonce is {} but no event matching that could be found!", our_cosmos_address, last_event_nonce)
}
