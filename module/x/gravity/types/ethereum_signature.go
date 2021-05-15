package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	_ EthereumSignature = &SignerSetTxSignature{}
	_ EthereumSignature = &ContractCallTxSignature{}
	_ EthereumSignature = &BatchTxSignature{}
)

///////////////
// GetSigner //
///////////////

func (u *SignerSetTxSignature) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *ContractCallTxSignature) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *BatchTxSignature) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

///////////////////
// GetStoreIndex //
///////////////////

func (sstx *SignerSetTxSignature) GetStoreIndex() []byte {
	return MakeSignerSetTxKey(sstx.Nonce)
}

func (btx *BatchTxSignature) GetStoreIndex() []byte {
	return MakeBatchTxKey(common.HexToAddress(btx.TokenContract), btx.Nonce)
}

func (cctx *ContractCallTxSignature) GetStoreIndex() []byte {
	return MakeContractCallTxKey(cctx.InvalidationScope, cctx.InvalidationNonce)
}

//////////////
// Validate //
//////////////

func (u *SignerSetTxSignature) Validate() error {
	if !(u.Nonce > 0) {
		return fmt.Errorf("nonce must be set")
	}
	if u.EthereumSigner == "" {
		return fmt.Errorf("ethereum signer must be set")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}

func (u *ContractCallTxSignature) Validate() error {
	if !(u.InvalidationNonce > 0) {
		return fmt.Errorf("invalidation nonce must be set")
	}
	if u.InvalidationScope == nil {
		return fmt.Errorf("invalidation scope must be set")
	}
	if u.EthereumSigner == "" {
		return fmt.Errorf("ethereum signer must be set")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}

func (u *BatchTxSignature) Validate() error {
	if !(u.Nonce > 0) {
		return fmt.Errorf("nonce must be set")
	}
	if !common.IsHexAddress(u.TokenContract) {
		return fmt.Errorf("token contract address must be valid ethereum address")
	}
	if u.EthereumSigner == "" {
		return fmt.Errorf("ethereum signer must be set")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	return nil
}
