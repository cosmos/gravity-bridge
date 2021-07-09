use clarity::abi::Token;
use clarity::Uint256;
use clarity::{abi::encode_tokens, Address as EthAddress};
use gravity_utils::error::GravityError;
use gravity_utils::types::*;
use sha3::{Digest, Keccak256};
use std::u128::MAX as U128MAX;
use std::u64::MAX as U64MAX;
use web30::{client::Web3, jsonrpc::error::Web3Error};

// pub fn get_correct_sig_for_address(
//     address: CosmosAddress,
//     confirms: &[ValsetConfirmResponse],
// ) -> (Uint256, Uint256, Uint256) {
//     for sig in confirms {
//         if sig.eth_signer == address {
//             return (
//                 sig.eth_signature.v.clone(),
//                 sig.eth_signature.r.clone(),
//                 sig.eth_signature.s.clone(),
//             );
//         }
//     }
//     panic!("Could not find that address!");
// }

pub fn get_checkpoint_abi_encode(
    valset: &Valset,
    gravity_id: &str,
) -> Result<Vec<u8>, GravityError> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    Ok(encode_tokens(&[
        Token::FixedString(gravity_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ]))
}

pub fn get_checkpoint_hash(valset: &Valset, gravity_id: &str) -> Result<Vec<u8>, GravityError> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &gravity_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode?);
    Ok(locally_computed_digest.to_vec())
}

pub fn downcast_uint256(input: Uint256) -> Option<u64> {
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

pub fn downcast_to_u128(input: Uint256) -> Option<u128> {
    if input >= U128MAX.into() {
        None
    } else {
        let mut val = input.to_bytes_be();
        // pad to 8 bytes
        while val.len() < 16 {
            val.insert(0, 0);
        }
        let mut lower_bytes: [u8; 16] = [0; 16];
        // get the 'lowest' 16 bytes from a 256 bit integer
        lower_bytes.copy_from_slice(&val[0..val.len()]);
        Some(u128::from_be_bytes(lower_bytes))
    }
}

#[test]
fn test_downcast_nonce() {
    let mut i = 0u64;
    while i < 100_000 {
        assert_eq!(i, downcast_uint256(i.into()).unwrap());
        i += 1
    }
    let mut i: u64 = std::u32::MAX.into();
    i -= 100;
    let end = i + 100_000;
    while i < end {
        assert_eq!(i, downcast_uint256(i.into()).unwrap());
        i += 1
    }
}

#[test]
fn test_downcast_to_u128() {
    let mut i = 0u128;
    while i < 100_000 {
        assert_eq!(i, downcast_to_u128(i.into()).unwrap());
        i += 1
    }
    let mut i: u128 = std::u64::MAX.into();
    i -= 100;
    let end = i + 100_000;
    while i < end {
        assert_eq!(i, downcast_to_u128(i.into()).unwrap());
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
            caller_address, None
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_uint256(real_num).expect("Valset nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_tx_batch_nonce(
    gravity_contract_address: EthAddress,
    erc20_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let val = web3
        .contract_call(
            gravity_contract_address,
            "lastBatchNonce(address)",
            &[erc20_contract_address.into()],
            caller_address, None
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_uint256(real_num).expect("TxBatch nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_logic_call_nonce(
    gravity_contract_address: EthAddress,
    invalidation_id: Vec<u8>,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let val = web3
        .contract_call(
            gravity_contract_address,
            "lastLogicCallNonce(bytes32)",
            &[Token::Bytes(invalidation_id)],
            caller_address, None
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_uint256(real_num).expect("LogicCall nonce overflow! Bridge Halt!"))
}

/// Gets the latest transaction batch nonce
pub async fn get_event_nonce(
    gravity_contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<u64, Web3Error> {
    let val = web3
        .contract_call(
            gravity_contract_address,
            "state_lastEventNonce()",
            &[],
            caller_address, None
        )
        .await?;
    // the go represents all nonces as u64, there's no
    // reason they should ever overflow without a user
    // submitting millions or tens of millions of dollars
    // worth of transactions. But we properly check and
    // handle that case here.
    let real_num = Uint256::from_bytes_be(&val);
    Ok(downcast_uint256(real_num).expect("EventNonce nonce overflow! Bridge Halt!"))
}

/// Gets the gravityID
pub async fn get_gravity_id(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<Vec<u8>, Web3Error> {
    let val = web3
        .contract_call(contract_address, "state_gravityId()", &[], caller_address, None)
        .await?;
    Ok(val)
}

/// Gets the ERC20 symbol, should maybe be upstreamed
pub async fn get_erc20_symbol(
    contract_address: EthAddress,
    caller_address: EthAddress,
    web3: &Web3,
) -> Result<String, GravityError> {
    let val_symbol = web3
        .contract_call(contract_address, "symbol()", &[], caller_address, None)
        .await?;
    // Pardon the unwrap, but this is temporary code, intended only for the tests, to help them
    // deal with a deprecated feature (the symbol), which will be removed soon
    Ok(String::from_utf8(val_symbol).unwrap())
}

/// Just a helper struct to represent the cost of actions on Ethereum
#[derive(Debug, Default, Clone)]
pub struct GasCost {
    pub gas: Uint256,
    pub gas_price: Uint256,
}

impl GasCost {
    pub fn get_total(&self) -> Uint256 {
        self.gas.clone() * self.gas_price.clone()
    }
}
