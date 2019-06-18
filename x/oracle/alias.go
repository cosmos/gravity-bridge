package oracle

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

// DefaultConsensusNeeded is the default fraction of validators needed to make claims on a prophecy in order for it to pass
const DefaultConsensusNeeded float64 = 0.7

type (
	Keeper = keeper.Keeper

	Prophecy = types.Prophecy

	Status = types.Status
)

var (
	NewKeeper = keeper.NewKeeper

	NewProphecy = types.NewProphecy
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

	TestID = types.TestID
)

var (
	ErrProphecyNotFound              = types.ErrProphecyNotFound
	ErrMinimumConsensusNeededInvalid = types.ErrMinimumConsensusNeededInvalid
	ErrInvalidIdentifier             = types.ErrInvalidIdentifier
)
