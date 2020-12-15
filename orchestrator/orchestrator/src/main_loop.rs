use crate::ethereum_event_watcher::check_for_events;
use clarity::PrivateKey as EthPrivateKey;
use clarity::{address::Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::{
    query::{get_oldest_unsigned_transaction_batch, get_oldest_unsigned_valset},
    send::{send_batch_confirm, send_valset_confirm},
};
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use ethereum_peggy::utils::get_peggy_id;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use relayer::batch_relaying::relay_batches;
use relayer::valset_relaying::relay_valsets;
use std::time::Duration;
use std::time::Instant;
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

//const BLOCK_DELAY: u128 = 50;

pub const LOOP_SPEED: Duration = Duration::from_secs(10);

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
pub async fn orchestrator_main_loop(
    cosmos_key: CosmosPrivateKey,
    ethereum_key: EthPrivateKey,
    web3: Web3,
    contact: Contact,
    grpc_client: PeggyQueryClient<Channel>,
    peggy_contract_address: EthAddress,
    pay_fees_in: String,
) {
    let our_cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
    let our_ethereum_address = ethereum_key.to_public_key().unwrap();
    let mut grpc_client = grpc_client;
    let mut last_checked_block: Uint256 = web3.eth_block_number().await.unwrap();
    let fee = Coin {
        denom: pay_fees_in.clone(),
        amount: 1u32.into(),
    };
    let peggy_id = get_peggy_id(peggy_contract_address, our_ethereum_address, &web3).await;
    if peggy_id.is_err() {
        error!("Failed to get PeggyID");
        return;
    }
    let peggy_id = peggy_id.unwrap();
    let peggy_id = String::from_utf8(peggy_id.clone()).expect("Invalid PeggyID");

    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_block_number().await;
        let latest_cosmos_block = contact.get_latest_block_number().await;
        if let (Ok(latest_eth_block), Ok(latest_cosmos_block)) =
            (latest_eth_block, latest_cosmos_block)
        {
            trace!(
                "Latest Eth block {} Latest Cosmos block {}",
                latest_eth_block,
                latest_cosmos_block
            );
        }

        //  Checks for new valsets to sign and relays validator sets from Cosmos -> Ethereum including
        relay_valsets(
            cosmos_key,
            ethereum_key,
            &web3,
            &contact,
            &mut grpc_client,
            peggy_contract_address,
            fee.clone(),
            LOOP_SPEED,
        )
        .await;

        // sign the last unsigned valset, TODO check if we already have signed this
        match get_oldest_unsigned_valset(&mut grpc_client, our_cosmos_address).await {
            Ok(Some(last_unsigned_valset)) => {
                info!("Sending valset confirm for {}", last_unsigned_valset.nonce);
                let res = send_valset_confirm(
                    &contact,
                    ethereum_key,
                    fee.clone(),
                    last_unsigned_valset,
                    cosmos_key,
                    peggy_id.clone(),
                )
                .await;
                info!("Valset confirm result is {:?}", res);
            }
            Ok(None) => trace!("No valset waiting to be signed!"),
            Err(e) => trace!("Failed to get unsigned valsets with {:?}", e),
        }

        // sign the last unsigned batch, TODO check if we already have signed this
        match get_oldest_unsigned_transaction_batch(&mut grpc_client, our_cosmos_address).await {
            Ok(Some(last_unsigned_batch)) => {
                info!("Sending batch confirm for {}", last_unsigned_batch.nonce);
                let res = send_batch_confirm(
                    &contact,
                    ethereum_key,
                    fee.clone(),
                    last_unsigned_batch,
                    cosmos_key,
                    peggy_id.clone(),
                )
                .await
                .unwrap();
                info!("Batch confirm result is {:?}", res);
            }
            Ok(None) => trace!("No unsigned batches! Everything good!"),
            Err(e) => trace!("Failed to get unsigned Batches with {:?}", e),
        }

        relay_batches(
            cosmos_key,
            ethereum_key,
            &web3,
            &contact,
            &mut grpc_client,
            peggy_contract_address,
            fee.clone(),
            LOOP_SPEED,
        )
        .await;

        // Relays events from Ethereum -> Cosmos
        match check_for_events(
            &web3,
            &contact,
            peggy_contract_address,
            cosmos_key,
            fee.clone(),
            last_checked_block.clone(),
        )
        .await
        {
            Ok(new_block) => last_checked_block = new_block,
            Err(e) => error!("Failed to get events for block range {:?}", e),
        }

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < LOOP_SPEED {
            delay_for(LOOP_SPEED - elapsed).await;
        }
    }
}
