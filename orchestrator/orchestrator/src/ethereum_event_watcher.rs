//! Ethereum Event watcher watches for events such as a deposit to the Peggy Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use clarity::{Address as EthAddress, Uint256};
use peggy_utils::{
    error::OrchestratorError,
    types::{SendToCosmosEvent, TransactionBatchExecutedEvent, ValsetUpdatedEvent},
};
use web30::jsonrpc::error::Web3Error;
use web30::{client::Web3, types::Log};

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
            vec![peggy_contract_address],
            "SendToCosmosEvent(address,address,bytes32,uint256)",
            Vec::new(),
        )
        .await;
    trace!("Deposits {:?}", deposits);
    // todo write a parser for each of these events to get the data out
    // then send a cosmos transaction to mint tokens

    let batches = web3
        .check_for_events(
            last_checked_block.clone(),
            Some(latest_block.clone()),
            vec![peggy_contract_address],
            "TransactionBatchExecutedEvent(uint256,address)",
            Vec::new(),
        )
        .await;
    trace!("Batches {:?}", batches);

    let valsets = web3
        .check_for_events(
            last_checked_block.clone(),
            Some(latest_block.clone()),
            vec![peggy_contract_address],
            "ValsetUpdatedEvent(uint256,address[],uint256[])",
            Vec::new(),
        )
        .await;
    trace!("Valsets {:?}", valsets);
    if let (Ok(valsets), Ok(batches), Ok(deposits)) = (valsets, batches, deposits) {
        let valsets = ValsetUpdatedEvent::from_logs(&valsets)?;
        let batches = TransactionBatchExecutedEvent::from_logs(&batches)?;
        let deposits = SendToCosmosEvent::from_logs(&deposits)?;
        info!("Parsed stuff!");
        if !deposits.is_empty() {
            info!(
                "Got deposit with sender {} and destination {} and amount {}",
                deposits[0].sender, deposits[0].destination, deposits[0].amount
            )
        }
        Ok(latest_block)
    } else {
        error!("Failed to get events");
        Err(OrchestratorError::EthereumRestErr(Web3Error::BadResponse(
            "Failed to get logs!".to_string(),
        )))
    }
}
