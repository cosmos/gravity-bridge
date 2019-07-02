package ethbridge

import (
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/querier"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
)

type (
	MsgCreateEthBridgeClaim = types.MsgCreateEthBridgeClaim
)

var (
	NewMsgCreateEthBridgeClaim = types.NewMsgCreateEthBridgeClaim
	NewEthBridgeClaim          = types.NewEthBridgeClaim

	NewQueryEthProphecyParams = types.NewQueryEthProphecyParams

	ErrInvalidEthNonce = types.ErrInvalidEthNonce

	RegisterCodec = types.RegisterCodec

	NewQuerier = querier.NewQuerier
)

const (
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace

	QueryEthProphecy = querier.QueryEthProphecy
)
