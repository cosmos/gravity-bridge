use crate::cosmos_interop::types::*;
use contact::client::Contact;
use contact::jsonrpc::error::JsonRpcError;
use contact::types::ResponseWrapper;
use contact::types::TypeWrapper;
use deep_space::address::Address;

/// Get the latest valset recorded by the peggy module. If no valset has ever been created
/// you will instead get a blank valset at height 0. Any value above this may or may not
/// be a complete valset and it's up to the caller to interpret the response.
pub async fn get_peggy_valset(contact: &Contact) -> Result<ResponseWrapper<Valset>, JsonRpcError> {
    let none: Option<bool> = None;
    let ret: Result<ResponseWrapper<ValsetUnparsed>, JsonRpcError> = contact
        .jsonrpc_client
        .request_method("peggy/current_valset", none, contact.timeout, None)
        .await;
    match ret {
        Ok(val) => Ok(ResponseWrapper {
            height: val.height,
            result: val.result.convert(),
        }),
        Err(e) => Err(e),
    }
}

/// get the valset for a given nonce (block) height
pub async fn get_peggy_valset_request(
    contact: &Contact,
    nonce: u128,
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
pub async fn get_last_valset_requests(
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
