use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use contact::client::Contact;
use cosmos_peggy::send::update_peggy_delegate_addresses;
use cosmos_peggy::utils::wait_for_next_cosmos_block;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use futures::future::join_all;
use std::process::Command;
use std::{fs::File, path::Path};
use std::{
    io::{BufRead, BufReader, Read, Write},
    process::ExitStatus,
};

use crate::COSMOS_NODE_ABCI;
use crate::ETH_NODE;
use crate::MINER_PRIVATE_KEY;
use crate::TOTAL_TIMEOUT;

/// Ethereum keys are generated for every validator inside
/// of this testing application and submitted to the blockchain
/// use the 'update eth address' message. In this case we generate
/// them based off of the Cosmos key as the seed so that we can run
/// the test runner multiple times against one chain and get the same keys.
///
/// There's no particular reason to use the public key except that the bytes
/// of the private key type are not public
pub fn generate_eth_private_key(seed: CosmosPrivateKey) -> EthPrivateKey {
    EthPrivateKey::from_slice(&seed.to_public_key().unwrap().as_bytes()[0..32]).unwrap()
}

/// Validator private keys are generated via the cosmoscli key add
/// command, from there they are used to create gentx's and start the
/// chain, these keys change every time the container is restarted.
/// The mnemonic phrases are dumped into a text file /validator-phrases
/// the phrases are in increasing order, so validator 1 is the first key
/// and so on. While validators may later fail to start it is guaranteed
/// that we have one key for each validator in this file.
pub fn parse_validator_keys() -> Vec<CosmosPrivateKey> {
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

pub fn get_keys() -> Vec<(CosmosPrivateKey, EthPrivateKey)> {
    let cosmos_keys = parse_validator_keys();
    let mut ret = Vec::new();
    for c_key in cosmos_keys {
        ret.push((c_key, generate_eth_private_key(c_key)))
    }
    ret
}

/// This function deploys the required contracts onto the Ethereum testnet
/// this runs only when the DEPLOY_CONTRACTS env var is set right after
/// the Ethereum test chain starts in the testing environment. We write
/// the stdout of this to a file for later test runs to parse
pub async fn deploy_contracts(
    contact: &Contact,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    fee: Coin,
) {
    // register all validator eth addresses, currently validators can just not do this
    // a full production version of Peggy would refuse to allow validators to enter the pool
    // without registering their address. It would also allow them to delegate their Cosmos addr
    //
    // Either way, validators need to setup their eth addresses out of band and it's not
    // the orchestrators job. So this isn't exactly where it needs to be in the final version
    // but neither is it really that different.
    let mut updates = Vec::new();
    for (c_key, e_key) in keys.iter() {
        info!(
            "Signing and submitting Delegate addresses {} for validator {}",
            e_key.to_public_key().unwrap(),
            c_key.to_public_key().unwrap().to_address(),
        );
        updates.push(update_peggy_delegate_addresses(
            &contact,
            e_key.to_public_key().unwrap(),
            c_key.to_public_key().unwrap().to_address(),
            *c_key,
            fee.clone(),
        ));
    }
    let update_results = join_all(updates).await;
    for i in update_results {
        i.expect("Failed to set delegate addresses!");
    }

    // prevents the node deployer from failing (rarely) when the chain has not
    // yet produced the next block after submitting each eth address
    wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;

    // these are the possible paths where we could find the contract deployer
    // and the peggy contract itself, feel free to expand this if it makes your
    // deployments more straightforward.

    // both files are just in the PWD
    const A: [&str; 2] = ["contract-deployer", "Peggy.json"];
    // files are placed in a root /solidity/ folder
    const B: [&str; 2] = ["/solidity/contract-deployer", "/solidity/Peggy.json"];
    // the default unmoved locations for the Gravity repo
    const C: [&str; 3] = [
        "/peggy/solidity/contract-deployer.ts",
        "/peggy/solidity/artifacts/contracts/Peggy.sol/Peggy.json",
        "/peggy/solidity/",
    ];
    let output = if all_paths_exist(&A) || all_paths_exist(&B) {
        let paths = return_existing(A, B);
        Command::new(paths[0])
            .args(&[
                &format!("--cosmos-node={}", COSMOS_NODE_ABCI),
                &format!("--eth-node={}", ETH_NODE),
                &format!("--eth-privkey={:#x}", *MINER_PRIVATE_KEY),
                &format!("--contract={}", paths[1]),
                "--test-mode=true",
            ])
            .output()
            .expect("Failed to deploy contracts!")
    } else if all_paths_exist(&C) {
        Command::new("npx")
            .args(&[
                "ts-node",
                C[0],
                &format!("--cosmos-node={}", COSMOS_NODE_ABCI),
                &format!("--eth-node={}", ETH_NODE),
                &format!("--eth-privkey={:#x}", *MINER_PRIVATE_KEY),
                &format!("--contract={}", C[1]),
                "--test-mode=true",
            ])
            .current_dir(C[2])
            .output()
            .expect("Failed to deploy contracts!")
    } else {
        panic!("Could not find Peggy.json contract artifact in any known location!")
    };

    info!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    info!("stderr: {}", String::from_utf8_lossy(&output.stderr));
    if !ExitStatus::success(&output.status) {
        panic!("Contract deploy failed!")
    }
    let mut file = File::create("/contracts").unwrap();
    file.write_all(&output.stdout).unwrap();
}

pub struct BootstrapContractAddresses {
    pub peggy_contract: EthAddress,
    pub erc20_addresses: Vec<EthAddress>,
    pub uniswap_liquidity_address: Option<EthAddress>,
}

/// Parses the ERC20 and Peggy contract addresses from the file created
/// in deploy_contracts()
pub fn parse_contract_addresses() -> BootstrapContractAddresses {
    let mut file =
        File::open("/contracts").expect("Failed to find contracts! did they not deploy?");
    let mut output = String::new();
    file.read_to_string(&mut output).unwrap();
    let mut maybe_peggy_address = None;
    let mut erc20_addresses = Vec::new();
    let mut uniswap_liquidity = None;
    for line in output.lines() {
        if line.contains("Peggy deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            maybe_peggy_address = Some(address_string.trim().parse().unwrap());
        } else if line.contains("ERC20 deployed at Address -") {
            let address_string = line.split('-').last().unwrap();
            erc20_addresses.push(address_string.trim().parse().unwrap());
        } else if line.contains("Uniswap Liquidity test deployed at Address - ") {
            let address_string = line.split('-').last().unwrap();
            uniswap_liquidity = Some(address_string.trim().parse().unwrap());
        }
    }
    let peggy_address: EthAddress = maybe_peggy_address.unwrap();
    BootstrapContractAddresses {
        peggy_contract: peggy_address,
        erc20_addresses,
        uniswap_liquidity_address: uniswap_liquidity,
    }
}

fn all_paths_exist(input: &[&str]) -> bool {
    for i in input {
        if !Path::new(i).exists() {
            return false;
        }
    }
    true
}

fn return_existing<'a>(a: [&'a str; 2], b: [&'a str; 2]) -> [&'a str; 2] {
    if all_paths_exist(&a) {
        a
    } else if all_paths_exist(&b) {
        b
    } else {
        panic!("No paths exist!")
    }
}
