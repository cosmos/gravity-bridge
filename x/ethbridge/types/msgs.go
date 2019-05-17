package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/common"
)

// MsgMakeEthBridgeClaim defines a message for creating claims on the ethereum bridge
type MsgMakeEthBridgeClaim struct {
	EthBridgeClaim `json:"eth_bridge_claim"`
}

// NewMsgMakeEthBridgeClaim is a constructor function for MsgMakeBridgeClaim
func NewMsgMakeEthBridgeClaim(ethBridgeClaim EthBridgeClaim) MsgMakeEthBridgeClaim {
	return MsgMakeEthBridgeClaim{ethBridgeClaim}
}

// Route should return the name of the module
func (msg MsgMakeEthBridgeClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMakeEthBridgeClaim) Type() string { return "make_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMakeEthBridgeClaim) ValidateBasic() sdk.Error {
	if msg.EthBridgeClaim.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String())
	}
	if msg.EthBridgeClaim.Nonce < 0 {
		return ErrInvalidEthNonce(DefaultCodespace)
	}
	if !common.IsValidEthAddress(msg.EthBridgeClaim.EthereumSender) {
		return ErrInvalidEthAddress(DefaultCodespace)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMakeEthBridgeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgMakeEthBridgeClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.EthBridgeClaim.Validator}
}
