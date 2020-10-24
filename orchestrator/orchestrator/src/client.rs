//! This file is the binary entry point for the Peggy client software, an easy to use cli utility that
//! allows anyone to send funds across the Peggy bridge. Currently this application only does anything
//! on the Ethereum side of the bridge since withdraw batches are incomplete.
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
use crate::main_loop::LOOP_SPEED;
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
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_fees: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
        "Usage: {} --cosmos-dest=<key> --erc20-address=<erc20> --amount=<amount> --contract-address=<addr> --ethereum-key=<key> --ethereum-rpc=<url>
        Options:
            -h --help                     Show this screen.
            --cosmos-dest=<ckey>   The Cosmos address of your destination
            --erc20-address=<erc20>       The Ethereum contract address for the ERC20 token you want to send
            --amount=<amount>             The number of ERC20 tokens you wish to send
            --contract-address=<addr>     The Ethereum contract address for Peggy, this is temporary
            --ethereum-key=<ekey>         The Ethereum private key with the funds you wish to send
            --ethereum-rpc=<eurl>         The Ethereum RPC url, should be a self hosted node
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
    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());
    let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
        .expect("Failed to parse validator key");
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
    )
    .await;
}
