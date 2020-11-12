use super::*;
use clarity::{abi::Token, Address as EthAddress};
use contact::types::parse_val;
use deep_space::address::Address as CosmosAddress;
use num256::Uint256;

/// This represents an individual transaction being bridged over to Ethereum
/// parallel is the OutgoingTransferTx in x/peggy/types/batch.go
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct BatchTransaction {
    #[serde(deserialize_with = "parse_val")]
    pub txid: Uint256,
    #[serde(deserialize_with = "parse_val")]
    pub sender: CosmosAddress,
    #[serde(deserialize_with = "parse_val", rename = "dest_address")]
    pub destination: EthAddress,
    pub send: ERC20Token,
    pub bridge_fee: ERC20Token,
}
/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TransactionBatchUnparsed {
    #[serde(deserialize_with = "parse_val")]
    pub nonce: Uint256,
    pub elements: Vec<BatchTransaction>,
    pub total_fee: ERC20Token,
    pub bridged_denominator: ERC20Denominator,
    pub valset: ValsetUnparsed,
    #[serde(deserialize_with = "parse_val")]
    pub token_contract: EthAddress,
}

impl TransactionBatchUnparsed {
    pub fn convert(self) -> TransactionBatch {
        TransactionBatch {
            nonce: self.nonce,
            elements: self.elements,
            total_fee: self.total_fee,
            bridged_denominator: self.bridged_denominator,
            valset: self.valset.convert(),
            token_contract: self.token_contract,
        }
    }
}

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct TransactionBatch {
    pub nonce: Uint256,
    pub elements: Vec<BatchTransaction>,
    pub total_fee: ERC20Token,
    pub bridged_denominator: ERC20Denominator,
    pub valset: Valset,
    pub token_contract: EthAddress,
}

impl TransactionBatch {
    /// extracts the amounts, destinations and fees as submitted to the Ethereum contract
    /// and used for signatures
    pub fn get_checkpoint_values(&self) -> (Token, Token, Token) {
        let mut amounts = Vec::new();
        let mut destinations = Vec::new();
        let mut fees = Vec::new();
        for item in self.elements.iter() {
            amounts.push(Token::Uint(item.send.amount.clone()));
            fees.push(Token::Uint(item.bridge_fee.amount.clone()));
            destinations.push(item.destination)
        }
        (
            Token::Dynamic(amounts),
            destinations.into(),
            Token::Dynamic(fees),
        )
    }
}
