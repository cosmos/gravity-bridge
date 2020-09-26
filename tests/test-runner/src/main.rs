//! Test runner is a testing script for the Peggy Cosmos module. It is built in Rust rather than python or bash
//! to maximize code and tooling shared with the validator-daemon and relayer binaries.
use clarity::{
    abi::encode_call, abi::encode_tokens, utils::bytes_to_hex_str, Address as EthAddress,
    Signature, Uint256,
};
use clarity::{abi::Token, PrivateKey as EthPrivateKey};
use contact::{client::test_rpc_calls, types::Valset};
use contact::{client::Contact, types::ValsetConfirmResponse};
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use sha3::{Digest, Keccak256};
use std::io::{BufRead, BufReader};
use std::process::Command;
use std::time::Duration;
use std::{fs::File, thread};
use tokio::time::timeout as future_timeout;
use web30::{client, jsonrpc::error::Web3Error, types::SendTxOption, types::TransactionRequest};
use web30::{client::Web3, types::Data, types::UnpaddedHex};

#[macro_use]
extern crate log;

const TIMEOUT: Duration = Duration::from_secs(30);

/// Ethereum keys are generated for every validator inside
/// of this testing application and submitted to the blockchain
/// use the 'update eth address' message.
fn generate_eth_private_key() -> EthPrivateKey {
    let key_buf: [u8; 32] = rand::random();
    EthPrivateKey::from_slice(&key_buf).unwrap()
}

/// Validator private keys are generated via the cosmoscli key add
/// command, from there they are used to create gentx's and start the
/// chain, these keys change every time the container is restarted.
/// The mnemonic phrases are dumped into a text file /validator-phrases
/// the phrases are in increasing order, so validator 1 is the first key
/// and so on. While validators may later fail to start it is guaranteed
/// that we have one key for each validator in this file.
fn parse_validator_keys() -> Vec<CosmosPrivateKey> {
    let filename = "/validator-phrases";
    let file = File::open(filename).expect("Failed to find phrases");
    let reader = BufReader::new(file);
    let mut ret = Vec::new();

    for line in reader.lines() {
        let phrase = line.expect("Error reading phrase file!");
        if phrase.is_empty()
            || phrase.contains("write this mnemonic phrase")
            || phrase.contains("recover your account if")
        {
            continue;
        }
        let key = CosmosPrivateKey::from_phrase(&phrase, "").expect("Bad phrase!");
        ret.push(key);
    }
    ret
}

fn get_keys() -> Vec<(CosmosPrivateKey, EthPrivateKey)> {
    let cosmos_keys = parse_validator_keys();
    let mut ret = Vec::new();
    for c_key in cosmos_keys {
        ret.push((c_key, generate_eth_private_key()))
    }
    ret
}

