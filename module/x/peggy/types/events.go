package types

const (
	EventTypeObservation               = "observation"
	EventTypeOutgoingBatch             = "outgoing_batch"
	EventTypeMultisigBootstrap         = "multisig_bootstrap"
	EventTypeMultisigUpdateRequest     = "multisig_update_request"
	EventTypeOutgoingBatchCanceled     = "outgoing_batch_canceled"
	EventTypeOutgoingLogicCallCanceled = "outgoing_logic_call_canceled"
	EventTypeBridgeWithdrawalReceived  = "withdrawal_received"
	EventTypeBridgeDepositReceived     = "deposit_received"

	AttributeKeyAttestationID     = "attestation_id"
	AttributeKeyAttestationIDs    = "attestation_ids"
	AttributeKeyBatchConfirmKey   = "batch_confirm_key"
	AttributeKeyValsetConfirmKey  = "valset_confirm_key"
	AttributeKeyMultisigID        = "multisig_id"
	AttributeKeyOutgoingBatchID   = "batch_id"
	AttributeKeyOutgoingTXID      = "outgoing_tx_id"
	AttributeKeyAttestationType   = "attestation_type"
	AttributeKeyContract          = "bridge_contract"
	AttributeKeyNonce             = "nonce"
	AttributeKeyValsetNonce       = "valset_nonce"
	AttributeKeyBatchNonce        = "batch_nonce"
	AttributeKeyBridgeChainID     = "bridge_chain_id"
	AttributeKeySetOperatorAddr   = "set_operator_address"
	AttributeKeyInvalidationID    = "logic_call_invalidation_id"
	AttributeKeyInvalidationNonce = "logic_call_invalidation_nonce"
)
