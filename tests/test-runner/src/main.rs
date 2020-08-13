//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.
use clarity::Address;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use std::time::Duration;

const TIMEOUT: Duration = Duration::from_secs(1);

fn generate_eth_private_key() -> EthPrivateKey {
    let key_buf: [u8; 32] = rand::random();
    EthPrivateKey::from_slice(&key_buf).unwrap()
}

#[actix_rt::main]
async fn main() {
    println!("Staring Peggy test-runner");
    let contact = Contact::new("http:://localhost:1317", TIMEOUT);

    // TODO
    // update-eth-addr
    // valset-request
    // submit-valset
    // get current valset / specific valset
}
