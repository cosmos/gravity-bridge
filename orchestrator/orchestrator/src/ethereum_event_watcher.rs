//! Ethereum Event watcher watches for events such as a deposit to the Gravity Ethereum contract or a validator set update
//! or a transaction batch update. It then responds to these events by performing actions on the Cosmos chain if required

use clarity::{utils::bytes_to_hex_str, Address as EthAddress, Uint256};
use cosmos_gravity::{query::get_last_event_nonce, send::send_ethereum_claims};
use deep_space::Contact;
use deep_space::{coin::Coin, private_key::PrivateKey as CosmosPrivateKey};
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_utils::{
    error::GravityError,
    types::{
        Erc20DeployedEvent, LogicCallExecutedEvent, SendToCosmosEvent,
        TransactionBatchExecutedEvent, ValsetUpdatedEvent,
    },
};
use tonic::transport::Channel;
use web30::client::Web3;
use web30::jsonrpc::error::Web3Error;

use crate::get_with_retry::get_block_number_with_retry;
use crate::get_with_retry::get_net_version_with_retry;

pub async fn check_for_events(
    web3: &Web3,
    contact: &Contact,
    grpc_client: &mut GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    our_private_key: CosmosPrivateKey,
    fee: Coin,
    starting_block: Uint256,
) -> Result<Uint256, GravityError> {
    let our_cosmos_address = our_private_key.to_address(&contact.get_prefix()).unwrap();
    let latest_block = get_block_number_with_retry(web3).await;
    let latest_block = latest_block - get_block_delay(web3).await;

    let deposits = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![gravity_contract_address],
            vec!["SendToCosmosEvent(address,address,bytes32,uint256,uint256)"],
        )
        .await;
    debug!("Deposit events detected {:?}", deposits);

    let batches = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![gravity_contract_address],
            vec!["TransactionBatchExecutedEvent(uint256,address,uint256)"],
        )
        .await;
    debug!("Batche events detected {:?}", batches);

    let valsets = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![gravity_contract_address],
            vec!["ValsetUpdatedEvent(uint256,uint256,address[],uint256[])"],
        )
        .await;
    debug!("Valset events detected {:?}", valsets);

    let erc20_deployed = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![gravity_contract_address],
            vec!["ERC20DeployedEvent(string,address,string,string,uint8,uint256)"],
        )
        .await;
    debug!("ERC20 events detected {:?}", erc20_deployed);

    let logic_calls = web3
        .check_for_events(
            starting_block.clone(),
            Some(latest_block.clone()),
            vec![gravity_contract_address],
            vec!["LogicCallEvent(bytes32,uint256,bytes,uint256)"],
        )
        .await;
    debug!("Logic call events detected {:?}", logic_calls);

    if let (Ok(valsets), Ok(batches), Ok(deposits), Ok(deploys), Ok(logic_calls)) =
        (valsets, batches, deposits, erc20_deployed, logic_calls)
    {
        let deposits = SendToCosmosEvent::from_logs(&deposits)?;
        debug!("parsed deposits {:?}", deposits);

        let batches = TransactionBatchExecutedEvent::from_logs(&batches)?;
        debug!("parsed batches {:?}", batches);

        let valsets = ValsetUpdatedEvent::from_logs(&valsets)?;
        debug!("parsed valsets {:?}", valsets);

        let erc20_deploys = Erc20DeployedEvent::from_logs(&deploys)?;
        debug!("parsed erc20 deploys {:?}", erc20_deploys);

        let logic_calls = LogicCallExecutedEvent::from_logs(&logic_calls)?;
        debug!("logic call executions {:?}", logic_calls);

        // note that starting block overlaps with our last checked block, because we have to deal with
        // the possibility that the relayer was killed after relaying only one of multiple events in a single
        // block, so we also need this routine so make sure we don't send in the first event in this hypothetical
        // multi event block again. In theory we only send all events for every block and that will pass of fail
        // atomicly but lets not take that risk.
        let last_event_nonce = get_last_event_nonce(grpc_client, our_cosmos_address).await?;

        let deposits = SendToCosmosEvent::filter_by_event_nonce(last_event_nonce, &deposits);
        let batches = TransactionBatchExecutedEvent::filter_by_event_nonce(last_event_nonce, &batches);
        let valsets = ValsetUpdatedEvent::filter_by_event_nonce(last_event_nonce, &valsets);
        let erc20_deploys = Erc20DeployedEvent::filter_by_event_nonce(last_event_nonce, &erc20_deploys);
        let logic_calls = LogicCallExecutedEvent::filter_by_event_nonce(last_event_nonce, &logic_calls);

        for deposit in deposits.iter() {
            info!(
                    "Oracle observed deposit with ethereum sender {}, cosmos_reciever {}, amount {}, and event nonce {}",
                    deposit.sender, deposit.destination, deposit.amount, deposit.event_nonce
            );
        }

        for batch in batches.iter() {
            info!(
                "Oracle observed batch with nonce {}, contract {}, and event nonce {}",
                batch.batch_nonce, batch.erc20, batch.event_nonce
            );
        }

        for valset in valsets.iter() {
            info!(
                "Oracle observed valset with nonce {}, event nonce {}, block height {} and members {:?}",
                valset.valset_nonce, valset.event_nonce, valset.block_height, valset.members,
            )
        }

        for erc20_deploy in erc20_deploys.iter() {
            info!(
                "Oracle observed ERC20 deployment with denom {} erc20 name {} and symbol {} and event nonce {}",
                erc20_deploy.cosmos_denom, erc20_deploy.name, erc20_deploy.symbol, erc20_deploy.event_nonce,
            )
        }

        for logic_call in logic_calls.iter() {
            info!(
                "Oracle observed logic call execution with ID {} Nonce {} and event nonce {}",
                bytes_to_hex_str(&logic_call.invalidation_id),
                logic_call.invalidation_nonce,
                logic_call.event_nonce
            );
        }

        if !deposits.is_empty()
            || !batches.is_empty()
            || !valsets.is_empty()
            || !erc20_deploys.is_empty()
            || !logic_calls.is_empty()
        {
            let res = send_ethereum_claims(
                contact,
                our_private_key,
                deposits,
                batches,
                erc20_deploys,
                logic_calls,
                valsets,
                fee,
            )
            .await?;
            info!("Claims response {:?}", res);

            let new_event_nonce = get_last_event_nonce(grpc_client, our_cosmos_address).await?;
            // since we can't actually trust that the above txresponse is correct we have to check here
            // we may be able to trust the tx response post grpc
            if new_event_nonce == last_event_nonce {
                return Err(GravityError::InvalidBridgeStateError(
                    format!("Claims did not process, trying to update but still on {}, trying again in a moment, check txhash {} for errors", last_event_nonce, res.txhash),
                ));
            } else {
                info!("Claims processed, new nonce {}", new_event_nonce);
            }
        }
        Ok(latest_block)
    } else {
        error!("Failed to get events");
        Err(GravityError::EthereumRestError(Web3Error::BadResponse(
            "Failed to get logs!".to_string(),
        )))
    }
}

