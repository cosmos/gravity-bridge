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

	// for backwards compatibility reasons  the V value of an Ethereum sig is presented
	// as 27 or 28, interally though it should be a 0-3 value due to changed formats.
	// It seems that go-ethereum expects this to be done before sigs actually reach it's
	// internal validation functions. In order to comply with this requirement we check
	// the sig an dif it's in standard format we correct it. If it's in go-ethereum's expected
	// format already we make no changes.
	//
	// We could attempt to break or otherwise exit early on obviously invalid values for this
	// byte, but that's a task best left to go-ethereum
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

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
