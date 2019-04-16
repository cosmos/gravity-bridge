package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgMakeBridgeEthClaim defines a SetName message
type MsgMakeBridgeEthClaim struct {
	Nonce          int
	EthereumSender string
	CosmosReceiver sdk.AccAddress
	Validator      sdk.AccAddress
	Amount         sdk.Coins
}

// NewMsgMakeEthBridgeClaim is a constructor function for MsgMakeBridgeClaim
func NewMsgMakeEthBridgeClaim(nonce int, ethereumSender string, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) MsgMakeBridgeEthClaim {
	return MsgMakeBridgeEthClaim{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Validator:      validator,
		Amount:         amount,
	}
}

// Route should return the name of the module
func (msg MsgMakeBridgeEthClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMakeBridgeEthClaim) Type() string { return "make_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMakeBridgeEthClaim) ValidateBasic() sdk.Error {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String())
	}
	//TODO: must have nonce/identifier
	//TODO: amount should be nonzero
	//TODO: investigate maybe the hacky mempool thing for offchain signature aggregation?
	//TODO: Check signer is in fact a validator (also work out if this check should be done here or in getsigners or in the handler?)
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMakeBridgeEthClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgMakeBridgeEthClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Validator}
}
