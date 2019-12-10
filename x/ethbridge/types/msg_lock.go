//nolint:dupl
package types

import (
	"encoding/json"
	"strconv"

	gethCommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgLock defines a message for locking coins and triggering a related event
type MsgLock struct {
	EthereumChainID  int             `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	TokenContract    EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	CosmosSender     sdk.AccAddress  `json:"cosmos_sender" yaml:"cosmos_sender"`
	EthereumReceiver EthereumAddress `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	Amount           sdk.Coins       `json:"amount" yaml:"amount"`
}

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(ethereumChainID int, tokenContract EthereumAddress, cosmosSender sdk.AccAddress, ethereumReceiver EthereumAddress, amount sdk.Coins) MsgLock {
	return MsgLock{
		EthereumChainID:  ethereumChainID,
		TokenContract:    tokenContract,
		CosmosSender:     cosmosSender,
		EthereumReceiver: ethereumReceiver,
		Amount:           amount,
	}
}

// Route should return the name of the module
func (msg MsgLock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgLock) Type() string { return "lock" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLock) ValidateBasic() sdk.Error {
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return ErrInvalidChainID(DefaultCodespace, strconv.Itoa(msg.EthereumChainID))
	}

	if msg.TokenContract.String() == "" {
		return ErrInvalidEthAddress(DefaultCodespace)
	}

	if !gethCommon.IsHexAddress(msg.TokenContract.String()) {
		return ErrInvalidEthAddress(DefaultCodespace)
	}

	if msg.CosmosSender.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosSender.String())
	}

	if msg.EthereumReceiver.String() == "" {
		return ErrInvalidEthAddress(DefaultCodespace)
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
		return ErrInvalidEthAddress(DefaultCodespace)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLock) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.CosmosSender)}
}
