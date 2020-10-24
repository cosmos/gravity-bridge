//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.

#[macro_use]
extern crate log;

mod ethereum_event_watcher;
mod main_loop;
mod tests;
mod valset_relaying;

use actix::Arbiter;
use clarity::{
    abi::{derive_signature, encode_call},
    utils::bytes_to_hex_str,
    PrivateKey as EthPrivateKey,
};
use clarity::{Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::send::send_valset_request;
use cosmos_peggy::send::update_peggy_eth_address;
use cosmos_peggy::utils::wait_for_cosmos_online;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use ethereum_peggy::send_to_cosmos::send_to_cosmos;
use ethereum_peggy::utils::get_valset_nonce;
use main_loop::orchestrator_main_loop;
use num::Bounded;
use std::io::{BufRead, BufReader};
use std::process::Command;
use std::time::Duration;
use std::{fs::File, time::Instant};
use tokio::time::delay_for;
use web30::{
    client::Web3,
    jsonrpc::error::Web3Error,
    types::NewFilter,
    types::{Log, SendTxOption},
};

const TIMEOUT: Duration = Duration::from_secs(60);

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

#[actix_rt::main]
async fn main() {
    env_logger::init();
    info!("Staring Peggy test-runner");
    const COSMOS_NODE: &str = "http://localhost:1317";
    const ETH_NODE: &str = "http://localhost:8545";
    const PEGGY_ID: &str = "foo";
    // this key is the private key for the public key defined in tests/assets/ETHGenesis.json
    // where the full node / miner sends its rewards. Therefore it's always going
    // to have a lot of ETH to pay for things like contract deployments
    let miner_private_key: EthPrivateKey =
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
            .parse()
            .unwrap();
    let miner_address = miner_private_key.to_public_key().unwrap();

    let contact = Contact::new(COSMOS_NODE, TIMEOUT);
    let web30 = web30::client::Web3::new(ETH_NODE, TIMEOUT);
    let keys = get_keys();
    let test_token_name = "footoken".to_string();

    let fee = Coin {
        denom: test_token_name.clone(),
        amount: 1u32.into(),
    };

    info!("Waiting for Cosmos chain to come online");
    wait_for_cosmos_online(&contact).await;

    // register all validator eth addresses, currently validators can just not do this
    // a full production version of Peggy would refuse to allow validators to enter the pool
    // without registering their address. It would also allow them to delegate their Cosmos addr
    //
    // Either way, validators need to setup their eth addresses out of band and it's not
    // the orchestrators job. So this isn't exactly where it needs to be in the final version
    // but neither is it really that different.
    for (c_key, e_key) in keys.iter() {
        info!(
            "Signing and submitting Eth address {} for validator {}",
            e_key.to_public_key().unwrap(),
            c_key.to_public_key().unwrap().to_address(),
        );
        update_peggy_eth_address(&contact, *e_key, *c_key, fee.clone(), None, None, None)
            .await
            .expect("Failed to update Eth address");
    }

    // wait for the orchestrators to finish registering their eth addresses
    let output = Command::new("npx")
        .args(&[
            "ts-node",
            "/peggy/solidity/contract-deployer.ts",
            &format!("--cosmos-node={}", COSMOS_NODE),
            &format!("--eth-node={}", ETH_NODE),
            &format!("--eth-privkey={:#x}", miner_private_key),
            &format!("--peggy-id={}", PEGGY_ID),
            "--contract=/peggy/solidity/artifacts/Peggy.json",
            "--erc20-contract=/peggy/solidity/artifacts/TestERC20.json",
            "--test-mode=true",
        ])
        .current_dir("/peggy/solidity/")
        .output()
        .expect("Failed to deploy contracts!");
    info!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    info!("stderr: {}", String::from_utf8_lossy(&output.stderr));

    let mut maybe_peggy_address = None;
    let mut maybe_contract_address = None;
    for line in String::from_utf8_lossy(&output.stdout).lines() {
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

    // before we start the orchestrators send them some funds so they can pay
    // for things
    for (_c_key, e_key) in keys.iter() {
        let validator_eth_address = e_key.to_public_key().unwrap();

        let balance = web30.eth_get_balance(miner_address).await.unwrap();
        info!(
            "Sending orchestrator 1 eth to pay for fees miner has {} WEI",
            balance
        );
        // send every orchestrator 1 eth to pay for fees
        let txid = web30
            .send_transaction(
                validator_eth_address,
                Vec::new(),
                1_000_000_000_000_000_000u128.into(),
                miner_address,
                miner_private_key,
                vec![],
            )
            .await
            .expect("Failed to send Eth to validator {}");
        web30
            .wait_for_transaction(txid, TIMEOUT, None)
            .await
            .unwrap();
    }

    // start orchestrators, send them some eth so that they can pay for things
    for (c_key, e_key) in keys.iter() {
        info!("Spawning Orchestrator");
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        Arbiter::spawn(orchestrator_main_loop(
            *c_key,
            *e_key,
            web30.clone(),
            contact.clone(),
            peggy_address,
            test_token_name.clone(),
        ));
    }

    // bootstrapping tests finish here and we move into operational tests

    // send 3 valset updates to make sure the process works back to back
    for _ in 0u32..2 {
        test_valset_update(
            &contact,
            &web30,
            &keys,
            peggy_address,
            miner_address,
            fee.clone(),
        )
        .await;
    }

    test_erc20_send(
        &contact,
        &web30,
        &keys,
        peggy_address,
        erc20_address,
        miner_private_key,
        miner_address,
        fee,
    )
    .await;
}

async fn test_valset_update(
    contact: &Contact,
    web30: &Web3,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    peggy_address: EthAddress,
    miner_address: EthAddress,
    fee: Coin,
) {
    // if we don't do this the orchestrators may run ahead of us and we'll be stuck here after
    // getting credit for two loops when we did one
    let starting_eth_valset_nonce = get_valset_nonce(peggy_address, miner_address, &web30)
        .await
        .expect("Failed to get starting eth valset");

    // now we send a valset request that the orchestrators will pick up on
    // in this case we send it as the first validator because they can pay the fee
    info!("Sending in valset request");
    let _res = send_valset_request(&contact, keys[0].0, fee, TIMEOUT * 2)
        .await
        .expect("Failed to send valset request");

    let mut current_eth_valset_nonce = get_valset_nonce(peggy_address, miner_address, &web30)
        .await
        .expect("Failed to get current eth valset");

    let start = Instant::now();
    while starting_eth_valset_nonce == current_eth_valset_nonce {
        info!(
            "Validator set is not yet updated to {}>, waiting",
            starting_eth_valset_nonce
        );
        current_eth_valset_nonce = get_valset_nonce(peggy_address, miner_address, &web30)
            .await
            .expect("Failed to get current eth valset");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TIMEOUT {
            panic!("Failed to update validator set");
        }
    }
    assert!(starting_eth_valset_nonce != current_eth_valset_nonce);
    info!("Validator set successfully updated!");
}

async fn test_erc20_send(
    contact: &Contact,
    web30: &Web3,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    peggy_address: EthAddress,
    erc20_address: EthAddress,
    miner_private_key: EthPrivateKey,
    miner_address: EthAddress,
    fee: Coin,
) {
    let dest = keys[0].0.to_public_key().unwrap().to_address();
    let amount: Uint256 = 1u64.into();
    info!(
        "Sending to Cosmos from {} to {} with amount {}",
        miner_address, dest, amount
    );
    // we send some erc20 tokens to the peggy contract to register a deposit
    let tx_id = send_to_cosmos(
        erc20_address,
        peggy_address,
        amount.clone(),
        dest,
        miner_private_key,
        Some(TIMEOUT),
        &web30,
        vec![],
    )
    .await
    .expect("Failed to send tokens to Cosmos");
    info!("Send to Cosmos txid: {:#066x}", tx_id);

    delay_for(TIMEOUT).await;

    let account_info = contact.get_account_info(dest).await.unwrap();
    let mut success = false;
    info!("Account info: {:?}", account_info);
    for coin in account_info.result.value.coins {
        if coin.denom.starts_with('p') && coin.denom.ends_with("maxx") && coin.amount == amount {
            success = true;
            break;
        }
    }
    assert!(success);
    info!("Successfully bridged ERC20 to Cosmos!");
}
