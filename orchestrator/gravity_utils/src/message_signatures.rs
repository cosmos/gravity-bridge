use crate::types::{LogicCall, TransactionBatch, Valset};
use clarity::abi::{encode_tokens, Token};
use clarity::utils::get_ethereum_msg_hash;

/// takes the required input data and produces the required signature to confirm a validator
/// set update on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_valset_confirm(gravity_id: String, valset: Valset) -> Vec<u8> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    encode_tokens(&[
        Token::FixedString(gravity_id),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ])
}

pub fn encode_valset_confirm_hashed(gravity_id: String, valset: Valset) -> Vec<u8> {
    let digest = encode_valset_confirm(gravity_id, valset);
    get_ethereum_msg_hash(&digest)
}

#[test]
fn test_valset_signature() {
    use crate::types::ValsetMember;
    use clarity::utils::hex_str_to_bytes;
    use sha3::{Digest, Keccak256};

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0x88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499")
            .unwrap();

    // a validator set
    let valset = Valset {
        nonce: 0,
        members: vec![
            ValsetMember {
                eth_address: Some(
                    "0xc783df8a850f42e7F7e57013759C285caa701eB6"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xE5904695748fe4A84b40b3fc79De2277660BD1D3"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
        ],
    };
    let checkpoint = encode_valset_confirm("foo".to_string(), valset);
    let checkpoint_hash = Keccak256::digest(&checkpoint);
    assert_eq!(correct_hash, checkpoint_hash.as_slice());

    // the same valset, except with an intentionally incorrect hash
    let valset = Valset {
        nonce: 1,
        members: vec![
            ValsetMember {
                eth_address: Some(
                    "0xc783df8a850f42e7F7e57013759C285caa701eB6"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
            ValsetMember {
                eth_address: Some(
                    "0xE5904695748fe4A84b40b3fc79De2277660BD1D3"
                        .parse()
                        .unwrap(),
                ),
                power: 3333,
            },
        ],
    };
    let checkpoint = encode_valset_confirm("foo".to_string(), valset);
    let checkpoint_hash = Keccak256::digest(&checkpoint);
    assert_ne!(correct_hash, checkpoint_hash.as_slice())
}

/// takes the required input data and produces the required signature to confirm a transaction
/// batch on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_tx_batch_confirm(gravity_id: String, batch: TransactionBatch) -> Vec<u8> {
    let (amounts, destinations, fees) = batch.get_checkpoint_values();
    encode_tokens(&[
        Token::FixedString(gravity_id),
        Token::FixedString("transactionBatch".to_string()),
        amounts,
        destinations,
        fees,
        batch.nonce.into(),
        batch.token_contract.into(),
        batch.batch_timeout.into(),
    ])
}

pub fn encode_tx_batch_confirm_hashed(gravity_id: String, batch: TransactionBatch) -> Vec<u8> {
    let digest = encode_tx_batch_confirm(gravity_id, batch);
    get_ethereum_msg_hash(&digest)
}

#[test]
fn test_batch_signature() {
    use crate::types::BatchTransaction;
    use crate::types::Erc20Token;
    use clarity::utils::hex_str_to_bytes;
    use clarity::PrivateKey as EthPrivateKey;
    use rand::Rng;
    use sha3::{Digest, Keccak256};

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0xa3a7ee0a363b8ad2514e7ee8f110d7449c0d88f3b0913c28c1751e6e0079a9b2")
            .unwrap();
    let erc20_addr = "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4"
        .parse()
        .unwrap();
    let sender_addr = "althea1c8nkaxk3d0p2gd7ummvmyqpdvqd6pkehqhwnnt"
        .parse()
        .unwrap();

    let token = Erc20Token {
        amount: 1u64.into(),
        token_contract_address: erc20_addr,
    };

    let batch = TransactionBatch {
        batch_timeout: 2111u64,
        nonce: 1u64,
        transactions: vec![BatchTransaction {
            id: 1u64,
            ethereum_recipient: "0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39"
                .parse()
                .unwrap(),
            sender: sender_addr,
            erc20_fee: token.clone(),
            erc20_token: token.clone(),
        }],
        total_fee: token,
        token_contract: erc20_addr,
    };

    let checkpoint = encode_tx_batch_confirm("foo".to_string(), batch.clone());
    let checkpoint_hash = Keccak256::digest(&checkpoint);
    assert_eq!(correct_hash.len(), checkpoint_hash.len());
    assert_eq!(correct_hash, checkpoint_hash.as_slice());

    // checkpoint is correct lets make sure our signature code works
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    let eth_key = EthPrivateKey::from_slice(&secret).unwrap();
    let eth_address = eth_key.to_public_key().unwrap();
    let checkpoint = encode_tx_batch_confirm_hashed("foo".to_string(), batch);

    let eth_signature = eth_key.sign_hash(&checkpoint);

    assert_eq!(eth_address, eth_signature.recover(&checkpoint).unwrap());
}

#[test]
fn test_specific_batch_signature() {
    use crate::types::BatchTransaction;
    use crate::types::Erc20Token;
    use clarity::PrivateKey as EthPrivateKey;
    use rand::Rng;

    let erc20_addr = "0x0635FF793Edf48cf5dB294916720A78e6e490E40"
        .parse()
        .unwrap();
    let sender_addr = "cosmos1g0etv93428tvxqftnmj25jn06mz6dtdasj5nz7"
        .parse()
        .unwrap();

    let token = Erc20Token {
        amount: 1u64.into(),
        token_contract_address: erc20_addr,
    };

    let batch = TransactionBatch {
        batch_timeout: 4427201u64,
        nonce: 15u64,
        transactions: vec![BatchTransaction {
            id: 1301u64,
            ethereum_recipient: "0x64D110e00064F2b428476cD64295d8E35836ffd6"
                .parse()
                .unwrap(),
            sender: sender_addr,
            erc20_fee: token.clone(),
            erc20_token: token.clone(),
        }],
        total_fee: token,
        token_contract: erc20_addr,
    };

    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    // the starting location of the funds
    let eth_key = EthPrivateKey::from_slice(&secret).unwrap();
    let eth_address = eth_key.to_public_key().unwrap();

    let checkpoint = encode_tx_batch_confirm_hashed("foo".to_string(), batch);

    let eth_signature = eth_key.sign_hash(&checkpoint);

    assert_eq!(eth_address, eth_signature.recover(&checkpoint).unwrap());
}

/// takes the required input data and produces the required signature to confirm a logic
/// call on the Gravity Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_logic_call_confirm(gravity_id: String, call: LogicCall) -> Vec<u8> {
    let mut transfer_amounts = Vec::new();
    let mut transfer_token_contracts = Vec::new();
    let mut fee_amounts = Vec::new();
    let mut fee_token_contracts = Vec::new();
    for item in call.transfers.iter() {
        transfer_amounts.push(Token::Uint(item.amount.clone()));
        transfer_token_contracts.push(item.token_contract_address);
    }
    for item in call.fees.iter() {
        fee_amounts.push(Token::Uint(item.amount.clone()));
        fee_token_contracts.push(item.token_contract_address);
    }

    encode_tokens(&[
        Token::FixedString(gravity_id),              // Gravity Instance ID
        Token::FixedString("logicCall".to_string()), //Function Name
        Token::Dynamic(transfer_amounts),            //Array of Transfer amounts
        transfer_token_contracts.into(),             //ERC-20 contract for transfers
        Token::Dynamic(fee_amounts),                 // Array of Fees
        fee_token_contracts.into(),                  // ERC-20 contract for fee payments
        call.logic_contract_address.into(),          // Address of a logic contract
        Token::UnboundedBytes(call.payload),         // Encoded arguments to logic contract
        call.timeout.into(),                         // Timeout on batch
        Token::Bytes(call.invalidation_id),          // Scope of logic batch
        call.invalidation_nonce.into(),              // Nonce of logic batch. See 2-d nonce scheme.
    ])
}

pub fn encode_logic_call_confirm_hashed(gravity_id: String, call: LogicCall) -> Vec<u8> {
    let digest = encode_logic_call_confirm(gravity_id, call);
    get_ethereum_msg_hash(&digest)
}

#[test]
fn test_logic_call_signature() {
    use crate::types::Erc20Token;
    use crate::types::LogicCall;
    use clarity::utils::hex_str_to_bytes;
    use sha3::{Digest, Keccak256};

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0x1de95c9ace999f8ec70c6dc8d045942da2612950567c4861aca959c0650194da")
            .unwrap();
    let token_contract_address = "0xC26eFfa98B8A2632141562Ae7E34953Cfe5B4888"
        .parse()
        .unwrap();
    let logic_contract_address = "0x17c1736CcF692F653c433d7aa2aB45148C016F68"
        .parse()
        .unwrap();
    let token = vec![Erc20Token {
        amount: 1u8.into(),
        token_contract_address,
    }];

    let logic_call = LogicCall {
        transfers: token.clone(),
        fees: token,
        logic_contract_address,
        payload: hex_str_to_bytes(
            "0x74657374696e675061796c6f6164000000000000000000000000000000000000",
        )
        .unwrap(),
        timeout: 4766922941000,
        invalidation_id: hex_str_to_bytes(
            "0x696e76616c69646174696f6e4964000000000000000000000000000000000000",
        )
        .unwrap(),
        invalidation_nonce: 1u8.into(),
    };
    let checkpoint = encode_logic_call_confirm("foo".to_string(), logic_call);
    println!("{}", checkpoint.len() / 32);

    let checkpoint_hash = Keccak256::digest(&checkpoint);

    assert_eq!(correct_hash.len(), checkpoint_hash.len());
    assert_eq!(correct_hash, checkpoint_hash.as_slice())
}
