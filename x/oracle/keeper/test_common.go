package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"

	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	TestAddress   = "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	TestValidator = "cosmos1xdp5tvt7lxh8rf9xx07wy2xlagzhq24ha48xtq"
	TestID        = "ethereumAddress0"
)

func CreateTestProphecy(t *testing.T) types.BridgeProphecy {
	testAddress, err1 := sdk.AccAddressFromBech32(TestAddress)
	testValidator, err2 := sdk.AccAddressFromBech32(TestValidator)
	amount, err3 := sdk.ParseCoins("1test")
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	bridgeClaim := types.NewBridgeClaim(TestID, testAddress, testValidator, amount)
	bridgeClaims := []types.BridgeClaim{bridgeClaim}
	newProphecy := types.NewBridgeProphecy(TestID, types.PendingStatus, 5, bridgeClaims)
	return newProphecy
}

// CreateTestKeepers greates an OracleKeeper, AccountKeeper and Context to be used for test input
func CreateTestKeepers(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, auth.AccountKeeper, Keeper) {
	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "testchainid"}, isCheckTx, log.NewNopLogger())
	cdc := MakeTestCodec()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	ck := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	keeper := NewKeeper(ck, keyOracle, cdc, types.DefaultCodespace)

	return ctx, accountKeeper, keeper
}

// MakeTestCodec creates a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	// Register Msgs
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
