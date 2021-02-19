use clarity::Address as EthAddress;
use contact::types::parse_val;
use num256::Uint256;
mod batches;
mod ethereum_events;
mod logic_call;
mod signatures;
mod valsets;
use crate::error::PeggyError;

pub use batches::*;
pub use ethereum_events::*;
pub use logic_call::*;
pub use signatures::*;
pub use valsets::*;

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20Token {
    pub amount: Uint256,
    #[serde(rename = "contract")]
    pub token_contract_address: EthAddress,
}

impl ERC20Token {
    pub fn from_proto(input: peggy_proto::peggy::Erc20Token) -> Result<Self, PeggyError> {
        Ok(ERC20Token {
            amount: input.amount.parse()?,
            token_contract_address: input.contract.parse()?,
        })
    }
}

#[derive(Serialize, Deserialize, Debug, Default, Clone, Eq, PartialEq, Hash)]
pub struct ERC20Denominator {
    #[serde(deserialize_with = "parse_val")]
    pub token_contract_address: EthAddress,
    pub symbol: String,
    pub cosmos_voucher_denom: String,
}
