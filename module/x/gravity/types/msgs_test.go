package types

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
)

const AddrLen = 20

func TestValidateMsgDelegateKeys(t *testing.T) {

	sdk.GetConfig().SetAddressVerifier(VerifyAddressFormat)
	var (
		ethAddress                   = "0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"
		cosmosAddress sdk.AccAddress = bytes.Repeat([]byte{0x1}, AddrLen)
		valAddress    sdk.ValAddress = bytes.Repeat([]byte{0x1}, AddrLen)
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
	for test_msg, spec := range specs {
		t.Run(test_msg, func(t *testing.T) {
			fmt.Println(test_msg)
			msg := NewMsgDelegateKeys(spec.srcValAddr, spec.srcCosmosAddr, spec.srcETHAddr)
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

func VerifyAddressFormat(bz []byte) error {

	if len(bz) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
	}

	if len(bz) > address.MaxAddrLen {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bz))
	}
	if len(bz) != 20 {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address length is %d, got %d", 20, len(bz))

	}

	return nil
}
