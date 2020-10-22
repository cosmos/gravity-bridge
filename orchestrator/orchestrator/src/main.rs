//! Orchestrator is a sort of specialized relayer for Althea-Peggy that runs on every validator.
//! Things this binary is responsible for
//!   * Performing all the Ethereum signing required to submit updates and generate batches
//!   * Progressing the validator set update generation process.
//!   * Observing events on the Ethereum chain and submitting oracle messages for validator consensus
//! Things this binary needs
//!   * Access to the validators signing Ethereum key
//!   * Access to the validators Cosmos key
//!   * Access to an Cosmos chain RPC server
//!   * Access to an Ethereum chain RPC server

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[macro_use]
extern crate log;

mod ethereum_event_watcher;
mod main_loop;
mod tests;
mod valset_relaying;

use crate::main_loop::orchestrator_main_loop;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use docopt::Docopt;
use std::time::Duration;
use url::Url;
use web30::client::Web3;

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_key: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_fees: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
        "Usage: {} --cosmos-key=<key> --ethereum-key=<key> --cosmos-rpc=<url> --ethereum-rpc=<url>
        Options:
            -h --help                 Show this screen.
            --cosmos-key=<ckey>       The Cosmos private key of the validator
            --ethereum-key=<ekey>     The Ethereum private key of the validator
            --cosmos-rpc=<curl>       The Cosmos RPC url, usually the validator
            --ethereum-rpc=<eurl>     The Ethereum RPC url, should be a self hosted node
            --fees=<denom>            The Cosmos Denom in which to pay Cosmos chain fees
            --contract-address=<addr> The Ethereum contract address for Peggy, this is temporary
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
    let contract_address: EthAddress = args
        .flag_contract_address
        .parse()
        .expect("Invalid contract address!");
    let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
    let eth_url = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");
    let fee_denom = args.flag_fees;

    let web3 = Web3::new(&eth_url.to_string(), LOOP_SPEED);
    let contact = Contact::new(&cosmos_url.to_string(), LOOP_SPEED);

    orchestrator_main_loop(
        cosmos_key,
        ethereum_key,
        web3,
        contact,
        contract_address,
        fee_denom,
        LOOP_SPEED,
    )
    .await;
}
