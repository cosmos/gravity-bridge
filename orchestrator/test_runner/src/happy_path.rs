use crate::get_fee;
use crate::get_test_token_name;
use crate::{
    utils::*, COSMOS_NODE_GRPC, MINER_ADDRESS, MINER_PRIVATE_KEY, STARTING_STAKE_PER_VALIDATOR,
    TOTAL_TIMEOUT,
};
use actix::Arbiter;
use clarity::PrivateKey as EthPrivateKey;
use clarity::{Address as EthAddress, Uint256};
use contact::client::Contact;
use cosmos_peggy::send::{send_request_batch, send_to_eth};
use cosmos_peggy::utils::wait_for_next_cosmos_block;
use cosmos_peggy::{query::get_oldest_unsigned_transaction_batch, send::send_ethereum_claims};
use deep_space::address::Address as CosmosAddress;
use deep_space::coin::Coin;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use ethereum_peggy::utils::get_valset_nonce;
use ethereum_peggy::{send_to_cosmos::send_to_cosmos, utils::get_tx_batch_nonce};
use orchestrator::main_loop::orchestrator_main_loop;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_utils::types::SendToCosmosEvent;
use rand::Rng;
use std::{env, process::Command, time::Duration};
use std::{process::ExitStatus, time::Instant};
use tokio::time::delay_for;
use tonic::transport::Channel;
use web30::client::Web3;

pub async fn happy_path_test(
    web30: &Web3,
    grpc_client: PeggyQueryClient<Channel>,
    contact: &Contact,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
    peggy_address: EthAddress,
    erc20_address: EthAddress,
    validator_out: bool,
) {
    let mut grpc_client = grpc_client;

    // used to break out of the loop early to simulate one validator
    // not running an Orchestrator
    let num_validators = keys.len();
    let mut count = 0;

    // start orchestrators
    #[allow(clippy::explicit_counter_loop)]
    for (c_key, e_key) in keys.iter() {
        info!("Spawning Orchestrator");
        let grpc_client = PeggyQueryClient::connect(COSMOS_NODE_GRPC).await.unwrap();
        // we have only one actual futures executor thread (see the actix runtime tag on our main function)
        // but that will execute all the orchestrators in our test in parallel
        Arbiter::spawn(orchestrator_main_loop(
            *c_key,
            *e_key,
            web30.clone(),
            contact.clone(),
            grpc_client,
            peggy_address,
            get_test_token_name(),
        ));

        // used to break out of the loop early to simulate one validator
        // not running an orchestrator
        count += 1;
        if validator_out && count == num_validators - 1 {
            break;
        }
    }

    // bootstrapping tests finish here and we move into operational tests

    // send 3 valset updates to make sure the process works back to back
    // don't do this in the validator out test because it changes powers
    // randomly and may actually make it impossible for that test to pass
    // by random re-allocation of powers. If we had 5 or 10 validators
    // instead of 3 this wouldn't be a problem. But with 3 not even 1% of
    // power can be reallocated to the down validator before things stop
    // working. We'll settle for testing that the initial valset (generated
    // with the first block) is successfully updated
    if !validator_out {
        for _ in 0u32..2 {
            test_valset_update(&web30, &keys, peggy_address).await;
        }
    } else {
        wait_for_nonzero_valset(&web30, peggy_address).await;
    }

    // generate an address for coin sending tests, this ensures test imdepotency
    let mut rng = rand::thread_rng();
    let secret: [u8; 32] = rng.gen();
    let dest_cosmos_private_key = CosmosPrivateKey::from_secret(&secret);
    let dest_cosmos_address = dest_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    let dest_eth_private_key = EthPrivateKey::from_slice(&secret).unwrap();
    let dest_eth_address = dest_eth_private_key.to_public_key().unwrap();

    // the denom and amount of the token bridged from Ethereum -> Cosmos
    // so the denom is the peggy<hash> token name
    // Send a token 3 times
    for _ in 0u32..3 {
        test_erc20_deposit(
            &web30,
            &contact,
            dest_cosmos_address,
            peggy_address,
            erc20_address,
            100u64.into(),
        )
        .await;
    }

    // We are going to submit a duplicate tx with nonce 1
    // This had better not increase the balance again
    // this test may have false positives if the timeout is not
    // long enough. TODO check for an error on the cosmos send response
    submit_duplicate_erc20_send(
        1u64.into(),
        &contact,
        erc20_address,
        1u64.into(),
        dest_cosmos_address,
        keys.clone(),
    )
    .await;

    // we test a batch by sending a transaction
    test_batch(
        &contact,
        &mut grpc_client,
        &web30,
        dest_eth_address,
        peggy_address,
        keys[0].0,
        dest_cosmos_private_key,
        erc20_address,
    )
    .await;
}

