use clarity::Address as EthAddress;
use deep_space::address::Address;
use deep_space::canonical_json::{to_canonical_json, CanonicalJsonError};
use deep_space::coin::Coin;
use deep_space::msg::DeepSpaceMsg;
use num256::Uint256;
use peggy_utils::types::ERC20Token;
/// Any arbitrary message
#[derive(Serialize, Deserialize, Debug, Clone, Eq, PartialEq)]
#[serde(tag = "type", content = "value")]
pub enum PeggyMsg {
    #[serde(rename = "peggy/MsgSetEthAddress")]
    SetEthAddressMsg(SetEthAddressMsg),

    #[serde(rename = "peggy/MsgValsetRequest")]
    ValsetRequestMsg(ValsetRequestMsg),

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
}

impl DeepSpaceMsg for PeggyMsg {
    fn to_sign_bytes(&self) -> Result<Vec<u8>, CanonicalJsonError> {
        Ok(to_canonical_json(self)?)
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct SetEthAddressMsg {
    #[serde(rename = "address")]
    pub eth_address: EthAddress,
    pub validator: Address,
    /// a hex encoded string representing the Ethereum signature
    #[serde(rename = "signature")]
    pub eth_signature: String,
}
#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ValsetRequestMsg {
    pub requester: Address,
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
    pub dest_address: EthAddress,
    pub send: Coin,
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
    pub address: EthAddress,
    /// a hex encoded string representing the Ethereum signature
    #[serde(rename = "signature")]
    pub eth_signature: String,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeDepositClaim {
    /// this was added on the cosmos side due to a poorly designed interface
    /// will always be zero
    pub nonce: Uint256,
    pub erc20_token: ERC20Token,
    pub ethereum_sender: EthAddress,
    pub cosmos_receiver: Address,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeWithdrawBatchClaim {
    pub nonce: Uint256,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct EthereumBridgeMultiSigUpdateClaim {
    pub nonce: Uint256,
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
    #[serde(rename = "peggy/EthereumBridgeDepositClaim")]
    EthereumBridgeDepositClaim(EthereumBridgeDepositClaim),
    #[serde(rename = "peggy/EthereumBridgeMultiSigUpdateClaim")]
    EthereumBridgeMultiSigUpdateClaim(EthereumBridgeMultiSigUpdateClaim),
    #[serde(rename = "peggy/EthereumBridgeWithdrawBatchClaim")]
    EthereumBridgeWithdrawBatchClaim(EthereumBridgeWithdrawBatchClaim),
    #[serde(rename = "peggy/EthereumBridgeBootstrappedClaim")]
    EthereumBridgeBootstrappedClaim(EthereumBridgeBootstrappedClaim),
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct CreateEthereumClaimsMsg {
    pub ethereum_chain_id: Uint256,
    pub bridge_contract_address: EthAddress,
    pub orchestrator: Address,
    pub claims: Vec<EthereumBridgeClaim>,
}
