package oracle

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/querier"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

type (
	Keeper = keeper.Keeper

	BridgeClaim    = types.BridgeClaim
	BridgeProphecy = types.BridgeProphecy

	MsgMakeBridgeEthClaim = types.MsgMakeBridgeEthClaim
)

var (
	NewKeeper = keeper.NewKeeper

	NewMsgMakeEthBridgeClaim = types.NewMsgMakeEthBridgeClaim
	NewBridgeClaim           = types.NewBridgeClaim
	NewBridgeProphecy        = types.NewBridgeProphecy

	RegisterCodec = types.RegisterCodec

	NewQuerier             = querier.NewQuerier
	NewQueryProphecyParams = querier.NewQueryProphecyParams
)

const (
	QueryProphecy = querier.QueryProphecy
	PendingStatus = types.PendingStatus
)

const (
	StoreKey               = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	RouterKey              = types.RouterKey
	DefaultCodespace       = types.DefaultCodespace
	DefaultConsensusNeeded = types.DefaultConsensusNeeded
)

var (
	ErrInvalidNonce       = types.ErrInvalidNonce
	ErrNotFound           = types.ErrNotFound
	ErrMinimumPowerTooLow = types.ErrMinimumPowerTooLow
	ErrInvalidIdentifier  = types.ErrInvalidIdentifier
)
