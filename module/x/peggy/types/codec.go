package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	proto "github.com/gogo/protobuf/proto"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterInterfaces regiesteres the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgValsetConfirm{},
		&MsgValsetRequest{},
		&MsgSetEthAddress{},
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgCreateEthereumClaims{},
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
	)

	registry.RegisterInterface(
		"peggy.v1beta1.EthereumClaim",
		(*EthereumClaim)(nil),
		&DepositClaim{},
		&WithdrawClaim{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// PackEthereumClaim constructs a new Any packed with the eth claim value. It returns
// an error if the client state can't be casted to a protobuf message or if the concrete
// implemention is not registered to the protobuf codec.
func PackEthereumClaim(ad EthereumClaim) (*types.Any, error) {
	msg, ok := ad.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", ad)
	}

	anyad, err := types.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyad, nil
}

// UnpackEthereumClaim unpacks an Any into a EthereumClaim. It returns an error if the
// attestation details can't be unpacked into a EthereumClaim.
func UnpackEthereumClaim(any *types.Any) (EthereumClaim, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	ethereumClaim, ok := any.GetCachedValue().(EthereumClaim)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into EthereumClaim %T", any)
	}

	return ethereumClaim, nil
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(&MsgSetEthAddress{}, "peggy/MsgSetEthAddress", nil)
	cdc.RegisterConcrete(&MsgValsetRequest{}, "peggy/MsgValsetRequest", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "peggy/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(&MsgSendToEth{}, "peggy/MsgSendToEth", nil)
	cdc.RegisterConcrete(&MsgRequestBatch{}, "peggy/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, "peggy/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(&Valset{}, "peggy/Valset", nil)
	cdc.RegisterConcrete(&MsgCreateEthereumClaims{}, "peggy/MsgCreateEthereumClaims", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "peggy/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "peggy/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&DepositClaim{}, "peggy/DepositClaim", nil)
	cdc.RegisterConcrete(&WithdrawClaim{}, "peggy/WithdrawClaim", nil)
	cdc.RegisterConcrete(&OutgoingTxBatch{}, "peggy/OutgoingTxBatch", nil)
	cdc.RegisterConcrete(&OutgoingTransferTx{}, "peggy/OutgoingTransferTx", nil)
	cdc.RegisterConcrete(&ERC20Token{}, "peggy/ERC20Token", nil)
	cdc.RegisterConcrete(&IDSet{}, "peggy/IDSet", nil)
	cdc.RegisterConcrete(&Attestation{}, "peggy/Attestation", nil)
}
