use clarity::Address as EthAddress;
use deep_space::address::Address;
use deep_space::canonical_json::{to_canonical_json, CanonicalJsonError};
use deep_space::coin::Coin;
use deep_space::msg::DeepSpaceMsg;
use ethereum_peggy::utils::downcast_uint256;
use num256::Uint256;
use peggy_utils::types::{
    ERC20DeployedEvent, LogicCallExecutedEvent, SendToCosmosEvent, TransactionBatchExecutedEvent,
};
/// Any arbitrary message
#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq)]
#[serde(tag = "type", content = "value")]
pub enum PeggyMsg {
    #[serde(rename = "peggy/MsgSetOrchestratorAddress")]
    SetOrchestratorAddressMsg(SetOrchestratorAddressMsg),

    #[serde(rename = "peggy/MsgValsetConfirm")]
    ValsetConfirmMsg(ValsetConfirmMsg),

    #[serde(rename = "peggy/MsgSendToEth")]
    SendToEthMsg(SendToEthMsg),

    #[serde(rename = "peggy/MsgRequestBatch")]
    RequestBatchMsg(RequestBatchMsg),

    #[serde(rename = "peggy/MsgConfirmBatch")]
    ConfirmBatchMsg(ConfirmBatchMsg),

    #[serde(rename = "peggy/MsgConfirmLogicCall")]
    ConfirmLogicCallMsg(ConfirmLogicCallMsg),

    #[serde(rename = "peggy/MsgDepositClaim")]
    DepositClaimMsg(DepositClaimMsg),

    #[serde(rename = "peggy/MsgWithdrawClaim")]
    WithdrawClaimMsg(WithdrawClaimMsg),

    #[serde(rename = "peggy/MsgERC20DeployedClaim")]
    ERC20DeployedClaimMsg(ERC20DeployedClaimMsg),

    #[serde(rename = "peggy/MsgLogicCallExecutedClaim")]
    LogicCallExecutedClaim(LogicCallExecutedClaim),
}

impl DeepSpaceMsg for PeggyMsg {
    fn to_sign_bytes(&self) -> Result<Vec<u8>, CanonicalJsonError> {
        Ok(to_canonical_json(self)?)
    }
}

/// This message sets both the Cosmos and Ethereum address being delegated for
/// Orchestrator operations. This allows a validator to use their highly valuable
/// valoper key to simply sign off on these addresses.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct SetOrchestratorAddressMsg {
    #[serde(rename = "eth_address")]
    // the Ethereum address being delegated to
    pub eth_address: EthAddress,
    // the valoper address
    pub validator: String,
    // the Cosmos address being delegated to
    pub orchestrator: Address,
}
/// a transaction we send to submit a valset confirmation signature
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ValsetConfirmMsg {
    pub orchestrator: Address,
    pub eth_address: EthAddress,
    pub nonce: Uint256,
    #[serde(rename = "signature")]
    pub eth_signature: String,
}

/// a transaction we send to move funds from Cosmos to Ethereum
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct SendToEthMsg {
    pub sender: Address,
    pub eth_dest: EthAddress,
    pub amount: Coin,
    pub bridge_fee: Coin,
}

/// This message requests that a batch be created on the Cosmos chain, this
/// may or may not actually trigger a batch to be created depending on the
/// internal batch creation rules. Said batch will be of arbitrary size also
/// depending on those rules. What this message does determine is the coin
/// type of the batch. Since all batches only move a single asset within them.
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct RequestBatchMsg {
    pub orchestrator: Address,
    pub denom: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ConfirmBatchMsg {
    pub nonce: Uint256,
    pub orchestrator: Address,
    pub token_contract: EthAddress,
    pub eth_signer: EthAddress,
    /// a hex encoded string representing the Ethereum signature
    #[serde(rename = "signature")]
    pub eth_signature: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ConfirmLogicCallMsg {
    pub invalidation_id: String,
    pub invalidation_nonce: Uint256,
    pub orchestrator: Address,
    pub eth_signer: EthAddress,
    /// a hex encoded string representing the Ethereum signature
    #[serde(rename = "signature")]
    pub eth_signature: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct WithdrawClaimMsg {
    pub event_nonce: Uint256,
    pub block_height: Uint256,
    pub batch_nonce: Uint256,
    pub token_contract: EthAddress,
    pub orchestrator: Address,
}

impl WithdrawClaimMsg {
    pub fn from_event(input: TransactionBatchExecutedEvent, sender: Address) -> Self {
        WithdrawClaimMsg {
            event_nonce: downcast_uint256(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            block_height: downcast_uint256(input.block_height)
                .expect("Block Height overflow! Bridge Halt!")
                .into(),
            batch_nonce: downcast_uint256(input.batch_nonce)
                .expect("Batch nonce overflow! Bridge halt!")
                .into(),
            token_contract: input.erc20,
            orchestrator: sender,
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct DepositClaimMsg {
    pub event_nonce: Uint256,
    pub block_height: Uint256,
    pub token_contract: EthAddress,
    pub amount: Uint256,
    pub ethereum_sender: EthAddress,
    pub cosmos_receiver: Address,
    pub orchestrator: Address,
}

impl DepositClaimMsg {
    pub fn from_event(input: SendToCosmosEvent, sender: Address) -> Self {
        DepositClaimMsg {
            event_nonce: downcast_uint256(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            block_height: downcast_uint256(input.block_height)
                .expect("Block number overflow! Bridge Halt!")
                .into(),
            amount: input.amount,
            token_contract: input.erc20,
            ethereum_sender: input.sender,
            cosmos_receiver: input.destination,
            orchestrator: sender,
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20DeployedClaimMsg {
    pub event_nonce: Uint256,
    pub block_height: Uint256,
    pub cosmos_denom: String,
    pub token_contract: EthAddress,
    pub name: String,
    pub symbol: String,
    pub decimals: Uint256,
    pub orchestrator: Address,
}

impl ERC20DeployedClaimMsg {
    pub fn from_event(input: ERC20DeployedEvent, sender: Address) -> Self {
        ERC20DeployedClaimMsg {
            event_nonce: downcast_uint256(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            block_height: downcast_uint256(input.block_height)
                .expect("Block number overflow! Bridge Halt!")
                .into(),
            cosmos_denom: input.cosmos_denom,
            token_contract: input.erc20_address,
            name: input.name,
            symbol: input.symbol,
            decimals: input.decimals.into(),
            orchestrator: sender,
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct LogicCallExecutedClaim {
    pub event_nonce: Uint256,
    pub block_height: Uint256,
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: Uint256,
    pub orchestrator: Address,
}

impl LogicCallExecutedClaim {
    pub fn from_event(input: LogicCallExecutedEvent, sender: Address) -> Self {
        LogicCallExecutedClaim {
            event_nonce: downcast_uint256(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            block_height: downcast_uint256(input.block_height)
                .expect("Block number overflow! Bridge Halt!")
                .into(),
            invalidation_nonce: downcast_uint256(input.invalidation_nonce)
                .expect("Invalidation nonce overflow! Bridge Halt!")
                .into(),
            invalidation_id: input.invalidation_id,
            orchestrator: sender,
        }
    }
}