/// for downstream test cases the binary name needs to be different, so we detect
/// if a runtime env var is set and if so use that as the chain binary name for
/// our re-delegation. This is actually not the best way to handle this problem
/// ideally we would send the tx staking delegate command ourselves and never need
/// to know what the binary name is
pub fn chain_binary() -> Option<String> {
    match env::var("CHAIN_BINARY") {
        Ok(s) => Some(s),
        _ => None,
    }
}

/// This deals with the horrible error behavior of the Cosmos sdk command
/// line, mainly that when the tx seq changes while creating the message
/// it doesn't produce a failed output status or even print to stderror
/// it logs to stdout and simply fails. This will retry every 5 seconds
/// until the delegation succeeds
pub async fn delegate_tokens(delegate_address: &str, amount: &str) {
    let args = [
        "tx",
        "staking",
        "delegate",
        // target address
        delegate_address,
        // amount, this should be about 4% of the total power to start
        amount,
        "--home=/validator1",
        // this is defined in /tests/container-scripts/setup-validator.sh
        "--chain-id=peggy-test",
        "--keyring-backend=test",
        "--yes",
        "--from=validator1",
    ];
    let mut cmd = if let Some(bin) = chain_binary() {
        Command::new(bin)
    } else {
        Command::new("peggy")
    };

    let cmd = cmd.args(&args);
    let mut output = cmd
        .current_dir("/")
        .output()
        .expect("Failed to update stake to trigger validator set change");
    let mut stdout = String::from_utf8_lossy(&output.stdout).to_string();
    while stdout.contains("account sequence mismatch") {
        if stdout.contains("insufficient funds") {
            panic!("Can't continue to produce validator sets! Not enough STAKE token")
        }
        output = cmd
            .current_dir("/")
            .output()
            .expect("Failed to update stake to trigger validator set change");
        stdout = String::from_utf8_lossy(&output.stdout).to_string();
        delay_for(Duration::from_secs(5)).await;
    }
    info!("stdout: {}", stdout);
    info!("stderr: {}", String::from_utf8_lossy(&output.stderr));
    if !ExitStatus::success(&output.status) {
        panic!("Delegation failed!")
    }
}

pub async fn wait_for_nonzero_valset(web30: &Web3, peggy_address: EthAddress) {
    let start = Instant::now();
    let mut current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
        .await
        .expect("Failed to get current eth valset");

    while 0 == current_eth_valset_nonce {
        info!("Validator set is not yet updated to 0>, waiting",);
        current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
            .await
            .expect("Failed to get current eth valset");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TOTAL_TIMEOUT {
            panic!("Failed to update validator set");
        }
    }
}

