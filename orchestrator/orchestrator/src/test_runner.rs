//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.

// there are several binaries for this crate if we allow dead code on all of them
// we will see functions not used in one binary as dead code. In order to fix that
// we forbid dead code in all but the 'main' binary
#![allow(dead_code)]

#[macro_use]
extern crate log;
#[macro_use]
extern crate lazy_static;

mod batch_relaying;
mod ethereum_event_watcher;
mod main_loop;
mod valset_relaying;

use actix::Arbiter;
use clarity::{Address as EthAddress, Uint256};
use clarity::{PrivateKey as EthPrivateKey, Transaction};
use contact::client::Contact;
use cosmos_peggy::send::send_valset_request;
use cosmos_peggy::send::{request_batch, send_to_eth, update_peggy_eth_address};
use cosmos_peggy::utils::wait_for_cosmos_online;
use cosmos_peggy::utils::wait_for_next_cosmos_block;
use cosmos_peggy::{query::get_oldest_unsigned_transaction_batch, send::send_ethereum_claims};
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::{send_to_cosmos::send_to_cosmos, utils::get_tx_batch_nonce};
use futures::future::join_all;
use main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::SendToCosmosEvent;
use rand::Rng;
use std::io::{BufRead, BufReader, Read, Write};
use std::process::Command;
use std::time::Duration;
use std::{fs::File, time::Instant};
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

/// the timeout for individual requests
const OPERATION_TIMEOUT: Duration = Duration::from_secs(30);
/// the timeout for the total system
const TOTAL_TIMEOUT: Duration = Duration::from_secs(300);

/// Ethereum keys are generated for every validator inside
/// of this testing application and submitted to the blockchain
/// use the 'update eth address' message. In this case we generate
/// them based off of the Cosmos key as the seed so that we can run
/// the test runner multiple times against one chain and get the same keys.
///
/// There's no particular reason to use the public key except that the bytes
/// of the private key type are not public
fn generate_eth_private_key(seed: CosmosPrivateKey) -> EthPrivateKey {
    EthPrivateKey::from_slice(&seed.to_public_key().unwrap().as_bytes()[0..32]).unwrap()
}

/// Validator private keys are generated via the cosmoscli key add
/// command, from there they are used to create gentx's and start the
/// chain, these keys change every time the container is restarted.
/// The mnemonic phrases are dumped into a text file /validator-phrases
/// the phrases are in increasing order, so validator 1 is the first key
/// and so on. While validators may later fail to start it is guaranteed
/// that we have one key for each validator in this file.
fn parse_validator_keys() -> Vec<CosmosPrivateKey> {
    let filename = "/validator-phrases";
    let file = File::open(filename).expect("Failed to find phrases");
    let reader = BufReader::new(file);
    let mut ret = Vec::new();

    for line in reader.lines() {
        let phrase = line.expect("Error reading phrase file!");
        if phrase.is_empty()
            || phrase.contains("write this mnemonic phrase")
            || phrase.contains("recover your account if")
        {
            continue;
        }
        let key = CosmosPrivateKey::from_phrase(&phrase, "").expect("Bad phrase!");
        ret.push(key);
    }
    ret
}

fn get_keys() -> Vec<(CosmosPrivateKey, EthPrivateKey)> {
    let cosmos_keys = parse_validator_keys();
    let mut ret = Vec::new();
    for c_key in cosmos_keys {
        ret.push((c_key, generate_eth_private_key(c_key)))
    }
    ret
}

const COSMOS_NODE: &str = "http://localhost:1317";
const COSMOS_NODE_GRPC: &str = "http://localhost:9090";
const COSMOS_NODE_ABCI: &str = "http://localhost:26657";
const ETH_NODE: &str = "http://localhost:8545";

lazy_static! {
    // this key is the private key for the public key defined in tests/assets/ETHGenesis.json
    // where the full node / miner sends its rewards. Therefore it's always going
    // to have a lot of ETH to pay for things like contract deployments
    static ref MINER_PRIVATE_KEY: EthPrivateKey =
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
            .parse()
            .unwrap();
    static ref MINER_ADDRESS: EthAddress = MINER_PRIVATE_KEY.to_public_key().unwrap();
}

