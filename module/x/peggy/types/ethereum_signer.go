package types

import (
	"crypto/ecdsa"
	"encoding/hex"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const signaturePrefix = "\x19Ethereum Signed Message:\n32"

func NewEthereumSignature(hash []byte, privateKey *ecdsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", sdkerrors.Wrap(ErrEmpty, "private key")
	}
	protectedHash := crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...))
	bytes, err := crypto.Sign(protectedHash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ValidateEthereumSignature(hash []byte, signature []byte, ethAddress string) error {
	if len(signature) < 65 {
		return sdkerrors.Wrap(ErrInvalid, "signature too short")
	}
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
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	protectedHash := crypto.Keccak256Hash(append([]uint8(signaturePrefix), hash...))

	pubkey, err := crypto.SigToPub(protectedHash.Bytes(), signature)
	if err != nil {
		return sdkerrors.Wrap(err, "signature to public key")
	}

	addr := crypto.PubkeyToAddress(*pubkey)

	if addr.Hex() != ethAddress {
		return sdkerrors.Wrap(ErrInvalid, "signature not matching")
	}

	return nil
}
