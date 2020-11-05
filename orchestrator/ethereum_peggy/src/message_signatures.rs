use clarity::abi::{encode_tokens, Token};
use peggy_utils::types::{TransactionBatch, Valset};
use sha3::{Digest, Keccak256};

/// takes the required input data and produces the required signature to confirm a validator
/// set update on the Peggy Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_valset_confirm(peggy_id: String, valset: Valset) -> Vec<u8> {
    let (eth_addresses, powers) = valset.filter_empty_addresses();
    encode_tokens(&[
        Token::FixedString(peggy_id),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        eth_addresses.into(),
        powers.into(),
    ])
}

#[test]
fn test_valset_signature() {
    use clarity::utils::hex_str_to_bytes;
    use peggy_utils::types::ValsetMember;

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
/// batch on the Peggy Ethereum contract. This value will then be signed before being
/// submitted to Cosmos, verified, and then relayed to Ethereum
/// Note: This is the message, you need to run Keccak256::digest() in order to get the 32byte
/// digest that is normally signed or may be used as a 'hash of the message'
pub fn encode_tx_batch_confirm(peggy_id: String, batch: TransactionBatch) -> Vec<u8> {
    // transaction batches include a validator set update, the way this is verified is that the valset checkpoint
    // (encoded ethereum data) is included within the batch signature, which is itself a checkpoint over the batch data
    let valset_checkpoint = encode_valset_confirm(peggy_id.clone(), batch.valset.clone());
    let valset_digest = Keccak256::digest(&valset_checkpoint).as_slice().to_vec();
    let (amounts, destinations, fees) = batch.get_checkpoint_values();
    encode_tokens(&[
        Token::FixedString(peggy_id),
        Token::FixedString("valsetAndTransactionBatch".to_string()),
        Token::Bytes(valset_digest),
        amounts,
        destinations,
        fees,
        batch.nonce.into(),
        batch.token_contract.into(),
    ])
}

#[test]
fn test_batch_signature() {
    use clarity::utils::hex_str_to_bytes;
    use peggy_utils::types::BatchTransaction;
    use peggy_utils::types::ERC20Denominator;
    use peggy_utils::types::ERC20Token;
    use peggy_utils::types::ValsetMember;

    let correct_hash: Vec<u8> =
        hex_str_to_bytes("0x746471abc2232c11039c2160365c4593110dbfbe25ff9a2dcf8b5b7376e9f346")
            .unwrap();
    let erc20_addr = "0x22474D350EC2dA53D717E30b96e9a2B7628Ede5b"
        .parse()
        .unwrap();
    let sender_addr = "0x527FBEE652609AB150F0AEE9D61A2F76CFC4A73E"
        .parse()
        .unwrap();

    let valset = Valset {
        nonce: 1,
        members: vec![ValsetMember {
            eth_address: Some(
                "0xc783df8a850f42e7F7e57013759C285caa701eB6"
                    .parse()
                    .unwrap(),
            ),
            power: 6670,
        }],
    };

    let token = ERC20Token {
        amount: 1u64.into(),
        symbol: "MAX".to_string(),
        token_contract_address: erc20_addr,
    };

    let batch = TransactionBatch {
        nonce: 1u64.into(),
        elements: vec![BatchTransaction {
            txid: 1u64.into(),
            destination: "0x9FC9C2DfBA3b6cF204C37a5F690619772b926e39"
                .parse()
                .unwrap(),
            sender: sender_addr,
            bridge_fee: token.clone(),
            send: token.clone(),
        }],
        created_at: "".to_string(),
        total_fee: token,
        bridged_denominator: ERC20Denominator {
            cosmos_voucher_denom: "peggy39b512461b".to_string(),
            symbol: "MAX".to_string(),
            token_contract_address: erc20_addr,
        },
        batch_status: 1,
        token_contract: erc20_addr,
        valset,
    };

    let checkpoint = encode_tx_batch_confirm("foo".to_string(), batch);
    let checkpoint_hash = Keccak256::digest(&checkpoint);
    assert_eq!(correct_hash.len(), checkpoint_hash.len());
    assert_eq!(correct_hash, checkpoint_hash.as_slice())
}
