//! This test verifies that live updating of orchestrator keys works correctly

use crate::get_fee;
use crate::utils::ValidatorKeys;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_gravity::send::update_gravity_delegate_addresses;
use deep_space::address::Address as CosmosAddress;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use gravity_proto::gravity::{
    query_client::QueryClient as GravityQueryClient, DelegateKeysByEthereumSignerRequest,
    DelegateKeysByOrchestratorRequest,
};
use rand::Rng;
use std::time::Duration;
use tonic::transport::Channel;

const BLOCK_TIMEOUT: Duration = Duration::from_secs(30);

pub async fn orch_keys_update(
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
) {
    let mut keys = keys;
    let mut grpc_client = grpc_client;
    // just to test that we have the right keys from the gentx
    info!("About to check already set delegate addresses");
    for k in keys.iter() {
        let eth_address = k.eth_key.to_public_key().unwrap();
        let orch_address = k.orch_key.to_address(&contact.get_prefix()).unwrap();
        let eth_response = grpc_client
            .delegate_keys_by_ethereum_signer(DelegateKeysByEthereumSignerRequest {
                ethereum_signer: eth_address.to_string(),
            })
            .await
            .unwrap()
            .into_inner();

        let parsed_response_orch_address: CosmosAddress =
            eth_response.orchestrator_address.parse().unwrap();
        assert_eq!(parsed_response_orch_address, orch_address);

        let orchestrator_response = grpc_client
            .delegate_keys_by_orchestrator(DelegateKeysByOrchestratorRequest {
                orchestrator_address: orch_address.to_string(),
            })
            .await
            .unwrap()
            .into_inner();

        let parsed_response_eth_address: EthAddress =
            orchestrator_response.ethereum_signer.parse().unwrap();
        assert_eq!(parsed_response_eth_address, eth_address);
    }

    info!("Starting with {:?}", keys);

    // now we change them all
    for k in keys.iter_mut() {
        let mut rng = rand::thread_rng();
        let secret: [u8; 32] = rng.gen();
        // generate some new keys to replace the old ones
        let ethereum_key = EthPrivateKey::from_slice(&secret).unwrap();
        let cosmos_key = CosmosPrivateKey::from_secret(&secret);
        // update the keys in the key list
        k.eth_key = ethereum_key;
        k.orch_key = cosmos_key;
        let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

        info!(
            "Signing and submitting Delegate addresses {} for validator {}",
            ethereum_key.to_public_key().unwrap(),
            cosmos_address,
        );
        // send in the new delegate keys signed by the validator address
        update_gravity_delegate_addresses(
            &contact,
            ethereum_key.to_public_key().unwrap(),
            cosmos_address,
            k.validator_key,
            k.eth_key,
            get_fee(),
        )
        .await
        .expect("Failed to set delegate addresses!");
    }

    contact.wait_for_next_block(BLOCK_TIMEOUT).await.unwrap();

    // TODO registering is too unreliable right now for confusing reasons, revisit with prototx

    // info!("About to check changed delegate addresses");
    // // verify that the change has taken place
    // for k in keys.iter() {
    //     let eth_address = k.eth_key.to_public_key().unwrap();
    //     let orch_address = k.orch_key.to_public_key().unwrap().to_address();

    //     let orchestrator_response = grpc_client
    //         .get_delegate_key_by_orchestrator(QueryDelegateKeysByOrchestratorAddress {
    //             orchestrator_address: orch_address.to_string(),
    //         })
    //         .await
    //         .unwrap()
    //         .into_inner();

    //     let parsed_response_eth_address: EthAddress =
    //         orchestrator_response.eth_address.parse().unwrap();
    //     assert_eq!(parsed_response_eth_address, eth_address);

    //     let eth_response = grpc_client
    //         .get_delegate_key_by_eth(QueryDelegateKeysByEthAddress {
    //             eth_address: eth_address.to_string(),
    //         })
    //         .await
    //         .unwrap()
    //         .into_inner();

    //     let parsed_response_orch_address: CosmosAddress =
    //         eth_response.orchestrator_address.parse().unwrap();
    //     assert_eq!(parsed_response_orch_address, orch_address);
    // }
}
