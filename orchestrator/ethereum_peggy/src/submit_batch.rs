use crate::utils::get_valset_nonce;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use peggy_utils::error::OrchestratorError;
use peggy_utils::types::*;
use std::time::Duration;
use web30::client::Web3;
use web30::types::SendTxOption;

/// this function generates an appropriate Ethereum transaction
/// to submit the provided transaction batch and validator set update.
pub async fn send_eth_transaction_batch(
    old_valset: Valset,
    batch: SignedTransactionBatch,
    web3: &Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<(), OrchestratorError> {
    let (old_addresses, old_powers) = old_valset.filter_empty_addresses();
    let (new_addresses, new_powers) = batch.batch.valset.filter_empty_addresses();
    let old_nonce = old_valset.nonce;
    let new_nonce = batch.batch.valset.nonce;
    assert!(new_nonce > old_nonce);
    let eth_address = our_eth_key.to_public_key().unwrap();
    info!(
        "Ordering signatures and submitting TransactionBatch {} -> {} update to Ethereum",
        old_nonce, new_nonce
    );

    let new_valset = batch.batch.valset.clone();
    let sig_data = new_valset.order_batch_sigs(batch.clone())?;
    let sig_arrays = to_arrays(sig_data);
    let (amounts, destinations, fees) = batch.batch.get_checkpoint_values();

    // Solidity function signature
    // function updateValsetAndSubmitBatch(
    // // The new version of the validator set
    // address[] memory _newValidators,
    // uint256[] memory _newPowers,
    // uint256 _newValsetNonce,
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
    // uint256 _nonce,
    // address _tokenContract
    let payload = clarity::abi::encode_call("updateValsetAndSubmitBatch(address[],uint256[],uint256,address[],uint256[],uint256,uint8[],bytes32[],bytes32[],uint256[],address[],uint256[],uint256,address)",
    &[new_addresses.into(), new_powers.into(), new_nonce.into(), old_addresses.into(), old_powers.into(), old_nonce.into(), sig_arrays.v, sig_arrays.r, sig_arrays.s, amounts, destinations, fees, new_nonce.into(), batch.batch.token_contract.into()]).unwrap();

    let before_nonce = get_valset_nonce(peggy_contract_address, eth_address, web3).await?;
    if before_nonce != old_nonce.into() {
        info!(
            "Someone else updated the valset to {}, exiting early",
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
    info!("Finished valset update with txid {:#066x}", tx);

    // TODO this segment of code works around the race condition for submitting valsets mostly
    // by not caring if our own submission reverts and only checking if the valset has been updated
    // period not if our update succeeded in particular. This will require some further consideration
    // in the future as many independent relayers racing to update the same thing will hopefully
    // be the common case.
    web3.wait_for_transaction(tx, timeout, None).await?;

    let last_nonce = get_valset_nonce(peggy_contract_address, eth_address, web3).await?;
    if last_nonce != new_nonce.into() {
        error!(
            "Current nonce is {} expected to update to nonce {}",
            last_nonce, new_nonce
        );
    } else {
        info!(
            "Successfully updated Valset with new Nonce {:?}",
            last_nonce
        );
    }
    Ok(())
}