pub async fn test_valset_update(
    web30: &Web3,
    keys: &[(CosmosPrivateKey, EthPrivateKey)],
    peggy_address: EthAddress,
) {
    // if we don't do this the orchestrators may run ahead of us and we'll be stuck here after
    // getting credit for two loops when we did one
    let starting_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
        .await
        .expect("Failed to get starting eth valset");
    let start = Instant::now();

    // now we send a valset request that the orchestrators will pick up on
    // in this case we send it as the first validator because they can pay the fee
    info!("Sending in valset request");

    // this is hacky and not really a good way to test validator set updates in a highly
    // repeatable fashion. What we really need to do is be aware of the total staking state
    // and manipulate the validator set very intentionally rather than kinda blindly like
    // we are here. For example the more your run this function the less this fixed amount
    // makes any difference, eventually it will fail because the change to the total staked
    // percentage is too small.
    let mut rng = rand::thread_rng();
    let validator_to_change = rng.gen_range(0..keys.len());
    let delegate_address = &keys[validator_to_change]
        .0
        .to_public_key()
        .unwrap()
        .to_address()
        .to_bech32("cosmosvaloper")
        .unwrap();
    let amount = &format!("{}stake", STARTING_STAKE_PER_VALIDATOR / 4);
    info!(
        "Delegating {} to {} in order to generate a validator set update",
        amount, delegate_address
    );
    delegate_tokens(delegate_address, amount).await;

    let mut current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
        .await
        .expect("Failed to get current eth valset");

    while starting_eth_valset_nonce == current_eth_valset_nonce {
        info!(
            "Validator set is not yet updated to {}>, waiting",
            starting_eth_valset_nonce
        );
        current_eth_valset_nonce = get_valset_nonce(peggy_address, *MINER_ADDRESS, &web30)
            .await
            .expect("Failed to get current eth valset");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TOTAL_TIMEOUT {
            panic!("Failed to update validator set");
        }
    }
    assert!(starting_eth_valset_nonce != current_eth_valset_nonce);
    info!("Validator set successfully updated!");
}

/// this function tests Ethereum -> Cosmos
async fn test_erc20_deposit(
    web30: &Web3,
    contact: &Contact,
    dest: CosmosAddress,
    peggy_address: EthAddress,
    erc20_address: EthAddress,
    amount: Uint256,
) {
    let start_coin = check_cosmos_balance("peggy", dest, &contact).await;
    info!(
        "Sending to Cosmos from {} to {} with amount {}",
        *MINER_ADDRESS, dest, amount
    );
    // we send some erc20 tokens to the peggy contract to register a deposit
    let tx_id = send_to_cosmos(
        erc20_address,
        peggy_address,
        amount.clone(),
        dest,
        *MINER_PRIVATE_KEY,
        Some(TOTAL_TIMEOUT),
        &web30,
        vec![],
    )
    .await
    .expect("Failed to send tokens to Cosmos");
    info!("Send to Cosmos txid: {:#066x}", tx_id);

    let start = Instant::now();
    while Instant::now() - start < TOTAL_TIMEOUT {
        match (
            start_coin.clone(),
            check_cosmos_balance("peggy", dest, &contact).await,
        ) {
            (Some(start_coin), Some(end_coin)) => {
                if start_coin.amount + amount.clone() == end_coin.amount
                    && start_coin.denom == end_coin.denom
                {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, start_coin.denom, end_coin.amount, end_coin.denom
                    );
                    return;
                }
            }
            (None, Some(end_coin)) => {
                if amount == end_coin.amount {
                    info!(
                        "Successfully bridged ERC20 {}{} to Cosmos! Balance is now {}{}",
                        amount, end_coin.denom, end_coin.amount, end_coin.denom
                    );
                    return;
                } else {
                    panic!("Failed to bridge ERC20!")
                }
            }
            _ => {}
        }
        info!("Waiting for ERC20 deposit");
        wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;
    }
    panic!("Failed to bridge ERC20!")
}

