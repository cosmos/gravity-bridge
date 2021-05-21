use crate::args::OrchestratorOpts;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use gravity_utils::connection_prep::{
    check_delegate_addresses, check_for_eth, wait_for_cosmos_node_ready,
};
use gravity_utils::connection_prep::{check_for_fee_denom, create_rpc_connections};
use orchestrator::main_loop::orchestrator_main_loop;
use orchestrator::main_loop::{ETH_ORACLE_LOOP_SPEED, ETH_SIGNER_LOOP_SPEED};
use relayer::main_loop::LOOP_SPEED as RELAYER_LOOP_SPEED;
use std::cmp::min;

pub async fn orchestrator(args: OrchestratorOpts, address_prefix: String) {
    let fee_denom = args.fees;
    let cosmos_grpc = args.cosmos_grpc;
    let ethereum_rpc = args.ethereum_rpc;
    let ethereum_key = args.ethereum_key;
    let cosmos_phrase = args.cosmos_phrase;
    let contract_address = args.gravity_contract_address;

    let cosmos_key =
        CosmosPrivateKey::from_phrase(&cosmos_phrase, "").expect("Failed to parse cosmos key");

    let timeout = min(
        min(ETH_SIGNER_LOOP_SPEED, ETH_ORACLE_LOOP_SPEED),
        RELAYER_LOOP_SPEED,
    );

    trace!("Probing RPC connections");
    // probe all rpc connections and see if they are valid
    let connections = create_rpc_connections(
        address_prefix,
        Some(cosmos_grpc),
        Some(ethereum_rpc),
        timeout,
    )
    .await;

    let mut grpc = connections.grpc.clone().unwrap();
    let contact = connections.contact.clone().unwrap();
    let web3 = connections.web3.clone().unwrap();

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    let public_cosmos_key = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    info!("Starting Gravity Validator companion binary Relayer + Oracle + Eth Signer");
    info!(
        "Ethereum Address: {} Cosmos Address {}",
        public_eth_key, public_cosmos_key
    );

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;

    // check if the delegate addresses are correctly configured
    check_delegate_addresses(
        &mut grpc,
        public_eth_key,
        public_cosmos_key,
        &contact.get_prefix(),
    )
    .await;

    // check if we actually have the promised balance of tokens to pay fees
    check_for_fee_denom(&fee_denom, public_cosmos_key, &contact).await;
    check_for_eth(public_eth_key, &web3).await;

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
