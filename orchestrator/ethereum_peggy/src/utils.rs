use clarity::abi::Token;
use clarity::Uint256;
use clarity::{abi::encode_tokens, Address as EthAddress};
use deep_space::address::Address as CosmosAddress;
use peggy_utils::error::OrchestratorError;
use peggy_utils::types::*;
use sha3::{Digest, Keccak256};
use web30::{client::Web3, jsonrpc::error::Web3Error};

pub fn get_correct_sig_for_address(
    address: CosmosAddress,
    confirms: &[ValsetConfirmResponse],
) -> (Uint256, Uint256, Uint256) {
    for sig in confirms {
        if sig.validator == address {
            return (
                sig.eth_signature.v.clone(),
                sig.eth_signature.r.clone(),
                sig.eth_signature.s.clone(),
            );
        }
    }
    panic!("Could not find that address!");
}

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    peggy_id: &str,
) -> Result<Vec<u8>, OrchestratorError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses()?;
    Ok(encode_tokens(&[
        Token::FixedString(peggy_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, peggy_id: &str) -> Result<Vec<u8>, OrchestratorError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode?);
    Ok(locally_computed_digest.to_vec())
}

/// Gets the latest validator set nonce
pub async fn get_valset_nonce(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<Uint256, Web3Error> {
    let val = web3
        .contract_call(
            contract_address,
            "state_lastValsetNonce()",
            &[],
            caller_address,
        )
        .await?;
    Ok(Uint256::from_bytes_be(&val))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<Uint256, Web3Error> {
    let val = web3
        .contract_call(
            contract_address,
            "state_lastBatchNonces()",
            &[],
            caller_address,
        )
        .await?;
    Ok(Uint256::from_bytes_be(&val))
}

/// Gets the peggyID
pub async fn get_peggy_id(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<Vec<u8>, Web3Error> {
    let val = web3
        .contract_call(contract_address, "state_peggyId()", &[], caller_address)
        .await?;
    Ok(val)
}
