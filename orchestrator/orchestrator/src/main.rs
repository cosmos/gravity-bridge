//! Orchestrator is a sort of specialized relayer for Althea-Gravity that runs on every validator.
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
mod get_with_retry;
mod main_loop;
mod metrics;
mod oracle_resync;

use crate::main_loop::orchestrator_main_loop;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use docopt::Docopt;
use env_logger::Env;
use gravity_utils::connection_prep::{
    check_delegate_addresses, check_for_eth, wait_for_cosmos_node_ready,
};
use gravity_utils::connection_prep::{check_for_fee_denom, create_rpc_connections};
use main_loop::{ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED};
use relayer::main_loop::LOOP_SPEED as RELAYER_LOOP_SPEED;
use std::cmp::min;

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_grpc: String,
    flag_address_prefix: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_fees: String,
    flag_metrics_listen: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --cosmos-phrase=<key> --ethereum-key=<key> --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --fees=<denom> --contract-address=<addr> --metrics-listen=<addr>
        Options:
            -h --help                    Show this screen.
            --cosmos-phrase=<ckey>       The mnenmonic of the Cosmos account key of the validator
            --ethereum-key=<ekey>        The Ethereum private key of the validator
            --cosmos-grpc=<gurl>         The Cosmos gRPC url, usually the validator
            --address-prefix=<prefix>    The prefix for addresses on this Cosmos chain
            --ethereum-rpc=<eurl>        The Ethereum RPC url, should be a self hosted node
            --fees=<denom>               The Cosmos Denom in which to pay Cosmos chain fees
            --contract-address=<addr>    The Ethereum contract address for Gravity, this is temporary
            --metrics-listen=<addr>      The address metrics server listens on [default: 127.0.0.1:3000].
        About:
            The Validator companion binary for Gravity. This must be run by all Gravity chain validators
            and is a mix of a relayer + oracle + ethereum signing infrastructure
            Written By: {}
            Version {}",
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
    let cosmos_key = CosmosPrivateKey::from_phrase(&args.flag_cosmos_phrase, "")
        .expect("Invalid Private Cosmos Key!");
    let ethereum_key: EthPrivateKey = args
        .flag_ethereum_key
        .parse()
        .expect("Invalid Ethereum private key!");
    let contract_address: EthAddress = args
        .flag_contract_address
        .parse()
        .expect("Invalid contract address!");
    let metrics_listen = args
        .flag_metrics_listen
        .parse()
        .expect("Invalid metrics listen address!");

    let fee_denom = args.flag_fees;

    let timeout = min(
        min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
        RELAYER_LOOP_SPEED,
    );

    trace!("Probing RPC connections");
    // probe all rpc connections and see if they are valid
    let connections = create_rpc_connections(
        args.flag_address_prefix,
        Some(args.flag_cosmos_grpc),
        Some(args.flag_ethereum_rpc),
        timeout,
    )
    .await;

    let mut grpc = connections.grpc.clone().unwrap();
    let contact = connections.contact.clone().unwrap();
    let web3 = connections.web3.clone().unwrap();

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    let public_cosmos_key = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    info!("Starting Gravity Validator companion binary Relayer + Oracle + Eth Signer");
    info!(
        "Ethereum Address: {} Cosmos Address {}",
        public_eth_key, public_cosmos_key
    );

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;

    // check if the delegate addresses are correctly configured
    check_delegate_addresses(
        &mut grpc,
        public_eth_key,
        public_cosmos_key,
        &contact.get_prefix(),
    )
    .await;

    // check if we actually have the promised balance of tokens to pay fees
    check_for_fee_denom(&fee_denom, public_cosmos_key, &contact).await;
    check_for_eth(public_eth_key, &web3).await;

    orchestrator_main_loop(
        cosmos_key,
        ethereum_key,
        connections.web3.unwrap(),
        connections.contact.unwrap(),
        connections.grpc.unwrap(),
        contract_address,
        (1f64, fee_denom.to_owned()),
        &metrics_listen,
    )
    .await;
}
