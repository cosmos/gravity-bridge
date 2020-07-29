package utils

import "github.com/ethereum/go-ethereum/crypto"

func ValidateEthSig(hash []byte, signature []byte, ethAddress string) (bool, error) {
	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.
	pubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return false, err
	}

	addr := crypto.PubkeyToAddress(*pubkey)

	if addr.Hex() != ethAddress {
		return false, nil
	}

	return true, nil
}
