use clarity::{Address, Uint256};
use cosmos_peggy::query::get_last_event_nonce;
use deep_space::address::Address as CosmosAddress;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::{SendToCosmosEvent, TransactionBatchExecutedEvent};
use tonic::transport::Channel;
use web30::client::Web3;

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
pub async fn get_last_checked_block(
    grpc_client: PeggyQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    peggy_contract_address: Address,
    web3: &Web3,
) -> Uint256 {
    let mut grpc_client = grpc_client;
    // all events from the peggy contract forever, TODO reduce scope of request to reduce load
    // on full node. Also response is limited to 5mbyte in size so if you have too many events
    // this will simply fail.
    let all_batch_events = web3
        .check_for_events(
            0u8.into(),
            None,
            vec![peggy_contract_address],
            vec!["TransactionBatchExecutedEvent(uint256,address,uint256)"],
        )
        .await
        .unwrap();
    let all_valset_events = web3
        .check_for_events(
            0u8.into(),
            None,
            vec![peggy_contract_address],
            vec!["SendToCosmosEvent(address,address,bytes32,uint256,uint256)"],
        )
        .await
        .unwrap();
    let last_event_nonce: Uint256 = get_last_event_nonce(&mut grpc_client, our_cosmos_address)
        .await
        .unwrap()
        .into();

    if last_event_nonce == 0u8.into() {
        return web3.eth_block_number().await.unwrap();
    }

    trace!(
        "Found events {:?} {:?}",
        all_batch_events,
        all_valset_events
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
    for event in all_valset_events {
        match SendToCosmosEvent::from_log(&event) {
            Ok(send) => {
                if send.event_nonce == last_event_nonce && event.block_number.is_some() {
                    return event.block_number.unwrap();
                }
            }
            Err(e) => error!("Got valset event that we can't parse {}", e),
        }
    }

    panic!("Could not find the last event relayed by {}, Last Event nonce is {} but no event matching that could be found!", our_cosmos_address, last_event_nonce)
}
