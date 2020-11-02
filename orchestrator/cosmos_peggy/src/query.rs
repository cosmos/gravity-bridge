use clarity::Uint256;
use contact::client::Contact;
use contact::jsonrpc::error::JsonRpcError;
use contact::types::ResponseWrapper;
use contact::types::TypeWrapper;
use deep_space::address::Address;
use peggy_utils::types::*;

/// get the valset for a given nonce (block) height
pub async fn get_valset(
    contact: &Contact,
    nonce: Uint256,
) -> Result<ResponseWrapper<Valset>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<TypeWrapper<ValsetUnparsed>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &format!("peggy/valset_request/{}", nonce),
            none,
            contact.timeout,
            None,
        )
        .await;
    trace!("Unparsed valset {:?}", ret);
    match ret {
        Ok(val) => Ok(ResponseWrapper {
            height: val.height,
            result: val.result.value.convert(),
        }),
        Err(e) => Err(e),
    }
}

/// This hits the /pending_valset_requests endpoint and will provide the oldest
/// validator set we have not yet signed.
pub async fn get_oldest_unsigned_valset(
    contact: &Contact,
    address: Address,
) -> Result<ResponseWrapper<Valset>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<TypeWrapper<ValsetUnparsed>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &format!("peggy/pending_valset_requests/{}", address),
            none,
            contact.timeout,
            None,
        )
        .await;
    match ret {
        Ok(val) => Ok(ResponseWrapper {
            height: val.height,
            result: val.result.value.convert(),
        }),
        Err(e) => Err(e),
    }
}

/// this input views the last five valest requests that have been made, useful if you're
/// a relayer looking to ferry confirmations
pub async fn get_latest_valsets(
    contact: &Contact,
) -> Result<ResponseWrapper<Vec<Valset>>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<Vec<ValsetUnparsed>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &"peggy/valset_requests".to_string(),
            none,
            contact.timeout,
            None,
        )
        .await;

    match ret {
        Ok(val) => {
            let mut converted_values = Vec::new();
            for item in val.result {
                converted_values.push(item.convert());
            }
            Ok(ResponseWrapper {
                height: val.height,
                result: converted_values,
            })
        }
        Err(e) => Err(e),
    }
}

/// get all valset confirmations for a given nonce
pub async fn get_all_valset_confirms(
    contact: &Contact,
    nonce: u64,
) -> Result<ResponseWrapper<Vec<ValsetConfirmResponse>>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<Vec<ValsetConfirmResponse>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &format!("peggy/valset_confirm/{}", nonce),
            none,
            contact.timeout,
            None,
        )
        .await;
    match ret {
        Ok(val) => Ok(val),
        Err(e) => Err(e),
    }
}
pub async fn get_oldest_unsigned_transaction_batch(
    contact: &Contact,
    address: Address,
) -> Result<ResponseWrapper<TransactionBatch>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<TransactionBatch>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &format!("peggy/pending_batch_requests/{}", address),
            none,
            contact.timeout,
            None,
        )
        .await;
    match ret {
        Ok(val) => Ok(val),
        Err(e) => Err(e),
    }
}

pub async fn get_latest_transaction_batches(
    contact: &Contact,
) -> Result<ResponseWrapper<Vec<TransactionBatch>>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<Vec<TransactionBatch>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &"peggy/transaction_batches/".to_string(),
            none,
            contact.timeout,
            None,
        )
        .await;
    match ret {
        Ok(val) => Ok(val),
        Err(e) => Err(e),
    }
}

/// get all batch confirmations for a given nonce
pub async fn get_signed_transaction_batches(
    contact: &Contact,
) -> Result<ResponseWrapper<Vec<SignedTransactionBatch>>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<Vec<SignedTransactionBatch>>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method(
            &"peggy/signed_batches".to_string(),
            none,
            contact.timeout,
            None,
        )
        .await;
    match ret {
        Ok(val) => Ok(val),
        Err(e) => Err(e),
    }
}
