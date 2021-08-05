use bytes::BytesMut;
use clarity::Address as EthAddress;
use clarity::PrivateKey as EthPrivateKey;
use deep_space::address::Address;
use deep_space::coin::Coin;
use deep_space::error::CosmosGrpcError;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::Contact;
use deep_space::Fee;
use deep_space::Msg;
use gravity_proto::cosmos_sdk_proto::cosmos::base::abci::v1beta1::TxResponse;
use gravity_proto::cosmos_sdk_proto::cosmos::tx::v1beta1::BroadcastMode;
use gravity_proto::gravity as proto;
use prost::Message;
use std::time::Duration;

pub const MEMO: &str = "Sent using Althea Orchestrator";
pub const TIMEOUT: Duration = Duration::from_secs(60);

/// Send a transaction updating the eth address for the sending
/// Cosmos address. The sending Cosmos address should be a validator
pub async fn update_gravity_delegate_addresses(
    contact: &Contact,
    delegate_eth_address: EthAddress,
    delegate_cosmos_address: Address,
    cosmos_key: CosmosPrivateKey,
    etheruem_key: EthPrivateKey,
    fee: Coin,
) -> Result<TxResponse, CosmosGrpcError> {
    let our_valoper_address = cosmos_key
        .to_address(&contact.get_prefix())
        .unwrap()
        // This works so long as the format set by the cosmos hub is maintained
        // having a main prefix followed by a series of titles for specific keys
        // this will not work if that convention is broken. This will be resolved when
        // GRPC exposes prefix endpoints (coming to upstream cosmos sdk soon)
        .to_bech32(format!("{}valoper", contact.get_prefix()))
        .unwrap();

    let nonce = contact
        .get_account_info(cosmos_key.to_address(&contact.get_prefix()).unwrap())
        .await?
        .sequence;

    let eth_sign_msg = proto::DelegateKeysSignMsg {
        validator_address: our_valoper_address.clone(),
        nonce,
    };

    let mut data = BytesMut::with_capacity(eth_sign_msg.encoded_len());
    Message::encode(&eth_sign_msg, &mut data).expect("encoding failed");

    let eth_signature = etheruem_key.sign_ethereum_msg(&data).to_bytes().to_vec();
    let msg = proto::MsgDelegateKeys {
        validator_address: our_valoper_address.to_string(),
        orchestrator_address: delegate_cosmos_address.to_string(),
        ethereum_address: delegate_eth_address.to_string(),
        eth_signature,
    };
    let msg = Msg::new("/gravity.v1.MsgDelegateKeys", msg);
    __send_messages(contact, cosmos_key, fee, vec![msg]).await
}

/// Sends tokens from Cosmos to Ethereum. These tokens will not be sent immediately instead
/// they will require some time to be included in a batch
pub async fn send_to_eth(
    cosmos_key: CosmosPrivateKey,
    destination: EthAddress,
    amount: Coin,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let msg = proto::MsgSendToEthereum {
        sender: cosmos_address.to_string(),
        ethereum_recipient: destination.to_string(),
        amount: Some(amount.into()),
        bridge_fee: Some(fee.clone().into()),
    };
    let msg = Msg::new("/gravity.v1.MsgSendToEthereum", msg);
    __send_messages(contact, cosmos_key, fee, vec![msg]).await
}

pub async fn send_request_batch_tx(
    cosmos_key: CosmosPrivateKey,
    denom: String,
    fee: Coin,
    contact: &Contact,
) -> Result<TxResponse, CosmosGrpcError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let msg_request_batch = proto::MsgRequestBatchTx {
        signer: cosmos_address.to_string(),
        denom,
    };
    let msg = Msg::new("/gravity.v1.MsgRequestBatchTx", msg_request_batch);
    __send_messages(contact, cosmos_key, fee, vec![msg]).await
}

// TODO(Levi) teach this branch to accept gas_prices
async fn __send_messages(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    fee: Coin,
    messages: Vec<Msg>,
) -> Result<TxResponse, CosmosGrpcError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let fee = Fee {
        amount: vec![fee],
        gas_limit: 500_000_000u64 * (messages.len() as u64),
        granter: None,
        payer: None,
    };

    let args = contact.get_message_args(cosmos_address, fee).await?;

    let msg_bytes = cosmos_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_messages(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    gas_price: (f64, String),
    messages: Vec<Msg>,
) -> Result<TxResponse, CosmosGrpcError> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    let gas_limit = 500_000_000 * messages.len();

    let fee_amount: f64 = (gas_limit as f64) * gas_price.0;
    let fee_amount: u64 = fee_amount.abs().ceil() as u64;

    let fee_amount = Coin {
        denom: gas_price.1,
        amount: fee_amount.into(),
    };

    let gas_limit = gas_limit as u64;
    let fee = Fee {
        amount: vec![fee_amount],
        gas_limit,
        granter: None,
        payer: None,
    };

    let args = contact.get_message_args(cosmos_address, fee).await?;

    let msg_bytes = cosmos_key.sign_std_msg(&messages, args, MEMO)?;

    let response = contact
        .send_transaction(msg_bytes, BroadcastMode::Sync)
        .await?;

    contact.wait_for_tx(response, TIMEOUT).await
}

pub async fn send_main_loop(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    gas_price: (f64, String),
    mut rx: tokio::sync::mpsc::Receiver<Vec<Msg>>,
) {
    while let Some(messages) = rx.recv().await {
        match send_messages(contact, cosmos_key, gas_price.to_owned(), messages).await {
            Ok(res) => trace!("okay: {:?}", res),
            Err(err) => error!("fail: {}", err),
        }
    }
}
