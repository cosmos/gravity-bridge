package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateMsgSetEthAddress(t *testing.T) {
	var (
		ethAddress                   = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		cosmosAddress sdk.AccAddress = bytes.Repeat([]byte{0x1}, sdk.AddrLen)
		// privKeyString = "0xb8662f35f9de8720424e82b232e8c98d15399490adae9ca993f5ef1dc4883690"
		correctSig = "46402c54b2a13f229560c5406db56fbd9b307a32ca31997955498f0df99f97cb471e8bdeb927551cbbc4d548a7739b5782c918ff9d56eed03f86b29a4bc722c400"
	)
	specs := map[string]struct {
		srcCosmosAddr sdk.AccAddress
		srcSignature  string
		srcETHAddr    string
		expErr        bool
	}{
		"all good": {
			srcCosmosAddr: cosmosAddress,
			srcSignature:  correctSig,
			srcETHAddr:    ethAddress,
		},
		"empty signature": {
			srcCosmosAddr: cosmosAddress,
			srcETHAddr:    ethAddress,
			expErr:        true,
		},
		"invalid signature": {
			srcCosmosAddr: cosmosAddress,
			srcSignature:  correctSig[2:],
			srcETHAddr:    ethAddress,
			expErr:        true,
		},
		"empty cosmos address": {
			srcSignature: correctSig,
			srcETHAddr:   ethAddress,
			expErr:       true,
		},
		"invalid cosmos address": {
			srcCosmosAddr: []byte{0x1},
			srcSignature:  correctSig,
			srcETHAddr:    ethAddress,
			expErr:        true,
		},
		"empty eth address": {
			srcCosmosAddr: cosmosAddress,
			srcSignature:  correctSig,
			expErr:        true,
		},
		"invalid eth address": {
			srcCosmosAddr: cosmosAddress,
			srcSignature:  correctSig,
			srcETHAddr:    "invalid",
			expErr:        true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			msg := NewMsgSetEthAddress(spec.srcETHAddr, spec.srcCosmosAddr, spec.srcSignature)
			// when
			err := msg.ValidateBasic()
			if spec.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}

}
