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
use std::{time::Duration, u128};
use url::Url;
use web30::client::Web3;

const TIMEOUT: Duration = Duration::from_secs(60);

pub fn one_eth() -> f64 {
    1000000000000000000f64
}

/// TODO revisit this for higher precision while
/// still representing the number to the user as a float
pub fn fraction_eth_to_wei(num: f64) -> Uint256 {
    let mut res = num;
    // in order to avoid floating point rounding issues we
    // multiply only by 10 each time. this reduces the rounding
    // errors enough to be ignored
    for _ in 0..18 {
        res *= 10f64
    }
    (res as u128).into()
}

pub fn print_eth(input: Uint256) -> String {
    let float: f64 = input.to_string().parse().unwrap();
    let res = float / one_eth();
    format!("{}", res)
}

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_rpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_amount: f64,
    flag_cosmos_destination: String,
    flag_erc20_address: String,
    flag_eth_destination: String,
    flag_no_batch: bool,
    flag_times: usize,
    cmd_eth_to_cosmos: bool,
    cmd_cosmos_to_eth: bool,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage:
        {} cosmos-to-eth --cosmos-phrase=<key> --cosmos-rpc=<url> --erc20-address=<addr> --amount=<amount> --eth-destination=<dest> [--no-batch] [--times=<number>]
        {} eth-to-cosmos --ethereum-key=<key> --ethereum-rpc=<url> --contract-address=<addr> --erc20-address=<addr> --amount=<amount> --cosmos-destination=<dest> [--times=<number>]
        Options:
            -h --help                   Show this screen.
            --cosmos-key=<ckey>         The Cosmos private key of the sender
            --ethereum-key=<ekey>       The Ethereum private key of the sender
            --cosmos-rpc=<curl>         The Cosmos Legacy RPC url, this will need to be manually enabled
            --ethereum-rpc=<eurl>       The Ethereum RPC url, should be a self hosted node
            --contract-address=<addr>   The Ethereum contract address for Peggy, this is temporary
            --erc20-address=<addr>      An erc20 address to send funds
            --amount=<amount>           The amount of tokens to send, for example 1.5
            --cosmos-destination=<dest> A cosmos address to send tokens to
            --eth-destination=<dest>    A cosmos address to send tokens to
            --no-batch                  Don't request a batch when sending to Ethereum
            --times=<number>            The number of times this send should be preformed, useful for stress testing
        About:
            Althea Gravity client software, moves tokens from Ethereum to Cosmos and back
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

    let amount = fraction_eth_to_wei(args.flag_amount);
    let times = if args.flag_times == 0 {
        1usize
    } else {
        args.flag_times
    };
    let erc20_address: EthAddress = args
        .flag_erc20_address
        .parse()
        .expect("Invalid contract address!");

    if args.cmd_cosmos_to_eth {
        let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
            .expect("Failed to parse cosmos key phrase, does it have a password?");
        let cosmos_address = cosmos_key.to_public_key().unwrap().to_address();

        println!("Sending from Cosmos address {}", cosmos_address);

        let cosmos_url = Url::parse(&args.flag_cosmos_rpc).expect("Invalid Cosmos RPC url");
        let cosmos_url = cosmos_url.to_string();
        let cosmos_url = cosmos_url.trim_end_matches('/');
        let peggy_denom = format!("peggy{}", erc20_address);
        let contact = Contact::new(&cosmos_url, TIMEOUT);
        let amount = Coin {
            amount,
            denom: peggy_denom.clone(),
        };
        let bridge_fee = Coin {
            denom: peggy_denom.clone(),
            amount: 1u64.into(),
        };
        let eth_dest: EthAddress = args.flag_eth_destination.parse().unwrap();

        let balances = contact
            .get_balances(cosmos_address)
            .await
            .expect("Failed to get balances!")
            .result;
        let mut found = None;
        for coin in balances.iter() {
            if coin.denom == peggy_denom {
                found = Some(coin);
            }
        }

        println!("Cosmos balances {:?}", balances);

        if found.is_none() {
            panic!("You don't have any {} tokens!", peggy_denom);
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount
            && times == 1
        {
            panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_eth(amount.amount), peggy_denom, print_eth(found.unwrap().amount.clone()));
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount {
            panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_eth(amount.amount), times, peggy_denom, print_eth(found.unwrap().amount.clone()));
        }

        for _ in 0..times {
            println!(
                "Locking {} / {} into the batch pool",
                args.flag_amount, erc20_address
            );
            let res = send_to_eth(
                cosmos_key,
                eth_dest,
                amount.clone(),
                bridge_fee.clone(),
                &contact,
            )
            .await;
            match res {
                Ok(tx_id) => println!("Send to Eth txid {}", tx_id.txhash),
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        }

        if !args.flag_no_batch {
            println!("Requesting a batch to push transaction along immediately");
            send_request_batch(cosmos_key, peggy_denom, bridge_fee, &contact)
                .await
                .expect("Failed to request batch");
        } else {
            println!("--no-batch specified, your transfer will wait until someone requests a batch for this token type")
        }
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

        let erc20_balance = web3
            .get_erc20_balance(erc20_address, ethereum_public_key)
            .await
            .expect("Failed to get balance, check ERC20 contract address");

        if erc20_balance == 0u8.into() {
            panic!(
                "You have zero {} tokens, please double check your sender and erc20 addresses!",
                contract_address
            );
        } else if amount.clone() * times.into() > erc20_balance {
            panic!(
                "Insufficient balance {} > {}",
                amount * times.into(),
                erc20_balance
            );
        }

        for _ in 0..times {
            println!(
                "Sending {} / {} to Cosmos from {} to {}",
                args.flag_amount, erc20_address, ethereum_public_key, cosmos_dest
            );
            // we send some erc20 tokens to the peggy contract to register a deposit
            let res = send_to_cosmos(
                erc20_address,
                contract_address,
                amount.clone(),
                cosmos_dest,
                ethereum_key,
                Some(TIMEOUT),
                &web3,
                vec![],
            )
            .await;
            match res {
                Ok(tx_id) => println!("Send to Cosmos txid: {:#066x}", tx_id),
                Err(e) => println!("Failed to send tokens! {:?}", e),
            }
        }
    }
}

#[test]
fn even_f32_rounding() {
    let one_eth: Uint256 = 1000000000000000000u128.into();
    let one_point_five_eth: Uint256 = 1500000000000000000u128.into();
    let one_point_one_five_eth: Uint256 = 1150000000000000000u128.into();
    let a_high_precision_number: Uint256 = 1150100000000000000u128.into();
    let res = fraction_eth_to_wei(1f64);
    assert_eq!(one_eth, res);
    let res = fraction_eth_to_wei(1.5f64);
    assert_eq!(one_point_five_eth, res);
    let res = fraction_eth_to_wei(1.15f64);
    assert_eq!(one_point_one_five_eth, res);
    let res = fraction_eth_to_wei(1.1501f64);
    assert_eq!(a_high_precision_number, res);
}
