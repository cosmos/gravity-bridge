package types

const (
	EventTypeOutgoingBatch             = "outgoing_batch"
	EventTypeMultisigUpdateRequest     = "multisig_update_request"
	EventTypeOutgoingBatchCanceled     = "outgoing_batch_canceled"
	EventTypeOutgoingLogicCallCanceled = "outgoing_logic_call_canceled"
	EventTypeTransferPooled            = "transfer_pooled"
	EventTypeTransferCanceled          = "transfer_canceled"

	AttributeKeyAttestationID         = "attestation_id"
	AttributeKeyEventID               = "event_id"
	AttributeKeyConfirmType           = "confirm_type"
	AttributeKeyTxID                  = "tx_id"
	AttributeKeyEthRecipient          = "eth_recipient"
	AttributeKeyDenom                 = "denom"
	AttributeKeyTokenContract         = "token_contract"
	AttributeKeyNonce                 = "nonce"
	AttributeKeyOrchestratorValidator = "orchestrator_validator"
	AttributeKeyValsetNonce           = "valset_nonce"
	AttributeKeySetOperatorAddr       = "set_operator_address"
	AttributeKeyInvalidationID        = "invalidation_id"
	AttributeKeyInvalidationNonce     = "invalidation_nonce"
)
