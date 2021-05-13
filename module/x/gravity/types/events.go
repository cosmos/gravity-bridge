package types

// TODO: check strings and check for dead code
const (
	EventTypeObservation              = "observation"
	EventTypeBatchTx                  = "batch_tx"
	EventTypeMultisigUpdateRequest    = "multisig_update_request"
	EventTypeBatchTxCanceled          = "batch_tx_canceled"
	EventTypeContractCallTxCanceled   = "contract_call_canceled"
	EventTypeBridgeWithdrawalReceived = "withdrawal_received"
	EventTypeBridgeDepositReceived    = "deposit_received"
	EventTypeBridgeWithdrawCanceled   = "withdraw_canceled"

	AttributeKeyEthereumEventVoteRecordID   = "ethereumEventVoteRecord_id"
	AttributeKeyBatchTxSignatureKey         = "batch_tx_signature_key"
	AttributeKeySignerSetTxSignatureKey     = "signer_set_tx_signature_key"
	AttributeKeyMultisigID                  = "multisig_id"
	AttributeKeyBatchTxID                   = "batch_id"
	AttributeKeySendToEthereumID            = "send_to_ethereum_id"
	AttributeKeyEthereumEventVoteRecordType = "ethereumEventVoteRecord_type"
	AttributeKeyContract                    = "bridge_contract"
	AttributeKeyNonce                       = "nonce"
	AttributeKeySignerSetTxNonce            = "signer_set_nonce"
	AttributeKeyBatchNonce                  = "batch_nonce"
	AttributeKeyBridgeChainID               = "bridge_chain_id"
	AttributeKeySetOperatorAddr             = "set_operator_address"
	AttributeKeyInvalidationID              = "contract_call_invalidation_id"
	AttributeKeyInvalidationNonce           = "contract_call_invalidation_nonce"
	AttributeKeyBadEthSignature             = "bad_eth_signature"
	AttributeKeyBadEthSignatureSubject      = "bad_eth_signature_subject"
)
