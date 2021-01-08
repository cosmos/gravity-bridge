use clarity::abi::Token;
use clarity::Uint256;
use clarity::{abi::encode_tokens, Address as EthAddress};
use deep_space::address::Address as CosmosAddress;
use peggy_utils::error::PeggyError;
use peggy_utils::types::*;
use sha3::{Digest, Keccak256};
use std::u64::MAX as U64MAX;
use web30::{client::Web3, jsonrpc::error::Web3Error};

pub fn get_correct_sig_for_address(
    address: CosmosAddress,
    confirms: &[ValsetConfirmResponse],
) -> (Uint256, Uint256, Uint256) {
    for sig in confirms {
        if sig.orchestrator == address {
            return (
                sig.eth_signature.v.clone(),
                sig.eth_signature.r.clone(),
                sig.eth_signature.s.clone(),
            );
        }
    }
    panic!("Could not find that address!");
}

pub fn get_checkpoint_abi_encode(valset: &Valset, peggy_id: &str) -> Result<Vec<u8>, PeggyError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    Ok(encode_tokens(&[
        Token::FixedString(peggy_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, peggy_id: &str) -> Result<Vec<u8>, PeggyError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode?);
    Ok(locally_computed_digest.to_vec())
}

pub fn downcast_nonce(input: Uint256) -> Option<u64> {
    if input >= U64MAX.into() {
        None
    } else {
        let mut val = input.to_bytes_be();
        // pad to 8 bytes
        while val.len() < 8 {
            val.insert(0, 0);
        }
        let mut lower_bytes: [u8; 8] = [0; 8];
        // get the 'lowest' 8 bytes from a 256 bit integer
        lower_bytes.copy_from_slice(&val[0..val.len()]);
        Some(u64::from_be_bytes(lower_bytes))
    }
}

#[test]
fn test_downcast_nonce() {
    let mut i = 0u64;
    while i < 100_000 {
        assert_eq!(i, downcast_nonce(i.into()).unwrap());
        i += 1
    }
    let mut i: u64 = std::u32::MAX.into();
    i -= 100;
    let end = i + 100_000;
    while i < end {
        assert_eq!(i, downcast_nonce(i.into()).unwrap());
        i += 1
    }
}

/// Gets the latest validator set nonce
pub async fn get_valset_nonce(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let val = web3
        .contract_call(
            contract_address,
            "state_lastValsetNonce()",
            &[],
            caller_address,
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_nonce(real_num).expect("Valset nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    peggy_contract_address: EthAddress,
    erc20_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let val = web3
        .contract_call(
            peggy_contract_address,
            "lastBatchNonce(address)",
            &[erc20_contract_address.into()],
            caller_address,
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_nonce(real_num).expect("TxBatch nonce overflow! Bridge Halt!"))
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

/// Gets the ERC20 symbol, should maybe be upstreamed
pub async fn get_erc20_symbol(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<String, PeggyError> {
    let val_symbol = web3
        .contract_call(contract_address, "symbol()", &[], caller_address)
        .await?;
    // Pardon the unwrap, but this is temporary code, intended only for the tests, to help them
    // deal with a deprecated feature (the symbol), which will be removed soon
    Ok(String::from_utf8(val_symbol).unwrap())
}
