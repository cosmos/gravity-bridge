use crate::utils::get_tx_batch_nonce;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use peggy_utils::error::PeggyError;
use peggy_utils::types::*;
use std::time::Duration;
use web30::client::Web3;
use web30::types::SendTxOption;

/// this function generates an appropriate Ethereum transaction
/// to submit the provided transaction batch and validator set update.
pub async fn send_eth_transaction_batch(
    current_valset: Valset,
    batch: TransactionBatch,
    confirms: &[BatchConfirmResponse],
    web3: &Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<(), PeggyError> {
    let (current_addresses, current_powers) = current_valset.filter_empty_addresses();
    let current_valset_nonce = current_valset.nonce;
    let new_batch_nonce = batch.nonce;
    //assert!(new_valset_nonce > old_valset_nonce);
    let eth_address = our_eth_key.to_public_key().unwrap();
    info!(
        "Ordering signatures and submitting TransacqtionBatch {}:{} to Ethereum",
        batch.token_contract, new_batch_nonce
    );
    trace!("Batch {:?}", batch);

    let sig_data = current_valset.order_batch_sigs(confirms)?;
    let sig_arrays = to_arrays(sig_data);
    let (amounts, destinations, fees) = batch.get_checkpoint_values();

    // Solidity function signature
    // function submitBatch(
    // // The validators that approve the batch and new valset
    // address[] memory _currentValidators,
    // uint256[] memory _currentPowers,
    // uint256 _currentValsetNonce,
    // // These are arrays of the parts of the validators signatures
    // uint8[] memory _v,
    // bytes32[] memory _r,
    // bytes32[] memory _s,
    // // The batch of transactions
    // uint256[] memory _amounts,
    // address[] memory _destinations,
    // uint256[] memory _fees,
    // uint256 _batchNonce,
    // address _tokenContract,
    // uint256 _batchTimeout
    let tokens = &[
        current_addresses.into(),
        current_powers.into(),
        current_valset_nonce.into(),
        sig_arrays.v,
        sig_arrays.r,
        sig_arrays.s,
        amounts,
        destinations,
        fees,
        new_batch_nonce.clone().into(),
        batch.token_contract.into(),
        batch.batch_timeout.into(),
    ];
    let payload = clarity::abi::encode_call("submitBatch(address[],uint256[],uint256,uint8[],bytes32[],bytes32[],uint256[],address[],uint256[],uint256,address,uint256)",
    tokens).unwrap();
    trace!("Tokens {:?}", tokens);

    let before_nonce = get_tx_batch_nonce(
        peggy_contract_address,
        batch.token_contract,
        eth_address,
        &web3,
    )
    .await?;
    if before_nonce >= new_batch_nonce {
        info!(
            "Someone else updated the batch to {}, exiting early",
            before_nonce
        );
        return Ok(());
    }

    let tx = web3
        .send_transaction(
            peggy_contract_address,
            payload,
            0u32.into(),
            eth_address,
            our_eth_key,
            vec![SendTxOption::GasLimit(1_000_000u32.into())],
        )
        .await?;
    info!("Sent batch update with txid {:#066x}", tx);

    // TODO this segment of code works around the race condition for submitting batches mostly
    // by not caring if our own submission reverts and only checking if the valset has been updated
    // period not if our update succeeded in particular. This will require some further consideration
    // in the future as many independent relayers racing to update the same thing will hopefully
    // be the common case.
    web3.wait_for_transaction(tx, timeout, None).await?;

    let last_nonce = get_tx_batch_nonce(
        peggy_contract_address,
        batch.token_contract,
        eth_address,
        &web3,
    )
    .await?;
    if last_nonce != new_batch_nonce {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_batch_nonce
        );
    } else {
        info!("Successfully updated Batch with new Nonce {:?}", last_nonce);
    }
    Ok(())
}
