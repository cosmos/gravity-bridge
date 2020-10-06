use clarity::abi::Token;
use clarity::private_key::PrivateKey as EthPrivateKey;
use clarity::Address as EthAddress;
use clarity::Signature;
use cosmos_peggy::types::*;
use deep_space::address::Address as CosmosAddress;
use deep_space::{private_key::PrivateKey as CosmosPrivateKey, utils::bytes_to_hex_str};
use ethereum_peggy::utils::filter_empty_eth_addresses;
use ethereum_peggy::utils::{get_checkpoint_hash, get_correct_power_for_address};
use ethereum_peggy::{contract_call, utils::get_checkpoint_abi_encode};
use web30::client::Web3;

async fn verify_signature_passing(
    web3: &Web3,
    peggy_contract_address: EthAddress,
    peggy_id: String,
    confirms: &[ValsetConfirmResponse],
    valset: Valset,
    eth_address_with_funds: EthAddress,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
) {
    info!("Checking validity of all validator signatures");
    let locally_computed_checkpoint_hash = get_checkpoint_hash(&valset, &peggy_id).unwrap();

    let mut addresses = Vec::new();
    let mut powers = Vec::new();
    let mut v = Vec::new();
    let mut r = Vec::new();
    let mut s = Vec::new();
    for sig in confirms {
        // we can get signatures and powers back from cosmos in any order, we must now order them properly
        let eth_private_key = get_eth_key_from_cosmos_addr(sig.validator, &keys);
        let eth_address = eth_private_key.to_public_key().unwrap();
        let (address, power) = get_correct_power_for_address(eth_address, &valset);
        addresses.push(address);
        powers.push(power);

        //v.push(format!("{}", sig.eth_signature.v.clone()).parse().unwrap());
        v.push(sig.eth_signature.v.clone());
        r.push(Token::Bytes(sig.eth_signature.r.clone().to_bytes_be()));
        s.push(Token::Bytes(sig.eth_signature.s.clone().to_bytes_be()));
    }

    let contract_output = contract_call(
        &web3,
        peggy_contract_address,
        "checkValidatorSignatures(address[],uint256[],uint8[],bytes32[],bytes32[],bytes32,uint256)",
        &[
            addresses.into(),
            powers.into(),
            v.into(),
            Token::Dynamic(r),
            Token::Dynamic(s),
            Token::Bytes(locally_computed_checkpoint_hash),
            500u64.into(),
        ],
        eth_address_with_funds,
    )
    .await
    .expect("Failed to get sig verifications from contract");
    // signature verification passed
    assert!(*contract_output.iter().last().unwrap() == 1u8);
    info!("Passed!");
}

async fn verify_contract_interactions(
    web3: &Web3,
    peggy_contract_address: EthAddress,
    eth_private_key: EthPrivateKey,
    cosmos_submitted_signature: Signature,
    valset: Valset,
    peggy_id: String,
    eth_address_with_funds: EthAddress,
) {
    let eth_address = eth_private_key.to_public_key().unwrap();
    // Solidity function signature
    // function getValsetNonce() public returns (uint256)
    let contract_computed_checkpoint = contract_call(
        &web3,
        peggy_contract_address,
        "makeCheckpoint(address[],uint256[],uint256,bytes32)",
        &[
            filter_empty_eth_addresses(&valset.eth_addresses)
                .unwrap()
                .into(),
            valset.powers.clone().into(),
            valset.nonce.into(),
            Token::FixedString(peggy_id.clone()),
        ],
        eth_address_with_funds,
    )
    .await
    .expect("Failed to get checkpoint hash from contract");
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id).unwrap();
    let locally_computed_digest = get_checkpoint_hash(&valset, &peggy_id).unwrap();
    assert_eq!(
        locally_computed_digest.to_vec(),
        contract_computed_checkpoint
    );
    trace!(
        "Correct checkpoint hash is {}",
        bytes_to_hex_str(&contract_computed_checkpoint)
    );

    let eth_signature = eth_private_key.sign_ethereum_msg(&locally_computed_abi_encode);
    assert_eq!(eth_signature, cosmos_submitted_signature);

    let contract_output = contract_call(
        &web3,
        peggy_contract_address,
        "verifySig(address,bytes32,uint8,bytes32,bytes32)",
        &[
            eth_address.into(),
            Token::Bytes(contract_computed_checkpoint.clone()),
            cosmos_submitted_signature.v.clone().into(),
            Token::Bytes(cosmos_submitted_signature.r.to_bytes_be()),
            Token::Bytes(cosmos_submitted_signature.s.to_bytes_be()),
        ],
        eth_address_with_funds,
    )
    .await
    .expect("Failed to get sig verification from contract");
    info!("Checking validity of individual validator signatures");
    info!(
        "Address: {} v: {:x} r: {:x} s: {:x} got: {}",
        eth_address,
        cosmos_submitted_signature.v,
        cosmos_submitted_signature.r,
        cosmos_submitted_signature.s,
        bytes_to_hex_str(&contract_output)
    );
    // signature verification passed
    assert!(*contract_output.iter().last().unwrap() == 1u8);
    info!("Passed!");
}

/// looks through our list of addresses to map a cosmos address to an
/// eth address
fn get_eth_key_from_cosmos_addr(
    search_key: CosmosAddress,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
) -> EthPrivateKey {
    for (c_key, e_key) in keys {
        if c_key.to_public_key().unwrap().to_address() == search_key {
            return *e_key;
        }
    }
    panic!("Did not find address!");
}

fn get_cosmos_address_from_eth_addr(
    search_key: EthAddress,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
) -> CosmosAddress {
    for (c_key, e_key) in keys {
        if e_key.to_public_key().unwrap() == search_key {
            return c_key.to_public_key().unwrap().to_address();
        }
    }
    panic!("Did not find address!");
}
