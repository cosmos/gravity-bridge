package types

// Ethbridge module event types
var (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"
	EventTypeBurn           = "burn"

	AttributeKeyEthereumSender = "ethereum_sender"
	AttributeKeyCosmosReceiver = "cosmos_receiver"
	AttributeKeyAmount         = "amount"
	AttributeKeyStatus         = "status"

	AttributeKeyCosmosSender     = "cosmos_sender"
	AttributeKeyEthereumReceiver = "ethereum_receiver"

	AttributeValueCategory = ModuleName
)