#[actix_rt::main]
async fn main() {
    println!("Staring Peggy test-runner");
    const COSMOS_NODE: &str = "http://localhost:1317";
    const ETH_NODE: &str = "http://localhost:8545";
    const PEGGY_ID: &str = "foo";
    // this key is the private key for the public key defined in tests/assets/ETHGenesis.json
    // where the full node / miner sends its rewards. Therefore it's always going
    // to have a lot of ETH to pay for things like contract deployments
    let miner_private_key: EthPrivateKey =
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
            .parse()
            .unwrap();

    let contact = Contact::new(COSMOS_NODE, TIMEOUT);
    let web30 = web30::client::Web3::new(ETH_NODE, TIMEOUT);
    let keys = get_keys();
    let test_token_name = "footoken".to_string();

    let fee = Coin {
        denom: test_token_name.clone(),
        amount: 1u32.into(),
    };

    for (c_key, e_key) in keys.iter() {
        // set the eth address for all the validators
        contact
            .update_peggy_eth_address(*e_key, *c_key, fee.clone(), None, None, None)
            .await
            .expect("Failed to update eth address!");
    }

    wait_for_next_cosmos_block(&contact).await;

    // get the first validator and have them send a valset request
    let (c_key, _e_key) = keys[0];
    let request_block = contact
        .send_valset_request(c_key, fee.clone(), None, None, None)
        .await
        .expect("Failed to send valset request!")
        .height;

    wait_for_next_cosmos_block(&contact).await;

    for (c_key, e_key) in keys.iter() {
        let valset = contact
            .get_oldest_unsigned_valset(c_key.to_public_key().unwrap().to_address())
            .await
            .expect("Failed to get valset!");
        // send in valset confirmation for all validators
        let res = contact
            .send_valset_confirm(
                *e_key,
                fee.clone(),
                valset.result,
                *c_key,
                PEGGY_ID.to_string(),
                None,
                None,
                None,
            )
            .await;
        res.expect("Failed to send valset confirm!");
    }

    wait_for_next_cosmos_block(&contact).await;

    let output = Command::new("npx")
        .args(&[
            "ts-node",
            "/peggy/solidity/contract-deployer.ts",
            &format!("--cosmos-node={}", COSMOS_NODE),
            &format!("--eth-node={}", ETH_NODE),
            &format!("--eth-privkey={:#x}", miner_private_key),
            &format!("--peggy-id={}", PEGGY_ID),
            "--contract=/peggy/solidity/artifacts/Peggy.json",
            "--erc20-contract=/peggy/solidity/artifacts/TestERC20.json",
            "--test-mode=true",
        ])
        .current_dir("/peggy/solidity/")
        .output()
        .expect("Failed to deploy contracts!");
    println!("status: {}", output.status);
    println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    println!("stderr: {}", String::from_utf8_lossy(&output.stderr));

    // TODO: here we need to bootstrap the chain with the new deployed contract
    // in the meantime we parse stdout
    let mut maybe_peggy_address = None;
    for line in String::from_utf8_lossy(&output.stdout).lines() {
        if line.contains("Peggy deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            maybe_peggy_address = Some(address_string.trim().parse().unwrap());
            break;
        }
    }
    let peggy_address = maybe_peggy_address.unwrap();

    let latest_valsets = contact
        .get_last_valset_requests()
        .await
        .expect("Failed to get latest valsets");
    // this will panic if there are no valsets in the response, but there must be
    // the one we have just signed and submitted above
    let latest = latest_valsets.result[0].clone();
    let confirms = contact
        .get_all_valset_confirms(latest.nonce)
        .await
        .expect("Failed to get valset confirms");

    for confirm in confirms.result.iter() {
        let eth_private_key = get_eth_key_from_cosmos_addr(confirm.validator, &keys);
        verify_contract_interactions(
            &web30,
            peggy_address,
            eth_private_key,
            confirm.eth_signature.clone(),
            latest.clone(),
            PEGGY_ID.to_string(),
            miner_private_key.to_public_key().unwrap(),
        )
        .await;
    }
    verify_signature_passing(
        &web30,
        peggy_address,
        PEGGY_ID.to_string(),
        &confirms.result,
        latest.clone(),
        miner_private_key.to_public_key().unwrap(),
        &keys,
    )
    .await;

    send_eth_valset_update(
        latest,
        &confirms.result,
        web30,
        TIMEOUT,
        peggy_address,
        miner_private_key,
        &keys,
    )
    .await;
}

