package types_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/gravity-bridge/module/app"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

func TestValidateMsgDelegateKeys(t *testing.T) {
	app.SetAddressConfig()

	var (
		ethAddress                   = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		cosmosAddress sdk.AccAddress = bytes.Repeat([]byte{0x1}, app.MaxAddrLen)
		valAddress    sdk.ValAddress = bytes.Repeat([]byte{0x1}, app.MaxAddrLen)
	)
	specs := map[string]struct {
		srcCosmosAddr sdk.AccAddress
		srcValAddr    sdk.ValAddress
		srcETHAddr    string
		expErr        bool
	}{
		"all good": {
			srcCosmosAddr: cosmosAddress,
			srcValAddr:    valAddress,
			srcETHAddr:    ethAddress,
		},
		"empty validator address": {
			srcETHAddr:    ethAddress,
			srcCosmosAddr: cosmosAddress,
			expErr:        true,
		},
		"invalid validator address": {
			srcValAddr:    []byte{0x1},
			srcCosmosAddr: cosmosAddress,
			srcETHAddr:    ethAddress,
			expErr:        true,
		},
		"empty cosmos address": {
			srcValAddr: valAddress,
			srcETHAddr: ethAddress,
			expErr:     true,
		},
		"invalid cosmos address": {
			srcCosmosAddr: []byte{0x1},
			srcValAddr:    valAddress,
			srcETHAddr:    ethAddress,
			expErr:        true,
		},
		"empty eth address": {
			srcValAddr:    valAddress,
			srcCosmosAddr: cosmosAddress,
			expErr:        true,
		},
		"invalid eth address": {
			srcValAddr:    valAddress,
			srcCosmosAddr: cosmosAddress,
			srcETHAddr:    "invalid",
			expErr:        true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			msg := types.NewMsgDelegateKeys(spec.srcValAddr, spec.srcCosmosAddr, spec.srcETHAddr, []byte{0x1})
			err := msg.ValidateBasic()
			if spec.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

}
