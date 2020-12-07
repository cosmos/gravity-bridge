//! This file is a single use binary that will allow you to register your validator ethereum key

// there are several binaries for this crate if we allow dead code on all of them
// we will see functions not used in one binary as dead code. In order to fix that
// we forbid dead code in all but the 'main' binary
#![allow(dead_code)]

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

mod batch_relaying;
mod ethereum_event_watcher;
mod main_loop;
mod valset_relaying;

use crate::main_loop::LOOP_SPEED;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use cosmos_peggy::send::update_peggy_eth_address;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use docopt::Docopt;
use url::Url;

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_fees: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
        "Usage: {} --cosmos-phrase=<key> --ethereum-key=<key> --cosmos-rpc=<url> --fees=<denom>
        Options:
            -h --help                     Show this screen.
            --cosmos-phrase=<ckey>    The Cosmos private key of the validator
            --ethereum-key=<ekey>     The Ethereum private key of the validator
            --cosmos-rpc=<curl>       The Cosmos RPC url, usually the validator
            --fees=<denom>            The Cosmos Denom in which to pay Cosmos chain fees
        About:
            Client software for the Althea-Peggy bridge
            Written By: {}
            Version {}",
        env!("CARGO_PKG_NAME"),
        env!("CARGO_PKG_AUTHORS"),
        env!("CARGO_PKG_VERSION"),
    );
}

#[actix_rt::main]
async fn main() {
    env_logger::init();

    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
        .expect("Failed to parse validator key");
    let ethereum_key: EthPrivateKey = args
        .flag_ethereum_key
        .parse()
        .expect("Invalid Ethereum private key!");
    let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
    let cosmos_url = cosmos_url.to_string();
    let cosmos_url = cosmos_url.trim_end_matches('/');
    let fee_denom = args.flag_fees;

    let contact = Contact::new(&cosmos_url, LOOP_SPEED);
    let fee = Coin {
        denom: fee_denom,
        amount: 1u64.into(),
    };

    update_peggy_eth_address(
        &contact,
        ethereum_key,
        cosmos_key,
        fee.clone(),
    )
    .await
    .expect("Failed to update Eth address");
}
