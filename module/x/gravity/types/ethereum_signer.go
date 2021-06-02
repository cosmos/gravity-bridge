package types

import (
	"crypto/ecdsa"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	signaturePrefix = "\x19Ethereum Signed Message:\n32"
)

// NewEthereumSignature creates a new signuature over a given byte array
func NewEthereumSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, sdkerrors.Wrap(ErrInvalid, "did not pass in private key")
	}
	protectedHash := crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...))
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

// ValidateEthereumSignature takes a message, an associated signature and public key and
// returns an error if the signature isn't valid
func ValidateEthereumSignature(hash []byte, signature []byte, ethAddress common.Address) error {
	if len(signature) < 65 {
		return sdkerrors.Wrap(ErrInvalid, "signature too short")
	}

	// Copy to avoid mutating signature slice by accident
	var sigCopy = make([]byte, len(signature))
	copy(sigCopy, signature)

	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.

	// for backwards compatibility reasons  the V value of an Ethereum sig is presented
	// as 27 or 28, internally though it should be a 0-3 value due to changed formats.
	// It seems that go-ethereum expects this to be done before sigs actually reach it's
	// internal validation functions. In order to comply with this requirement we check
	// the sig an dif it's in standard format we correct it. If it's in go-ethereum's expected
	// format already we make no changes.
	//
	// We could attempt to break or otherwise exit early on obviously invalid values for this
	// byte, but that's a task best left to go-ethereum
	if sigCopy[64] == 27 || sigCopy[64] == 28 {
		sigCopy[64] -= 27
	}

	pubkey, err := crypto.SigToPub(crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...)).Bytes(), sigCopy)
	if err != nil {
		return sdkerrors.Wrap(err, "signature to public key")
	}

	if addr := crypto.PubkeyToAddress(*pubkey); addr != ethAddress {
		return sdkerrors.Wrap(ErrInvalid, "signature not matching")
	}

	return nil
}
