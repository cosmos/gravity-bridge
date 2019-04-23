package oracle

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

// DefaultConsensusNeeded is the default fraction of validators needed to make claims on a prophecy in order for it to pass
const DefaultConsensusNeeded float64 = 0.7

type (
	Keeper = keeper.Keeper

	Claim    = types.Claim
	Prophecy = types.Prophecy
)

var (
	NewKeeper = keeper.NewKeeper

	NewClaim    = types.NewClaim
	NewProphecy = types.NewProphecy
)

const (
	PendingStatus = types.PendingStatus
	SuccessStatus = types.SuccessStatus
	FailedStatus  = types.FailedStatus
)

const (
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace

	TestID           = types.TestID
	TestMinimumPower = types.TestMinimumPower
)

var (
	ErrProphecyNotFound   = types.ErrProphecyNotFound
	ErrMinimumPowerTooLow = types.ErrMinimumPowerTooLow
	ErrInvalidIdentifier  = types.ErrInvalidIdentifier
)