async fn send_basic_eth_transaction(
    sending_eth_private_key: EthPrivateKey,
    dest: EthAddress,
    web3: Web3,
) {
    let eth_address = sending_eth_private_key.to_public_key().unwrap();
    println!(
        "Our balance is {:?}",
        web3.eth_get_balance(eth_address).await
    );
    println!("Our gas price is {:?}", web3.eth_gas_price().await);
    println!("Our chain id is {:?}", web3.net_version().await);
    let tx = web3
        .send_transaction(
            dest,
            Vec::new(),
            10000000u32.into(),
            eth_address,
            sending_eth_private_key,
            vec![SendTxOption::GasLimit(27000u32.into())],
        )
        .await;
    println!("Our tx result is {:?}", tx);
    panic!("exiting!");
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

fn get_correct_power_for_address(address: EthAddress, valset: &Valset) -> (EthAddress, u64) {
    for (a, p) in valset.eth_addresses.iter().zip(valset.powers.iter()) {
        if let Some(a) = a {
            if *a == address {
                return (*a, *p);
            }
        }
    }
    panic!("Could not find that address!");
}

fn get_correct_sig_for_address(
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

/// this function generates an appropriate Ethereum transaction
/// to submit the provided validator set and signatures.
/// TODO this function uses the same validator set as the old and
/// new validator set, this is because there's no actual changes to
/// the set in testing and because there's no oracle to tell us what
/// the old set was anyways.
/// TODO TODO should we have an oracle for the old set or look in the chain?
async fn send_eth_valset_update(
    valset: Valset,
    confirms: &[ValsetConfirmResponse],
    web3: Web3,
    timeout: Duration,
    peggy_contract_address: EthAddress,
    sending_eth_private_key: EthPrivateKey,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
) {
    let old_addresses = filter_empty_addresses(&valset.eth_addresses);
    let old_powers = valset.powers;
    let new_addresses = old_addresses.clone();
    let new_powers = old_powers.clone();
    let old_nonce = 0u64;
    let new_nonce = valset.nonce;
    let mut v = Vec::new();
    let mut r = Vec::new();
    let mut s = Vec::new();
    for address in old_addresses.iter() {
        let cosmos_address = get_cosmos_address_from_eth_addr(*address, &keys);
        let (sig_v, sig_r, sig_s) = get_correct_sig_for_address(cosmos_address, confirms);
        v.push(sig_v.clone());
        r.push(Token::Bytes(sig_r.clone().to_bytes_be()));
        s.push(Token::Bytes(sig_s.clone().to_bytes_be()));
    }
    let eth_address = sending_eth_private_key.to_public_key().unwrap();

    // Solidity function signature
    // function getValsetNonce() public returns (uint256)
    let first_nonce = contract_call(
        &web3,
        peggy_contract_address,
        "getValsetNonce()",
        &[],
        eth_address,
    )
    .await
    .expect("Failed to get the first nonce");
    println!("First nonce {:?}", first_nonce);

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
            sending_eth_private_key,
            vec![SendTxOption::GasLimit(1_000_000u32.into())],
        ),
    )
    .await
    .expect("Valset update timed out")
    .expect("Valset update failed for other reasons");
    println!("Finished valset update with txid {:#066x}", tx);

    let mut not_in_chain = true;
    while not_in_chain {
        let res = web3.eth_get_transaction_by_hash(tx.clone()).await.unwrap();
        if let Some(val) = res {
            if let Some(_block) = val.block_number {
                not_in_chain = false;
            }
        }
    }

    // Solidity function signature
    // function getValsetNonce() public returns (uint256)
    let last_nonce = contract_call(
        &web3,
        peggy_contract_address,
        "getValsetNonce()",
        &[],
        eth_address,
    )
    .await
    .expect("Failed to get the last nonce");
    println!("Last nonce {:?}", last_nonce);
}

/// Takes an array of Option<EthAddress> and converts to EthAddress erroring when
/// an None is found, in a prod environment you would replace with zeros if a sig
/// or address is missing, this is test so we want to exit with an error
pub fn filter_empty_addresses(input: &[Option<EthAddress>]) -> Vec<EthAddress> {
    let mut output = Vec::new();
    for val in input.iter() {
        match val {
            Some(a) => output.push(*a),
            None => panic!("This should be impossible!"),
        }
    }
    output
}

