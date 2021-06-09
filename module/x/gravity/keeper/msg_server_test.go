package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/require"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

func TestMsgServer_SubmitEthereumSignature(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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

func TestMsgServer_SendToEthereum(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testDenom = "stake"

		balance = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10000),
		}
		amount = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(1000),
		}
		fee = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10),
		}
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	{ // add balance to bank
		err = env.AddBalanceToBank(ctx, orcAddr1, sdk.Coins{balance})
		require.NoError(t, err)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, testDenom, "testcontractstring")

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEthereum{
		Sender:            orcAddr1.String(),
		EthereumRecipient: ethAddr1.String(),
		Amount:            amount,
		BridgeFee:         fee,
	}

	_, err = msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_CancelSendToEthereum(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testDenom = "stake"

		balance = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10000),
		}
		amount = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(1000),
		}
		fee = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10),
		}
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	{ // add balance to bank
		err = env.AddBalanceToBank(ctx, orcAddr1, sdk.Coins{balance})
		require.NoError(t, err)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, testDenom, "testcontractstring")

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEthereum{
		Sender:            orcAddr1.String(),
		EthereumRecipient: ethAddr1.String(),
		Amount:            amount,
		BridgeFee:         fee,
	}

	response, err := msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	cancelMsg := &types.MsgCancelSendToEthereum{
		Id:     response.Id,
		Sender: orcAddr1.String(),
	}
	_, err = msgServer.CancelSendToEthereum(sdk.WrapSDKContext(ctx), cancelMsg)
	require.NoError(t, err)
}

func TestMsgServer_RequestBatchTx(t *testing.T) {
	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		//ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testDenom = "stake"
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, testDenom, "testcontractstring")

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgRequestBatchTx{
		Signer: orcAddr1.String(),
		Denom:  testDenom,
	}

	_, err := msgServer.RequestBatchTx(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_SubmitEthereumEvent(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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

	sendToCosmosEvent := &types.SendToCosmosEvent{
		EventNonce: 1,
		TokenContract: "test-token-contract-string",
		Amount: sdk.NewInt(1000),
		EthereumSender: ethAddr1.String(),
		CosmosReceiver: orcAddr1.String(),
		EthereumHeight: 200,
	}

	event, err := types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEthereumEvent{
		Event:  event,
		Signer: orcAddr1.String(),
	}

	_, err = msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

// TODO(levi) ensure coverage for:
// SubmitEthereumEvent(context.Context, *MsgSubmitEthereumEvent) (*MsgSubmitEthereumEventResponse, error)
// SetDelegateKeys(context.Context, *MsgDelegateKeys) (*MsgDelegateKeysResponse, error)
