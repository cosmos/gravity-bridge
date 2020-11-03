use clarity::Address as EthAddress;
use contact::types::parse_val;
use num256::Uint256;

mod batches;
mod ethereum_events;
mod signatures;
mod valsets;

pub use batches::*;
pub use ethereum_events::*;
pub use signatures::*;
pub use valsets::*;

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20Token {
    #[serde(deserialize_with = "parse_val")]
    pub amount: Uint256,
    pub symbol: String,
    #[serde(deserialize_with = "parse_val")]
    pub token_contract_address: EthAddress,
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20Denominator {
    #[serde(deserialize_with = "parse_val")]
    pub token_contract_address: EthAddress,
    pub symbol: String,
    pub cosmos_voucher_denom: String,
}
