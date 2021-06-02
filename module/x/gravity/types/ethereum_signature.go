package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ EthereumTxConfirmation = &SignerSetTxConfirmation{}
	_ EthereumTxConfirmation = &ContractCallTxConfirmation{}
	_ EthereumTxConfirmation = &BatchTxConfirmation{}
)

///////////////
// GetSigner //
///////////////

func (u *SignerSetTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *ContractCallTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *BatchTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

///////////////////
// GetStoreIndex //
///////////////////

func (sstx *SignerSetTxConfirmation) GetStoreIndex() []byte {
	return MakeSignerSetTxKey(sstx.SignerSetNonce)
}

func (btx *BatchTxConfirmation) GetStoreIndex() []byte {
	return MakeBatchTxKey(common.HexToAddress(btx.TokenContract), btx.BatchNonce)
}

func (cctx *ContractCallTxConfirmation) GetStoreIndex() []byte {
	return MakeContractCallTxKey(cctx.InvalidationScope, cctx.InvalidationNonce)
}

//////////////
// Validate //
//////////////

func (u *SignerSetTxConfirmation) Validate() error {
	if u.SignerSetNonce == 0 {
		return fmt.Errorf("nonce must be set")
	}
	if !common.IsHexAddress(u.EthereumSigner) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}

func (u *ContractCallTxConfirmation) Validate() error {
	if u.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce must be set")
	}
	if u.InvalidationScope == nil {
		return fmt.Errorf("invalidation scope must be set")
	}
	if !common.IsHexAddress(u.EthereumSigner) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}

func (u *BatchTxConfirmation) Validate() error {
	if u.BatchNonce == 0 {
		return fmt.Errorf("nonce must be set")
	}
	if !common.IsHexAddress(u.TokenContract) {
		return fmt.Errorf("token contract address must be valid ethereum address")
	}
	if !common.IsHexAddress(u.EthereumSigner) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}
