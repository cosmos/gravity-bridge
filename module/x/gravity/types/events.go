package types

const (
	EventTypeObservation               = "observation"
	EventTypeOutgoingBatch             = "outgoing_batch"
	EventTypeMultisigUpdateRequest     = "multisig_update_request"
	EventTypeOutgoingBatchCanceled     = "outgoing_batch_canceled"
	EventTypeContractCallTxCanceled = "outgoing_logic_call_canceled"
	EventTypeBridgeWithdrawalReceived  = "withdrawal_received"
	EventTypeBridgeDepositReceived     = "deposit_received"
	EventTypeBridgeWithdrawCanceled    = "withdraw_canceled"

	AttributeKeyEthereumEventVoteRecordID = "ethereum_event_vote_record_id"
	AttributeKeyBatchConfirmKey           = "batch_confirm_key"
	AttributeKeyEthereumSignatureKey      = "ethereum_signature_key"
	AttributeKeyMultisigID                = "multisig_id"
	AttributeKeyOutgoingBatchID           = "batch_id"
	AttributeKeyOutgoingTXID              = "outgoing_tx_id"
	AttributeKeyEthereumEventType         = "ethereum_event_type"
	AttributeKeyContract                  = "bridge_contract"
	AttributeKeyNonce                     = "nonce"
	AttributeKeyValsetNonce               = "valset_nonce"
	AttributeKeyBatchNonce                = "batch_nonce"
	AttributeKeyBridgeChainID             = "bridge_chain_id"
	AttributeKeySetOperatorAddr           = "set_operator_address"
	AttributeKeyInvalidationID            = "logic_call_invalidation_id"
	AttributeKeyInvalidationNonce         = "logic_call_invalidation_nonce"
)
