use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use peggy_utils::error::OrchestratorError;
use peggy_utils::types::*;
use std::time::Duration;
use tokio::time::timeout as future_timeout;
use web30::client::Web3;
use web30::types::SendTxOption;

/// this function generates an appropriate Ethereum transaction
/// to submit the provided validator set and signatures.
/// TODO this function uses the same validator set as the old and
/// new validator set, this is because there's no actual changes to
/// the set in testing and because there's no oracle to tell us what
/// the old set was anyways.
/// TODO TODO should we have an oracle for the old set or look in the chain?
pub async fn send_eth_valset_update(
    new_valset: Valset,
    old_valset: Valset,
    confirms: &[ValsetConfirmResponse],
    web3: &Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    our_eth_key: EthPrivateKey,
) -> Result<(), OrchestratorError> {
    let (old_addresses, old_powers) = old_valset.filter_empty_addresses()?;
    let (new_addresses, new_powers) = new_valset.filter_empty_addresses()?;
    let old_nonce = old_valset.nonce;
    let new_nonce = new_valset.nonce;
    let eth_address = our_eth_key.to_public_key().unwrap();
    info!(
        "Ordering signatures and submitting validator set {} -> {} update to Ethereum",
        old_nonce, new_nonce
    );

    let sig_data = new_valset.order_sigs(confirms)?;
    let sig_arrays = to_arrays(sig_data);

    // Solidity function signature
    // function getValsetNonce() public returns (uint256)
    let first_nonce = web3
        .contract_call(peggy_contract_address, "getValsetNonce()", &[], eth_address)
        .await
        .expect("Failed to get the first nonce");
    info!("Current valset nonce {:?}", first_nonce);

    // Solidity function signature
    // function updateValset(
    // // The new version of the validator set
    // address[] memory _newValidators,
    // uint256[] memory _newPowers,
    // uint256 _newValsetNonce,
    // // The current validators that approve the change
    // address[] memory _currentValidators,
    // uint256[] memory _currentPowers,
    // uint256 _currentValsetNonce,
    // // These are arrays of the parts of the current validator's signatures
    // uint8[] memory _v,
    // bytes32[] memory _r,
    // bytes32[] memory _s
    let payload = clarity::abi::encode_call("updateValset(address[],uint256[],uint256,address[],uint256[],uint256,uint8[],bytes32[],bytes32[])",
    &[new_addresses.into(), new_powers.into(), new_nonce.into(), old_addresses.into(), old_powers.into(), old_nonce.into(), sig_arrays.v, sig_arrays.r, sig_arrays.s]).unwrap();

    let tx = future_timeout(
        timeout,
        web3.send_transaction(
            peggy_contract_address,
            payload,
            0u32.into(),
            eth_address,
            our_eth_key,
            vec![SendTxOption::GasLimit(1_000_000u32.into())],
        ),
    )
    .await
    .expect("Valset update timed out")
    .expect("Valset update failed for other reasons");
    info!("Finished valset update with txid {:#066x}", tx);

    web3.wait_for_transaction(tx, timeout, None).await.unwrap();

    // Solidity function signature
    // function getValsetNonce() public returns (uint256)
    let last_nonce = web3
        .contract_call(peggy_contract_address, "getValsetNonce()", &[], eth_address)
        .await
        .expect("Failed to get the last nonce");
    assert!(first_nonce != last_nonce);
    info!(
        "Successfully updated Valset with new Nonce {:?}",
        last_nonce
    );
    Ok(())
}