/// The number of blocks behind the 'latest block' on Ethereum our event checking should be.
/// Ethereum does not have finality and as such is subject to chain reorgs and temporary forks
/// if we check for events up to the very latest block we may process an event which did not
/// 'actually occur' in the longest POW chain.
///
/// Obviously we must chose some delay in order to prevent incorrect events from being claimed
///
/// For EVM chains with finality the correct value for this is zero. As there's no need
/// to concern ourselves with re-orgs or forking. This function checks the netID of the
/// provided Ethereum RPC and adjusts the block delay accordingly
///
/// The value used here for Ethereum is a balance between being reasonably fast and reasonably secure
/// As you can see on https://etherscan.io/blocks_forked uncles (one block deep reorgs)
/// occur once every few minutes. Two deep once or twice a day.
/// https://etherscan.io/chart/uncles
/// Let's make a conservative assumption of 1% chance of an uncle being a two block deep reorg
/// (actual is closer to 0.3%) and assume that continues as we increase the depth.
/// Given an uncle every 2.8 minutes, a 6 deep reorg would be 2.8 minutes * (100^4) or one
/// 6 deep reorg every 53,272 years.
///
pub async fn get_block_delay(web3: &Web3) -> Uint256 {
    let net_version = get_net_version_with_retry(web3).await;

    match net_version {
        // Mainline Ethereum, Ethereum classic, or the Ropsten, Mordor testnets
        // all POW Chains
        1 | 3 | 7 => 6u8.into(),
        // Rinkeby, Goerli, Dev, our own Gravity Ethereum testnet, and Kotti respectively
        // all non-pow chains
        4 | 5 | 2018 | 15 | 6 => 0u8.into(),
        // assume the safe option (POW) where we don't know
        _ => 6u8.into(),
    }
}
