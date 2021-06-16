use clarity::Address as EthAddress;
use deep_space::address::Address;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use gravity_proto::gravity::*;
use gravity_utils::error::GravityError;
use gravity_utils::types::*;
use tonic::transport::Channel;

/// get the valset for a given nonce (block) height
pub async fn get_valset(
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
) -> Result<Option<Valset>, GravityError> {
    let response = client
        .signer_set_tx(SignerSetTxRequest {
            signer_set_nonce: nonce,
        })
        .await?;
    let valset = response.into_inner().signer_set;
    let valset = match valset {
        Some(v) => Some(v.into()),
        None => None,
    };
    Ok(valset)
}

/// This hits the /pending_valset_requests endpoint and will provide
/// an array of validator sets we have not already signed
pub async fn get_oldest_unsigned_valsets(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Vec<Valset>, GravityError> {
    let response = client
        .unsigned_signer_set_txs(UnsignedSignerSetTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let valsets = response.into_inner().signer_sets;
    // convert from proto valset type to rust valset type
    let valsets = valsets.iter().map(|v| v.clone().into()).collect();
    Ok(valsets)
}

/// this input views the last five signer set txs that have been made, useful if you're
/// a relayer looking to ferry confirmations
pub async fn get_latest_valset(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Option<Valset>, GravityError> {
    let response = client
        .latest_signer_set_tx(LatestSignerSetTxRequest {})
        .await?;
    let valset = response.into_inner().signer_set;
    let valset = match valset {
        Some(v) => Some(v.into()),
        None => None,
    };
    Ok(valset)
}

/// get all valset confirmations for a given nonce
pub async fn get_all_valset_confirms(
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
) -> Result<Vec<ValsetConfirmResponse>, GravityError> {
    let request = client
        .signer_set_tx_confirmations(SignerSetTxConfirmationsRequest {
            signer_set_nonce: nonce,
        })
        .await?;
    let confirms = request.into_inner().signatures;
    let mut parsed_confirms = Vec::new();
    for item in confirms {
        parsed_confirms.push(ValsetConfirmResponse::from_proto(item)?)
    }
    Ok(parsed_confirms)
}

pub async fn get_oldest_unsigned_transaction_batch(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Option<TransactionBatch>, GravityError> {
    let request = client
        .unsigned_batch_txs(UnsignedBatchTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let batches = request.into_inner().batches;
    let batch = batches.get(0);
    match batch {
        Some(batch) => Ok(Some(TransactionBatch::from_proto(batch.clone())?)),
        None => Ok(None),
    }
}

/// gets the latest 100 transaction batches, regardless of token type
/// for relayers to consider relaying
pub async fn get_latest_transaction_batches(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Vec<TransactionBatch>, GravityError> {
    let request = client
        .batch_txs(BatchTxsRequest { pagination: None })
        .await?;
    let batches = request.into_inner().batches;
    let mut out = Vec::new();
    for batch in batches {
        out.push(TransactionBatch::from_proto(batch)?)
    }
    Ok(out)
}

/// get all batch confirmations for a given nonce and denom
pub async fn get_transaction_batch_signatures(
    client: &mut GravityQueryClient<Channel>,
    nonce: u64,
    contract_address: EthAddress,
) -> Result<Vec<BatchConfirmResponse>, GravityError> {
    let request = client
        .batch_tx_confirmations(BatchTxConfirmationsRequest {
            batch_nonce: nonce,
            token_contract: contract_address.to_string(),
        })
        .await?;
    let batch_confirms = request.into_inner().signatures;
    let mut out = Vec::new();
    for confirm in batch_confirms {
        out.push(BatchConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<u64, GravityError> {
    let request = client
        .last_submitted_ethereum_event(LastSubmittedEthereumEventRequest {
            address: address.to_string(),
        })
        .await?;
    Ok(request.into_inner().event_nonce)
}

/// Gets the 100 latest logic calls for a relayer to consider relaying
pub async fn get_latest_logic_calls(
    client: &mut GravityQueryClient<Channel>,
) -> Result<Vec<LogicCall>, GravityError> {
    let request = client
        .contract_call_txs(ContractCallTxsRequest { pagination: None })
        .await?;
    let calls = request.into_inner().calls;
    let mut out = Vec::new();
    for call in calls {
        out.push(LogicCall::from_proto(call)?);
    }
    Ok(out)
}

pub async fn get_logic_call_signatures(
    client: &mut GravityQueryClient<Channel>,
    invalidation_scope: Vec<u8>,
    invalidation_nonce: u64,
) -> Result<Vec<LogicCallConfirmResponse>, GravityError> {
    let request = client
        .contract_call_tx_confirmations(ContractCallTxConfirmationsRequest {
            invalidation_scope,
            invalidation_nonce,
        })
        .await?;
    let call_confirms = request.into_inner().signatures;
    let mut out = Vec::new();
    for confirm in call_confirms {
        out.push(LogicCallConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

pub async fn get_oldest_unsigned_logic_call(
    client: &mut GravityQueryClient<Channel>,
    address: Address,
) -> Result<Vec<LogicCall>, GravityError> {
    let request = client
        .unsigned_contract_call_txs(UnsignedContractCallTxsRequest {
            address: address.to_string(),
        })
        .await?;
    let calls = request.into_inner().calls;
    let mut out = Vec::new();
    for call in calls {
        out.push(LogicCall::from_proto(call)?)
    }
    Ok(out)
}
