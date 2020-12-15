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

use crate::main_loop::orchestrator_main_loop;
use crate::main_loop::LOOP_SPEED;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use docopt::Docopt;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use url::Url;
use web30::client::Web3;

#[derive(Debug, Deserialize)]
struct Args {
    flag_cosmos_phrase: String,
    flag_ethereum_key: String,
    flag_cosmos_legacy_rpc: String,
    flag_cosmos_grpc: String,
    flag_ethereum_rpc: String,
    flag_contract_address: String,
    flag_fees: String,
}

lazy_static! {
    pub static ref USAGE: String = format!(
    "Usage: {} --cosmos-phrase=<key> --ethereum-key=<key> --cosmos-legacy-rpc=<url> --cosmos-grpc=<url> --ethereum-rpc=<url> --fees=<denom> --contract-address=<addr>
        Options:
            -h --help                    Show this screen.
            --cosmos-key=<ckey>          The Cosmos private key of the validator
            --ethereum-key=<ekey>        The Ethereum private key of the validator
            --cosmos-legacy-rpc=<curl>   The Cosmos RPC url, usually the validator
            --cosmos-grpc=<gurl>         The Cosmos gRPC url, usually the validator
            --ethereum-rpc=<eurl>        The Ethereum RPC url, should be a self hosted node
            --fees=<denom>               The Cosmos Denom in which to pay Cosmos chain fees
            --contract-address=<addr>    The Ethereum contract address for Peggy, this is temporary
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

#[actix_rt::main]
async fn main() {
    env_logger::init();

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

    let _ = Url::parse(&args.flag_cosmos_legacy_rpc).expect("Invalid Cosmos legacy RPC url");
    let cosmos_legacy_url = args.flag_cosmos_legacy_rpc.trim_end_matches('/');

    let _ = Url::parse(&args.flag_cosmos_grpc).expect("Invalid Cosmos gRPC url");
    let cosmos_grpc_url = args.flag_cosmos_grpc.trim_end_matches('/').to_string();

    let _ = Url::parse(&args.flag_ethereum_rpc).expect("Invalid Ethereum RPC url");
    let eth_url = args.flag_ethereum_rpc.trim_end_matches('/');

    let fee_denom = args.flag_fees;

    let grpc_client = PeggyQueryClient::connect(cosmos_grpc_url).await.unwrap();
    let web3 = Web3::new(&eth_url, LOOP_SPEED);
    let contact = Contact::new(&cosmos_legacy_url, LOOP_SPEED);

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    let public_cosmos_key = cosmos_key
        .to_public_key()
        .expect("Invalid Cosmos Phrase!")
        .to_address();
    info!("Starting Peggy Relayer + Eth Signer");
    info!(
        "Ethereum Address: {} Cosmos Address {}",
        public_eth_key, public_cosmos_key
    );

    orchestrator_main_loop(
        cosmos_key,
        ethereum_key,
        web3,
        contact,
        grpc_client,
        contract_address,
        fee_denom,
    )
    .await;
}
