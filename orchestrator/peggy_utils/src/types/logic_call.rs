use super::*;
use crate::error::PeggyError;
use clarity::Address as EthAddress;
use clarity::Signature as EthSignature;
use deep_space::address::Address as CosmosAddress;

/// the response we get when querying for a valset confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct LogicCall {
    pub transfers: Vec<ERC20Token>,
    pub fees: Vec<ERC20Token>,
    pub logic_contract_address: EthAddress,
    pub payload: Vec<u8>,
    pub timeout: u64,
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: u64,
}

impl LogicCall {
    pub fn from_proto(input: peggy_proto::peggy::OutgoingLogicCall) -> Result<Self, PeggyError> {
        let mut transfers: Vec<ERC20Token> = Vec::new();
        let mut fees: Vec<ERC20Token> = Vec::new();
        for transfer in input.transfers {
            transfers.push(ERC20Token {
                amount: transfer.amount.parse()?,
                token_contract_address: transfer.contract.parse()?,
            })
        }
        for fee in input.fees {
            fees.push(ERC20Token {
                amount: fee.amount.parse()?,
                token_contract_address: fee.contract.parse()?,
            })
        }
        if transfers.is_empty() || fees.is_empty() {
            return Err(PeggyError::InvalidBridgeStateError(
                "Transaction batch containing no transactions!".to_string(),
            ));
        }

        Ok(LogicCall {
            transfers,
            fees,
            logic_contract_address: input.logic_contract_address.parse()?,
            payload: input.payload,
            timeout: input.timeout,
            invalidation_id: input.invalidation_id,
            invalidation_nonce: input.invalidation_nonce,
        })
    }
}

/// the response we get when querying for a logic call confirmation
#[derive(Serialize, Deserialize, Debug, Default, Clone)]
pub struct LogicCallConfirmResponse {
    pub invalidation_id: Vec<u8>,
    pub invalidation_nonce: u64,
    pub ethereum_signer: EthAddress,
    pub orchestrator: CosmosAddress,
    pub eth_signature: EthSignature,
}

impl LogicCallConfirmResponse {
    pub fn from_proto(input: peggy_proto::peggy::MsgConfirmLogicCall) -> Result<Self, PeggyError> {
        Ok(LogicCallConfirmResponse {
            invalidation_id: input.invalidation_id,
            invalidation_nonce: input.invalidation_nonce,
            orchestrator: input.orchestrator.parse()?,
            ethereum_signer: input.eth_signer.parse()?,
            eth_signature: input.signature.parse()?,
        })
    }
}

impl Confirm for LogicCallConfirmResponse {
    fn get_eth_address(&self) -> EthAddress {
        self.ethereum_signer
    }
    fn get_signature(&self) -> EthSignature {
        self.eth_signature.clone()
    }
}
