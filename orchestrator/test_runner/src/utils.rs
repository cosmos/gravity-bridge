use crate::get_test_token_name;
use crate::COSMOS_NODE_GRPC;
use crate::ETH_NODE;
use crate::TOTAL_TIMEOUT;
use crate::{one_eth, MINER_PRIVATE_KEY};
use crate::{MINER_ADDRESS, OPERATION_TIMEOUT};
use actix::System;
use clarity::{Address as EthAddress, Uint256};
use clarity::{PrivateKey as EthPrivateKey, Transaction};
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use futures::future::join_all;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use orchestrator::main_loop::orchestrator_main_loop;
use rand::Rng;
use std::thread;
use web30::{client::Web3, types::SendTxOption};

/// This overly complex function primarily exists to parallelize the sending of Eth to the
/// orchestrators, waiting for these there transactions takes up nearly a minute of test time
/// and it seemed like low hanging fruit. It in fact was not, mostly because we are sending
/// these tx's from the same address and we therefore need to take into account the correct
/// nonce given the other transactions in flight. This means we need to build the transactions
/// ourselves with that info right here. If you have to modify this seriously consider
/// just calling send_one_eth in a loop.
pub async fn send_eth_to_orchestrators(keys: &[ValidatorKeys], web30: &Web3) {
    let balance = web30.eth_get_balance(*MINER_ADDRESS).await.unwrap();
    info!(
        "Sending orchestrators 100 eth to pay for fees miner has {} ETH",
        balance / one_eth()
    );
    let mut eth_addresses = Vec::new();
    for k in keys {
        let e_key = k.eth_key;
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
            value: one_eth() * 100u16.into(),
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

/// This function pays the piper for the strange concurrency model that we use for the tests
/// we launch a thread, create an actix executor and then start the orchestrator within that scope
/// previously we could just throw around futures to spawn despite not having 'send' newer versions
/// of Actix rightly forbid this and we have to take the time to handle it here.
pub async fn start_orchestrators(
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
    validator_out: bool,
) {
    // used to break out of the loop early to simulate one validator
    // not running an Orchestrator
    let num_validators = keys.len();
    let mut count = 0;

    #[allow(clippy::explicit_counter_loop)]
    for k in keys {
        info!(
            "Spawning Orchestrator with delegate keys {} {} and validator key {}",
            k.eth_key.to_public_key().unwrap(),
            k.orch_key
                .to_address(CosmosAddress::DEFAULT_PREFIX)
                .unwrap(),
            k.validator_key.to_address("cosmosvaloper").unwrap()
        );
        let grpc_client = GravityQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        thread::spawn(move || {
            let web30 = web30::client::Web3::new(ETH_NODE, OPERATION_TIMEOUT);
            let contact = Contact::new(
                COSMOS_NODE_GRPC,
                OPERATION_TIMEOUT,
                CosmosAddress::DEFAULT_PREFIX,
            )
            .unwrap();
            let fut = orchestrator_main_loop(
                k.orch_key,
                k.eth_key,
                web30,
                contact,
                grpc_client,
                gravity_address,
                get_test_token_name(),
            );
            let system = System::new();
            system.block_on(fut);
        });
        // used to break out of the loop early to simulate one validator
        // not running an orchestrator
        count += 1;
        if validator_out && count == num_validators - 1 {
            break;
        }
    }
}
