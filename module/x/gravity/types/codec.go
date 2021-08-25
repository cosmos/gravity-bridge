package types

import (
	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDelegateKeys{}, "gravity-bridge/", nil)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendToEthereum{},
		&MsgCancelSendToEthereum{},
		&MsgRequestBatchTx{},
		&MsgSubmitEthereumEvent{},
		&MsgSubmitEthereumTxConfirmation{},
		&MsgDelegateKeys{},
	)

	registry.RegisterInterface(
		"gravity.v1.EthereumEvent",
		(*EthereumEvent)(nil),
		&SendToCosmosEvent{},
		&BatchExecutedEvent{},
		&ERC20DeployedEvent{},
		&ContractCallExecutedEvent{},
		&SignerSetTxExecutedEvent{},
	)

	registry.RegisterInterface(
		"gravity.v1.EthereumSignature",
		(*EthereumTxConfirmation)(nil),
		&BatchTxConfirmation{},
		&ContractCallTxConfirmation{},
		&SignerSetTxConfirmation{},
	)

	registry.RegisterInterface(
		"gravity.v1.OutgoingTx",
		(*OutgoingTx)(nil),
		&SignerSetTx{},
		&BatchTx{},
		&ContractCallTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func PackEvent(event EthereumEvent) (*types.Any, error) {
	msg, ok := event.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", event)
	}

	anyEvent, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyEvent, nil
}

// UnpackEvent unpacks an Any into an EthereumEvent. It returns an error if the
// event can't be unpacked.
func UnpackEvent(any *types.Any) (EthereumEvent, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	event, ok := any.GetCachedValue().(EthereumEvent)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into EthereumEvent %T", any)
	}

	return event, nil
}

func PackConfirmation(confirmation EthereumTxConfirmation) (*types.Any, error) {
	msg, ok := confirmation.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", confirmation)
	}

	anyEvent, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyEvent, nil
}

// UnpackConfirmation unpacks an Any into a Confirm interface. It returns an error if the
// confirm can't be unpacked.
func UnpackConfirmation(any *types.Any) (EthereumTxConfirmation, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	confirm, ok := any.GetCachedValue().(EthereumTxConfirmation)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into EthereumSignature %T", any)
	}

	return confirm, nil
}

func PackOutgoingTx(outgoing OutgoingTx) (*types.Any, error) {
	msg, ok := outgoing.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", outgoing)
	}

	anyEvent, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyEvent, nil
}

func UnpackOutgoingTx(any *types.Any) (OutgoingTx, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	confirm, ok := any.GetCachedValue().(OutgoingTx)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into OutgoingTx %T", any)
	}

	return confirm, nil
}
