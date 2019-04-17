package oracle

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

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
	PendingStatus  = types.PendingStatus
	CompleteStatus = types.CompleteStatus
	FailedStatus   = types.FailedStatus
)

const (
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace

	DefaultConsensusNeeded = types.DefaultConsensusNeeded

	TestID           = types.TestID
	TestMinimumPower = types.TestMinimumPower
)

var (
	ErrProphecyNotFound   = types.ErrProphecyNotFound
	ErrMinimumPowerTooLow = types.ErrMinimumPowerTooLow
	ErrInvalidIdentifier  = types.ErrInvalidIdentifier
)
