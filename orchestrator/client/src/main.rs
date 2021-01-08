//! This file is the binary entry point for the Peggy client software, an easy to use cli utility that
//! allows anyone to send funds across the Peggy bridge. Currently this application only does anything
//! on the Ethereum side of the bridge since withdraw batches are incomplete.

// there are several binaries for this crate if we allow dead code on all of them
// we will see functions not used in one binary as dead code. In order to fix that
// we forbid dead code in all but the 'main' binary
#![allow(dead_code)]

#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;

use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use clarity::Uint256;
use contact::client::Contact;
use cosmos_peggy::send::{send_request_batch, send_to_eth};
use deep_space::address::Address as CosmosAddress;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use docopt::Docopt;
use ethereum_peggy::send_to_cosmos::send_to_cosmos;
use std::time::Duration;
use url::Url;
use web30::client::Web3;

const TIMEOUT: Duration = Duration::from_secs(60);

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_fees: String,
    flag_amount: String,
    flag_cosmos_destination: String,
    flag_erc20_address: String,
    flag_eth_destination: String,
    cmd_eth_to_cosmos: bool,
    cmd_cosmos_to_eth: bool,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage:
        {} cosmos-to-eth --cosmos-phrase=<key> --cosmos-rpc=<url> --fees=<denom> --erc20-address=<addr> --amount=<amount> --eth-destination=<dest>
        {} eth-to-cosmos --ethereum-key=<key> --ethereum-rpc=<url> --contract-address=<addr> --erc20-address=<addr> --amount=<amount> --cosmos-destination=<dest>
        Options:
            -h --help                   Show this screen.
            --cosmos-key=<ckey>         The Cosmos private key of the sender
            --ethereum-key=<ekey>       The Ethereum private key of the sender
            --cosmos-rpc=<curl>         The Cosmos Legacy RPC url, this will need to be manually enabled
            --ethereum-rpc=<eurl>       The Ethereum RPC url, should be a self hosted node
            --fees=<denom>              The Cosmos Denom in which to pay Cosmos chain fees
            --contract-address=<addr>   The Ethereum contract address for Peggy, this is temporary
            --erc20-address=<addr>      An erc20 address to send funds
            --amount=<amount>           The amount of tokens to send
            --cosmos-destination=<dest> A cosmos address to send tokens to
            --eth-destination=<dest> A cosmos address to send tokens to
        About:
            Althea Peggy client software, moves tokens from Ethereum to Cosmos and back
            Written By: {}
            Version {}",
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_NAME"),
            env!("CARGO_PKG_AUTHORS"),
            env!("CARGO_PKG_VERSION"),
        );
}

#[actix_rt::main]
async fn main() {
    env_logger::init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();
    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    let amount: Uint256 = args.flag_amount.parse().unwrap();
    let erc20_address: EthAddress = args
        .flag_erc20_address
        .parse()
        .expect("Invalid contract address!");
    if args.cmd_cosmos_to_eth {
        let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
            .expect("Failed to parse cosmos key phrase, does it have a password?");
        let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
        let cosmos_url = cosmos_url.to_string();
        let cosmos_url = cosmos_url.trim_end_matches('/');
        let fee_denom = args.flag_fees;
        let fee = Coin {
            denom: fee_denom,
            amount: 1u64.into(),
        };
        let peggy_denom = format!("peggy/{}", erc20_address);
        let contact = Contact::new(&cosmos_url, TIMEOUT);
        let amount = Coin {
            amount,
            denom: peggy_denom.clone(),
        };
        let eth_dest: EthAddress = args.flag_eth_destination.parse().unwrap();

        println!("Locking funds into the batch pool");
        send_to_eth(cosmos_key, eth_dest, amount, fee.clone(), &contact)
            .await
            .expect("Failed to Send to ETH");

        println!("Requesting a batch to push transaction along immediately");
        send_request_batch(cosmos_key, peggy_denom, fee, &contact)
            .await
            .expect("Failed to request batch");
    } else if args.cmd_eth_to_cosmos {
        let ethereum_key: EthPrivateKey = args
            .flag_ethereum_key
            .parse()
            .expect("Invalid Ethereum private key!");
        let contract_address: EthAddress = args
            .flag_contract_address
            .parse()
            .expect("Invalid contract address!");
        let eth_url = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");
        let eth_url = eth_url.to_string();
        let eth_url = eth_url.trim_end_matches('/');
        let web3 = Web3::new(&eth_url, TIMEOUT);
        let cosmos_dest: CosmosAddress = args.flag_cosmos_destination.parse().unwrap();

        let ethereum_public_key = ethereum_key.to_public_key().unwrap();

        println!(
            "Sending to Cosmos from {} to {} with amount {}",
            ethereum_public_key, cosmos_dest, amount
        );
        // we send some erc20 tokens to the peggy contract to register a deposit
        let tx_id = send_to_cosmos(
            erc20_address,
            contract_address,
            amount.clone(),
            cosmos_dest,
            ethereum_key,
            Some(TIMEOUT),
            &web3,
            vec![],
        )
        .await
        .expect("Failed to send tokens to Cosmos");
        println!("Send to Cosmos txid: {:#066x}", tx_id);
    }
}
