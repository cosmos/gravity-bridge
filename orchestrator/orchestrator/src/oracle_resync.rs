use std::time::Duration;

use clarity::{Address, Uint256};
use cosmos_peggy::query::get_last_event_nonce;
use deep_space::address::Address as CosmosAddress;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::{
    ERC20DeployedEvent, LogicCallExecutedEvent, SendToCosmosEvent, TransactionBatchExecutedEvent,
    ValsetUpdatedEvent,
};
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

const RETRY_TIME: Duration = Duration::from_secs(5);

/// This function retrieves the last event nonce this oracle has relayed to Cosmos
/// it then uses the Ethereum indexes to determine what block the last entry
pub async fn get_last_checked_block(
    grpc_client: PeggyQueryClient<Channel>,
    our_cosmos_address: CosmosAddress,
    peggy_contract_address: Address,
    web3: &Web3,
) -> Uint256 {
    let mut grpc_client = grpc_client;
    const BLOCKS_TO_SEARCH: u128 = 50_000u128;

    let latest_block = get_block_number_with_retry(web3).await;
    let mut last_event_nonce: Uint256 =
        get_last_event_nonce_with_retry(&mut grpc_client, our_cosmos_address)
            .await
            .into();

    // zero indicates this oracle has never submitted an event before since there is no
    // zero event nonce (it's pre-incremented in the solidity contract) we have to go
    // and look for event nonce one.
    if last_event_nonce == 0u8.into() {
        last_event_nonce = 1u8.into();
    }

    let mut current_block: Uint256 = latest_block.clone();

    while current_block.clone() > 0u8.into() {
        info!(
            "Oracle is resyncing, looking back into the history to find our last event nonce {}, on block {}",
            last_event_nonce, current_block
        );
        let end_search = if current_block.clone() < BLOCKS_TO_SEARCH.into() {
            0u8.into()
        } else {
            current_block.clone() - BLOCKS_TO_SEARCH.into()
        };
        let batch_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["TransactionBatchExecutedEvent(uint256,address,uint256)"],
            )
            .await;
        let send_to_cosmos_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["SendToCosmosEvent(address,address,bytes32,uint256,uint256)"],
            )
            .await;
        let erc20_deployed_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["ERC20DeployedEvent(string,address,string,string,uint8,uint256)"],
            )
            .await;
        let logic_call_executed_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["LogicCallEvent(bytes32,uint256,bytes,uint256)"],
            )
            .await;

        // valset events do not have an event nonce (because they are not relayed to cosmos)
        // and therefore they are mostly useless to us. But they do have one special property
        // that is useful to us in this handler a valset update event for nonce 0 is emitted
        // in the contract constructor meaning once you find that event you can exit the search
        // with confidence that you have not missed any events without searching the entire blockchain
        // history
        let valset_events = web3
            .check_for_events(
                end_search.clone(),
                Some(current_block.clone()),
                vec![peggy_contract_address],
                vec!["ValsetUpdatedEvent(uint256,address[],uint256[])"],
            )
            .await;
        if batch_events.is_err()
            || send_to_cosmos_events.is_err()
            || valset_events.is_err()
            || erc20_deployed_events.is_err()
            || logic_call_executed_events.is_err()
        {
            error!("Failed to get blockchain events while resyncing, is your Eth node working?");
            delay_for(RETRY_TIME).await;
            continue;
        }
        let batch_events = batch_events.unwrap();
        let send_to_cosmos_events = send_to_cosmos_events.unwrap();
        let valset_events = valset_events.unwrap();
        let erc20_deployed_events = erc20_deployed_events.unwrap();
        let logic_call_executed_events = logic_call_executed_events.unwrap();

        // look for and return the block number of the event last seen on the Cosmos chain
        // then we will play events from that block (including that block, just in case
        // there is more than one event there) onwards. We use valset nonce 0 as an indicator
        // of what block the contract was deployed on.
        for event in batch_events {
            match TransactionBatchExecutedEvent::from_log(&event) {
                Ok(batch) => {
                    if batch.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got batch event that we can't parse {}", e),
            }
        }
        for event in send_to_cosmos_events {
            match SendToCosmosEvent::from_log(&event) {
                Ok(send) => {
                    if send.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got SendToCosmos event that we can't parse {}", e),
            }
        }
        for event in erc20_deployed_events {
            match ERC20DeployedEvent::from_log(&event) {
                Ok(deploy) => {
                    if deploy.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got ERC20Deployed event that we can't parse {}", e),
            }
        }
        for event in logic_call_executed_events {
            match LogicCallExecutedEvent::from_log(&event) {
                Ok(call) => {
                    if call.event_nonce == last_event_nonce && event.block_number.is_some() {
                        return event.block_number.unwrap();
                    }
                }
                Err(e) => error!("Got ERC20Deployed event that we can't parse {}", e),
            }
        }
        for event in valset_events {
            match ValsetUpdatedEvent::from_log(&event) {
                Ok(valset) => {
                    // if we've found this event it is the first possible event from the contract
                    // no other events can come before it, therefore either there's been a parsing error
                    // or no events have been submitted on this chain yet.
                    if valset.nonce == 0 && last_event_nonce == 1u8.into() {
                        return latest_block;
                    }
                    // if we're looking for a later event nonce and we find the deployment of the contract
                    // we must have failed to parse the event we're looking for. The oracle can not start
                    if valset.nonce == 0 && last_event_nonce > 1u8.into() {
                        panic!("Could not find the last event relayed by {}, Last Event nonce is {} but no event matching that could be found!", our_cosmos_address, last_event_nonce)
                    }
                }
                Err(e) => error!("Got valset event that we can't parse {}", e),
            }
        }
        current_block = end_search;
    }

    panic!("You have reached the end of block history without finding the Peggy contract deploy event! You must have the wrong contract address!");
}

/// gets the current block number, no matter how long it takes
async fn get_block_number_with_retry(web3: &Web3) -> Uint256 {
    let mut res = web3.eth_block_number().await;
    while res.is_err() {
        error!("Failed to get latest block! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = web3.eth_block_number().await;
    }
    res.unwrap()
}

/// gets the last event nonce, no matter how long it takes.
async fn get_last_event_nonce_with_retry(
    client: &mut PeggyQueryClient<Channel>,
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
