package types

import (
	proto "github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgSubmitEvent{},
		&MsgSubmitConfirm{},
		&MsgCancelSendToEth{},
	)

	registry.RegisterInterface(
		"gravity.v1.EthereumEvent",
		(*EthereumEvent)(nil),
		&DepositEvent{},
		&WithdrawEvent{},
		&CosmosERC20DeployedEvent{},
		&LogicCallExecutedEvent{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func PackEvent(claim EthereumEvent) (*types.Any, error) {
	msg, ok := claim.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", claim)
	}

	anyEvent, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyEvent, nil
}

// UnpackEvent unpacks an Any into an EthereumEvent. It returns an error if the
// claim can't be unpacked.
func UnpackEvent(any *types.Any) (EthereumEvent, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	claim, ok := any.GetCachedValue().(EthereumEvent)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into EthereumEvent %T", any)
	}

	return claim, nil
}
