package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	gethCommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
)

// MsgLock defines a message for locking coins and triggering a related event
type MsgLockNFT struct {
	EthereumChainID  int                       `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	TokenContract    ethbridge.EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	CosmosSender     sdk.AccAddress            `json:"cosmos_sender" yaml:"cosmos_sender"`
	EthereumReceiver ethbridge.EthereumAddress `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	Denom            string                    `json:"denom" yaml:"denom"`
	ID               string                    `json:"id" yaml:"id"`
}

// NewMsgLock is a constructor function for MsgLock
func NewMsgLockNFT(ethereumChainID int, tokenContract ethbridge.EthereumAddress, cosmosSender sdk.AccAddress, ethereumReceiver ethbridge.EthereumAddress, denom, id string) MsgLockNFT {
	return MsgLockNFT{
		EthereumChainID:  ethereumChainID,
		TokenContract:    tokenContract,
		CosmosSender:     cosmosSender,
		EthereumReceiver: ethereumReceiver,
		Denom:            denom,
		ID:               id,
	}
}

// Route should return the name of the module
func (msg MsgLockNFT) Route() string { return RouterKey }

// Type should return the action
func (msg MsgLockNFT) Type() string { return "lock_nft" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLockNFT) ValidateBasic() error {
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainID)
	}

	if msg.TokenContract.String() == "" {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.TokenContract.String()) {
		return ErrInvalidEthAddress
	}

	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}

	if msg.EthereumReceiver.String() == "" {
		return ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
		return ErrInvalidEthAddress
	}

	if msg.Denom == "" {
		return ErrInvalidDenom
	}

	if msg.ID == "" {
		return ErrInvalidID
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLockNFT) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgLockNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.CosmosSender)}
}

// MsgBurnNFT defines a message for burning coins and triggering a related event
type MsgBurnNFT struct {
	EthereumChainID  int                       `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	TokenContract    ethbridge.EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	CosmosSender     sdk.AccAddress            `json:"cosmos_sender" yaml:"cosmos_sender"`
	EthereumReceiver ethbridge.EthereumAddress `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	Denom            string                    `json:"denom" yaml:"denom"`
	ID               string                    `json:"id" yaml:"id"`
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(ethereumChainID int, tokenContract ethbridge.EthereumAddress, cosmosSender sdk.AccAddress, ethereumReceiver ethbridge.EthereumAddress, denom, id string) MsgBurnNFT {
	return MsgBurnNFT{
		EthereumChainID:  ethereumChainID,
		TokenContract:    tokenContract,
		CosmosSender:     cosmosSender,
		EthereumReceiver: ethereumReceiver,
		Denom:            denom,
		ID:               id,
	}
}

// Route should return the name of the module
func (msg MsgBurnNFT) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBurnNFT) Type() string { return "burn_nft" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBurnNFT) ValidateBasic() error {
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return sdkerrors.Wrapf(ErrInvalidEthereumChainID, "%d", msg.EthereumChainID)
	}
	if msg.TokenContract.String() == "" {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.TokenContract.String()) {
		return ErrInvalidEthAddress
	}
	if msg.CosmosSender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosSender.String())
	}
	if msg.EthereumReceiver.String() == "" {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
		return ErrInvalidEthAddress
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBurnNFT) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.CosmosSender)}
}

// MsgCreateNFTBridgeClaim defines a message for creating claims on the ethereum bridge
type MsgCreateNFTBridgeClaim NFTBridgeClaim

// NewMsgCreateNFTBridgeClaim is a constructor function for MsgCreateBridgeClaim
func NewMsgCreateNFTBridgeClaim(nftBridgeClaim NFTBridgeClaim) MsgCreateNFTBridgeClaim {
	return MsgCreateNFTBridgeClaim(nftBridgeClaim)
}

// Route should return the name of the module
func (msg MsgCreateNFTBridgeClaim) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateNFTBridgeClaim) Type() string { return "create_nft_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateNFTBridgeClaim) ValidateBasic() error {
	if msg.CosmosReceiver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.CosmosReceiver.String())
	}

	if msg.ValidatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress.String())
	}

	if msg.Nonce < 0 {
		return ErrInvalidEthNonce
	}

	if !gethCommon.IsHexAddress(msg.EthereumSender.String()) {
		return ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.BridgeContractAddress.String()) {
		return ErrInvalidEthAddress
	}
	if msg.TokenContractAddress == ethbridge.NewEthereumAddress("0x0000000000000000000000000000000000000000") {
		return ErrInvalidTokenAddress
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateNFTBridgeClaim) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgCreateNFTBridgeClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// MapOracleClaimsToNFTBridgeClaims maps a set of generic oracle claim data into NFTBridgeClaim objects
func MapOracleClaimsToNFTBridgeClaims(ethereumChainID int, bridgeContract ethbridge.EthereumAddress, nonce int, symbol string, tokenContract ethbridge.EthereumAddress, ethereumSender ethbridge.EthereumAddress, oracleValidatorClaims map[string]string, f func(int, ethbridge.EthereumAddress, int, string, ethbridge.EthereumAddress, ethbridge.EthereumAddress, sdk.ValAddress, string) (NFTBridgeClaim, error)) ([]NFTBridgeClaim, error) {
	mappedClaims := make([]NFTBridgeClaim, len(oracleValidatorClaims))
	i := 0
	for validatorBech32, validatorClaim := range oracleValidatorClaims {
		validatorAddress, parseErr := sdk.ValAddressFromBech32(validatorBech32)
		if parseErr != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("failed to parse claim: %s", parseErr))
		}
		mappedClaim, err := f(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender, validatorAddress, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}
	return mappedClaims, nil
}
