package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func newCosmosAddress() (crypto.Address, error) {
	kb, err := keyring.New("keybasename", "memory", "", nil)
	if err != nil {
		return nil, err
	}

	info, _, err := kb.NewMnemonic("john", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	if err != nil {
		return nil, err
	}

	return info.GetPubKey().Address(), nil
}

func TestMsgDelegateKey_ValidateBasic(t *testing.T) {
	valCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	valAddr, err := sdk.ValAddressFromHex(valCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to validator address")
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")

	testCases := []struct {
		name string
		val string
		orch string
		eth string
		expError bool
	}{
		{"valid input", valAddr.String(), orchAddr.String(), ethAddr.String(), false},
		{"invalid eth addr", valAddr.String(), orchAddr.String(), "not an addr", true},
		{"invalid orchestrator", valAddr.String(), "not an addr", ethAddr.String(), true},
		{"invalid validator", "not an addr", orchAddr.String(), ethAddr.String(), true},
	}

	for _, tc := range testCases {
		mdk := MsgDelegateKey{tc.val, tc.orch, tc.eth}
		err := mdk.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestMsgTransfer_ValidateBasic(t *testing.T) {
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")

	testCases := []struct{
		name string
		sender string
		eth string
		amount sdk.Coin
		bridgeFee sdk.Coin
		expError bool
	}{
		{"valid input", orchAddr.String(), ethAddr.String(),
			sdk.Coin{"testcoin", sdk.NewInt(10)}, sdk.Coin{"testcoin", sdk.NewInt(2)}, false},
		{"no sender", "not an addr", ethAddr.String(),
			sdk.Coin{"testcoin", sdk.NewInt(10)}, sdk.Coin{"testcoin", sdk.NewInt(2)}, true},
		{"no eth address", orchAddr.String(), "not an addr",
			sdk.Coin{"testcoin", sdk.NewInt(10)}, sdk.Coin{"testcoin", sdk.NewInt(2)}, true},
		{"unmatched denominations", orchAddr.String(), ethAddr.String(),
			sdk.Coin{"testcoin", sdk.NewInt(10)}, sdk.Coin{"othercoin", sdk.NewInt(2)}, true},
		{"negative amount", orchAddr.String(), ethAddr.String(),
			sdk.Coin{"testcoin", sdk.NewInt(-10)}, sdk.Coin{"testcoin", sdk.NewInt(2)}, true},
	}

	for _, tc := range testCases {
		mt := MsgTransfer{tc.sender, tc.eth, tc.amount, tc.bridgeFee}
		err := mt.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestMsgRequestBatch_ValidateBasic(t *testing.T) {
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")

	testCases := []struct{
		name string
		orch string
		denom string
		expError bool
	}{
		{"valid input", orchAddr.String(), "testcoin", false},
		{"invalid denom", orchAddr.String(), "gravity/broken", true},
		{"no orchestrator", "not an addr", "testcoin", true},
	}

	for _, tc := range testCases {
		mrb := MsgRequestBatch{tc.orch, tc.denom}
		err := mrb.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestMsgCancelTransfer_ValidateBasic(t *testing.T) {
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")

	testCases := []struct{
		name string
		sender string
		txid []byte
		expError bool
	}{
		{"valid input", orchAddr.String(), []byte("10"), false},
		{"invalid transaction ID", orchAddr.String(), []byte(""), true},
		{"no sender", "not an addr", []byte("10"), true},
	}

	for _, tc := range testCases {
		mct := MsgCancelTransfer{tc.txid, tc.sender}
		err := mct.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestMsgSubmitConfirm_ValidateBasic(t *testing.T) {
	orchCryptoAddr, err := newCosmosAddress()
	require.NoError(t, err, "unable to generate cosmos address")
	orchAddr, err := sdk.AccAddressFromHex(orchCryptoAddr.String())
	require.NoError(t, err, "unable to cast cosmos address to orchestrator address")
	ethAddr, err := newEthAddress()
	require.NoError(t, err, "unable to generate ethereum address")

	testCases := []struct{
		name string
		signer string
		confirm interface{}
		expError bool
	}{
		{"valid input", orchAddr.String(), ConfirmSignerSet{12, ethAddr.String(), orchAddr.String(), []byte("signature")}, false},
		{"no confirm", orchAddr.String(), nil, true},
		{"no signer", "not an addr", ConfirmSignerSet{12, ethAddr.String(), orchAddr.String(), []byte("signature")}, true},
	}

	for _, tc := range testCases {
		any, err := types.NewAnyWithValue(tc.confirm.(proto.Message))
		require.NoError(t, err, tc.name)
		msc := MsgSubmitConfirm{any, tc.signer}
		err = msc.ValidateBasic()
		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}
