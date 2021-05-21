use crate::args::RelayerOpts;
use gravity_utils::connection_prep::{
    check_for_eth, create_rpc_connections, wait_for_cosmos_node_ready,
};
use relayer::main_loop::relayer_main_loop;
use relayer::main_loop::LOOP_SPEED;

pub async fn relayer(args: RelayerOpts, address_prefix: String) {
    let cosmos_grpc = args.cosmos_grpc;
    let ethereum_rpc = args.ethereum_rpc;
    let ethereum_key = args.ethereum_key;
    let gravity_contract_address = args.gravity_contract_address;
    let connections = create_rpc_connections(
        address_prefix,
        Some(cosmos_grpc),
        Some(ethereum_rpc),
        LOOP_SPEED,
    )
    .await;

    let public_eth_key = ethereum_key
        .to_public_key()
        .expect("Invalid Ethereum Private Key!");
    info!("Starting Gravity Relayer");
    info!("Ethereum Address: {}", public_eth_key);

    let contact = connections.contact.clone().unwrap();
    let web3 = connections.web3.clone().unwrap();

    // check if the cosmos node is syncing, if so wait for it
    // we can't move any steps above this because they may fail on an incorrect
    // historic chain state while syncing occurs
    wait_for_cosmos_node_ready(&contact).await;
    check_for_eth(public_eth_key, &web3).await;

    relayer_main_loop(
        ethereum_key,
        connections.web3.unwrap(),
        connections.grpc.unwrap(),
        gravity_contract_address,
    )
    .await
}
