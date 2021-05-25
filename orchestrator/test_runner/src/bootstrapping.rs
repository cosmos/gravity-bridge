use crate::{utils::ValidatorKeys};
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use std::{fs::File};
use std::{
    io::{BufRead, BufReader, Read},
};

/// Ethereum private keys for the validators are generated using the gravity eth_keys add command
/// and dumped into a file /validator-eth-keys in the container, from there they are then used by
/// the orchestrator on startup
pub fn parse_ethereum_keys() -> Vec<EthPrivateKey> {
    // TODO - grab this from the json files in /testdata instead
    let filename = "/validator-eth-keys";
    let file = File::open(filename).expect("Failed to find eth keys");
    let reader = BufReader::new(file);
    let mut ret = Vec::new();

    for line in reader.lines() {
        let key = line.expect("Error reading eth key file!");
        if key.is_empty() || key.contains("public") || key.contains("address") {
            continue;
        }
        let key = key.split(':').last().unwrap().trim();
        ret.push(key.parse().unwrap());
    }
    ret
}

/// Validator private keys are generated via the gravity key add
/// command, from there they are used to create gentx's and start the
/// chain, these keys change every time the container is restarted.
/// The mnemonic phrases are dumped into a text file /validator-phrases
/// the phrases are in increasing order, so validator 1 is the first key
/// and so on. While validators may later fail to start it is guaranteed
/// that we have one key for each validator in this file.
pub fn parse_validator_keys() -> Vec<CosmosPrivateKey> {
    // TODO - grab this from the json files in /testdata instead
    // probably: find /testdata/testchain -name validator_key.json
    let filename = "/validator-phrases";
    parse_phrases(filename)
}

/// Orchestrator private keys are generated via the gravity key add
/// command just like the validator keys themselves and stored in a
/// similar file /orchestrator-phrases
pub fn parse_orchestrator_keys() -> Vec<CosmosPrivateKey> {
    // TODO - grab this from the json files in /testdata instead
    // probably: find /testdata/testchain -name orchestrator_key.json
    let filename = "/orchestrator-phrases";
    parse_phrases(filename)
}

/// Parses the output of the cosmoscli keys add command to import the private key
fn parse_phrases(filename: &str) -> Vec<CosmosPrivateKey> {
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

pub fn get_keys() -> Vec<ValidatorKeys> {
    let cosmos_keys = parse_validator_keys();
    let orch_keys = parse_orchestrator_keys();
    let eth_keys = parse_ethereum_keys();
    let mut ret = Vec::new();
    for ((c_key, o_key), e_key) in cosmos_keys.into_iter().zip(orch_keys).zip(eth_keys) {
        ret.push(ValidatorKeys {
            eth_key: e_key,
            validator_key: c_key,
            orch_key: o_key,
        })
    }
    ret
}

pub struct BootstrapContractAddresses {
    pub gravity_contract: EthAddress,
    pub erc20_addresses: Vec<EthAddress>,
    pub uniswap_liquidity_address: Option<EthAddress>,
}

/// Parses the ERC20 and Gravity contract addresses from the file created
/// in deploy_contracts()
pub fn parse_contract_addresses() -> BootstrapContractAddresses {
    let mut file =
        File::open("/contracts").expect("Failed to find contracts! did they not deploy?");
    let mut output = String::new();
    file.read_to_string(&mut output).unwrap();
    let mut maybe_gravity_address = None;
    let mut erc20_addresses = Vec::new();
    let mut uniswap_liquidity = None;
    for line in output.lines() {
        if line.contains("Gravity deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            maybe_gravity_address = Some(address_string.trim().parse().unwrap());
        } else if line.contains("ERC20 deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            erc20_addresses.push(address_string.trim().parse().unwrap());
        } else if line.contains("Uniswap Liquidity test deployed at Address - ") {
            let address_string = line.split('-').last().unwrap();
            uniswap_liquidity = Some(address_string.trim().parse().unwrap());
        }
    }
    let gravity_contract: EthAddress = maybe_gravity_address.unwrap();
    BootstrapContractAddresses {
        gravity_contract,
        erc20_addresses,
        uniswap_liquidity_address: uniswap_liquidity,
    }
}

// fn all_paths_exist(input: &[&str]) -> bool {
//     for i in input {
//         if !Path::new(i).exists() {
//             return false;
//         }
//     }
//     true
// }

// fn return_existing<'a>(a: [&'a str; 2], b: [&'a str; 2]) -> [&'a str; 2] {
//     if all_paths_exist(&a) {
//         a
//     } else if all_paths_exist(&b) {
//         b
//     } else {
//         panic!("No paths exist!")
//     }
// }
