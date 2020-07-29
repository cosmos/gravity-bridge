package utils

import (
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

func ValidateEthSig(hash []byte, signature []byte, ethAddress string) error {
	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.
	pubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return err
	}

	addr := crypto.PubkeyToAddress(*pubkey)

	if addr.Hex() != ethAddress {
		return errors.New("Signature is not valid")
	}

	return nil
}
