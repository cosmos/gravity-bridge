//! This test verifies that live updating of orchestrator keys works correctly

use crate::get_fee;
use crate::utils::ValidatorKeys;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use cosmos_gravity::send::update_gravity_delegate_addresses;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use futures::future::join_all;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::connection_prep::check_delegate_addresses;
use rand::Rng;
use tonic::transport::Channel;

pub async fn orch_keys_update(
    grpc_client: GravityQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<ValidatorKeys>,
) {
    let mut keys = keys;
    let mut grpc_client = grpc_client;
    // just to test that we have the right keys from the gentx
    for k in keys.iter() {
        check_delegate_addresses(
            &mut grpc_client,
            k.eth_key.to_public_key().unwrap(),
            k.orch_key.to_public_key().unwrap().to_address(),
        )
        .await;
    }

    // now we change them all
    let mut updates = Vec::new();
    for k in keys.iter_mut() {
        let mut rng = rand::thread_rng();
        let secret: [u8; 32] = rng.gen();
        // generate some new keys to replace the old ones
        let eth_key = EthPrivateKey::from_slice(&secret).unwrap();
        let cosmos_key = CosmosPrivateKey::from_secret(&secret);
        // update the keys in the key list
        k.eth_key = eth_key;
        k.orch_key = cosmos_key;

        info!(
            "Signing and submitting Delegate addresses {} for validator {}",
            eth_key.to_public_key().unwrap(),
            cosmos_key.to_public_key().unwrap().to_address(),
        );
        // send in the new delegate keys signed by the validator address
        updates.push(update_gravity_delegate_addresses(
            &contact,
            eth_key.to_public_key().unwrap(),
            cosmos_key.to_public_key().unwrap().to_address(),
            k.validator_key,
            get_fee(),
        ));
    }
    let update_results = join_all(updates).await;
    for i in update_results {
        i.expect("Failed to set delegate addresses!");
    }

    // verify that the change has taken place
    for k in keys.iter() {
        check_delegate_addresses(
            &mut grpc_client,
            k.eth_key.to_public_key().unwrap(),
            k.orch_key.to_public_key().unwrap().to_address(),
        )
        .await;
    }
}
