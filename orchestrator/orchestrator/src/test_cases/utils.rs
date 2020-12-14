use clarity::{Address as EthAddress};
use clarity::{PrivateKey as EthPrivateKey, Transaction};
use contact::client::Contact;
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use futures::future::join_all;
use web30::client::Web3;

use crate::test_cases::MINER_PRIVATE_KEY;
use crate::test_cases::MINER_ADDRESS;
use crate::test_cases::TOTAL_TIMEOUT;


/// This overly complex function primarily exists to parallelize the sending of Eth to the
/// orchestrators, waiting for these there transactions takes up nearly a minute of test time
/// and it seemed like low hanging fruit. It in fact was not, mostly because we are sending
/// these tx's from the same address and we therefore need to take into account the correct
/// nonce given the other transactions in flight. This means we need to build the transactions
/// ourselves with that info right here. If you have to modify this seriously consider
/// just calling send_one_eth in a loop.
pub async fn send_eth_to_orchestrators(keys: &[(CosmosPrivateKey, EthPrivateKey)], web30: &Web3) {
    let balance = web30.eth_get_balance(*MINER_ADDRESS).await.unwrap();
    info!(
        "Sending orchestrators 1 eth to pay for fees miner has {} WEI",
        balance
    );
    let mut eth_addresses = Vec::new();
    for (_, e_key) in keys {
        eth_addresses.push(e_key.to_public_key().unwrap())
    }
    let net_version = web30.net_version().await.unwrap();
    let mut nonce = web30
        .eth_get_transaction_count(*MINER_ADDRESS)
        .await
        .unwrap();
    let mut transactions = Vec::new();
    for address in eth_addresses {
        let t = Transaction {
            to: address,
            nonce: nonce.clone(),
            gas_price: 1_000_000_000u64.into(),
            gas_limit: 24000u64.into(),
            value: 1_000_000_000_000_000_000u128.into(),
            data: Vec::new(),
            signature: None,
        };
        let t = t.sign(&*MINER_PRIVATE_KEY, Some(net_version));
        transactions.push(t);
        nonce += 1u64.into();
    }
    let mut sends = Vec::new();
    for tx in transactions {
        sends.push(web30.eth_send_raw_transaction(tx.to_bytes().unwrap()));
    }
    let txids = join_all(sends).await;
    let mut wait_for_txid = Vec::new();
    for txid in txids {
        let wait = web30.wait_for_transaction(txid.unwrap(), TOTAL_TIMEOUT, None);
        wait_for_txid.push(wait);
    }
    join_all(wait_for_txid).await;
}

pub async fn send_one_eth(dest: EthAddress, web30: &Web3) {
    let txid = web30
        .send_transaction(
            dest,
            Vec::new(),
            1_000_000_000_000_000_000u128.into(),
            *MINER_ADDRESS,
            *MINER_PRIVATE_KEY,
            vec![],
        )
        .await
        .expect("Failed to send Eth to validator {}");
    web30
        .wait_for_transaction(txid, TOTAL_TIMEOUT, None)
        .await
        .unwrap();
}

pub async fn check_cosmos_balance(address: CosmosAddress, contact: &Contact) -> Option<Coin> {
    let account_info = contact.get_balances(address).await.unwrap();
    trace!("Cosmos balance {:?}", account_info.result);
    for coin in account_info.result {
        // make sure the name and amount is correct
        if coin.denom.starts_with("peggy") {
            return Some(coin);
        }
    }
    None
}

pub async fn test_check_cosmos_balance(address: CosmosAddress, contact: &Contact) -> Option<Coin> {
    let account_info = contact.get_balances(address).await.unwrap();
    trace!("Cosmos balance {:?}", account_info.result);
    for coin in account_info.result {
        // make sure the name and amount is correct
        if coin.denom.starts_with("footoken") {
            return Some(coin);
        }
    }
    None
}