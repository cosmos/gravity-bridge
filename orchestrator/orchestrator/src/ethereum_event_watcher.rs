//! Ethereum Event watcher watches for events such as a deposit to the Peggy Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use clarity::{Address as EthAddress, Uint256};
use peggy_utils::error::OrchestratorError;
use web30::client::Web3;

pub async fn check_for_events(
    web3: &Web3,
    peggy_contract_address: EthAddress,
    last_checked_block: Uint256,
) -> Result<Uint256, OrchestratorError> {
    let latest_block = web3.eth_block_number().await?;
    let deposits = web3
        .check_for_events(
            last_checked_block.clone(),
            Some(latest_block.clone()),
            peggy_contract_address,
            "SendToCosmos(address,address,bytes32,uint256)",
            None,
            None,
        )
        .await;
    info!("Deposits {:?}", deposits);

    let batches = web3
        .check_for_events(
            last_checked_block.clone(),
            Some(latest_block.clone()),
            peggy_contract_address,
            "TransactionBatchExecutedEvent(uint256,address)",
            None,
            None,
        )
        .await;
    info!("Batches {:?}", deposits);

    let valsets = web3
        .check_for_events(
            last_checked_block.clone(),
            Some(latest_block.clone()),
            peggy_contract_address,
            "ValsetUpdatedEvent(address[],uint256[])",
            None,
            None,
        )
        .await;
    info!("Valsets {:?}", valsets);

    Ok(latest_block)
}
