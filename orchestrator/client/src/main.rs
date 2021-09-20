//! This file is the binary entry point for the Gravity client software, an easy to use cli utility that
//! allows anyone to send funds across the Gravity bridge. Currently this application only does anything
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
use cosmos_gravity::send::{send_request_batch_tx, send_to_eth};
use deep_space::address::Address as CosmosAddress;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use docopt::Docopt;
use env_logger::Env;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use gravity_proto::gravity::DenomToErc20Request;
use gravity_utils::connection_prep::{check_for_eth, check_for_fee_denom, create_rpc_connections};
use std::{process::exit, time::Duration};

const TIMEOUT: Duration = Duration::from_secs(60);

pub fn one_eth() -> f64 {
    1000000000000000000f64
}

pub fn one_atom() -> f64 {
    1000000f64
}

pub fn print_eth(input: Uint256) -> String {
    let float: f64 = input.to_string().parse().unwrap();
    let res = float / one_eth();
    format!("{}", res)
}

pub fn print_atom(input: Uint256) -> String {
    let float: f64 = input.to_string().parse().unwrap();
    let res = float / one_atom();
    format!("{}", res)
}

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_grpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_cosmos_denom: String,
    flag_amount: Option<String>,
    flag_cosmos_destination: String,
    flag_erc20_address: String,
    flag_eth_destination: String,
    flag_no_batch: bool,
    flag_times: usize,
    flag_cosmos_prefix: String,
    cmd_eth_to_cosmos: bool,
    cmd_cosmos_to_eth: bool,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage:
        {} cosmos-to-eth --cosmos-phrase=<key> --cosmos-grpc=<url> --cosmos-prefix=<prefix> --cosmos-denom=<denom> --amount=<amount> --eth-destination=<dest> [--no-batch] [--times=<number>]
        {} eth-to-cosmos --ethereum-key=<key> --ethereum-rpc=<url> --cosmos-prefix=<prefix> --contract-address=<addr> --erc20-address=<addr> --amount=<amount> --cosmos-destination=<dest> [--times=<number>]
        Options:
            -h --help                   Show this screen.
            --cosmos-phrase=<ckey>      The mnenmonic of the Cosmos account key of the validator
            --ethereum-key=<ekey>       The Ethereum private key of the sender
            --cosmos-legacy-rpc=<curl>  The Cosmos Legacy RPC url, this will need to be manually enabled
            --cosmos-grpc=<curl>        The Cosmos gRPC url
            --cosmos-prefix=<prefix>    The Bech32 Prefix used for the Cosmos chain's addresses
            --ethereum-rpc=<eurl>       The Ethereum RPC url, should be a self hosted node
            --contract-address=<addr>   The Ethereum contract address for Gravity, this is temporary
            --erc20-address=<addr>      An erc20 address on Ethereum to send funds from
            --cosmos-denom=<amount>     The Cosmos denom that you intend to send to Ethereum
            --amount=<amount>           The amount of tokens to send, for example 1.5
            --cosmos-destination=<dest> A cosmos address to send tokens to
            --eth-destination=<dest>    A cosmos address to send tokens to
            --no-batch                  Don't request a batch when sending to Ethereum
            --times=<number>            The number of times this send should be preformed, useful for stress testing
        Description:
            cosmos-to-eth               Locks up a Cosmos asset in the batch pool. Optionally this command will also request a batch.
            eth-to-cosmos               Sends an Ethereum ERC20 asset to a Cosmos destination address
            deploy-erc20-representation Deploys an ERC20 representation for a Cosmos asset, required to bridge a Cosmos native asset with 'cosmos-to-eth'
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
    env_logger::Builder::from_env(Env::default().default_filter_or("info")).init();
    // On Linux static builds we need to probe ssl certs path to be able to
    // do TLS stuff.
    openssl_probe::init_ssl_cert_env_vars();
    let args: Args = Docopt::new(USAGE.as_str())
        .and_then(|d| d.deserialize())
        .unwrap_or_else(|e| e.exit());

    let times = if args.flag_times == 0 {
        1usize
    } else {
        args.flag_times
    };

    if args.cmd_cosmos_to_eth {
        let gravity_denom = args.flag_cosmos_denom;
        // todo actually query metadata for this
        let is_cosmos_originated = !gravity_denom.starts_with("gravity");
        let amount = args.flag_amount.unwrap().parse().unwrap();
        let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
            .expect("Failed to parse cosmos key phrase, does it have a password?");
        let cosmos_address = cosmos_key.to_address(&args.flag_cosmos_prefix).unwrap();

        println!("Sending from Cosmos address {}", cosmos_address);
        let connections = create_rpc_connections(
            args.flag_cosmos_prefix,
            Some(args.flag_cosmos_grpc),
            None,
            TIMEOUT,
        )
        .await;
        let contact = connections.contact.unwrap();
        let mut grpc = connections.grpc.unwrap();

        let res = grpc
            .denom_to_erc20(DenomToErc20Request {
                denom: gravity_denom.clone(),
            })
            .await;
        match res {
            Ok(val) => println!(
                "Asset {} has ERC20 representation {}",
                gravity_denom,
                val.into_inner().erc20
            ),
            Err(_e) => {
                println!(
                    "Asset {} has no ERC20 representation, you may need to deploy an ERC20 for it!",
                    gravity_denom
                );
                exit(1);
            }
        }

        let amount = Coin {
            amount,
            denom: gravity_denom.clone(),
        };
        let bridge_fee = Coin {
            denom: gravity_denom.clone(),
            amount: 1u64.into(),
        };
        let eth_dest: EthAddress = args.flag_eth_destination.parse().unwrap();
        check_for_fee_denom(&gravity_denom, cosmos_address, &contact).await;

        let balances = contact
            .get_balances(cosmos_address)
            .await
            .expect("Failed to get balances!");
        let mut found = None;
        for coin in balances.iter() {
            if coin.denom == gravity_denom {
                found = Some(coin);
            }
        }

        println!("Cosmos balances {:?}", balances);

        if found.is_none() {
            panic!("You don't have any {} tokens!", gravity_denom);
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount
            && times == 1
        {
            if is_cosmos_originated {
                panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_atom(amount.amount), gravity_denom, print_atom(found.unwrap().amount.clone()));
            } else {
                panic!("Your transfer of {} {} tokens is greater than your balance of {} tokens. Remember you need some to pay for fees!", print_eth(amount.amount), gravity_denom, print_eth(found.unwrap().amount.clone()));
            }
        } else if amount.amount.clone() * times.into() >= found.clone().unwrap().amount {
            if is_cosmos_originated {
                panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_atom(amount.amount), times, gravity_denom, print_atom(found.unwrap().amount.clone()));
            } else {
                panic!("Your transfer of {} * {} {} tokens is greater than your balance of {} tokens. Try to reduce the amount or the --times parameter", print_eth(amount.amount), times, gravity_denom, print_eth(found.unwrap().amount.clone()));
            }
        }

        for _ in 0..times {
            println!(
                "Locking {} / {} into the batch pool",
                amount.clone(),
                gravity_denom
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
            send_request_batch_tx(cosmos_key, gravity_denom, bridge_fee, &contact)
                .await
                .expect("Failed to request batch");
        } else {
            println!("--no-batch specified, your transfer will wait until someone requests a batch for this token type")
        }
    } else if args.cmd_eth_to_cosmos {
        let erc20_address: EthAddress = args
            .flag_erc20_address
            .parse()
            .expect("Invalid ERC20 contract address!");
        let ethereum_key: EthPrivateKey = args
            .flag_ethereum_key
            .parse()
            .expect("Invalid Ethereum private key!");
        let contract_address: EthAddress = args
            .flag_contract_address
            .parse()
            .expect("Invalid contract address!");
        let connections = create_rpc_connections(
            args.flag_cosmos_prefix,
            None,
            Some(args.flag_ethereum_rpc),
            TIMEOUT,
        )
        .await;
        let web3 = connections.web3.unwrap();
        let cosmos_dest: CosmosAddress = args.flag_cosmos_destination.parse().unwrap();
        let ethereum_public_key = ethereum_key.to_public_key().unwrap();
        check_for_eth(ethereum_public_key, &web3).await;

        let amount: Uint256 = args.flag_amount.unwrap().parse().unwrap();

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
                amount.clone(),
                erc20_address,
                ethereum_public_key,
                cosmos_dest
            );
            // we send some erc20 tokens to the gravity contract to register a deposit
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
