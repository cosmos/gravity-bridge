package keeper

import (
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var TestingStakeParams = staking.Params{
	UnbondingTime:     100,
	MaxValidators:     10,
	MaxEntries:        10,
	HistoricalEntries: 10,
	BondDenom:         "stake",
}

type TestKeepers struct {
	AccountKeeper auth.AccountKeeper
	StakingKeeper staking.Keeper
	DistKeeper    distribution.Keeper
	SupplyKeeper  supply.Keeper
	GovKeeper     gov.Keeper
	BankKeeper    bank.Keeper
}

func CreateTestEnv(t *testing.T) (Keeper, sdk.Context, TestKeepers) {
	t.Helper()
	peggyKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyDistro := sdk.NewKVStoreKey(distribution.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyGov := sdk.NewKVStoreKey(govtypes.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(peggyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistro, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	const isCheckTx = false
	ctx := sdk.NewContext(ms, abci.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, isCheckTx, log.TestingLogger())

	cdc := MakeTestCodec()

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	// this is also used to initialize module accounts for all the map keys
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          {supply.Minter, supply.Burner},
	}
	blockedAddr := make(map[string]bool, len(maccPerms))
	for acc := range maccPerms {
		blockedAddr[supply.NewModuleAddress(acc).String()] = true
	}
	bankKeeper := bank.NewBaseKeeper(
		accountKeeper,
		paramsKeeper.Subspace(bank.DefaultParamspace),
		blockedAddr,
	)
	bankKeeper.SetSendEnabled(ctx, true)

	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(cdc, keyStaking, supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace))
	stakingKeeper.SetParams(ctx, TestingStakeParams)

	distKeeper := distribution.NewKeeper(cdc, keyDistro, paramsKeeper.Subspace(distribution.DefaultParamspace), stakingKeeper, supplyKeeper, auth.FeeCollectorName, nil)
	distKeeper.SetParams(ctx, distribution.DefaultParams())
	stakingKeeper.SetHooks(distKeeper.Hooks())

	// set genesis items required for distribution
	distKeeper.SetFeePool(ctx, distribution.InitialFeePool())

	// total supply to track this
	totalSupply := sdk.NewCoins(sdk.NewInt64Coin("stake", 100000000))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// set up initial accounts
	for name, perms := range maccPerms {
		mod := supply.NewEmptyModuleAccount(name, perms...)
		if name == staking.NotBondedPoolName {
			err = mod.SetCoins(totalSupply)
			require.NoError(t, err)
		} else if name == distribution.ModuleName {
			// some big pot to pay out
			err = mod.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("stake", 500000)))
			require.NoError(t, err)
		}
		supplyKeeper.SetModuleAccount(ctx, mod)
	}

	stakeAddr := supply.NewModuleAddress(staking.BondedPoolName)
	moduleAcct := accountKeeper.GetAccount(ctx, stakeAddr)
	require.NotNil(t, moduleAcct)

	router := baseapp.NewRouter()
	bh := bank.NewHandler(bankKeeper)
	router.AddRoute(bank.RouterKey, bh)
	sh := staking.NewHandler(stakingKeeper)
	router.AddRoute(staking.RouterKey, sh)
	dh := distribution.NewHandler(distKeeper)
	router.AddRoute(distribution.RouterKey, dh)

	// Load default wasm config

	govRouter := gov.NewRouter().
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler)

	govKeeper := gov.NewKeeper(
		cdc, keyGov, paramsKeeper.Subspace(govtypes.DefaultParamspace).WithKeyTable(gov.ParamKeyTable()), supplyKeeper, stakingKeeper, govRouter,
	)

	govKeeper.SetProposalID(ctx, govtypes.DefaultStartingProposalID)
	govKeeper.SetDepositParams(ctx, govtypes.DefaultDepositParams())
	govKeeper.SetVotingParams(ctx, govtypes.DefaultVotingParams())
	govKeeper.SetTallyParams(ctx, govtypes.DefaultTallyParams())

	keepers := TestKeepers{
		AccountKeeper: accountKeeper,
		SupplyKeeper:  supplyKeeper,
		StakingKeeper: stakingKeeper,
		DistKeeper:    distKeeper,
		GovKeeper:     govKeeper,
		BankKeeper:    bankKeeper,
	}

	k := NewKeeper(cdc, peggyKey, paramsKeeper.Subspace(types.DefaultParamspace), stakingKeeper, supplyKeeper)
	k.setParams(ctx, types.Params{
		// any eth address
		BridgeContractAddress: types.NewEthereumAddress("0x8858eeb3dfffa017d4bce9801d340d36cf895ccf"),
		BridgeChainID:         11,
	})
	return k, ctx, keepers
}

