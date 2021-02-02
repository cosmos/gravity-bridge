//! This is the happy path test for Cosmos to Ethereum asset transfers, meaning assets originated on Cosmos

use crate::get_test_token_name;
use crate::{COSMOS_NODE_GRPC, TOTAL_TIMEOUT};
use actix::Arbiter;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use ethereum_peggy::{deploy_erc20::deploy_erc20, utils::get_event_nonce};
use orchestrator::main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub async fn happy_path_test_v2(
    web30: &Web3,
    grpc_client: PeggyQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    peggy_address: EthAddress,
    validator_out: bool,
) {
    let starting_event_nonce =
        get_event_nonce(peggy_address, keys[0].1.to_public_key().unwrap(), web30)
            .await
            .unwrap();
    deploy_erc20(
        "footoken-a".to_string(),
        "footoken-b".to_string(),
        "foo".to_string(),
        18,
        peggy_address,
        web30,
        Some(TOTAL_TIMEOUT),
        keys[0].1,
        vec![],
    )
    .await
    .unwrap();
    let ending_event_nonce =
        get_event_nonce(peggy_address, keys[0].1.to_public_key().unwrap(), web30)
            .await
            .unwrap();

    assert!(starting_event_nonce != ending_event_nonce);
    info!(
        "Successfully deployed new ERC20 representing FooToken on Cosmos with event nonce {}",
        ending_event_nonce
    );

    // used to break out of the loop early to simulate one validator
    // not running an Orchestrator
    let num_validators = keys.len();
    let mut count = 0;

    // start orchestrators
    #[allow(clippy::explicit_counter_loop)]
    for (c_key, e_key) in keys.iter() {
        info!("Spawning Orchestrator");
        let grpc_client = PeggyQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        Arbiter::spawn(orchestrator_main_loop(
            *c_key,
            *e_key,
            web30.clone(),
            contact.clone(),
            grpc_client,
            peggy_address,
            get_test_token_name(),
        ));

        // used to break out of the loop early to simulate one validator
        // not running an orchestrator
        count += 1;
        if validator_out && count == num_validators - 1 {
            break;
        }
    }

    delay_for(TOTAL_TIMEOUT).await;

    // TODO make sure that Cosmos adopts the contract
    // TODO generate batch of footoken to send over to ethereum, check that it gets there
}
