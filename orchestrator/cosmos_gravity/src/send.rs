use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::address::Address;
use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey;
use deep_space::Contact;
use deep_space::Fee;
use deep_space::Msg;
use deep_space::{coin::Coin, utils::bytes_to_hex_str};
use ethereum_gravity::utils::downcast_uint256;
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::service_client::ServiceClient as TxServiceClient;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastMode;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastTxRequest;
use gravity_proto::gravity as proto;
use bytes::{Buf, BufMut};


use gravity_utils::message_signatures::{
    encode_logic_call_confirm, encode_tx_batch_confirm, encode_valset_confirm,
};
use gravity_utils::types::*;
use std::{collections::HashMap, time::Duration};

use bytes::BytesMut;
use prost::Message;
use prost_types::Any;

pub const MEMO: &str = "Sent using Althea Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
pub async fn update_gravity_delegate_addresses(
    contact: &Contact,
    delegate_eth_address: EthAddress,
    delegate_cosmos_address: Address,
    private_key: PrivateKey,
    eth_private_key: EthPrivateKey,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    trace!("Updating Gravity Delegate addresses");
    let our_valoper_address = private_key
        .to_address(&contact.get_prefix())
        .unwrap()
        // This works so long as the format set by the cosmos hub is maintained
        // having a main prefix followed by a series of titles for specific keys
        // this will not work if that convention is broken. This will be resolved when
        // GRPC exposes prefix endpoints (coming to upstream cosmos sdk soon)
        .to_bech32(format!("{}valoper", contact.get_prefix()))
        .unwrap();
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();

    let sequence = &contact.get_account_info(private_key
        .to_address(&contact.get_prefix()).unwrap()).await?.sequence;

    let eth_sign_msg = proto::DelegateKeysSignMsg{
        validator_address: our_valoper_address.clone(),
        nonce:*sequence,
    };
    let size = Message::encoded_len(&eth_sign_msg);
    let mut buf = BytesMut::with_capacity(size);
    Message::encode(&eth_sign_msg, &mut buf).expect("Failed to encode DelegateKeysSignMsg!");
    
    let eth_signature = eth_private_key.sign_ethereum_msg(&buf).to_bytes().to_vec();



    let msg_set_orch_address = proto::MsgDelegateKeys {
        validator_address: our_valoper_address.to_string(),
        orchestrator_address: delegate_cosmos_address.to_string(),
        ethereum_address: delegate_eth_address.to_string(),
        eth_signature
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgDelegateKeys", msg_set_orch_address);

    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&[msg], args, MEMO)?;

    let mut txrpc = TxServiceClient::connect(contact.get_url()).await?;
    let response = txrpc
        .broadcast_tx(BroadcastTxRequest {
            tx_bytes: msg_bytes,
            mode: BroadcastMode::Block.into(),
        })
        .await?;
    let response = response.into_inner();

    contact
        .wait_for_tx(response.tx_response.unwrap(), TIMEOUT)
        .await
}