func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.AppModuleBasic{}.RegisterCodec(cdc)
	bank.AppModuleBasic{}.RegisterCodec(cdc)
	supply.AppModuleBasic{}.RegisterCodec(cdc)
	staking.AppModuleBasic{}.RegisterCodec(cdc)
	distribution.AppModuleBasic{}.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	params.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

func MintVouchersFromAir(t *testing.T, ctx sdk.Context, k Keeper, dest sdk.AccAddress, amount types.ERC20Token) sdk.Coin {
	if !k.HasCounterpartDenominator(ctx, types.NewVoucherDenom(amount.TokenContractAddress, amount.Symbol)) {
		k.StoreCounterpartDenominator(ctx, amount.TokenContractAddress, amount.Symbol)
	}
	coin := amount.AsVoucherCoin()
	vouchers := sdk.Coins{coin}
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, vouchers)
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, dest, vouchers)
	require.NoError(t, err)
	return coin
}

var _ types.StakingKeeper = &StakingKeeperMock{}

type StakingKeeperMock struct {
	BondedValidators []staking.Validator
	ValidatorPower   map[string]int64
}

func NewStakingKeeperMock(operators ...sdk.ValAddress) *StakingKeeperMock {
	r := &StakingKeeperMock{
		BondedValidators: make([]staking.Validator, 0),
		ValidatorPower:   make(map[string]int64, 0),
	}
	const defaultTestPower = 100
	for _, a := range operators {
		r.BondedValidators = append(r.BondedValidators, staking.Validator{
			OperatorAddress: a,
		})
		r.ValidatorPower[a.String()] = defaultTestPower
	}
	return r
}

type MockStakingValidatorData struct {
	Operator sdk.ValAddress
	Power    int64
}

func NewStakingKeeperWeightedMock(t ...MockStakingValidatorData) *StakingKeeperMock {
	r := &StakingKeeperMock{
		BondedValidators: make([]staking.Validator, len(t)),
		ValidatorPower:   make(map[string]int64, len(t)),
	}

	for i, a := range t {
		r.BondedValidators[i] = staking.Validator{
			OperatorAddress: a.Operator,
		}
		r.ValidatorPower[a.Operator.String()] = a.Power
	}
	return r
}

func (s *StakingKeeperMock) GetBondedValidatorsByPower(ctx sdk.Context) []staking.Validator {
	return s.BondedValidators
}

func (s *StakingKeeperMock) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64 {
	v, ok := s.ValidatorPower[operator.String()]
	if !ok {
		panic("unknown address")
	}
	return v
}

func (s *StakingKeeperMock) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	var total int64
	for _, v := range s.ValidatorPower {
		total += v
	}
	return sdk.NewInt(total)
}

type AlwaysPanicStakingMock struct{}

func (s AlwaysPanicStakingMock) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	panic("unexpected call")
}

func (s AlwaysPanicStakingMock) GetBondedValidatorsByPower(ctx sdk.Context) []staking.Validator {
	panic("unexpected call")
}

func (s AlwaysPanicStakingMock) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64 {
	panic("unexpected call")
}
