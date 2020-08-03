package peggy

import (
	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.ModuleName
	QuerierRoute      = types.QuerierRoute
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
	// NewMsgBuyName       = types.NewMsgBuyName
	// NewMsgSetName       = types.NewMsgSetName
	// NewMsgDeleteName    = types.NewMsgDeleteName
	NewMsgSetEthAddress = types.NewMsgSetEthAddress
	// NewWhois            = types.NewWhois
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
	// MsgSetName       = types.MsgSetName
	// MsgBuyName       = types.MsgBuyName
	// MsgDeleteName    = types.MsgDeleteName
	MsgSetEthAddress = types.MsgSetEthAddress
	MsgValsetConfirm = types.MsgValsetConfirm
	MsgValsetRequest = types.MsgValsetRequest
	// QueryResResolve  = types.QueryResResolve
	// QueryResNames    = types.QueryResNames
	// Whois            = types.Whois
)
