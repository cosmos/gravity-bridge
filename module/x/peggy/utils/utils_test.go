package utils

import (
	"encoding/hex"
	"testing"
)

func TestValsetConfirmSig(t *testing.T) {
	correctSig := "e108a7776de6b87183b0690484a74daef44aa6daf907e91abaf7bbfa426ae7706b12e0bd44ef7b0634710d99c2d81087a2f39e075158212343a3b2948ecf33d01c"
	ethAddress := "0xc783df8a850f42e7F7e57013759C285caa701eB6"
	hash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"

	hashBytes, hexErr := hex.DecodeString(hash)
	if hexErr != nil {
		panic("Hash hex decoding error")
	}
	sigBytes, hexErr := hex.DecodeString(correctSig)
	if hexErr != nil {
		panic("Signature hex decoding error")
	}
	validationErr := ValidateEthSig(hashBytes, sigBytes, ethAddress)
	if validationErr != nil {
		panic("Failed to validate signature!")
	}
}
