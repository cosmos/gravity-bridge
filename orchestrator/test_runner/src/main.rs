//! this crate, namely runs all up integration tests of the Peggy code against
//! several scenarios, happy path and non happy path. This is essentially meant
//! to be executed in our specific CI docker container and nowhere else. If you
//! find some function useful pull it up into the more general peggy_utils or the like

#[macro_use]
extern crate log;
#[macro_use]
extern crate lazy_static;

use crate::bootstrapping::*;
use crate::utils::*;
use clarity::PrivateKey as EthPrivateKey;
use clarity::{Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::utils::wait_for_cosmos_online;
use deep_space::coin::Coin;
use happy_path::happy_path_test;
use happy_path_v2::happy_path_test_v2;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use std::{env, time::Duration};
use transaction_stress_test::transaction_stress_test;
use valset_stress::validator_set_stress_test;

mod bootstrapping;
mod happy_path;
mod happy_path_v2;
mod transaction_stress_test;
mod utils;
mod valset_stress;

/// the timeout for individual requests
const OPERATION_TIMEOUT: Duration = Duration::from_secs(30);
/// the timeout for the total system
const TOTAL_TIMEOUT: Duration = Duration::from_secs(300);

pub const COSMOS_NODE: &str = "http://localhost:1317";
pub const COSMOS_NODE_GRPC: &str = "http://localhost:9090";
pub const COSMOS_NODE_ABCI: &str = "http://localhost:26657";
pub const ETH_NODE: &str = "http://localhost:8545";

/// this value reflects the contents of /tests/container-scripts/setup-validator.sh
/// and is used to compute if a stake change is big enough to trigger a validator set
/// update since we want to make several such changes intentionally
pub const STAKE_SUPPLY_PER_VALIDATOR: u128 = 1000000000;
/// this is the amount each validator bonds at startup
pub const STARTING_STAKE_PER_VALIDATOR: u128 = STAKE_SUPPLY_PER_VALIDATOR / 2;

lazy_static! {
    // this key is the private key for the public key defined in tests/assets/ETHGenesis.json
    // where the full node / miner sends its rewards. Therefore it's always going
    // to have a lot of ETH to pay for things like contract deployments
    static ref MINER_PRIVATE_KEY: EthPrivateKey =
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
            .parse()
            .unwrap();
    static ref MINER_ADDRESS: EthAddress = MINER_PRIVATE_KEY.to_public_key().unwrap();
}

/// Gets the standard non-token fee for the testnet. We deploy the test chain with STAKE
/// and FOOTOKEN balances by default, one footoken is sufficient for any Cosmos tx fee except
/// fees for send_to_eth messages which have to be of the same bridged denom so that the relayers
/// on the Ethereum side can be paid in that token.
pub fn get_fee() -> Coin {
    Coin {
        denom: get_test_token_name(),
        amount: 1u32.into(),
    }
}

pub fn get_test_token_name() -> String {
    "footoken".to_string()
}

pub fn one_eth() -> Uint256 {
    1000000000000000000u128.into()
}

pub fn should_deploy_contracts() -> bool {
    match env::var("DEPLOY_CONTRACTS") {
        Ok(s) => s == "1" || s.to_lowercase() == "yes" || s.to_lowercase() == "true",
        _ => false,
    }
}

#[actix_rt::main]
pub async fn main() {
    env_logger::init();
    info!("Staring Peggy test-runner");
    let contact = Contact::new(COSMOS_NODE, OPERATION_TIMEOUT);

    info!("Waiting for Cosmos chain to come online");
    wait_for_cosmos_online(&contact, TOTAL_TIMEOUT).await;

    let grpc_client = PeggyQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
    let web30 = web30::client::Web3::new(ETH_NODE, OPERATION_TIMEOUT);
    let keys = get_keys();

    // if we detect this env var we are only deploying contracts, do that then exit.
    if should_deploy_contracts() {
        info!("test-runner in contract deploying mode, deploying contracts, then exiting");
        deploy_contracts(&contact, &keys, get_fee()).await;
        return;
    }

    let contracts = parse_contract_addresses();
    // the address of the deployed Peggy contract
    let peggy_address = contracts.peggy_contract;
    // addresses of deployed ERC20 token contracts to be used for testing
    let erc20_addresses = contracts.erc20_addresses;

    // before we start the orchestrators send them some funds so they can pay
    // for things
    send_eth_to_orchestrators(&keys, &web30).await;

    assert!(check_cosmos_balance(
        &get_test_token_name(),
        keys[0].0.to_public_key().unwrap().to_address(),
        &contact
    )
    .await
    .is_some());

    // This segment contains optional tests, by default we run a happy path test
    // this tests all major functionality of Peggy once or twice.
    // VALSET_STRESS sends in 1k valsets to sign and update
    // BATCH_STRESS fills several batches and executes an out of order batch
    // VALIDATOR_OUT simulates a validator not participating in the happy path test
    let test_type = env::var("TEST_TYPE");
    info!("Starting tests with {:?}", test_type);
    if let Ok(test_type) = test_type {
        if test_type == "VALIDATOR_OUT" {
            info!("Starting Validator out test");
            happy_path_test(
                &web30,
                grpc_client,
                &contact,
                keys,
                peggy_address,
                erc20_addresses[0],
                true,
            )
            .await;
            return;
        } else if test_type == "BATCH_STRESS" {
            let contact = Contact::new(COSMOS_NODE, TOTAL_TIMEOUT);
            transaction_stress_test(&web30, &contact, keys, peggy_address, erc20_addresses).await;
            return;
        } else if test_type == "VALSET_STRESS" {
            info!("Starting Valset update stress test");
            validator_set_stress_test(&web30, &contact, keys, peggy_address).await;
            return;
        } else if test_type == "V2_HAPPY_PATH" {
            info!("Starting happy path for Gravity v2");
            happy_path_test_v2(&web30, grpc_client, &contact, keys, peggy_address, false).await;
            return;
        }
    }
    info!("Starting Happy path test");
    happy_path_test(
        &web30,
        grpc_client,
        &contact,
        keys,
        peggy_address,
        erc20_addresses[0],
        false,
    )
    .await;
}
