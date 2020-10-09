use crate::utils::filter_empty_eth_addresses;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use clarity::{abi::Token, Uint256};
use cosmos_peggy::types::*;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use peggy_utils::error::OrchestratorError;
use std::time::Duration;
use tokio::time::timeout as future_timeout;
use web30::types::SendTxOption;
use web30::{client::Web3, jsonrpc::error::Web3Error};

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
    our_comsmos_key: CosmosPrivateKey,
) -> Result<(), OrchestratorError> {
    info!("Ordering signatures and submitting validator set update to Ethereum");
    let old_addresses = filter_empty_eth_addresses(&old_valset.eth_addresses)?;
    let old_powers = old_valset.powers;
    let new_addresses = filter_empty_eth_addresses(&new_valset.eth_addresses.clone())?;
    let new_powers = new_valset.powers.clone();
    let old_nonce = old_valset.nonce;
    let new_nonce = new_valset.nonce;
    let mut v: Vec<u8> = Vec::new();
    let mut r = Vec::new();
    let mut s = Vec::new();
    //replace this with a function to get ordered addresses and sigs
    // for address in old_addresses.iter() {
    //     let cosmos_address = get_cosmos_address_from_eth_addr(*address, &keys);
    //     let (sig_v, sig_r, sig_s) = get_correct_sig_for_address(cosmos_address, confirms);
    //     v.push(sig_v.clone());
    //     r.push(Token::Bytes(sig_r.clone().to_bytes_be()));
    //     s.push(Token::Bytes(sig_s.clone().to_bytes_be()));
    // }
    let eth_address = our_eth_key.to_public_key().unwrap();

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
    &[new_addresses.into(), new_powers.into(), new_nonce.into(), old_addresses.into(), old_powers.into(), old_nonce.into(), v.into(), Token::Dynamic(r), Token::Dynamic(s)]).unwrap();

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
