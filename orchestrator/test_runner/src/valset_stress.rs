use crate::get_test_token_name;
use crate::happy_path::test_valset_update;
use crate::utils::ValidatorKeys;
use crate::COSMOS_NODE_GRPC;
use actix::Arbiter;
use clarity::Address as EthAddress;
use contact::client::Contact;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use orchestrator::main_loop::orchestrator_main_loop;
use web30::client::Web3;

pub async fn validator_set_stress_test(
    web30: &Web3,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
    gravity_address: EthAddress,
) {
    // start orchestrators
    for k in keys.iter() {
        info!("Spawning Orchestrator");
        let grpc_client = GravityQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        Arbiter::spawn(orchestrator_main_loop(
            k.orch_key,
            k.eth_key,
            web30.clone(),
            contact.clone(),
            grpc_client,
            gravity_address,
            get_test_token_name(),
        ));
    }

    for _ in 0u32..10 {
        test_valset_update(&web30, &keys, gravity_address).await;
    }
}
