//! Ethereum Event watcher watches for events such as a deposit to the Peggy Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use clarity::{Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::send::send_ethereum_claims;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use peggy_utils::{
    error::PeggyError,
    types::{SendToCosmosEvent, TransactionBatchExecutedEvent, ValsetUpdatedEvent},
};
use web30::client::Web3;
use web30::jsonrpc::error::Web3Error;

pub async fn check_for_events(
    web3: &Web3,
    contact: &Contact,
    peggy_contract_address: EthAddress,
    our_private_key: CosmosPrivateKey,
    fee: Coin,
    last_checked_block: Uint256,
) -> Result<Uint256, PeggyError> {
    let starting_block = last_checked_block + 1u8.into();
    let latest_block = web3.eth_block_number().await?;
    let deposits = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![peggy_contract_address],
            "SendToCosmosEvent(address,address,bytes32,uint256,uint256)",
            Vec::new(),
        )
        .await;
    trace!("Deposits {:?}", deposits);
    // todo write a parser for each of these events to get the data out
    // then send a cosmos transaction to mint tokens

    let batches = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![peggy_contract_address],
            "TransactionBatchExecutedEvent(uint256,address,uint256)",
            Vec::new(),
        )
        .await;
    trace!("Batches {:?}", batches);

    let valsets = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![peggy_contract_address],
            "ValsetUpdatedEvent(uint256,address[],uint256[])",
            Vec::new(),
        )
        .await;
    trace!("Valsets {:?}", valsets);

    if let (Ok(valsets), Ok(batches), Ok(deposits)) = (valsets, batches, deposits) {
        let valsets = ValsetUpdatedEvent::from_logs(&valsets)?;
        trace!("parsed valsets {:?}", valsets);
        let withdraws = TransactionBatchExecutedEvent::from_logs(&batches)?;
        trace!("parsed batches {:?}", batches);
        let deposits = SendToCosmosEvent::from_logs(&deposits)?;
        trace!("parsed deposits {:?}", deposits);
        if !deposits.is_empty() {
            info!(
                "Oracle observed deposit with sender {}, destination {}, amount {}, and event nonce {}",
                deposits[0].sender, deposits[0].destination, deposits[0].amount, deposits[0].event_nonce
            )
        }

        if !deposits.is_empty() || !withdraws.is_empty() {
            // todo get eth chain id from the chain
            let res =
                send_ethereum_claims(contact, our_private_key, deposits, withdraws, fee).await?;
            trace!("Sent in Oracle claims response: {:?}", res);
        }

        Ok(latest_block)
    } else {
        error!("Failed to get events");
        Err(PeggyError::EthereumRestError(Web3Error::BadResponse(
            "Failed to get logs!".to_string(),
        )))
    }
}
