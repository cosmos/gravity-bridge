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
mod oracle_resync;

use crate::main_loop::orchestrator_main_loop;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::address::Address as CosmosAddress;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use docopt::Docopt;
use main_loop::{ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED};
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_proto::peggy::QueryDelegateKeysByEthAddress;
use peggy_proto::peggy::QueryDelegateKeysByOrchestratorAddress;
use peggy_utils::connection_prep::create_rpc_connections;
use relayer::main_loop::LOOP_SPEED as RELAYER_LOOP_SPEED;
use std::{cmp::min, process::exit};
use tonic::transport::Channel;

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
            The Validator companion binary for Peggy. This must be run by all Peggy chain validators
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
    env_logger::init();
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

    let fee_denom = args.flag_fees;

    let timeout = min(
        min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
        RELAYER_LOOP_SPEED,
    );

    let connections = create_rpc_connections(
        Some(args.flag_cosmos_grpc),
        Some(args.flag_cosmos_legacy_rpc),
        Some(args.flag_ethereum_rpc),
        timeout,
    )
    .await;

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    let public_cosmos_key = cosmos_key
        .to_public_key()
        .expect("Invalid Cosmos Phrase!")
        .to_address();
    info!("Starting Peggy Validator companion binary Relayer + Oracle + Eth Signer");
    info!(
        "Ethereum Address: {} Cosmos Address {}",
        public_eth_key, public_cosmos_key
    );

    let mut grpc = connections.grpc.clone().unwrap();
    check_delegate_addresses(&mut grpc, public_eth_key, public_cosmos_key).await;

    // TODO this should wait here if the cosmos node is still syncing.

    orchestrator_main_loop(
        cosmos_key,
        ethereum_key,
        connections.web3.unwrap(),
        connections.contact.unwrap(),
        connections.grpc.unwrap(),
        contract_address,
        fee_denom,
    )
    .await;
}

/// This function checks the orchestrator delegate addresses
/// for consistency what this means is that it takes the Ethereum
/// address and Orchestrator address from the Orchestrator and checks
/// that both are registered and internally consistent.
async fn check_delegate_addresses(
    client: &mut PeggyQueryClient<Channel>,
    delegate_eth_address: EthAddress,
    delegate_orchestrator_address: CosmosAddress,
) {
    let eth_response = client
        .get_delegate_key_by_eth(QueryDelegateKeysByEthAddress {
            eth_address: delegate_eth_address.to_string(),
        })
        .await;
    let orchestrator_response = client
        .get_delegate_key_by_orchestrator(QueryDelegateKeysByOrchestratorAddress {
            orchestrator_address: delegate_orchestrator_address.to_string(),
        })
        .await;
    match (eth_response, orchestrator_response) {
        (Ok(e), Ok(o)) => {
            let e = e.into_inner();
            let o = o.into_inner();
            let req_delegate_orchestrator_address: CosmosAddress =
                e.orchestrator_address.parse().unwrap();
            let req_delegate_eth_address: EthAddress = o.eth_address.parse().unwrap();
            if req_delegate_eth_address != delegate_eth_address
                && req_delegate_orchestrator_address != delegate_orchestrator_address
            {
                error!("Your Delegate Ethereum and Orchestrator addresses are both incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_orchestrator_address, req_delegate_orchestrator_address
                );
                error!("In order to resolve this issue you should double check your input value or re-register your delegate keys");
                exit(1);
            } else if req_delegate_eth_address != delegate_eth_address {
                error!("Your Delegate Ethereum address is incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!("In order to resolve this issue you should double check how you input your eth private key");
                exit(1);
            } else if req_delegate_orchestrator_address != delegate_orchestrator_address {
                error!("Your Delegate Orchestrator address is incorrect!");
                error!(
                    "You provided {}  Correct Value {}",
                    delegate_eth_address, req_delegate_eth_address
                );
                error!("In order to resolve this issue you should double check how you input your Orchestrator address phrase, make sure you didn't use your Validator phrase!");
                exit(1);
            }

            if e.validator_address != o.validator_address {
                error!("You are using delegate keys from two different validator addresses!");
                error!("If you get this error message I would just blow everything away and start again");
                exit(1);
            }
        }
        (Err(_), Ok(_)) | (Ok(_), Err(_)) => {
            panic!("Failed to check delegate Eth address. Maybe try running the program again? If that doesn't work try registering delegate keys again")
        }
        (Err(_), Err(_)) => {
            panic!("Delegate addresses are not set! Please Register your delegate keys and make sure your Althea binary is updated!")
        }
    }
}
