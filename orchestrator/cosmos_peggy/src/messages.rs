use clarity::Address as EthAddress;
use deep_space::address::Address;
use deep_space::canonical_json::{to_canonical_json, CanonicalJsonError};
use deep_space::coin::Coin;
use deep_space::msg::DeepSpaceMsg;
use ethereum_peggy::utils::downcast_nonce;
use num256::Uint256;
use peggy_utils::types::{ERC20Token, SendToCosmosEvent, TransactionBatchExecutedEvent};
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

    #[serde(rename = "peggy/MsgCreateEthereumClaims")]
    CreateEthereumClaimsMsg(CreateEthereumClaimsMsg),

    #[serde(rename = "peggy/MsgDepositClaim")]
    DepositClaimMsg(DepositClaimMsg),

    #[serde(rename = "peggy/MsgWithdrawClaim")]
    WithdrawClaimMsg(WithdrawClaimMsg),
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
    pub validator: Address,
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
    pub requester: Address,
    pub denom: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ConfirmBatchMsg {
    pub nonce: Uint256,
    pub validator: Address,
    pub token_contract: EthAddress,
    pub eth_signer: EthAddress,
    /// a hex encoded string representing the Ethereum signature
    #[serde(rename = "signature")]
    pub eth_signature: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeDepositClaim {
    #[serde(rename = "nonce")]
    pub event_nonce: Uint256,
    pub erc20_token: ERC20Token,
    pub ethereum_sender: EthAddress,
    pub cosmos_receiver: Address,
}

impl EthereumBridgeDepositClaim {
    pub fn from_event(input: SendToCosmosEvent) -> Self {
        EthereumBridgeDepositClaim {
            erc20_token: ERC20Token {
                amount: input.amount,
                token_contract_address: input.erc20,
            },
            ethereum_sender: input.sender,
            cosmos_receiver: input.destination,
            event_nonce: input.event_nonce,
        }
    }
    // used for enum typing
    pub fn into_enum(self) -> EthereumBridgeClaim {
        EthereumBridgeClaim::EthereumBridgeDepositClaim(self)
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeWithdrawBatchClaim {
    pub batch_nonce: Uint256,
    pub event_nonce: Uint256,
    pub erc20_token: EthAddress,
}

impl EthereumBridgeWithdrawBatchClaim {
    pub fn from_event(input: TransactionBatchExecutedEvent) -> Self {
        EthereumBridgeWithdrawBatchClaim {
            batch_nonce: input.batch_nonce,
            event_nonce: input.event_nonce,
            erc20_token: input.erc20,
        }
    }
    // used for enum typing
    pub fn into_enum(self) -> EthereumBridgeClaim {
        EthereumBridgeClaim::EthereumBridgeWithdrawBatchClaim(self)
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeBootstrappedClaim {
    /// the claim nonce, in case multiple claims are made before one passes
    pub nonce: Uint256,
    /// the validator set in the contract being claimed
    pub allowed_validator_set: Vec<EthAddress>,
    /// the powers of the validator set in the contract being claimed
    pub validator_powers: Vec<Uint256>,
    /// the peggy ID a 32 byte unique value encoded as a hex string
    pub peggy_id: String,
    /// the amount of voting power (measured by the bridge, not cosmos) required
    /// to start the bridge, remember bridge powers are normalized to u32 max so
    /// this would be computed as some percentage of that with no bearing on what
    /// Cosmos would consider the power number to be.
    pub start_threshold: Uint256,
}

#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq, Hash)]
#[serde(tag = "type", content = "value")]
pub enum EthereumBridgeClaim {
    #[serde(rename = "peggy/DepositClaim")]
    EthereumBridgeDepositClaim(EthereumBridgeDepositClaim),
    #[serde(rename = "peggy/WithdrawClaim")]
    EthereumBridgeWithdrawBatchClaim(EthereumBridgeWithdrawBatchClaim),
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct CreateEthereumClaimsMsg {
    pub ethereum_chain_id: Uint256,
    pub bridge_contract_address: EthAddress,
    pub orchestrator: Address,
    pub deposits: Vec<EthereumBridgeClaim>,
    pub withdraws: Vec<EthereumBridgeClaim>,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct WithdrawClaimMsg {
    pub event_nonce: Uint256,
    pub batch_nonce: Uint256,
    pub token_contract: EthAddress,
    pub orchestrator: Address,
}

impl WithdrawClaimMsg {
    pub fn from_event(input: TransactionBatchExecutedEvent, sender: Address) -> Self {
        WithdrawClaimMsg {
            event_nonce: downcast_nonce(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            batch_nonce: downcast_nonce(input.batch_nonce)
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
    pub token_contract: EthAddress,
    pub amount: Uint256,
    pub ethereum_sender: EthAddress,
    pub cosmos_receiver: Address,
    pub orchestrator: Address,
}

impl DepositClaimMsg {
    pub fn from_event(input: SendToCosmosEvent, sender: Address) -> Self {
        DepositClaimMsg {
            event_nonce: downcast_nonce(input.event_nonce)
                .expect("Event nonce overflow! Bridge Halt!")
                .into(),
            amount: input.amount,
            token_contract: input.erc20,
            ethereum_sender: input.sender,
            cosmos_receiver: input.destination,
            orchestrator: sender,
        }
    }
}
