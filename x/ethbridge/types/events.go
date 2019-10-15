package types

// Ethbridge module event types
var (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"
	EventTypeBurn           = "burn"
	EventTypeLock           = "lock"

	AttributeKeyEthereumSender = "ethereum_sender"
	AttributeKeyCosmosReceiver = "cosmos_receiver"
	AttributeKeyAmount         = "amount"
	AttributeKeyStatus         = "status"
	AttributeKeyClaimType      = "claim_type"

	AttributeKeyCosmosSender     = "cosmos_sender"
	AttributeKeyEthereumReceiver = "ethereum_receiver"

	AttributeValueCategory = ModuleName
)
