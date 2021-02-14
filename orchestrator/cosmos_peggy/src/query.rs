use clarity::Address as EthAddress;
use deep_space::address::Address;
use peggy_proto::peggy::query_client::QueryClient as PeggyQueryClient;
use peggy_proto::peggy::QueryBatchConfirmsRequest;
use peggy_proto::peggy::QueryCurrentValsetRequest;
use peggy_proto::peggy::QueryLastEventNonceByAddrRequest;
use peggy_proto::peggy::QueryLastPendingBatchRequestByAddrRequest;
use peggy_proto::peggy::QueryLastPendingValsetRequestByAddrRequest;
use peggy_proto::peggy::QueryLastValsetRequestsRequest;
use peggy_proto::peggy::QueryOutgoingTxBatchesRequest;
use peggy_proto::peggy::QueryValsetConfirmsByNonceRequest;
use peggy_proto::peggy::QueryValsetRequestRequest;
use peggy_utils::error::PeggyError;
use peggy_utils::types::*;
use tonic::transport::Channel;

/// get the valset for a given nonce (block) height
pub async fn get_valset(
    client: &mut PeggyQueryClient<Channel>,
    nonce: u64,
) -> Result<Option<Valset>, PeggyError> {
    let request = client
        .valset_request(QueryValsetRequestRequest { nonce })
        .await?;
    let valset = request.into_inner().valset;
    let valset = match valset {
        Some(v) => Some(v.into()),
        None => None,
    };
    Ok(valset)
}

/// get the current valset. You should never sign this valset
/// valset requests create a consensus point around the block height
/// that transaction got in. Without that consensus point everyone trying
/// to sign the 'current' valset would run into slight differences and fail
/// to produce a viable update.
pub async fn get_current_valset(
    client: &mut PeggyQueryClient<Channel>,
) -> Result<Valset, PeggyError> {
    let request = client.current_valset(QueryCurrentValsetRequest {}).await?;
    let valset = request.into_inner().valset;
    if let Some(valset) = valset {
        Ok(valset.into())
    } else {
        error!("Current valset returned None? This should be impossible");
        Err(PeggyError::InvalidBridgeStateError(
            "Must have a current valset!".to_string(),
        ))
    }
}

/// This hits the /pending_valset_requests endpoint and will provide
/// an array of validator sets we have not already signed
pub async fn get_oldest_unsigned_valsets(
    client: &mut PeggyQueryClient<Channel>,
    address: Address,
) -> Result<Vec<Valset>, PeggyError> {
    let request = client
        .last_pending_valset_request_by_addr(QueryLastPendingValsetRequestByAddrRequest {
            address: address.to_string(),
        })
        .await?;
    let valsets = request.into_inner().valsets;
    // convert from proto valset type to rust valset type
    let valsets = valsets.iter().map(|v| v.into()).collect();
    Ok(valsets)
}

/// this input views the last five valset requests that have been made, useful if you're
/// a relayer looking to ferry confirmations
pub async fn get_latest_valsets(
    client: &mut PeggyQueryClient<Channel>,
) -> Result<Vec<Valset>, PeggyError> {
    let request = client
        .last_valset_requests(QueryLastValsetRequestsRequest {})
        .await?;
    let valsets = request.into_inner().valsets;
    Ok(valsets.iter().map(|v| v.into()).collect())
}

/// get all valset confirmations for a given nonce
pub async fn get_all_valset_confirms(
    client: &mut PeggyQueryClient<Channel>,
    nonce: u64,
) -> Result<Vec<ValsetConfirmResponse>, PeggyError> {
    let request = client
        .valset_confirms_by_nonce(QueryValsetConfirmsByNonceRequest { nonce })
        .await?;
    let confirms = request.into_inner().confirms;
    let mut parsed_confirms = Vec::new();
    for item in confirms {
        parsed_confirms.push(ValsetConfirmResponse::from_proto(item)?)
    }
    Ok(parsed_confirms)
}

pub async fn get_oldest_unsigned_transaction_batch(
    client: &mut PeggyQueryClient<Channel>,
    address: Address,
) -> Result<Option<TransactionBatch>, PeggyError> {
    let request = client
        .last_pending_batch_request_by_addr(QueryLastPendingBatchRequestByAddrRequest {
            address: address.to_string(),
        })
        .await?;
    let batch = request.into_inner().batch;
    match batch {
        Some(batch) => Ok(Some(TransactionBatch::from_proto(batch)?)),
        None => Ok(None),
    }
}

pub async fn get_latest_transaction_batches(
    client: &mut PeggyQueryClient<Channel>,
) -> Result<Vec<TransactionBatch>, PeggyError> {
    let request = client
        .outgoing_tx_batches(QueryOutgoingTxBatchesRequest {})
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
    client: &mut PeggyQueryClient<Channel>,
    nonce: u64,
    contract_address: EthAddress,
) -> Result<Vec<BatchConfirmResponse>, PeggyError> {
    let request = client
        .batch_confirms(QueryBatchConfirmsRequest {
            nonce,
            contract_address: contract_address.to_string(),
        })
        .await?;
    let batch_confirms = request.into_inner().confirms;
    let mut out = Vec::new();
    for confirm in batch_confirms {
        out.push(BatchConfirmResponse::from_proto(confirm)?)
    }
    Ok(out)
}

/// Gets the last event nonce that a given validator has attested to, this lets us
/// catch up with what the current event nonce should be if a oracle is restarted
pub async fn get_last_event_nonce(
    client: &mut PeggyQueryClient<Channel>,
    address: Address,
) -> Result<u64, PeggyError> {
    let request = client
        .last_event_nonce_by_addr(QueryLastEventNonceByAddrRequest {
            address: address.to_string(),
        })
        .await?;
    Ok(request.into_inner().event_nonce)
}
