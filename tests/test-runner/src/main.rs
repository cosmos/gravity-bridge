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
use std::time::Duration;

#[macro_use]
extern crate log;

const TIMEOUT: Duration = Duration::from_secs(1);

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
        println!("the phrase is {}", phrase);
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
    let contact = Contact::new("http:://localhost:1317/", TIMEOUT);
    let keys = get_keys();

    // runs through a full set of rpc tests for each validator, this includes the basics like
    // sending and querying transactions and creating, getting, and setting validator set info
    // this is a little big overkill for a *Peggy* test as opposed to testing the Contact library
    // but it's pretty much free to run.
    for (c_key, e_key) in keys {
        test_rpc_calls(contact.clone(), c_key, e_key)
            .await
            .expect("Failed to test endpoints")
    }

    // TODO
    // update-eth-addr
    // valset-request
    // submit-valset
    // get current valset / specific valset
}
