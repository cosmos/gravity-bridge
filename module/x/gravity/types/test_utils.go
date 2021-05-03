package types

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateTestEthAddress() (*common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, nil
}

func GenerateTestCosmosAddress() (types.Address, error) {
	kb, err := keyring.New("keybasename", "memory", "", nil)
	if err != nil {
		return nil, err
	}

	info, _, err := kb.NewMnemonic("john", keyring.English, types2.FullFundraiserPath, hd.Secp256k1)
	if err != nil {
		return nil, err
	}

	return info.GetPubKey().Address(), nil
}

