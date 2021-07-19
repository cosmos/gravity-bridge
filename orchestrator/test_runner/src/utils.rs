use crate::TOTAL_TIMEOUT;
use crate::{MINER_PRIVATE_KEY};
use crate::{MINER_ADDRESS, OPERATION_TIMEOUT};
use clarity::{Address as EthAddress, Uint256};
use clarity::{PrivateKey as EthPrivateKey, Transaction};
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use futures::future::join_all;
use rand::Rng;
use web30::{client::Web3, types::SendTxOption};

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

pub async fn check_cosmos_balance(
    denom: &str,
    address: CosmosAddress,
    contact: &Contact,
) -> Option<Coin> {
    let account_info = contact.get_balances(address).await.unwrap();
    trace!("Cosmos balance {:?}", account_info);
    for coin in account_info {
        // make sure the name and amount is correct
        if coin.denom.starts_with(denom) {
            return Some(coin);
        }
    }
    None
}

/// This function efficiently distributes ERC20 tokens to a large number of provided Ethereum addresses
/// the real problem here is that you can't do more than one send operation at a time from a
/// single address without your sequence getting out of whack. By manually setting the nonce
/// here we can send thousands of transactions in only a few blocks
pub async fn send_erc20_bulk(
    amount: Uint256,
    erc20: EthAddress,
    destinations: &[EthAddress],
    web3: &Web3,
) {
    let miner_balance = web3.get_erc20_balance(erc20, *MINER_ADDRESS).await.unwrap();
    assert!(miner_balance > amount.clone() * destinations.len().into());
    let mut nonce = web3
        .eth_get_transaction_count(*MINER_ADDRESS)
        .await
        .unwrap();
    let mut transactions = Vec::new();
    for address in destinations {
        let send = web3.erc20_send(
            amount.clone(),
            *address,
            erc20,
            *MINER_PRIVATE_KEY,
            Some(OPERATION_TIMEOUT),
            vec![
                SendTxOption::Nonce(nonce.clone()),
                SendTxOption::GasLimit(100_000u32.into()),
            ],
        );
        transactions.push(send);
        nonce += 1u64.into();
    }
    let _txids = join_all(transactions).await;
    for address in destinations {
        let new_balance = web3.get_erc20_balance(erc20, *address).await.unwrap();
        assert!(new_balance >= amount.clone());
    }
}

/// This function efficiently distributes ETH to a large number of provided Ethereum addresses
/// the real problem here is that you can't do more than one send operation at a time from a
/// single address without your sequence getting out of whack. By manually setting the nonce
/// here we can quickly send thousands of transactions in only a few blocks
pub async fn send_eth_bulk(amount: Uint256, destinations: &[EthAddress], web3: &Web3) {
    let net_version = web3.net_version().await.unwrap();
    let mut nonce = web3
        .eth_get_transaction_count(*MINER_ADDRESS)
        .await
        .unwrap();
    let mut transactions = Vec::new();
    for address in destinations {
        let t = Transaction {
            to: *address,
            nonce: nonce.clone(),
            gas_price: 1_000_000_000u64.into(),
            gas_limit: 24000u64.into(),
            value: amount.clone(),
            data: Vec::new(),
            signature: None,
        };
        let t = t.sign(&*MINER_PRIVATE_KEY, Some(net_version));
        transactions.push(t);
        nonce += 1u64.into();
    }
    let mut sends = Vec::new();
    for tx in transactions {
        sends.push(web3.eth_send_raw_transaction(tx.to_bytes().unwrap()));
    }
    let txids = join_all(sends).await;
    let mut wait_for_txid = Vec::new();
    for txid in txids {
        let wait = web3.wait_for_transaction(txid.unwrap(), TOTAL_TIMEOUT, None);
        wait_for_txid.push(wait);
    }
    join_all(wait_for_txid).await;
}

pub fn get_user_key() -> BridgeUserKey {
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the starting location of the funds
    let eth_key = EthPrivateKey::from_slice(&secret).unwrap();
    let eth_address = eth_key.to_public_key().unwrap();
    // the destination on cosmos that sends along to the final ethereum destination
    let cosmos_key = CosmosPrivateKey::from_secret(&secret);
    let cosmos_address = cosmos_key
        .to_address(CosmosAddress::DEFAULT_PREFIX)
        .unwrap();
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the final destination of the tokens back on Ethereum
    let eth_dest_key = EthPrivateKey::from_slice(&secret).unwrap();
    let eth_dest_address = eth_key.to_public_key().unwrap();
    BridgeUserKey {
        eth_address,
        eth_key,
        cosmos_address,
        cosmos_key,
        eth_dest_key,
        eth_dest_address,
    }
}
#[derive(Debug)]
pub struct BridgeUserKey {
    // the starting addresses that get Eth balances to send across the bridge
    pub eth_address: EthAddress,
    pub eth_key: EthPrivateKey,
    // the cosmos addresses that get the funds and send them on to the dest eth addresses
    pub cosmos_address: CosmosAddress,
    pub cosmos_key: CosmosPrivateKey,
    // the location tokens are sent back to on Ethereum
    pub eth_dest_address: EthAddress,
    pub eth_dest_key: EthPrivateKey,
}

#[derive(Debug, Clone)]
pub struct ValidatorKeys {
    /// The Ethereum key used by this validator to sign Gravity bridge messages
    pub eth_key: EthPrivateKey,
    /// The Orchestrator key used by this validator to submit oracle messages and signatures
    /// to the cosmos chain
    pub orch_key: CosmosPrivateKey,
    /// The validator key used by this validator to actually sign and produce blocks
    pub validator_key: CosmosPrivateKey,
}
