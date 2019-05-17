package keeper

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeperLib "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

// CreateTestKeepers greates an OracleKeeper, AccountKeeper and Context to be used for test input
func CreateTestKeepers(t *testing.T, consensusNeeded float64, validatorPowers []int64) (sdk.Context, auth.AccountKeeper, Keeper, bank.Keeper, []sdk.ValAddress, sdk.Error) {
	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "testchainid"}, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := MakeTestCodec()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	stakingKeeper := staking.NewKeeper(cdc, keyStaking, tkeyStaking, bankKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	stakingKeeper.SetPool(ctx, staking.InitialPool())
	stakingKeeper.SetParams(ctx, staking.DefaultParams())

	keeper, keeperErr := NewKeeper(stakingKeeper, keyOracle, cdc, types.DefaultCodespace, consensusNeeded)

	//construct the validators
	numValidators := len(validatorPowers)
	validators := make([]staking.Validator, numValidators)
	accountAddresses, valAddresses := CreateTestAddrs(len(validatorPowers))
	publicKeys := createTestPubKeys(len(validatorPowers))

	// create the validators addresses desired and fill them with the expected amount of coins
	for i, power := range validatorPowers {
		coins := cmtypes.TokensFromTendermintPower(power)
		pool := stakingKeeper.GetPool(ctx)
		err := error(nil)
		_, _, err = bankKeeper.AddCoins(ctx, accountAddresses[i], sdk.Coins{
			{stakingKeeper.BondDenom(ctx), coins},
		})
		require.Nil(t, err)
		pool.NotBondedTokens = pool.NotBondedTokens.Add(coins)
		stakingKeeper.SetPool(ctx, pool)
	}
	pool := stakingKeeper.GetPool(ctx)
	for i, power := range validatorPowers {
		validators[i] = staking.NewValidator(valAddresses[i], publicKeys[i], staking.Description{})
		validators[i].Status = sdk.Bonded
		validators[i].Tokens = sdk.ZeroInt()
		tokens := sdk.TokensFromTendermintPower(power)
		validators[i], pool, _ = validators[i].AddTokensFromDel(pool, tokens)
		stakingKeeper.SetPool(ctx, pool)
		stakingKeeperLib.TestingUpdateValidator(stakingKeeper, ctx, validators[i], true)
	}

	return ctx, accountKeeper, keeper, bankKeeper, valAddresses, keeperErr
}

// nolint: unparam
func CreateTestAddrs(numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	var addresses []sdk.AccAddress
	var valAddresses []sdk.ValAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		address := stakingKeeperLib.TestAddr(buffer.String(), bech)
		valAddress := sdk.ValAddress(address)
		addresses = append(addresses, address)
		valAddresses = append(valAddresses, valAddress)
		buffer.Reset()
	}
	return addresses, valAddresses
}

// nolint: unparam
func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, stakingKeeperLib.NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

// MakeTestCodec creates a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	// Register Msgs
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