/// Send in a confirmation for an array of validator sets, it's far more efficient to send these
/// as a single message
#[allow(clippy::too_many_arguments)]
pub async fn send_valset_confirms(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    valsets: Vec<Valset>,
    private_key: PrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();
    let our_eth_address = eth_private_key.to_public_key().unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let mut messages = Vec::new();

    for valset in valsets {
        trace!("Submitting signature for valset {:?}", valset);
        let message = encode_valset_confirm(gravity_id.clone(), valset.clone());
        let eth_signature = eth_private_key.sign_ethereum_msg(&message);
        info!(
            "Sending valset update address {} sig {} hash {}",
            our_eth_address,
            bytes_to_hex_str(&eth_signature.to_bytes()),
            bytes_to_hex_str(&message),
        );
        let confirm = proto::SignerSetTxConfirmation {
            ethereum_signer: our_eth_address.to_string(),
            signer_set_nonce: valset.nonce,
            signature: eth_signature.to_bytes().to_vec(),
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumTxConfirmation {
            signer: our_address.to_string(),
            confirmation: Some(Any {
                type_url: "/gravity.v1.SignerSetTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&messages, args, MEMO)?;

    let mut txrpc = TxServiceClient::connect(contact.get_url()).await?;
    let response = txrpc
        .broadcast_tx(BroadcastTxRequest {
            tx_bytes: msg_bytes,
            mode: BroadcastMode::Block.into(),
        })
        .await?;
    let response = response.into_inner();

    contact
        .wait_for_tx(response.tx_response.unwrap(), TIMEOUT)
        .await
}

/// Send in a confirmation for a specific transaction batch
pub async fn send_batch_confirm(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    transaction_batches: Vec<TransactionBatch>,
    private_key: PrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();
    let our_eth_address = eth_private_key.to_public_key().unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let mut messages = Vec::new();

    for batch in transaction_batches {
        info!("Submitting signature for batch {:?}", batch);
        let message = encode_tx_batch_confirm(gravity_id.clone(), batch.clone());
        let eth_signature = eth_private_key.sign_ethereum_msg(&message);
        info!(
            "Sending batch update address {} sig {} hash {}",
            our_eth_address,
            bytes_to_hex_str(&eth_signature.to_bytes()),
            bytes_to_hex_str(&message),
        );
        let confirm = proto::BatchTxConfirmation {
            token_contract: batch.token_contract.to_string(),
            batch_nonce: batch.nonce,
            ethereum_signer: our_eth_address.to_string(),
            signature: eth_signature.to_bytes().to_vec(),
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.BatchTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    let args = contact.get_message_args(our_address, fee).await?;
    info!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;
    contact.wait_for_tx(response, TIMEOUT).await
}

/// Send in a confirmation for a specific logic call
pub async fn send_logic_call_confirm(
    contact: &Contact,
    eth_private_key: EthPrivateKey,
    fee: Coin,
    logic_calls: Vec<LogicCall>,
    private_key: PrivateKey,
    gravity_id: String,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();
    let our_eth_address = eth_private_key.to_public_key().unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let mut messages = Vec::new();

    for call in logic_calls {
        trace!("Submitting signature for LogicCall {:?}", call);
        let message = encode_logic_call_confirm(gravity_id.clone(), call.clone());
        let eth_signature = eth_private_key.sign_ethereum_msg(&message);
        info!(
            "Sending LogicCall update address {} sig {} hash {}",
            our_eth_address,
            bytes_to_hex_str(&eth_signature.to_bytes()),
            bytes_to_hex_str(&message),
        );
        let confirm = proto::ContractCallTxConfirmation {
            ethereum_signer: our_eth_address.to_string(),
            signature: eth_signature.to_bytes().to_vec(),
            // TODO JEHAN: this will break
            invalidation_scope: bytes_to_hex_str(&call.invalidation_id).as_bytes().to_vec(),
            invalidation_nonce: call.invalidation_nonce,
        };
        let size = Message::encoded_len(&confirm);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&confirm, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumTxConfirmation {
            signer: our_address.to_string(),
            confirmation: Some(Any {
                type_url: "/gravity.v1.ContractCallTxConfirmation".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", wrapper);
        messages.push(msg);
    }
    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;
    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_ethereum_claims(
    contact: &Contact,
    private_key: PrivateKey,
    deposits: Vec<SendToCosmosEvent>,
    withdraws: Vec<TransactionBatchExecutedEvent>,
    erc20_deploys: Vec<Erc20DeployedEvent>,
    logic_calls: Vec<LogicCallExecutedEvent>,
    valsets: Vec<ValsetUpdatedEvent>,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();

    // This sorts oracle messages by event nonce before submitting them. It's not a pretty implementation because
    // we're missing an intermediary layer of abstraction. We could implement 'EventTrait' and then implement sort
    // for it, but then when we go to transform 'EventTrait' objects into GravityMsg enum values we'll have all sorts
    // of issues extracting the inner object from the TraitObject. Likewise we could implement sort of GravityMsg but that
    // would require a truly horrendous (nearly 100 line) match statement to deal with all combinations. That match statement
    // could be reduced by adding two traits to sort against but really this is the easiest option.
    //
    // We index the events by event nonce in an unordered hashmap and then play them back in order into a vec
    let mut unordered_msgs = HashMap::new();
    for deposit in deposits {
        let event = proto::SendToCosmosEvent {
            event_nonce: downcast_uint256(deposit.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(deposit.block_height).unwrap(),
            token_contract: deposit.erc20.to_string(),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.destination.to_string(),
            ethereum_sender: deposit.sender.to_string(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.SendToCosmosEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(deposit.event_nonce, msg);
    }
    for withdraw in withdraws {
        let event = proto::BatchExecutedEvent {
            event_nonce: downcast_uint256(withdraw.event_nonce.clone()).unwrap(),
            batch_nonce: downcast_uint256(withdraw.batch_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(withdraw.block_height).unwrap(),
            token_contract: withdraw.erc20.to_string(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.BatchExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(withdraw.event_nonce, msg);
    }
    for deploy in erc20_deploys {
        let event = proto::Erc20DeployedEvent {
            event_nonce: downcast_uint256(deploy.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(deploy.block_height).unwrap(),
            cosmos_denom: deploy.cosmos_denom,
            token_contract: deploy.erc20_address.to_string(),
            erc20_name: deploy.name,
            erc20_symbol: deploy.symbol,
            erc20_decimals: deploy.decimals as u64,
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.ERC20DeployedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(deploy.event_nonce, msg);
    }
    for call in logic_calls {
        let event = proto::ContractCallExecutedEvent {
            event_nonce: downcast_uint256(call.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(call.block_height).unwrap(),
            invalidation_id: call.invalidation_id,
            invalidation_nonce: downcast_uint256(call.invalidation_nonce).unwrap(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.ContractCallExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(call.event_nonce, msg);
    }
    for valset in valsets {
        let event = proto::SignerSetTxExecutedEvent {
            event_nonce: downcast_uint256(valset.event_nonce.clone()).unwrap(),
            signer_set_tx_nonce: downcast_uint256(valset.valset_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(valset.block_height).unwrap(),
            members: valset.members.iter().map(|v| v.into()).collect(),
        };
        let size = Message::encoded_len(&event);
        let mut buf = BytesMut::with_capacity(size);
        Message::encode(&event, &mut buf).expect("Failed to encode!"); // encoding should never fail so long as the buffer is big enough
        let wrapper = proto::MsgSubmitEthereumEvent {
            signer: our_address.to_string(),
            event: Some(Any {
                type_url: "/gravity.v1.SignerSetTxExecutedEvent".into(),
                value: buf.to_vec(),
            }),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", wrapper);
        unordered_msgs.insert(valset.event_nonce, msg);
    }

    let mut keys = Vec::new();
    for (key, _) in unordered_msgs.iter() {
        keys.push(key.clone());
    }
    keys.sort();

    let mut msgs = Vec::new();
    for i in keys {
        msgs.push(unordered_msgs.remove_entry(&i).unwrap().1);
    }

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64 * (msgs.len() as u64),
        granter: None,
        payer: None,
    };

    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&msgs, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;
    contact.wait_for_tx(response, TIMEOUT).await
}

/// Sends tokens from Cosmos to Ethereum. These tokens will not be sent immediately instead
/// they will require some time to be included in a batch
pub async fn send_to_eth(
    private_key: PrivateKey,
    destination: EthAddress,
    amount: Coin,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();
    if amount.denom != fee.denom {
        return Err(CosmosGrpcError::BadInput(format!(
            "{} {} is an invalid denom set for SendToEth you must pay fees in the same token your sending",
            amount.denom, fee.denom,
        )));
    }
    let balances = contact.get_balances(our_address).await.unwrap();
    let mut found = false;
    for balance in balances {
        if balance.denom == amount.denom {
            let total_amount = amount.amount.clone() + (fee.amount.clone() * 2u8.into());
            if balance.amount < total_amount {
                return Err(CosmosGrpcError::BadInput(format!(
                    "Insufficient balance of {} to send {}",
                    amount.denom, total_amount,
                )));
            }
            found = true;
        }
    }
    if !found {
        return Err(CosmosGrpcError::BadInput(format!(
            "No balance of {} to send",
            amount.denom,
        )));
    }

    let msg_send_to_eth = proto::MsgSendToEthereum {
        sender: our_address.to_string(),
        ethereum_recipient: destination.to_string(),
        amount: Some(amount.into()),
        bridge_fee: Some(fee.clone().into()),
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgSendToEthereum", msg_send_to_eth);

    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&[msg], args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;
    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_request_batch(
    private_key: PrivateKey,
    denom: String,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_address = private_key.to_address(&contact.get_prefix()).unwrap();

    let msg_request_batch = proto::MsgRequestBatchTx {
        signer: our_address.to_string(),
        denom,
    };

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64,
        granter: None,
        payer: None,
    };

    let msg = Msg::new("/gravity.v1.MsgRequestBatchTx", msg_request_batch);

    let args = contact.get_message_args(our_address, fee).await?;
    trace!("got optional tx info");

    let msg_bytes = private_key.sign_std_msg(&[msg], args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;
    contact.wait_for_tx(response, TIMEOUT).await
}
