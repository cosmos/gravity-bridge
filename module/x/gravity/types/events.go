package types

const (
	EventTypeOutgoingBatch             = "outgoing_batch"
	EventTypeMultisigUpdateRequest     = "multisig_update_request"
	EventTypeOutgoingBatchCanceled     = "outgoing_batch_canceled"
	EventTypeOutgoingLogicCallCanceled = "outgoing_logic_call_canceled"

	AttributeKeyAttestationID    = "attestation_id"
	AttributeKeyBatchConfirmKey  = "batch_confirm_key"
	AttributeKeyValsetConfirmKey = "valset_confirm_key"
	AttributeKeyTxID             = "tx_id"
	AttributeKeyNonce            = "nonce"
	AttributeKeyValsetNonce      = "valset_nonce"
	AttributeKeySetOperatorAddr  = "set_operator_address"
)
