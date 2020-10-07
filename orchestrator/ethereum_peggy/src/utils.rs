use clarity::Uint256;
use clarity::{abi::encode_tokens, Address as EthAddress};
use clarity::{abi::Token, PrivateKey as EthPrivateKey};
use cosmos_peggy::{send::filter_empty_addresses, types::*};
use deep_space::address::Address as CosmosAddress;
use peggy_utils::error::OrchestratorError;
use sha3::{Digest, Keccak256};
use web30::client::Web3;
use web30::types::SendTxOption;

pub fn get_correct_power_for_address(address: EthAddress, valset: &Valset) -> (EthAddress, u64) {
    for (a, p) in valset.eth_addresses.iter().zip(valset.powers.iter()) {
        if let Some(a) = a {
            if *a == address {
                return (*a, *p);
            }
        }
    }
    panic!("Could not find that address!");
}

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

/// this function is used to determine if all validator eth addresses are set. Since
/// this is required for many bridge operations it will return with an error if any
/// validator eth address is unset.
pub fn filter_empty_eth_addresses(
    input: &[Option<EthAddress>],
) -> Result<Vec<EthAddress>, OrchestratorError> {
    let mut res = Vec::new();
    for val in input {
        if let Some(addr) = val {
            res.push(*addr);
        } else {
            return Err(OrchestratorError::InvalidBridgeStateError(
                "Validator without registered Ethereum key".to_string(),
            ));
        }
    }
    Ok(res)
}

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    peggy_id: &str,
) -> Result<Vec<u8>, OrchestratorError> {
    Ok(encode_tokens(&[
        Token::FixedString(peggy_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        filter_empty_addresses(&valset.eth_addresses)?.into(),
        valset.powers.clone().into(),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, peggy_id: &str) -> Result<Vec<u8>, OrchestratorError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode?);
    Ok(locally_computed_digest.to_vec())
}
