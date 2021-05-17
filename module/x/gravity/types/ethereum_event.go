package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ EthereumEvent = &SendToCosmosEvent{}
	_ EthereumEvent = &BatchExecutedEvent{}
	_ EthereumEvent = &ContractCallExecutedEvent{}
	_ EthereumEvent = &ERC20DeployedEvent{}
	_ EthereumEvent = &SignerSetTxExecutedEvent{}
)

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *EthereumEventVoteRecord) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var event EthereumEvent
	return unpacker.UnpackAny(m.Event, &event)
}

//////////////
// GetNonce //
//////////////

func (stce *SendToCosmosEvent) GetNonce() uint64 {
	return stce.EventNonce
}

func (bee *BatchExecutedEvent) GetNonce() uint64 {
	return bee.EventNonce
}

func (ccee *ContractCallExecutedEvent) GetNonce() uint64 {
	return ccee.EventNonce
}

func (e20de *ERC20DeployedEvent) GetNonce() uint64 {
	return e20de.EventNonce
}

func (sse *SignerSetTxExecutedEvent) GetNonce() uint64 {
	return sse.EventNonce
}

//////////
// Hash //
//////////

// TODO: rewrite all hashes to match this one
func (stce *SendToCosmosEvent) Hash() tmbytes.HexBytes {
	rcv, _ := sdk.AccAddressFromBech32(stce.CosmosReceiver)
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(stce.EventNonce),
			stce.Amount.BigInt().Bytes(),
			common.Hex2Bytes(stce.TokenContract),
			common.Hex2Bytes(stce.EthereumSender),
			rcv.Bytes(),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (bee *BatchExecutedEvent) Hash() tmbytes.HexBytes {
	path := append(common.HexToAddress(bee.TokenContract).Bytes(), sdk.Uint64ToBigEndian(bee.EventNonce)...)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (ccee *ContractCallExecutedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			ccee.InvalidationId,
			sdk.Uint64ToBigEndian(ccee.InvalidationNonce),
			sdk.Uint64ToBigEndian(ccee.EventNonce),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (e20de *ERC20DeployedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			[]byte(e20de.CosmosDenom),
			common.HexToAddress(e20de.TokenContract).Bytes(),
			[]byte(e20de.Erc20Name),
			[]byte(e20de.Erc20Symbol),
			sdk.Uint64ToBigEndian(e20de.Erc20Decimals),
			sdk.Uint64ToBigEndian(e20de.EventNonce),
		},
		[]byte{},
	)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (sse *SignerSetTxExecutedEvent) Hash() tmbytes.HexBytes {
	path := bytes.Join(
		[][]byte{
			sdk.Uint64ToBigEndian(sse.SignerSetTxNonce),
			sdk.Uint64ToBigEndian(sse.EventNonce),
			EthereumSigners(sse.Members).Hash(),
		},
		[]byte{},
	)
	hash := sha256.Sum256(([]byte(path)))
	return hash[:]
}

//////////////
// Validate //
//////////////

func (stce *SendToCosmosEvent) Validate() error {
	if stce.EventNonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if !common.IsHexAddress(stce.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum contract address")
	}
	if !stce.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}
	if !common.IsHexAddress(stce.EthereumSender) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum sender")
	}
	if _, err := sdk.AccAddressFromBech32(stce.CosmosReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, stce.CosmosReceiver)
	}
	if stce.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

func (bee *BatchExecutedEvent) Validate() error {
	if bee.EventNonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if !common.IsHexAddress(bee.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum contract address")
	}
	if bee.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

func (ccee *ContractCallExecutedEvent) Validate() error {
	if ccee.EventNonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if len(ccee.InvalidationId) == 0 {
		return fmt.Errorf("invalidation id cannot be empty")
	}
	if ccee.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce cannot be 0")
	}
	if ccee.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

func (e20de *ERC20DeployedEvent) Validate() error {
	if e20de.EventNonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if !common.IsHexAddress(e20de.TokenContract) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum contract address")
	}
	if err := sdk.ValidateDenom(e20de.CosmosDenom); err != nil {
		return err
	}
	if strings.TrimSpace(e20de.Erc20Name) == "" {
		return fmt.Errorf("token name cannot be blank")
	}
	if strings.TrimSpace(e20de.Erc20Symbol) == "" {
		return fmt.Errorf("token symbol cannot be blank")
	}
	if e20de.Erc20Decimals <= 0 {
		return fmt.Errorf("decimal precision must be positive")
	}
	if e20de.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	return nil
}

func (sse *SignerSetTxExecutedEvent) Validate() error {
	if sse.EventNonce == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if sse.GetSignerSetTxNonce() == 0 {
		return fmt.Errorf("nonce cannot be 0")
	}
	if sse.EthereumHeight == 0 {
		return fmt.Errorf("ethereum height cannot be 0")
	}
	if sse.Members == nil {
		return fmt.Errorf("members cannot be nil")
	}
	return nil
}
