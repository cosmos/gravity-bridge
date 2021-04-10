package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgSubmitClaim{},
		&MsgSubmitConfirm{},
		&MsgCancelSendToEth{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.EthereumClaim",
		(*EthereumClaim)(nil),
		&DepositClaim{},
		&WithdrawClaim{},
		&ERC20DeployedClaim{},
		&LogicCallExecutedClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
