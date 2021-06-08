package keeper

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestMsgServer_SubmitEthereumSignature(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	ethPubKey := fmt.Sprintf("0x%s", hex.EncodeToString(ethCrypto.CompressPubkey(&ethPrivKey.PublicKey)))

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = common.HexToAddress(ethPubKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr2, orcAddr2)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr3, orcAddr3)
	}

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	// setup for GetOutgoingTx
	signerSetTx := gk.CreateSignerSetTx(ctx)

	// setup for ValidateEthereumSignature
	gravityId := gk.getGravityID(ctx)
	checkpoint := signerSetTx.GetCheckpoint([]byte(gravityId))
	signature, err := types.NewEthereumSignature(checkpoint, ethPrivKey)
	require.NoError(t, err)

	signerSetTxConfirmation := &types.SignerSetTxConfirmation{
		SignerSetNonce: signerSetTx.Nonce,
		EthereumSigner: ethAddr1.Hex(),
		Signature:      signature,
	}

	confirmation, err := types.PackConfirmation(signerSetTxConfirmation)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEthereumTxConfirmation{
		Confirmation: confirmation,
		Signer:       orcAddr1.String(),
	}

	_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

// TODO(levi) ensure coverage for:
// SendToEthereum(context.Context, *MsgSendToEthereum) (*MsgSendToEthereumResponse, error)
// CancelSendToEthereum(context.Context, *MsgCancelSendToEthereum) (*MsgCancelSendToEthereumResponse, error)
// RequestBatchTx(context.Context, *MsgRequestBatchTx) (*MsgRequestBatchTxResponse, error)
// SubmitEthereumTxConfirmation(context.Context, *MsgSubmitEthereumTxConfirmation) (*MsgSubmitEthereumTxConfirmationResponse, error)
// SubmitEthereumEvent(context.Context, *MsgSubmitEthereumEvent) (*MsgSubmitEthereumEventResponse, error)
// SetDelegateKeys(context.Context, *MsgDelegateKeys) (*MsgDelegateKeysResponse, error)
