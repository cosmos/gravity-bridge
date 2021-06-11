package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/bytes"
)

func TestKeeper_Params(t *testing.T) {
	env := CreateTestEnv(t)
	ctx := sdk.WrapSDKContext(env.Context)
	gk := env.GravityKeeper

	req := &types.ParamsRequest{}
	res, err := gk.Params(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestKeeper_LatestSignerSetTx(t *testing.T) {
	t.Run("read before there's anything in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		req := &types.LatestSignerSetTxRequest{}
		res, err := gk.LatestSignerSetTx(sdk.WrapSDKContext(ctx), req)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper
		{ // setup
			sstx := gk.CreateSignerSetTx(env.Context)
			require.NotNil(t, sstx)
		}
		{ // validate
			req := &types.LatestSignerSetTxRequest{}
			res, err := gk.LatestSignerSetTx(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, res)
		}
	})
}

func TestKeeper_SignerSetTx(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		var signerSetNonce uint64
		{ // setup
			sstx := gk.CreateSignerSetTx(env.Context)
			require.NotNil(t, sstx)
			signerSetNonce = sstx.Nonce
		}
		{ // validate
			req := &types.SignerSetTxRequest{SignerSetNonce: signerSetNonce}
			res, err := gk.SignerSetTx(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.NotNil(t, res.SignerSet)
		}
	})
}

func TestKeeper_BatchTx(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		const (
			batchNonce    = 55
			tokenContract = "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4"
		)

		{ // setup
			gk.SetOutgoingTx(ctx, &types.BatchTx{
				BatchNonce:    batchNonce,
				Timeout:       1000,
				Transactions:  nil,
				TokenContract: tokenContract,
				Height:        100,
			})
		}
		{ // validate
			req := &types.BatchTxRequest{
				BatchNonce:    batchNonce,
				TokenContract: tokenContract,
			}

			res, err := gk.BatchTx(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.NotNil(t, res.Batch)
		}
	})
}

func TestKeeper_ContractCallTx(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		const (
			invalidationNonce = 100
			invalidationScope = "an-invalidation-scope"
		)

		{ // setup
			gk.SetOutgoingTx(ctx, &types.ContractCallTx{
				InvalidationNonce: invalidationNonce,
				InvalidationScope: bytes.HexBytes(invalidationScope),
			})
		}
		{ // validate
			req := &types.ContractCallTxRequest{
				InvalidationNonce: invalidationNonce,
				InvalidationScope: bytes.HexBytes(invalidationScope),
			}

			res, err := gk.ContractCallTx(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.NotNil(t, res.LogicCall)
		}
	})
}

func TestKeeper_SignerSetTxs(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			require.NotNil(t, gk.CreateSignerSetTx(env.Context))
			require.NotNil(t, gk.CreateSignerSetTx(env.Context))
		}
		{ // validate
			req := &types.SignerSetTxsRequest{}
			res, err := gk.SignerSetTxs(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Len(t, res.SignerSets, 2)
		}
	})
}

func TestKeeper_BatchTxs(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.SetOutgoingTx(ctx, &types.BatchTx{
				BatchNonce:    1000,
				Timeout:       1000,
				Transactions:  nil,
				TokenContract: "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4",
				Height:        1000,
			})
			gk.SetOutgoingTx(ctx, &types.BatchTx{
				BatchNonce:    1001,
				Timeout:       1000,
				Transactions:  nil,
				TokenContract: "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4",
				Height:        1001,
			})
		}
		{ // validate
			req := &types.BatchTxsRequest{}
			got, err := gk.BatchTxs(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.Batches, 2)
		}
	})
}

func TestKeeper_ContractCallTxs(t *testing.T) {
	t.Run("read after there's something in state", func(t *testing.T) {
		env := CreateTestEnv(t)
		ctx := env.Context
		gk := env.GravityKeeper

		{ // setup
			gk.SetOutgoingTx(ctx, &types.ContractCallTx{
				InvalidationNonce: 5,
				InvalidationScope: []byte("an-invalidation-scope"),
				// TODO
			})
			gk.SetOutgoingTx(ctx, &types.ContractCallTx{
				InvalidationNonce: 6,
				InvalidationScope: []byte("an-invalidation-scope"),
			})
		}
		{ // validate
			req := &types.ContractCallTxsRequest{}
			got, err := gk.ContractCallTxs(sdk.WrapSDKContext(ctx), req)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.Calls, 2)
		}
	})
}

// TODO(levi) ensure coverage for:
// ContractCallTx(context.Context, *ContractCallTxRequest) (*ContractCallTxResponse, error)
// ContractCallTxs(context.Context, *ContractCallTxsRequest) (*ContractCallTxsResponse, error)

// SignerSetTxConfirmations(context.Context, *SignerSetTxConfirmationsRequest) (*SignerSetTxConfirmationsResponse, error)
// BatchTxConfirmations(context.Context, *BatchTxConfirmationsRequest) (*BatchTxConfirmationsResponse, error)
// ContractCallTxConfirmations(context.Context, *ContractCallTxConfirmationsRequest) (*ContractCallTxConfirmationsResponse, error)

// UnsignedSignerSetTxs(context.Context, *UnsignedSignerSetTxsRequest) (*UnsignedSignerSetTxsResponse, error)
// UnsignedBatchTxs(context.Context, *UnsignedBatchTxsRequest) (*UnsignedBatchTxsResponse, error)
// UnsignedContractCallTxs(context.Context, *UnsignedContractCallTxsRequest) (*UnsignedContractCallTxsResponse, error)

// BatchTxFees(context.Context, *BatchTxFeesRequest) (*BatchTxFeesResponse, error)
// ERC20ToDenom(context.Context, *ERC20ToDenomRequest) (*ERC20ToDenomResponse, error)
// DenomToERC20(context.Context, *DenomToERC20Request) (*DenomToERC20Response, error)
// BatchedSendToEthereums(context.Context, *BatchedSendToEthereumsRequest) (*BatchedSendToEthereumsResponse, error)
// UnbatchedSendToEthereums(context.Context, *UnbatchedSendToEthereumsRequest) (*UnbatchedSendToEthereumsResponse, error)
// DelegateKeysByValidator(context.Context, *DelegateKeysByValidatorRequest) (*DelegateKeysByValidatorResponse, error)
// DelegateKeysByEthereumSigner(context.Context, *DelegateKeysByEthereumSignerRequest) (*DelegateKeysByEthereumSignerResponse, error)
// DelegateKeysByOrchestrator(context.Context, *DelegateKeysByOrchestratorRequest) (*DelegateKeysByOrchestratorResponse, error)
