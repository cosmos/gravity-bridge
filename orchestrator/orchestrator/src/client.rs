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

use crate::main_loop::LOOP_SPEED;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use clarity::Uint256;
use contact::client::Contact;
use cosmos_peggy::send::send_valset_request;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use docopt::Docopt;
use ethereum_peggy::send_to_cosmos::send_to_cosmos;
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
    flag_amount: String,
    flag_erc20: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --cosmos-phrase=<key> --ethereum-key=<key> --cosmos-rpc=<url> --ethereum-rpc=<url> --fees=<denom> --contract-address=<addr> --erc20-address=<addr> --amount=<amount>
        Options:
            -h --help                 Show this screen.
            --cosmos-key=<ckey>       The Cosmos private key of the validator
            --ethereum-key=<ekey>     The Ethereum private key of the validator
            --cosmos-rpc=<curl>       The Cosmos RPC url, usually the validator
            --ethereum-rpc=<eurl>     The Ethereum RPC url, should be a self hosted node
            --fees=<denom>            The Cosmos Denom in which to pay Cosmos chain fees
            --contract-address=<addr> The Ethereum contract address for Peggy, this is temporary
            --erc20-address=<addr>    An erc20 address to send funds
            --amount=<amount>
        About:
            Client software for Althea-Peggy.
            Written By: {}
            Version {}",
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_AUTHORS"),
            env!("CARGO_PKG_VERSION"),
        );
}

const TIMEOUT: Duration = Duration::from_secs(60);

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
    let erc20_address: EthAddress = args.flag_erc20.parse().expect("Invalid contract address!");
    let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
    let cosmos_url = cosmos_url.to_string();
    let cosmos_url = cosmos_url.trim_end_matches('/');
    let eth_url = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");
    let eth_url = eth_url.to_string();
    let eth_url = eth_url.trim_end_matches('/');
    let fee_denom = args.flag_fees;
    let amount: Uint256 = args.flag_amount.parse().unwrap();

    let web3 = Web3::new(&eth_url, LOOP_SPEED);
    let contact = Contact::new(&cosmos_url, LOOP_SPEED);
    let fee = Coin {
        denom: fee_denom,
        amount: 1u64.into(),
    };

    let cosmos_public_key = cosmos_key.to_public_key().unwrap().to_address();
    let ethereum_public_key = ethereum_key.to_public_key().unwrap();

    info!("Sending in valset request");
    let _res = send_valset_request(&contact, cosmos_key, fee, None, None, None)
        .await
        .expect("Failed to send valset request");

    let dest = cosmos_public_key;
    info!(
        "Sending to Cosmos from {} to {} with amount {}",
        ethereum_public_key, dest, amount
    );
    // we send some erc20 tokens to the peggy contract to register a deposit
    let tx_id = send_to_cosmos(
        erc20_address,
        contract_address,
        amount.clone(),
        dest,
        ethereum_key,
        Some(TIMEOUT),
        &web3,
        vec![],
    )
    .await
    .expect("Failed to send tokens to Cosmos");
    info!("Send to Cosmos txid: {:#066x}", tx_id);
}