#[actix_rt::main]
async fn main() {
    env_logger::init();
    info!("Staring Peggy test-runner");

    let contact = Contact::new(COSMOS_NODE, OPERATION_TIMEOUT);
    let mut grpc_client = PeggyQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
    let web30 = web30::client::Web3::new(ETH_NODE, OPERATION_TIMEOUT);
    let keys = get_keys();
    let test_token_name = "footoken".to_string();

    let fee = Coin {
        denom: test_token_name.clone(),
        amount: 1u32.into(),
    };

    info!("Waiting for Cosmos chain to come online");
    wait_for_cosmos_online(&contact, TOTAL_TIMEOUT).await;

    // if we detect this env var we are only deploying contracts, do that then exit.
    if option_env!("DEPLOY_CONTRACTS").is_some() {
        info!("test-runner in contract deploying mode, deploying contracts, then exiting");
        deploy_contracts(&contact, &keys, fee).await;
        return;
    }

    let (peggy_address, erc20_address) = parse_contract_addresses();

    // before we start the orchestrators send them some funds so they can pay
    // for things
    send_eth_to_orchestrators(&keys, &web30).await;

    assert!(
        test_check_cosmos_balance(keys[0].0.to_public_key().unwrap().to_address(), &contact)
            .await
            .is_some()
    );

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

    // bootstrapping tests finish here and we move into operational tests

    // send 3 valset updates to make sure the process works back to back
    for _ in 0u32..2 {
        test_valset_update(&contact, &web30, &keys, peggy_address, fee.clone()).await;
    }

    // generate an address for coin sending tests, this ensures test imdepotency
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    let dest_cosmos_private_key = CosmosPrivateKey::from_secret(&secret);
    let dest_cosmos_address = dest_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    let dest_eth_private_key = EthPrivateKey::from_slice(&secret).unwrap();
    let dest_eth_address = dest_eth_private_key.to_public_key().unwrap();

    // the denom and amount of the token bridged from Ethereum -> Cosmos
    // so the denom is the peggy<hash> token name
    // Send a token 3 times
    for _ in 0u32..3 {
        test_erc20_send(
            &web30,
            &contact,
            dest_cosmos_address,
            peggy_address,
            erc20_address,
            100u64.into(),
        )
        .await;
    }

    // We are going to submit a duplicate tx with nonce 1
    // This had better not increase the balance again
    // this test may have false positives if the timeout is not
    // long enough. TODO check for an error on the cosmos send response
    submit_duplicate_erc20_send(
        1u64.into(),
        &contact,
        erc20_address,
        1u64.into(),
        dest_cosmos_address,
        keys.clone(),
        fee.clone(),
    )
    .await;

    // TODO this test is incomplete, the cosmos module is not currently in a state
    // where it will allow it to complete. We send a tx into the Cosmos -> Eth tx pool
    // create a batch with it, sign that batch, and then can not submit it due to failures
    // in the code handling that.
    test_batch(
        &contact,
        &mut grpc_client,
        &web30,
        dest_eth_address,
        peggy_address,
        fee,
        keys[0].0,
        dest_cosmos_private_key,
        erc20_address,
    )
    .await;
}

/// This function deploys the required contracts onto the Ethereum testnet
/// this runs only when the DEPLOY_CONTRACTS env var is set right after
/// the Ethereum test chain starts in the testing environment. We write
/// the stdout of this to a file for later test runs to parse
async fn deploy_contracts(
    contact: &Contact,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    fee: Coin,
) {
    // register all validator eth addresses, currently validators can just not do this
    // a full production version of Peggy would refuse to allow validators to enter the pool
    // without registering their address. It would also allow them to delegate their Cosmos addr
    //
    // Either way, validators need to setup their eth addresses out of band and it's not
    // the orchestrators job. So this isn't exactly where it needs to be in the final version
    // but neither is it really that different.
    let mut updates = Vec::new();
    for (c_key, e_key) in keys.iter() {
        info!(
            "Signing and submitting Eth address {} for validator {}",
            e_key.to_public_key().unwrap(),
            c_key.to_public_key().unwrap().to_address(),
        );
        updates.push(update_peggy_eth_address(
            &contact,
            *e_key,
            *c_key,
            fee.clone(),
        ));
    }
    let update_results = join_all(updates).await;
    for i in update_results {
        i.expect("Failed to update eth address!");
    }

    // prevents the node deployer from failing (rarely) when the chain has not
    // yet produced the next block after submitting each eth address
    wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;

    // wait for the orchestrators to finish registering their eth addresses
    let output = Command::new("npx")
        .args(&[
            "ts-node",
            "/peggy/solidity/contract-deployer.ts",
            &format!("--cosmos-node={}", COSMOS_NODE_ABCI),
            &format!("--eth-node={}", ETH_NODE),
            &format!("--eth-privkey={:#x}", *MINER_PRIVATE_KEY),
            "--contract=/peggy/solidity/artifacts/Peggy.json",
            "--erc20-contract=/peggy/solidity/artifacts/TestERC20.json",
            "--test-mode=true",
        ])
        .current_dir("/peggy/solidity/")
        .output()
        .expect("Failed to deploy contracts!");
    info!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    info!("stderr: {}", String::from_utf8_lossy(&output.stderr));
    let mut file = File::create("/contracts").unwrap();
    file.write_all(&output.stdout).unwrap();
}

