package types

const (
	EventTypeObservation              = "observation"
	EventTypeOutgoingBatch            = "outgoing_batch"
	EventTypeOutgoingBatchCanceled    = "outgoing_batch_canceled"
	EventTypeBridgeWithdrawalReceived = "withdrawal_received"
	EventTypeBridgeDepositReceived    = "deposit_received"
)
const (
	AttributeKeyAttestationID   = "attestation_id"
	AttributeKeyOutgoingBatchID = "batch_id"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeKeyAttestationType = "attestation_type"
	AttributeKeyContract        = "bridge_contract"
	AttributeKeyNonce           = "nonce"
	AttributeKeyBridgeChainID   = "bridge_chain_id"
)
