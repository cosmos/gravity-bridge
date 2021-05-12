package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.NewLegacyAmino()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterInterfaces registers the interfaces for the proto stuff
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSignerSetTxSignature{},
		&MsgSendToEthereum{},
		&MsgRequestBatchTx{},
		&MsgBatchTxSignature{},
		&MsgContractCallTxSignature{},
		&MsgSendToCosmosEvent{},
		&MsgBatchExecutedEvent{},
		&MsgERC20DeployedEvent{},
		&MsgDelegateKeys{},
		&MsgContractCallTxExecutedEvent{},
		&MsgSignerSetUpdatedEvent{},
		&MsgCancelSendToEthereum{},
		&MsgSubmitBadSignatureEvidence{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.EthereumClaim",
		(*EthereumClaim)(nil),
		&MsgSendToCosmosEvent{},
		&MsgBatchExecutedEvent{},
		&MsgERC20DeployedEvent{},
		&MsgContractCallTxExecutedEvent{},
		&MsgSignerSetUpdatedEvent{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(&MsgDelegateKeys{}, "gravity/MsgDelegateKeys", nil)
	cdc.RegisterConcrete(&MsgSignerSetTxSignature{}, "gravity/MsgSignerSetTxSignature", nil)
	cdc.RegisterConcrete(&MsgSendToEthereum{}, "gravity/MsgSendToEthereum", nil)
	cdc.RegisterConcrete(&MsgRequestBatchTx{}, "gravity/MsgRequestBatchTx", nil)
	cdc.RegisterConcrete(&MsgBatchTxSignature{}, "gravity/MsgBatchTxSignature", nil)
	cdc.RegisterConcrete(&MsgContractCallTxSignature{}, "gravity/MsgContractCallTxSignature", nil)
	cdc.RegisterConcrete(&Valset{}, "gravity/Valset", nil)
	cdc.RegisterConcrete(&MsgSendToCosmosEvent{}, "gravity/MsgSendToCosmosEvent", nil)
	cdc.RegisterConcrete(&MsgBatchExecutedEvent{}, "gravity/MsgBatchExecutedEvent", nil)
	cdc.RegisterConcrete(&MsgERC20DeployedEvent{}, "gravity/MsgERC20DeployedEvent", nil)
	cdc.RegisterConcrete(&MsgContractCallTxExecutedEvent{}, "gravity/MsgContractCallTxExecutedEvent", nil)
	cdc.RegisterConcrete(&MsgSignerSetUpdatedEvent{}, "gravity/MsgSignerSetUpdatedEvent", nil)
	cdc.RegisterConcrete(&BatchTx{}, "gravity/BatchTx", nil)
	cdc.RegisterConcrete(&MsgCancelSendToEthereum{}, "gravity/MsgCancelSendToEthereum", nil)
	cdc.RegisterConcrete(&SendToEthereum{}, "gravity/SendToEthereum", nil)
	cdc.RegisterConcrete(&ERC20Token{}, "gravity/ERC20Token", nil)
	cdc.RegisterConcrete(&IDSet{}, "gravity/IDSet", nil)
	cdc.RegisterConcrete(&EthereumEventVoteRecord{}, "gravity/EthereumEventVoteRecord", nil)
	cdc.RegisterConcrete(&MsgSubmitBadSignatureEvidence{}, "gravity/MsgSubmitBadSignatureEvidence", nil)
}
