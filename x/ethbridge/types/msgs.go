package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	common "github.com/swishlabsco/peggy_fork/x/ethbridge/common"
)

// MsgCreateEthBridgeClaim defines a message for creating claims on the ethereum bridge
type MsgCreateEthBridgeClaim EthBridgeClaim

// NewMsgCreateEthBridgeClaim is a constructor function for MsgCreateBridgeClaim
func NewMsgCreateEthBridgeClaim(ethBridgeClaim EthBridgeClaim) MsgCreateEthBridgeClaim {
	return MsgCreateEthBridgeClaim(ethBridgeClaim)
}

// Route should return the name of the module
func (msg MsgCreateEthBridgeClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateEthBridgeClaim) Type() string { return "create_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateEthBridgeClaim) ValidateBasic() sdk.Error {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String())
	}
	if msg.Nonce < 0 {
		return ErrInvalidEthNonce(DefaultCodespace)
	}
	if !common.IsValidEthAddress(msg.EthereumSender) {
		return ErrInvalidEthAddress(DefaultCodespace)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateEthBridgeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgCreateEthBridgeClaim) GetSigners() []sdk.AccAddress {
	var validatorAddress sdk.AccAddress = sdk.AccAddress(msg.ValidatorAddress)
	return []sdk.AccAddress{validatorAddress}
}
