//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.
use clarity::Address;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::test_rpc_calls;
use contact::client::Contact;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::process::Command;
use std::time::Duration;
use web30::client;

#[macro_use]
extern crate log;

const TIMEOUT: Duration = Duration::from_secs(30);

/// Ethereum keys are generated for every validator inside
/// of this testing application and submitted to the blockchain
/// use the 'update eth address' message.
fn generate_eth_private_key() -> EthPrivateKey {
    let key_buf: [u8; 32] = rand::random();
    EthPrivateKey::from_slice(&key_buf).unwrap()
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
        ret.push((c_key, generate_eth_private_key()))
    }
    ret
}

#[actix_rt::main]
async fn main() {
    println!("Staring Peggy test-runner");
    const COSMOS_NODE: &str = "http://localhost:1317";
    const ETH_NODE: &str = "http://localhost:8545";
    // this key is the private key for the public key defined in tests/assets/ETHGenesis.json
    // where the full node / miner sends its rewards. Therefore it's always going
    // to have a lot of ETH to pay for things like contract deployments
    let miner_private_key: EthPrivateKey =
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
            .parse()
            .unwrap();

    let contact = Contact::new(COSMOS_NODE, TIMEOUT);
    let web30 = web30::client::Web3::new(ETH_NODE, TIMEOUT);
    let keys = get_keys();
    let test_token_name = "footoken".to_string();

    let fee = Coin {
        denom: test_token_name.clone(),
        amount: 1u32.into(),
    };

    for (c_key, e_key) in keys.iter() {
        // set the eth address for all the validators
        contact
            .update_peggy_eth_address(*e_key, *c_key, fee.clone(), None, None, None)
            .await
            .expect("Failed to update eth address!");
    }
    // get the first validator and have them send a valset request
    let (c_key, _e_key) = keys[0];
    let request_block = contact
        .send_valset_request(c_key, fee.clone(), None, None, None)
        .await
        .expect("Failed to send valset request!")
        .height;

    let valset = contact
        .get_peggy_valset_request(request_block)
        .await
        .expect("Failed to get valset!");

    for (c_key, e_key) in keys.iter() {
        // send in valset confirmation for all validators
        let res = contact
            .send_valset_confirm(
                *e_key,
                fee.clone(),
                valset.result.clone(),
                *c_key,
                "foo".to_string(),
                None,
                None,
                None,
            )
            .await;
        res.expect("Failed to send valset confirm!");
    }

    /// TODO valset confirm goes here
    // now we can deploy the test peggy contract, this must come after the
    // first valset is created because the constructor requires this first
    // valset to be submitted.
    let output = Command::new("npx")
        .args(&[
            "ts-node",
            "/peggy/solidity/contract-deployer.ts",
            &format!("--cosmos-node={}", COSMOS_NODE),
            &format!("--eth-node={}", ETH_NODE),
            &format!("--eth-privkey={:#x}", miner_private_key),
            "--contract=/peggy/solidity/artifacts/Peggy.json",
            "--erc20-contract=/peggy/solidity/artifacts/TestERC20.json",
            "--test-mode=true",
        ])
        .current_dir("/peggy/solidity/")
        .output()
        .expect("Failed to deploy contracts!");
    println!("status: {}", output.status);
    println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    println!("stderr: {}", String::from_utf8_lossy(&output.stderr));

    // TODO
    // valset-request-confirm
    // submit-valset
    // get current valset / specific valset
}