/// Parses the ERC20 and Peggy contract addresses from the file created
/// in deploy_contracts()
fn parse_contract_addresses() -> (EthAddress, EthAddress) {
    let mut file =
        File::open("/contracts").expect("Failed to find contracts! did they not deploy?");
    let mut output = String::new();
    file.read_to_string(&mut output).unwrap();
    let mut maybe_peggy_address = None;
    let mut maybe_contract_address = None;
    for line in output.lines() {
        if line.contains("Peggy deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            maybe_peggy_address = Some(address_string.trim().parse().unwrap());
        } else if line.contains("ERC20 deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            maybe_contract_address = Some(address_string.trim().parse().unwrap());
        }
    }
    let peggy_address: EthAddress = maybe_peggy_address.unwrap();
    let erc20_address: EthAddress = maybe_contract_address.unwrap();
    (peggy_address, erc20_address)
}

async fn test_valset_update(
    contact: &Contact,
    web30: &Web3,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    peggy_address: EthAddress,
    fee: Coin,
) {
    // if we don't do this the orchestrators may run ahead of us and we'll be stuck here after
    // getting credit for two loops when we did one
    let starting_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
        .await
        .expect("Failed to get starting eth valset");
    let start = Instant::now();

    // now we send a valset request that the orchestrators will pick up on
    // in this case we send it as the first validator because they can pay the fee
    info!("Sending in valset request");

    // reset here since we might confirm a validator set while sending the next one resulting in an bad sequence
    while Instant::now() - start < TOTAL_TIMEOUT {
        let res = send_valset_request(&contact, keys[0].0, fee.clone()).await;
        if let Ok(res) = res {
            delay_for(Duration::from_secs(2)).await;
            if contact.get_tx_by_hash(&res.txhash).await.is_ok() {
                break;
            }
        }
    }

    let mut current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
        .await
        .expect("Failed to get current eth valset");

    while starting_eth_valset_nonce == current_eth_valset_nonce {
        info!(
            "Validator set is not yet updated to {}>, waiting",
            starting_eth_valset_nonce
        );
        current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
            .await
            .expect("Failed to get current eth valset");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TOTAL_TIMEOUT {
            panic!("Failed to update validator set");
        }
    }
    assert!(starting_eth_valset_nonce != current_eth_valset_nonce);
    info!("Validator set successfully updated!");
}

/// this function tests Ethereum -> Cosmos
async fn test_erc20_send(
    web30: &Web3,
    contact: &Contact,
    dest: CosmosAddress,
    peggy_address: EthAddress,
    erc20_address: EthAddress,
    amount: Uint256,
) {
    let start_coin = check_cosmos_balance(dest, &contact).await;
    info!(
        "Sending to Cosmos from {} to {} with amount {}",
        *MINER_ADDRESS, dest, amount
    );
    // we send some erc20 tokens to the peggy contract to register a deposit
    let tx_id = send_to_cosmos(
        erc20_address,
        peggy_address,
        amount.clone(),
        dest,
        *MINER_PRIVATE_KEY,
        Some(TOTAL_TIMEOUT),
        &web30,
        vec![],
    )
    .await
    .expect("Failed to send tokens to Cosmos");
    info!("Send to Cosmos txid: {:#066x}", tx_id);

    let start = Instant::now();
    while Instant::now() - start < TOTAL_TIMEOUT {
        match (
            start_coin.clone(),
            check_cosmos_balance(dest, &contact).await,
        ) {
            (Some(start_coin), Some(end_coin)) => {
                if start_coin.amount + amount.clone() == end_coin.amount
                    && start_coin.denom == end_coin.denom
                {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, start_coin.denom, end_coin.amount, end_coin.denom
                    );
                    return;
                }
            }
            (None, Some(end_coin)) => {
                if amount == end_coin.amount {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, end_coin.denom, end_coin.amount, end_coin.denom
                    );
                    return;
                } else {
                    panic!("Failed to bridge ERC20!")
                }
            }
            _ => {}
        }
        info!("Waiting for ERC20 deposit");
        wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;
    }
    panic!("Failed to bridge ERC20!")
}

