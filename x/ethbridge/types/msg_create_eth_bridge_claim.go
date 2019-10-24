package types

import (
	"encoding/json"
	"fmt"

	gethCommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	if msg.ValidatorAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.ValidatorAddress.String())
	}

	if msg.Nonce < 0 {
		return ErrInvalidEthNonce(DefaultCodespace)
	}

	if !gethCommon.IsHexAddress(msg.EthereumSender.String()) {
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
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// MapOracleClaimsToEthBridgeClaims maps a set of generic oracle claim data into EthBridgeClaim objects
func MapOracleClaimsToEthBridgeClaims(nonce int, ethereumSender EthereumAddress, oracleValidatorClaims map[string]string, f func(int, EthereumAddress, sdk.ValAddress, string) (EthBridgeClaim, sdk.Error)) ([]EthBridgeClaim, sdk.Error) {
	mappedClaims := make([]EthBridgeClaim, len(oracleValidatorClaims))
	i := 0
	for validatorBech32, validatorClaim := range oracleValidatorClaims {
		validatorAddress, parseErr := sdk.ValAddressFromBech32(validatorBech32)
		if parseErr != nil {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", parseErr))
		}
		mappedClaim, err := f(nonce, ethereumSender, validatorAddress, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}
	return mappedClaims, nil
}
