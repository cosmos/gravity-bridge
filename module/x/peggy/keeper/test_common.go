package keeper

import (
	"testing"
	"time"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func CreateTestEnv(t *testing.T) (Keeper, sdk.Context) {
	t.Helper()
	peggyKey := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(peggyKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, false, log.NewNopLogger())

	cdc := MakeTestCodec()
	return NewKeeper(cdc, peggyKey, AlwaysPanicStakingMock{}), ctx
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

func (s AlwaysPanicStakingMock) GetBondedValidatorsByPower(ctx sdk.Context) []staking.Validator {
	panic("unexpected call")
}

func (s AlwaysPanicStakingMock) GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64 {
	panic("unexpected call")
}