pub fn to_uint_vec(input: &[u64]) -> Vec<Uint256> {
    let mut new_vec = Vec::new();
    for value in input {
        let v: u64 = *value;
        new_vec.push(v.into())
    }
    new_vec
}

/// TODO modify code in web30 if this works
pub async fn contract_call(
    web30: &Web3,
    contract_address: EthAddress,
    sig: &str,
    tokens: &[Token],
    own_address: EthAddress,
) -> Result<Vec<u8>, Web3Error> {
    let gas_price = match web30.eth_gas_price().await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };

    let nonce = match web30.eth_get_transaction_count(own_address).await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };

    let payload = encode_call(sig, tokens).unwrap();

    let transaction = TransactionRequest {
        from: Some(own_address),
        to: contract_address,
        nonce: Some(UnpaddedHex(nonce)),
        gas: Some(UnpaddedHex(1_000_000u64.into())),
        gas_price: Some(UnpaddedHex(gas_price)),
        value: Some(UnpaddedHex(0u64.into())),
        data: Some(Data(payload)),
    };

    let bytes = match web30.eth_call(transaction).await {
        Ok(val) => val,
        Err(e) => return Err(e),
    };
    Ok(bytes.0)
}

async fn wait_for_next_cosmos_block(contact: &Contact) {
    let current_block = contact
        .get_latest_block()
        .await
        .unwrap()
        .block
        .last_commit
        .height;
    while current_block
        == contact
            .get_latest_block()
            .await
            .unwrap()
            .block
            .last_commit
            .height
    {
        thread::sleep(Duration::from_secs(1))
    }
}

fn get_checkpoint_abi_encode(valset: &Valset, peggy_id: &str) -> Vec<u8> {
    encode_tokens(&[
        Token::FixedString(peggy_id.to_string()),
        Token::FixedString("checkpoint".to_string()),
        valset.nonce.into(),
        filter_empty_addresses(&valset.eth_addresses).into(),
        valset.powers.clone().into(),
    ])
}

fn get_checkpoint_hash(valset: &Valset, peggy_id: &str) -> Vec<u8> {
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id);
    let locally_computed_digest = Keccak256::digest(&locally_computed_abi_encode);
    locally_computed_digest.to_vec()
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
            filter_empty_addresses(&valset.eth_addresses).into(),
            valset.powers.clone().into(),
            valset.nonce.into(),
            Token::FixedString(peggy_id.clone()),
        ],
        eth_address_with_funds,
    )
    .await
    .expect("Failed to get checkpoint hash from contract");
    let locally_computed_abi_encode = get_checkpoint_abi_encode(&valset, &peggy_id);
    let locally_computed_digest = get_checkpoint_hash(&valset, &peggy_id);
    assert_eq!(
        locally_computed_digest.to_vec(),
        contract_computed_checkpoint
    );
    println!(
        "Correct hash is {}",
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
    println!(
        "Address: {} v: {:x} r: {:x} s: {:x} got: {}",
        eth_address,
        cosmos_submitted_signature.v,
        cosmos_submitted_signature.r,
        cosmos_submitted_signature.s,
        bytes_to_hex_str(&contract_output)
    );
    // signature verification passed
    assert!(*contract_output.iter().last().unwrap() == 1u8);
}

struct OrderedSignatures {
    addresses: Vec<EthAddress>,
    powers: Vec<u64>,
    v: Vec<Uint256>,
    r: Vec<Uint256>,
    s: Vec<Uint256>,
}

/// takes a valset and signatures as input and returns a sorted
/// array of each in order of descending validator power with the
/// appropriate signatures lined up correctly
fn prepare_sigs() {}

async fn verify_signature_passing(
    web3: &Web3,
    peggy_contract_address: EthAddress,
    peggy_id: String,
    confirms: &[ValsetConfirmResponse],
    valset: Valset,
    eth_address_with_funds: EthAddress,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
) {
    let locally_computed_checkpoint_hash = get_checkpoint_hash(&valset, &peggy_id);

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
}
