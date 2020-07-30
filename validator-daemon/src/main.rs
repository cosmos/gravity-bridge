//! Validator-daemon is a sort of specialized relayer for Althea-Peggy that runs on every validator.
//! Things this binary is responsible for
//!   * Performing all the Ethereum signing required to submit updates and generate batches
//!   * Progressing the validator set update generation process.
//!   * Observing events on the Ethereum chain and submitting oracle messages for validator consensus
//! Things this binary needs
//!   * Access to the validators signing Ethereum key
//!   * Access to the validators Cosmos key
//!   * Access to an Cosmos chain RPC server
//!   * Access to an Ethereum chain RPC server

use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use docopt::Docopt;
use num256::Int256;
use std::thread;
use std::time::Duration;
use std::time::Instant;
use url::Url;
use web30::client::Web3;

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_key: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_ethereum_rpc: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
        "Usage: {} --cosmos-key=<key> --ethereum-key=<key> --cosmos-rpc=<url> --ethereum-rpc=<url>
        Options:
            -h --help              Show this screen.
            --cosmos-key=<ckey>    The Cosmos private key of the validator
            --ethereum-key=<ekey>  The Ethereum private key of the validator
            --cosmos-rpc=<curl>    The Cosmos RPC url, usually the validator
            --ethereum-rpc=<eurl>  The Ethereum RPC url, should be a self hosted node
        About:
            The Validator companion relayer and Ethereum network observer.
            for Althea-Peggy.
            Written By: {}
            Version {}",
        env!("CARGO_PKG_NAME"),
        env!("CARGO_PKG_AUTHORS"),
        env!("CARGO_PKG_VERSION"),
    );
}

const LOOP_SPEED: Duration = Duration::from_secs(5);

#[actix_rt::main]
async fn main() {
    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    let cosmos_key: CosmosPrivateKey = args
        .flag_cosmos_key
        .parse()
        .expect("Invalid Private Cosmos Key!");
    let ethereum_key: EthPrivateKey = args
        .flag_ethereum_key
        .parse()
        .expect("Invalid Ethereum private key!");
    let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
    let eth_url = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");

    let web3 = Web3::new(&eth_url.to_string(), LOOP_SPEED);
    let contact = Contact::new(&cosmos_url.to_string(), LOOP_SPEED);

    loop {
        let loop_start = Instant::now();

        let latest_eth_block = web3.eth_get_latest_block().await.unwrap();
        let latest_cosmos_block = contact.get_latest_block().await.unwrap();
        println!(
            "Latest Eth block {} Latest Cosmos block {}",
            latest_eth_block.number, latest_cosmos_block.block.header.version.block
        );

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < LOOP_SPEED {
            thread::sleep(LOOP_SPEED - elapsed)
        }
    }
}
