use crate::args::SetOrchestratorAddressOpts;
use crate::utils::TIMEOUT;
use clarity::PrivateKey as EthPrivateKey;
use cosmos_gravity::send::update_gravity_delegate_addresses;
use deep_space::{mnemonic::Mnemonic, private_key::PrivateKey as CosmosPrivateKey};
use gravity_utils::connection_prep::check_for_fee;
use gravity_utils::connection_prep::{create_rpc_connections, wait_for_cosmos_node_ready};
use rand::{thread_rng, Rng};

pub async fn set_orchestrator_address(args: SetOrchestratorAddressOpts, prefix: String) {
    let fee = args.fees;
    let cosmos_grpc = args.cosmos_grpc;
    let validator_phrase = args.validator_phrase;
    let cosmos_phrase = args.cosmos_phrase;
    let mut generated_cosmos = None;
    let mut generated_eth = false;

    let connections = create_rpc_connections(prefix, Some(cosmos_grpc), None, TIMEOUT).await;
    let contact = connections.contact.unwrap();
    wait_for_cosmos_node_ready(&contact).await;

    let validator_key = CosmosPrivateKey::from_phrase(&validator_phrase, "")
        .expect("Failed to parse validator key");
    let validator_addr = validator_key.to_address(&contact.get_prefix()).unwrap();
    check_for_fee(&fee, validator_addr, &contact).await;

    let cosmos_key = if let Some(cosmos_phrase) = cosmos_phrase {
        CosmosPrivateKey::from_phrase(&cosmos_phrase, "").expect("Failed to parse cosmos key")
    } else {
        let new_phrase = Mnemonic::generate(24).unwrap();
        let key = CosmosPrivateKey::from_phrase(new_phrase.as_str(), "").unwrap();
        generated_cosmos = Some(new_phrase);
        key
    };
    let ethereum_key = if let Some(key) = args.ethereum_key {
        key
    } else {
        generated_eth = true;
        let mut rng = thread_rng();
        let key: [u8; 32] = rng.gen();
        EthPrivateKey::from_slice(&key).unwrap()
    };

    let ethereum_address = ethereum_key.to_public_key().unwrap();
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    update_gravity_delegate_addresses(
        &contact,
        ethereum_address,
        cosmos_address,
        validator_key,
        fee.clone(),
    )
    .await
    .expect("Failed to update Eth address");

    if let Some(phrase) = generated_cosmos {
        info!(
            "No Cosmos key provided, your generated key is\n {} -> {}",
            phrase.as_str(),
            cosmos_key.to_address(&contact.get_prefix()).unwrap()
        );
    }
    if generated_eth {
        info!(
            "No Ethereum key provided, your generated key is\n Private: {} -> Address: {}",
            ethereum_key,
            ethereum_key.to_public_key().unwrap()
        );
    }

    let eth_address = ethereum_key.to_public_key().unwrap();
    info!(
        "Registered Delegate Ethereum address {} and Cosmos address {}",
        eth_address, cosmos_address
    )
}
