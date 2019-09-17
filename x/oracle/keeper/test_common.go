package keeper

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/cosmos/peggy/x/oracle/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	TestID                     = "oracleID"
	AlternateTestID            = "altOracleID"
	TestString                 = "{value: 5}"
	AlternateTestString        = "{value: 7}"
	AnotherAlternateTestString = "{value: 9}"
)

// CreateTestKeepers greates an Mock App, OracleKeeper, BankKeeper and ValidatorAddresses to be used for test input
func CreateTestKeepers(t *testing.T, consensusNeeded float64, validatorPowers []int64) (*mock.App, Keeper, bank.Keeper, []sdk.ValAddress) {
	mApp := mock.NewApp()
	RegisterTestCodecs(mApp.Cdc)

	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)

	feeCollector := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(stakingtypes.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollector.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	bankKeeper := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:          nil,
		stakingtypes.NotBondedPoolName: {supply.Burner, supply.Staking},
		stakingtypes.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(mApp.Cdc, keyStaking, tkeyStaking, supplyKeeper, mApp.ParamsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	keeper := NewKeeper(mApp.Cdc, keyOracle, stakingKeeper, types.DefaultCodespace, consensusNeeded)

	mApp.Router().AddRoute(staking.RouterKey, staking.NewHandler(stakingKeeper))
	mApp.SetEndBlocker(getEndBlocker(stakingKeeper))
	mApp.SetInitChainer(getInitChainer(mApp, stakingKeeper, mApp.AccountKeeper, supplyKeeper,
		[]supplyexported.ModuleAccountI{feeCollector, notBondedPool, bondPool}))

	require.NoError(t, mApp.CompleteSetup(keyStaking, tkeyStaking, keySupply))

	// create the validators addresses desired and fill them with the expected amount of coins
	commissionRates := staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	accounts := []*auth.BaseAccount{}
	addresses := []sdk.ValAddress{}
	for _, power := range validatorPowers {
		coinAmount := sdk.TokensFromConsensusPower(power)
		coin := sdk.Coin{"stake", coinAmount}
		coins := sdk.Coins{coin}
		priv := secp256k1.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())

		account := &auth.BaseAccount{
			Address: addr,
			Coins:   coins,
		}
		accounts = append(accounts, account)
		addresses = append(addresses, sdk.ValAddress(addr))
		description := staking.NewDescription("foo_moniker", "", "", "", "")
		createValidatorMsg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr), priv.PubKey(), coin, description, commissionRates, sdk.OneInt(),
		)
		header := abci.Header{Height: mApp.LastBlockHeight() + 1}
		mock.SignCheckDeliver(t, mApp.Cdc, mApp.BaseApp, header, []sdk.Msg{createValidatorMsg}, []uint64{0}, []uint64{0}, true, true, priv)
	}
	return mApp, keeper, bankKeeper, addresses

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
		address := stakingkeeper.TestAddr(buffer.String(), bech)
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
		publicKeys = append(publicKeys, stakingkeeper.NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

// RegisterTestCodecs registers codecs used only for testing
func RegisterTestCodecs(cdc *codec.Codec) {
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper staking.Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := staking.EndBlocker(ctx, keeper)

		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper staking.Keeper, accountKeeper stakingtypes.AccountKeeper, supplyKeeper stakingtypes.SupplyKeeper,
	blacklistedAddrs []supplyexported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		stakingGenesis := staking.DefaultGenesisState()
		validators := staking.InitGenesis(ctx, keeper, accountKeeper, supplyKeeper, stakingGenesis)
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}