async fn check_cosmos_balance(address: CosmosAddress, contact: &Contact) -> Option<Coin> {
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

async fn test_check_cosmos_balance(address: CosmosAddress, contact: &Contact) -> Option<Coin> {
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

#[allow(clippy::too_many_arguments)]
async fn test_batch(
    contact: &Contact,
    grpc_client: &mut PeggyQueryClient<Channel>,
    web30: &Web3,
    dest_eth_address: EthAddress,
    peggy_address: EthAddress,
    fee: Coin,
    requester_cosmos_private_key: CosmosPrivateKey,
    dest_cosmos_private_key: CosmosPrivateKey,
    erc20_contract: EthAddress,
) {
    let dest_cosmos_address = dest_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    let coin = check_cosmos_balance(dest_cosmos_address, &contact)
        .await
        .unwrap();
    let token_name = coin.denom;
    let amount = coin.amount;

    let bridge_denom_fee = Coin {
        denom: token_name.clone(),
        amount: 1u64.into(),
    };
    let amount = amount - 5u64.into();
    info!(
        "Sending {}{} from {} on Cosmos back to Ethereum",
        amount, token_name, dest_cosmos_address
    );
    let res = send_to_eth(
        dest_cosmos_private_key,
        dest_eth_address,
        Coin {
            denom: token_name.clone(),
            amount: amount.clone(),
        },
        bridge_denom_fee.clone(),
        &contact,
    )
    .await
    .unwrap();
    info!("Sent tokens to Ethereum with {:?}", res);

    info!("Requesting transaction batch");
    request_batch(
        requester_cosmos_private_key,
        token_name.clone(),
        fee.clone(),
        &contact,
    )
    .await
    .unwrap();

    wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;
    let requester_address = requester_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    get_oldest_unsigned_transaction_batch(grpc_client, requester_address)
        .await
        .expect("Failed to get batch to sign");

    let mut current_eth_batch_nonce =
        get_tx_batch_nonce(peggy_address, erc20_contract, *MINER_ADDRESS, &web30)
            .await
            .expect("Failed to get current eth valset");
    let starting_batch_nonce = current_eth_batch_nonce;

    let start = Instant::now();
    while starting_batch_nonce == current_eth_batch_nonce {
        info!(
            "Batch is not yet submitted {}>, waiting",
            starting_batch_nonce
        );
        current_eth_batch_nonce =
            get_tx_batch_nonce(peggy_address, erc20_contract, *MINER_ADDRESS, &web30)
                .await
                .expect("Failed to get current eth tx batch nonce");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TOTAL_TIMEOUT {
            panic!("Failed to submit transaction batch set");
        }
    }

    //
    let txid = web30
        .send_transaction(
            dest_eth_address,
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

    // we have to send this address one eth so that it can perform contract calls
    send_one_eth(dest_eth_address, web30).await;
    assert_eq!(
        web30
            .get_erc20_balance(erc20_contract, dest_eth_address)
            .await
            .unwrap(),
        amount
    );
    info!(
        "Successfully updated txbatch nonce to {} and sent {}{} tokens to Ethereum!",
        current_eth_batch_nonce, amount, token_name
    );
}

// this function submits a EthereumBridgeDepositClaim to the module with a given nonce. This can be set to be a nonce that has
// already been submitted to test the nonce functionality.
#[allow(clippy::too_many_arguments)]
async fn submit_duplicate_erc20_send(
    nonce: Uint256,
    contact: &Contact,
    erc20_address: EthAddress,
    amount: Uint256,
    receiver: CosmosAddress,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    fee: Coin,
) {
    let start_coin = check_cosmos_balance(receiver, &contact)
        .await
        .expect("Did not find coins!");

    let ethereum_sender = "0x912fd21d7a69678227fe6d08c64222db41477ba0"
        .parse()
        .unwrap();

    let event = SendToCosmosEvent {
        event_nonce: nonce,
        erc20: erc20_address,
        sender: ethereum_sender,
        destination: receiver,
        amount,
    };

    // iterate through all validators and try to send an event with duplicate nonce
    for (c_key, _) in keys.iter() {
        let res = send_ethereum_claims(contact, *c_key, vec![event.clone()], vec![], fee.clone())
            .await
            .unwrap();
        trace!("Submitted duplicate sendToCosmos event: {:?}", res);
    }

    if let Some(end_coin) = check_cosmos_balance(receiver, &contact).await {
        if start_coin.amount == end_coin.amount && start_coin.denom == end_coin.denom {
            info!("Successfully failed to duplicate ERC20!");
        } else {
            panic!("Duplicated ERC20!")
        }
    } else {
        panic!("Duplicate test failed for unknown reasons!");
    }
}

/// This overly complex function primarily exists to parallelize the sending of Eth to the
/// orchestrators, waiting for these there transactions takes up nearly a minute of test time
/// and it seemed like low hanging fruit. It in fact was not, mostly because we are sending
/// these tx's from the same address and we therefore need to take into account the correct
/// nonce given the other transactions in flight. This means we need to build the transactions
/// ourselves with that info right here. If you have to modify this seriously consider
/// just calling send_one_eth in a loop.
async fn send_eth_to_orchestrators(keys: &[(CosmosPrivateKey, EthPrivateKey)], web30: &Web3) {
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

async fn send_one_eth(dest: EthAddress, web30: &Web3) {
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
