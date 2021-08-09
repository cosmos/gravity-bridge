use clarity::PrivateKey as EthPrivateKey;
use deep_space::private_key::PrivateKey as CosmosPrivateKey;
use deep_space::utils::bytes_to_hex_str;
use deep_space::Contact;
use deep_space::Msg;
use ethereum_gravity::utils::downcast_uint256;
use gravity_proto::gravity as proto;
use gravity_proto::ToAny;
use gravity_utils::message_signatures::{
    encode_logic_call_confirm, encode_tx_batch_confirm, encode_valset_confirm,
};
use gravity_utils::types::*;

pub fn signer_set_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    valsets: Vec<Valset>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let ethereum_address = ethereum_key.to_public_key().unwrap();

    let mut msgs = Vec::new();
    for valset in valsets {
        let data = encode_valset_confirm(gravity_id.clone(), valset.clone());
        let signature = ethereum_key.sign_ethereum_msg(&data);
        let confirmation = proto::SignerSetTxConfirmation {
            ethereum_signer: ethereum_address.to_string(),
            signer_set_nonce: valset.nonce,
            signature: signature.to_bytes().to_vec(),
        };
        let msg = proto::MsgSubmitEthereumTxConfirmation {
            signer: cosmos_address.to_string(),
            confirmation: confirmation.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", msg);
        msgs.push(msg);
    }
    msgs
}

pub fn batch_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    batches: Vec<TransactionBatch>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let ethereum_address = ethereum_key.to_public_key().unwrap();

    let mut msgs = Vec::new();
    for batch in batches {
        let data = encode_tx_batch_confirm(gravity_id.clone(), batch.clone());
        let signature = ethereum_key.sign_ethereum_msg(&data);
        let confirmation = proto::BatchTxConfirmation {
            token_contract: batch.token_contract.to_string(),
            batch_nonce: batch.nonce,
            ethereum_signer: ethereum_address.to_string(),
            signature: signature.to_bytes().to_vec(),
        };
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: confirmation.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", msg);
        msgs.push(msg);
    }
    msgs
}

pub fn contract_call_tx_confirmation_messages(
    contact: &Contact,
    ethereum_key: EthPrivateKey,
    logic_calls: Vec<LogicCall>,
    cosmos_key: CosmosPrivateKey,
    gravity_id: String,
) -> Vec<Msg> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();
    let ethereum_address = ethereum_key.to_public_key().unwrap();

    let mut msgs = Vec::new();
    for logic_call in logic_calls {
        let data = encode_logic_call_confirm(gravity_id.clone(), logic_call.clone());
        let signature = ethereum_key.sign_ethereum_msg(&data);
        let confirmation = proto::ContractCallTxConfirmation {
            ethereum_signer: ethereum_address.to_string(),
            signature: signature.to_bytes().to_vec(),
            invalidation_scope: bytes_to_hex_str(&logic_call.invalidation_id)
                .as_bytes()
                .to_vec(),
            invalidation_nonce: logic_call.invalidation_nonce,
        };
        let msg = proto::MsgSubmitEthereumTxConfirmation {
            signer: cosmos_address.to_string(),
            confirmation: confirmation.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumTxConfirmation", msg);
        msgs.push(msg);
    }
    msgs
}

pub fn ethereum_event_messages(
    contact: &Contact,
    cosmos_key: CosmosPrivateKey,
    deposits: Vec<SendToCosmosEvent>,
    batches: Vec<TransactionBatchExecutedEvent>,
    erc20_deploys: Vec<Erc20DeployedEvent>,
    logic_calls: Vec<LogicCallExecutedEvent>,
    valsets: Vec<ValsetUpdatedEvent>,
) -> Vec<Msg> {
    let cosmos_address = cosmos_key.to_address(&contact.get_prefix()).unwrap();

    // This sorts oracle messages by event nonce before submitting them. It's not a pretty implementation because
    // we're missing an intermediary layer of abstraction. We could implement 'EventTrait' and then implement sort
    // for it, but then when we go to transform 'EventTrait' objects into GravityMsg enum values we'll have all sorts
    // of issues extracting the inner object from the TraitObject. Likewise we could implement sort of GravityMsg but that
    // would require a truly horrendous (nearly 100 line) match statement to deal with all combinations. That match statement
    // could be reduced by adding two traits to sort against but really this is the easiest option.
    //
    // We index the events by event nonce in an unordered hashmap and then play them back in order into a vec
    let mut unordered_msgs = std::collections::HashMap::new();
    for deposit in deposits {
        let event = proto::SendToCosmosEvent {
            event_nonce: downcast_uint256(deposit.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(deposit.block_height).unwrap(),
            token_contract: deposit.erc20.to_string(),
            amount: deposit.amount.to_string(),
            cosmos_receiver: deposit.destination.to_string(),
            ethereum_sender: deposit.sender.to_string(),
        };
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: event.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", msg);
        unordered_msgs.insert(deposit.event_nonce, msg);
    }
    for batch in batches {
        let event = proto::BatchExecutedEvent {
            event_nonce: downcast_uint256(batch.event_nonce.clone()).unwrap(),
            batch_nonce: downcast_uint256(batch.batch_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(batch.block_height).unwrap(),
            token_contract: batch.erc20.to_string(),
        };
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: event.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", msg);
        unordered_msgs.insert(batch.event_nonce, msg);
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
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: event.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", msg);
        unordered_msgs.insert(deploy.event_nonce, msg);
    }
    for logic_call in logic_calls {
        let event = proto::ContractCallExecutedEvent {
            event_nonce: downcast_uint256(logic_call.event_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(logic_call.block_height).unwrap(),
            invalidation_id: logic_call.invalidation_id,
            invalidation_nonce: downcast_uint256(logic_call.invalidation_nonce).unwrap(),
        };
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: event.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", msg);
        unordered_msgs.insert(logic_call.event_nonce, msg);
    }
    for valset in valsets {
        let event = proto::SignerSetTxExecutedEvent {
            event_nonce: downcast_uint256(valset.event_nonce.clone()).unwrap(),
            signer_set_tx_nonce: downcast_uint256(valset.valset_nonce.clone()).unwrap(),
            ethereum_height: downcast_uint256(valset.block_height).unwrap(),
            members: valset.members.iter().map(|v| v.into()).collect(),
        };
        let msg = proto::MsgSubmitEthereumEvent {
            signer: cosmos_address.to_string(),
            event: event.to_any(),
        };
        let msg = Msg::new("/gravity.v1.MsgSubmitEthereumEvent", msg);
        unordered_msgs.insert(valset.event_nonce, msg);
    }

    let mut keys = Vec::new();
    for (key, _) in unordered_msgs.iter() {
        keys.push(key.clone());
    }
    keys.sort();

    let mut msgs = Vec::new();
    for i in keys.iter() {
        msgs.push(unordered_msgs.remove_entry(&i).unwrap().1);
    }

    msgs
}
