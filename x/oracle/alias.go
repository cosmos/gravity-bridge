package oracle

import (
	"github.com/swishlabsco/peggy/x/oracle/keeper"
	"github.com/swishlabsco/peggy/x/oracle/types"
)

// DefaultConsensusNeeded is the default fraction of validators needed to create claims on a prophecy in order for it to pass
const DefaultConsensusNeeded float64 = 0.7

type (
	Keeper = keeper.Keeper

	Prophecy = types.Prophecy

	Status = types.Status

	Claim = types.Claim
)

var (
	NewKeeper = keeper.NewKeeper

	NewProphecy = types.NewProphecy

	NewClaim = types.NewClaim
)

const (
	PendingStatus = types.PendingStatusText
	SuccessStatus = types.SuccessStatusText
	FailedStatus  = types.FailedStatusText
)

const (
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace

	TestID = keeper.TestID
)

var (
	ErrProphecyNotFound              = types.ErrProphecyNotFound
	ErrMinimumConsensusNeededInvalid = types.ErrMinimumConsensusNeededInvalid
	ErrInvalidIdentifier             = types.ErrInvalidIdentifier
)
