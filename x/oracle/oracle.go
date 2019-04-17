package oracle

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/querier"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

type (
	Keeper = keeper.Keeper

	Claim    = types.Claim
	Prophecy = types.Prophecy

	MsgMakeEthBridgeClaim = types.MsgMakeEthBridgeClaim
)

var (
	NewKeeper = keeper.NewKeeper

	NewMsgMakeEthBridgeClaim = types.NewMsgMakeEthBridgeClaim
	NewClaim                 = types.NewClaim
	NewProphecy              = types.NewProphecy

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
	ErrProphecyNotFound   = types.ErrProphecyNotFound
	ErrMinimumPowerTooLow = types.ErrMinimumPowerTooLow
	ErrInvalidIdentifier  = types.ErrInvalidIdentifier

	ErrInvalidEthereumNonce = types.ErrInvalidEthereumNonce
)
