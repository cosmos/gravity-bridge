package keeper

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
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
		EventNonce:     1,
		TokenContract:  "test-token-contract-string",
		Amount:         sdk.NewInt(1000),
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

func TestMsgServer_SetDelegateKeys(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env         = CreateTestEnv(t)
		ctx         = env.Context
		gk          = env.GravityKeeper
		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)
	)

	// setup for getSignerValidator
	gk.StakingKeeper = NewStakingKeeperMock(valAddr1)

	// Set the sequence to 1 because the antehandler will do this in the full
	// chain.
	acc := env.AccountKeeper.NewAccountWithAddress(ctx, orcAddr1)
	acc.SetSequence(1)
	env.AccountKeeper.SetAccount(ctx, acc)

	msgServer := NewMsgServerImpl(gk)

	ethMsg := types.DelegateKeysSignMsg{
		ValidatorAddress: valAddr1.String(),
		Nonce:            0,
	}
	signMsgBz := env.Marshaler.MustMarshalBinaryBare(&ethMsg)

	sig, err := types.NewEthereumSignature(signMsgBz, ethPrivKey)
	require.NoError(t, err)

	msg := &types.MsgDelegateKeys{
		ValidatorAddress:    valAddr1.String(),
		OrchestratorAddress: orcAddr1.String(),
		EthereumAddress:     ethAddr1.String(),
		EthSignature:        sig,
	}

	_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestEthVerify(t *testing.T) {
	// Replace privKeyHexStr and addrHexStr with your own private key and address
	// HEX values.
	privKeyHexStr := "0x9a86de8a78c5a8f9787ecdd611494550b37690f6eff354533357386d73812664"
	addrHexStr := "0xCe7A018732f60Ad707595302bA64A711cbd5b658"

	// ==========================================================================
	// setup
	// ==========================================================================
	privKeyBz, err := hexutil.Decode(privKeyHexStr)
	require.NoError(t, err)

	privKey, err := crypto.ToECDSA(privKeyBz)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	require.True(t, bytes.Equal(privKeyBz, crypto.FromECDSA(privKey)))
	require.Equal(t, privKeyHexStr, hexutil.Encode(crypto.FromECDSA(privKey)))

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	require.Equal(t, addrHexStr, address.Hex())

	// ==========================================================================
	// signature verification
	// ==========================================================================
	cdc := MakeTestMarshaler()

	valAddr := "cosmosvaloper1dmly9yyhd5lyhyl8qhs7wtcd4xt7gyxlesgvmc"
	signMsgBz, err := cdc.MarshalBinaryBare(&types.DelegateKeysSignMsg{
		ValidatorAddress: valAddr,
		Nonce:            0,
	})
	require.NoError(t, err)

	fmt.Println("MESSAGE BYTES TO SIGN:", hexutil.Encode(signMsgBz))

	sig, err := types.NewEthereumSignature(signMsgBz, privKey)
	sig[64] += 27 // change the V value
	require.NoError(t, err)

	err = types.ValidateEthereumSignature(signMsgBz, sig, address)
	require.NoError(t, err)

	// replace gorcSig with what the following command produces:
	// $ gorc sign-delegate-keys <your-eth-key-name> cosmosvaloper1dmly9yyhd5lyhyl8qhs7wtcd4xt7gyxlesgvmc 0
	gorcSig := "0xd34881c746b8498926bdea191529d5af66aa34938349e789aefab90cf0fc4ffe3cbffb85313cef0107d49b17af83f3175c63db00cd5edee58d2369bd507410551c"
	require.Equal(t, hexutil.Encode(sig), gorcSig)
}
