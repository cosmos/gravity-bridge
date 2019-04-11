package oracle

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgMakeBridgeClaim defines a SetName message
type MsgMakeBridgeClaim struct {
	Nonce          int
	EthereumSender string
	CosmosReceiver sdk.AccAddress
	Validator      sdk.AccAddress
	Amount         sdk.Coins
}

// NewMsgMakeBridgeClaim is a constructor function for MsgMakeBridgeClaim
func NewMsgMakeBridgeClaim(nonce int, ethereumSender string, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) MsgMakeBridgeClaim {
	return MsgMakeBridgeClaim{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Validator:      validator,
		Amount:         amount,
	}
}

// Route should return the name of the module
func (msg MsgMakeBridgeClaim) Route() string { return "oracle" }

// Type should return the action
func (msg MsgMakeBridgeClaim) Type() string { return "make_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMakeBridgeClaim) ValidateBasic() sdk.Error {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String())
	}
	//must have nonce
	//amount should be nonzero
	//maybe the hacky mempool thing?
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMakeBridgeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgMakeBridgeClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}
