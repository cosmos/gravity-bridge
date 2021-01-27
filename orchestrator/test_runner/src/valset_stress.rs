use crate::{get_test_token_name, happy_path::test_valset_update, COSMOS_NODE_GRPC};
use actix::Arbiter;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use orchestrator::main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use web30::client::Web3;

#[allow(clippy::too_many_arguments)]
pub async fn validator_set_stress_test(
    web30: &Web3,
    contact: &Contact,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    peggy_address: EthAddress,
) {
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
            get_test_token_name(),
        ));
    }

    // TODO have some external system send hundreds of valset updates in parallel
    // to do this you need to generate a non-orchestrator address, send it funds
    // then use that to send the requests or your sequence gets all messed up
    for _ in 0u32..10 {
        test_valset_update(&web30, &keys, peggy_address).await;
    }
}
