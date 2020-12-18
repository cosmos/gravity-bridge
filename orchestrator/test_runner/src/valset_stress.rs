use crate::{happy_path::test_valset_update, COSMOS_NODE_GRPC};
use actix::Arbiter;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use orchestrator::main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use tonic::transport::Channel;
use web30::client::Web3;

#[allow(clippy::too_many_arguments)]
pub async fn validator_set_stress_test(
    web30: &Web3,
    grpc_client: PeggyQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    peggy_address: EthAddress,
    test_token_name: String,
    erc20_address: EthAddress,
    fee: Coin,
) {
    let mut grpc_client = grpc_client;
    // start orchestrators
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
            test_token_name.clone(),
        ));
    }

    for _ in 0u32..25 {
        test_valset_update(&contact, &web30, &keys, peggy_address, fee.clone()).await;
    }
}
