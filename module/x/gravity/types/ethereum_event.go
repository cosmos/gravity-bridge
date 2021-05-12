package types

import (
	"crypto/sha256"
	"fmt"
	"strings"

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
)

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

//////////
// Hash //
//////////

func (stce *SendToCosmosEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%s/%s/", stce.TokenContract, stce.EthereumSender, stce.CosmosReceiver)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (bee *BatchExecutedEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%d/", bee.TokenContract, bee.EventNonce)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (ccee *ContractCallExecutedEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%d/", ccee.InvalidationId, ccee.InvalidationNonce)
	hash := sha256.Sum256([]byte(path))
	return hash[:]
}

func (e20de *ERC20DeployedEvent) Hash() tmbytes.HexBytes {
	path := fmt.Sprintf("%s/%s/%s/%s/%d/", e20de.CosmosDenom, e20de.TokenContract, e20de.Erc20Name, e20de.Erc20Symbol, e20de.Erc20Decimals)
	hash := sha256.Sum256([]byte(path))
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
