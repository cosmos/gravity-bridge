package types

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateGravityDenomIncorrectValues(t *testing.T) {
	err := ValidateGravityDenom("fake/denom")
	require.Errorf(t, err, "invalid denom prefix accepted")

	err = ValidateGravityDenom("gravity/fakedenom")
	require.Errorf(t, err, "invalid denom accepted")
}

func TestValidateGravityDenomCorrectValues(t *testing.T) {
	err := ValidateGravityDenom("gravity/test-denom")
	require.Errorf(t, err, "valid denom prefix not accepted")

	err = ValidateGravityDenom("testdenom")
	require.NoError(t, err, "valid existing denom not accepted")

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate new crypto address")
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Errorf("failed casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	denom := fmt.Sprintf("gravity/%s", address.Hex())
	err = ValidateGravityDenom(denom)
	require.NoError(t, err, "valid external token with eth address not accepted")
}