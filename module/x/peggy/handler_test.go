package peggy

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleValsetRequest(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr       sdk.AccAddress = bytes.Repeat([]byte{1}, 12)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
		myBlockHeight      int64          = 200
	)

	k, ctx, _ := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	h := NewHandler(k)
	msg := types.MsgValsetRequest{Requester: myCosmosAddr}
	ctx = ctx.WithBlockTime(myBlockTime).WithBlockHeight(myBlockHeight)
	res, err := h(ctx, msg)
	// then
	require.NoError(t, err)
	nonce := types.UInt64NonceFromBytes(res.Data)
	require.False(t, nonce.IsEmpty())
	require.Equal(t, types.NewUInt64Nonce(uint64(myBlockHeight)), nonce)
	// and persisted
	valset := k.GetValsetRequest(ctx, nonce)
	require.NotNil(t, valset)
	assert.Equal(t, nonce, valset.Nonce)
	require.Len(t, valset.Members, 1)
	assert.Equal(t, []uint64{math.MaxUint32}, valset.Members.GetPowers())
	assert.Equal(t, types.NewEthereumAddress(""), valset.Members[0].EthereumAddress)
}

func TestHandleCreateEthereumClaims(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myCosmosAddr       sdk.AccAddress = bytes.Repeat([]byte{1}, 12)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myNonce                           = types.NewUInt64Nonce(1)
		anyETHAddr                        = types.NewEthereumAddress("any-address")
		tokenETHAddr                      = types.NewEthereumAddress("any-erc20-token-addr")
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)
	k, ctx, keepers := keeper.CreateTestEnv(t)
	k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
	h := NewHandler(k)

	msg := MsgCreateEthereumClaims{
		EthereumChainID:       0,
		BridgeContractAddress: types.NewEthereumAddress(""),
		Orchestrator:          myOrchestratorAddr,
		Claims: []EthereumClaim{
			EthereumBridgeDepositClaim{
				Nonce: myNonce,
				ERC20Token: types.ERC20Token{
					Amount:               12,
					Symbol:               "ALX",
					TokenContractAddress: tokenETHAddr,
				},
				EthereumSender: anyETHAddr,
				CosmosReceiver: myCosmosAddr,
			},
		},
	}
	// when
	ctx = ctx.WithBlockTime(myBlockTime)
	_, err := h(ctx, msg)
	// then
	require.NoError(t, err)
	// and claim persisted
	claimFound := k.HasClaim(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce, myValAddr, msg.Claims[0].Details())
	assert.True(t, claimFound)
	// and attestation persisted
	a := k.GetAttestation(ctx, types.ClaimTypeEthereumBridgeDeposit, myNonce)
	require.NotNil(t, a)
	// and vouchers added to the account
	balance := keepers.BankKeeper.GetCoins(ctx, myCosmosAddr)
	assert.Equal(t, sdk.Coins{sdk.NewInt64Coin("peggy96dde7db38", 12)}, balance)
}

func TestHandleBridgeSignatureSubmission(t *testing.T) {
	var (
		myOrchestratorAddr sdk.AccAddress = make([]byte, sdk.AddrLen)
		myValAddr                         = sdk.ValAddress(myOrchestratorAddr) // revisit when proper mapping is impl in keeper
		myBlockTime                       = time.Date(2020, 9, 14, 15, 20, 10, 0, time.UTC)
	)

	privKey, err := ethCrypto.HexToECDSA("0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf57"[2:])
	require.NoError(t, err)

	specs := map[string]struct {
		setup  func(ctx sdk.Context, k Keeper) MsgBridgeSignatureSubmission
		expErr bool
	}{
		"SignedMultiSigUpdate good": {
			setup: func(ctx sdk.Context, k Keeper) MsgBridgeSignatureSubmission {
				v := k.SetValsetRequest(ctx)
				validSig, err := types.NewEthereumSignature(v.GetCheckpoint(), privKey)
				require.NoError(t, err)
				return MsgBridgeSignatureSubmission{
					ClaimType:         types.ClaimTypeOrchestratorSignedMultiSigUpdate,
					Nonce:             v.Nonce,
					Orchestrator:      myOrchestratorAddr,
					EthereumSignature: validSig,
				}
			},
		},
		"SignedWithdrawBatch good": {
			setup: func(ctx sdk.Context, k Keeper) MsgBridgeSignatureSubmission {
				vouchers := keeper.MintVouchersFromAir(t, ctx, k, myOrchestratorAddr, types.NewERC20Token(12, "any", types.NewEthereumAddress("0x4251ed140bf791c4112bb61fcb6e72f927e8fef2")))
				require.NoError(t, err)
				// with a transaction
				k.AddToOutgoingPool(ctx, myOrchestratorAddr, types.NewEthereumAddress("0xb5f728530fe1477ba8b780823a2d48f367fc9fc2"), vouchers, sdk.NewInt64Coin(vouchers.Denom, 0))
				voucherDenom, err := types.AsVoucherDenom(vouchers.Denom)
				require.NoError(t, err)
				// in a batch
				b, err := k.BuildOutgoingTXBatch(ctx, voucherDenom, 10)
				require.NoError(t, err)
				// and a multisig observed
				v := k.SetValsetRequest(ctx)
				att, err := k.AddClaim(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate, v.Nonce, myValAddr, types.SignedCheckpoint{Checkpoint: v.GetCheckpoint()})
				require.NoError(t, err)
				require.Equal(t, types.ProcessStatusProcessed, att.Status)
				// create signature
				checkpoint, err := b.GetCheckpoint(&v)
				require.NoError(t, err)
				validSig, err := types.NewEthereumSignature(checkpoint, privKey)
				require.NoError(t, err)
				return MsgBridgeSignatureSubmission{
					ClaimType:         types.ClaimTypeOrchestratorSignedWithdrawBatch,
					Nonce:             b.Nonce,
					Orchestrator:      myOrchestratorAddr,
					EthereumSignature: validSig,
				}
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			k, ctx, _ := keeper.CreateTestEnv(t)
			k.StakingKeeper = keeper.NewStakingKeeperMock(myValAddr)
			h := NewHandler(k)
			k.SetEthAddress(ctx, myValAddr, types.NewEthereumAddress("0xbd5d7df0349ff9671e36ec5545e849cbb93ac7fa"))

			// when
			ctx = ctx.WithBlockTime(myBlockTime)
			msg := spec.setup(ctx, k)
			_, err = h(ctx, msg)
			if spec.expErr {
				assert.Error(t, err)
				return
			}
			// then
			require.NoError(t, err)

			// and claim persisted
			checkPoint, err := getCheckpoint(ctx, k, msg.ClaimType, msg.Nonce)
			require.NoError(t, err)
			claimFound := k.HasClaim(ctx, msg.ClaimType, msg.Nonce, myValAddr, types.SignedCheckpoint{
				Checkpoint: checkPoint,
			})
			assert.True(t, claimFound)
			// and approval persisted
			sigFound := k.HasBridgeApprovalSignature(ctx, msg.ClaimType, msg.Nonce, myValAddr)
			assert.True(t, sigFound)

			// and attestation persisted
			a := k.GetAttestation(ctx, msg.ClaimType, msg.Nonce)
			require.NotNil(t, a)
		})
	}
}
