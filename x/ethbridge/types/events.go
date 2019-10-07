package types

// Ethbridge module event types
var (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"

	AttributeKeyEthereumSender = "ethereum_sender"
	AttributeKeyCosmosReceiver = "cosmos_receiver"
	AttributeKeyAmount         = "amount"
	AttributeKeyStatus         = "status"

	AttributeValueCategory = ModuleName
)
