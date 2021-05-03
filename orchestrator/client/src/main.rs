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
use cosmos_gravity::send::{send_request_batch, send_to_eth};
use deep_space::address::Address as CosmosAddress;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use docopt::Docopt;
use env_logger::Env;
use ethereum_gravity::deploy_erc20::deploy_erc20;
use ethereum_gravity::send_to_cosmos::send_to_cosmos;
use gravity_proto::gravity::QueryDenomToErc20Request;
use gravity_utils::connection_prep::{check_for_eth, check_for_fee_denom, create_rpc_connections};
use std::time::Instant;
use std::{process::exit, time::Duration, u128};
use tokio::time::sleep as delay_for;
use web30::{client::Web3, jsonrpc::error::Web3Error};

const TIMEOUT: Duration = Duration::from_secs(60);

pub fn one_eth() -> f64 {
    1000000000000000000f64
}

pub fn one_atom() -> f64 {
    1000000f64
}

/// TODO revisit this for higher precision while
/// still representing the number to the user as a float
/// this takes a number like 0.37 eth and turns it into wei
/// or any erc20 with arbitrary decimals
pub fn fraction_to_exponent(num: f64, exponent: u8) -> Uint256 {
    let mut res = num;
    // in order to avoid floating point rounding issues we
    // multiply only by 10 each time. this reduces the rounding
    // errors enough to be ignored
    for _ in 0..exponent {
        res *= 10f64
    }
    (res as u128).into()
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
    flag_amount: Option<f64>,
    flag_cosmos_destination: String,
    flag_erc20_address: String,
    flag_eth_destination: String,
    flag_no_batch: bool,
    flag_times: usize,
    flag_erc20_name: String,
    flag_erc20_symbol: String,
    flag_erc20_decimals: u8,
    flag_cosmos_prefix: String,
    cmd_eth_to_cosmos: bool,
    cmd_cosmos_to_eth: bool,
    cmd_deploy_erc20_representation: bool,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage:
        {} cosmos-to-eth --cosmos-phrase=<key> --cosmos-grpc=<url> --cosmos-prefix=<prefix> --cosmos-denom=<denom> --amount=<amount> --eth-destination=<dest> [--no-batch] [--times=<number>]
        {} eth-to-cosmos --ethereum-key=<key> --ethereum-rpc=<url> --cosmos-prefix=<prefix> --contract-address=<addr> --erc20-address=<addr> --amount=<amount> --cosmos-destination=<dest> [--times=<number>]
        {} deploy-erc20-representation --cosmos-grpc=<url> --cosmos-prefix=<prefix> --cosmos-denom=<denom> --ethereum-key=<key> --ethereum-rpc=<url> --contract-address=<addr> --erc20-name=<name> --erc20-symbol=<symbol> --erc20-decimals=<decimals>
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
            --erc20-name=<name>         The 'name' value for the deployed ERC20 contract, must match Cosmos denom metadata
            --erc20-symbol=<symbol>     The 'symbol 'value for the deployed ERC20 contract, must match the Cosmos denom metadata
            --erc20-decimals=<decimals> The number of decimals the deployed ERC20 token will have, must match the resolution of the Cosmos asset to be adopted by the chain  
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
        let amount = if is_cosmos_originated {
            fraction_to_exponent(args.flag_amount.unwrap(), 6)
        } else {
            fraction_to_exponent(args.flag_amount.unwrap(), 18)
        };
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
            .denom_to_erc20(QueryDenomToErc20Request {
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
                args.flag_amount.unwrap(),
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
            send_request_batch(cosmos_key, gravity_denom, bridge_fee, &contact)
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

        let res = get_erc20_decimals(&web3, erc20_address, ethereum_public_key)
            .await
            .expect("Failed to query ERC20 contract");
        let decimals: u8 = res.to_string().parse().unwrap();
        let amount = fraction_to_exponent(args.flag_amount.unwrap(), decimals);

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
                args.flag_amount.unwrap(),
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
    } else if args.cmd_deploy_erc20_representation {
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
            Some(args.flag_cosmos_grpc),
            Some(args.flag_ethereum_rpc),
            TIMEOUT,
        )
        .await;
        let web3 = connections.web3.unwrap();
        let mut grpc = connections.grpc.unwrap();
        let ethereum_public_key = ethereum_key.to_public_key().unwrap();
        check_for_eth(ethereum_public_key, &web3).await;

        let denom = args.flag_cosmos_denom;

        let res = grpc
            .denom_to_erc20(QueryDenomToErc20Request {
                denom: denom.clone(),
            })
            .await;
        if let Ok(val) = res {
            println!(
                "Asset {} already has ERC20 representation {}",
                denom,
                val.into_inner().erc20
            );
            exit(1);
        }

        println!("Starting deploy of ERC20");
        let res = deploy_erc20(
            denom.clone(),
            args.flag_erc20_name,
            args.flag_erc20_symbol,
            args.flag_erc20_decimals,
            contract_address,
            &web3,
            Some(TIMEOUT),
            ethereum_key,
            vec![],
        )
        .await
        .unwrap();

        println!("We have deployed ERC20 contract {:#066x}, waiting to see if the Cosmos chain choses to adopt it", res);

        let start = Instant::now();
        loop {
            let res = grpc
                .denom_to_erc20(QueryDenomToErc20Request {
                    denom: denom.clone(),
                })
                .await;

            if let Ok(val) = res {
                println!(
                    "Asset {} has accepted new ERC20 representation {}",
                    denom,
                    val.into_inner().erc20
                );
                exit(0);
            }

            if Instant::now() - start > Duration::from_secs(100) {
                println!(
                    "Your ERC20 contract was not adopted, double check the metadata and try again"
                );
                exit(1);
            }
            delay_for(Duration::from_secs(1)).await;
        }
    }
}

#[test]
fn even_f32_rounding() {
    let one_eth: Uint256 = 1000000000000000000u128.into();
    let one_point_five_eth: Uint256 = 1500000000000000000u128.into();
    let one_point_one_five_eth: Uint256 = 1150000000000000000u128.into();
    let a_high_precision_number: Uint256 = 1150100000000000000u128.into();
    let res = fraction_to_exponent(1f64, 18);
    assert_eq!(one_eth, res);
    let res = fraction_to_exponent(1.5f64, 18);
    assert_eq!(one_point_five_eth, res);
    let res = fraction_to_exponent(1.15f64, 18);
    assert_eq!(one_point_one_five_eth, res);
    let res = fraction_to_exponent(1.1501f64, 18);
    assert_eq!(a_high_precision_number, res);
}

pub async fn get_erc20_decimals(
    web3: &Web3,
    erc20: EthAddress,
    caller_address: EthAddress,
) -> Result<Uint256, Web3Error> {
    let decimals = web3
        .contract_call(erc20, "decimals()", &[], caller_address)
        .await?;

    Ok(Uint256::from_bytes_be(match decimals.get(0..32) {
        Some(val) => val,
        None => {
            return Err(Web3Error::ContractCallError(
                "Bad response from ERC20 decimals".to_string(),
            ))
        }
    }))
}