#[allow(clippy::too_many_arguments)]
async fn test_batch(
    contact: &Contact,
    grpc_client: &mut PeggyQueryClient<Channel>,
    web30: &Web3,
    dest_eth_address: EthAddress,
    peggy_address: EthAddress,
    requester_cosmos_private_key: CosmosPrivateKey,
    dest_cosmos_private_key: CosmosPrivateKey,
    erc20_contract: EthAddress,
) {
    let dest_cosmos_address = dest_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    let coin = check_cosmos_balance("peggy", dest_cosmos_address, &contact)
        .await
        .unwrap();
    let token_name = coin.denom;
    let amount = coin.amount;

    let bridge_denom_fee = Coin {
        denom: token_name.clone(),
        amount: 1u64.into(),
    };
    let amount = amount - 5u64.into();
    info!(
        "Sending {}{} from {} on Cosmos back to Ethereum",
        amount, token_name, dest_cosmos_address
    );
    let res = send_to_eth(
        dest_cosmos_private_key,
        dest_eth_address,
        Coin {
            denom: token_name.clone(),
            amount: amount.clone(),
        },
        bridge_denom_fee.clone(),
        &contact,
    )
    .await
    .unwrap();
    info!("Sent tokens to Ethereum with {:?}", res);

    info!("Requesting transaction batch");
    send_request_batch(
        requester_cosmos_private_key,
        token_name.clone(),
        get_fee(),
        &contact,
    )
    .await
    .unwrap();

    wait_for_next_cosmos_block(contact, TOTAL_TIMEOUT).await;
    let requester_address = requester_cosmos_private_key
        .to_public_key()
        .unwrap()
        .to_address();
    get_oldest_unsigned_transaction_batch(grpc_client, requester_address)
        .await
        .expect("Failed to get batch to sign");

    let mut current_eth_batch_nonce =
        get_tx_batch_nonce(peggy_address, erc20_contract, *MINER_ADDRESS, &web30)
            .await
            .expect("Failed to get current eth valset");
    let starting_batch_nonce = current_eth_batch_nonce;

    let start = Instant::now();
    while starting_batch_nonce == current_eth_batch_nonce {
        info!(
            "Batch is not yet submitted {}>, waiting",
            starting_batch_nonce
        );
        current_eth_batch_nonce =
            get_tx_batch_nonce(peggy_address, erc20_contract, *MINER_ADDRESS, &web30)
                .await
                .expect("Failed to get current eth tx batch nonce");
        delay_for(Duration::from_secs(4)).await;
        if Instant::now() - start > TOTAL_TIMEOUT {
            panic!("Failed to submit transaction batch set");
        }
    }

    let txid = web30
        .send_transaction(
            dest_eth_address,
            Vec::new(),
            1_000_000_000_000_000_000u128.into(),
            *MINER_ADDRESS,
            *MINER_PRIVATE_KEY,
            vec![],
        )
        .await
        .expect("Failed to send Eth to validator {}");
    web30
        .wait_for_transaction(txid, TOTAL_TIMEOUT, None)
        .await
        .unwrap();

    // we have to send this address one eth so that it can perform contract calls
    send_one_eth(dest_eth_address, web30).await;
    assert_eq!(
        web30
            .get_erc20_balance(erc20_contract, dest_eth_address)
            .await
            .unwrap(),
        amount
    );
    info!(
        "Successfully updated txbatch nonce to {} and sent {}{} tokens to Ethereum!",
        current_eth_batch_nonce, amount, token_name
    );
}

// this function submits a EthereumBridgeDepositClaim to the module with a given nonce. This can be set to be a nonce that has
// already been submitted to test the nonce functionality.
#[allow(clippy::too_many_arguments)]
async fn submit_duplicate_erc20_send(
    nonce: Uint256,
    contact: &Contact,
    erc20_address: EthAddress,
    amount: Uint256,
    receiver: CosmosAddress,
    keys: Vec<(CosmosPrivateKey, EthPrivateKey)>,
) {
    let start_coin = check_cosmos_balance("peggy", receiver, &contact)
        .await
        .expect("Did not find coins!");

    let ethereum_sender = "0x912fd21d7a69678227fe6d08c64222db41477ba0"
        .parse()
        .unwrap();

    let event = SendToCosmosEvent {
        event_nonce: nonce,
        block_height: 500u16.into(),
        erc20: erc20_address,
        sender: ethereum_sender,
        destination: receiver,
        amount,
    };

    // iterate through all validators and try to send an event with duplicate nonce
    for (c_key, _) in keys.iter() {
        let res = send_ethereum_claims(
            contact,
            *c_key,
            vec![event.clone()],
            vec![],
            vec![],
            vec![],
            get_fee(),
        )
        .await
        .unwrap();
        trace!("Submitted duplicate sendToCosmos event: {:?}", res);
    }

    if let Some(end_coin) = check_cosmos_balance("peggy", receiver, &contact).await {
        if start_coin.amount == end_coin.amount && start_coin.denom == end_coin.denom {
            info!("Successfully failed to duplicate ERC20!");
        } else {
            panic!("Duplicated ERC20!")
        }
    } else {
        panic!("Duplicate test failed for unknown reasons!");
    }
}
