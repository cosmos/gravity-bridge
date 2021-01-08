use crate::{utils::*, COSMOS_NODE_GRPC, TOTAL_TIMEOUT};
use actix::{clock::delay_for, Arbiter};
use clarity::PrivateKey as EthPrivateKey;
use clarity::{Address as EthAddress, Uint256};
use contact::client::Contact;
use deep_space::address::Address as CosmosAddress;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use ethereum_peggy::send_to_cosmos::send_to_cosmos;
use futures::future::join_all;
use orchestrator::main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use rand::Rng;
use std::time::{Duration, Instant};
use web30::client::Web3;

const TIMEOUT: Duration = Duration::from_secs(120);

pub fn one_eth() -> Uint256 {
    1000000000000000000u128.into()
}

pub struct BridgeUserKey {
    pub eth_address: EthAddress,
    pub eth_key: EthPrivateKey,
    pub cosmos_address: CosmosAddress,
    pub cosmos_key: CosmosPrivateKey,
}

/// Perform a stress test by sending thousands of
/// transactions and producing large batches
#[allow(clippy::too_many_arguments)]
pub async fn transaction_stress_test(
    web30: &Web3,
    contact: &Contact,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    peggy_address: EthAddress,
    test_token_name: String,
    erc20_addresses: Vec<EthAddress>,
) {
    // start orchestrators
    for (c_key, e_key) in keys.iter() {
        info!("Spawning Orchestrator");
        let grpc_client = PeggyQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        Arbiter::spawn(orchestrator_main_loop(
            *c_key,
            *e_key,
            web30.clone(),
            contact.clone(),
            grpc_client,
            peggy_address,
            test_token_name.clone(),
        ));
    }

    // Generate 100 user keys to send ETH and multiple types of tokens
    let mut user_keys = Vec::new();
    for _ in 0..100 {
        let mut rng = rand::thread_rng();
        let secret: [u8; 32] = rng.gen();
        let cosmos_key = CosmosPrivateKey::from_secret(&secret);
        let cosmos_address = cosmos_key.to_public_key().unwrap().to_address();
        let eth_key = EthPrivateKey::from_slice(&secret).unwrap();
        let eth_address = eth_key.to_public_key().unwrap();
        user_keys.push(BridgeUserKey {
            eth_address,
            eth_key,
            cosmos_address,
            cosmos_key,
        })
    }
    let eth_destinations: Vec<EthAddress> = user_keys.iter().map(|i| i.eth_address).collect();
    send_eth_bulk(one_eth(), &eth_destinations, web30).await;
    info!("Sent {} addresses 1 ETH", user_keys.len());
    for token in erc20_addresses.iter() {
        send_erc20_bulk(one_eth(), *token, &eth_destinations, web30).await;
        info!("Sent {} addresses 1 {}", user_keys.len(), token);
    }
    for token in erc20_addresses.iter() {
        let mut sends = Vec::new();
        for keys in user_keys.iter() {
            let fut = send_to_cosmos(
                *token,
                peggy_address,
                one_eth(),
                keys.cosmos_address,
                keys.eth_key,
                Some(TIMEOUT),
                web30,
                Vec::new(),
            );
            sends.push(fut);
        }
        let txids = join_all(sends).await;
        let mut wait_for_txid = Vec::new();
        for txid in txids {
            let wait = web30.wait_for_transaction(txid.unwrap(), TIMEOUT, None);
            wait_for_txid.push(wait);
        }
        let results = join_all(wait_for_txid).await;
        for result in results {
            let result = result.unwrap();
            result.block_number.unwrap();
        }
        info!("Locked 1 {} from {} into Peggy", token, user_keys.len());
    }

    let start = Instant::now();
    while Instant::now() - start < TOTAL_TIMEOUT {
        let mut good = true;
        for keys in user_keys.iter() {
            let c_addr = keys.cosmos_address;
            let balances = contact.get_balances(c_addr).await.unwrap().result;
            for token in erc20_addresses.iter() {
                let mut found = false;
                for balance in balances.iter() {
                    if balance.denom.contains(&token.to_string()) {
                        found = true;
                    }
                }
                if !found {
                    good = false;
                }
            }
        }
        if good {
            info!("All {} deposits bridged successfully!", user_keys.len());
            break;
        }
        delay_for(Duration::from_secs(1)).await;
    }
}
