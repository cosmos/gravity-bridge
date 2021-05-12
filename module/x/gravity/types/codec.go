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
		&MsgConfirmLogicCall{},
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgDelegateKeys{},
		&MsgLogicCallExecutedClaim{},
		&MsgValsetUpdatedClaim{},
		&MsgCancelSendToEthereum{},
		&MsgSubmitBadSignatureEvidence{},
	)

	registry.RegisterInterface(
		"gravity.v1beta1.EthereumClaim",
		(*EthereumClaim)(nil),
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgLogicCallExecutedClaim{},
		&MsgValsetUpdatedClaim{},
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
	cdc.RegisterConcrete(&MsgConfirmLogicCall{}, "gravity/MsgConfirmLogicCall", nil)
	cdc.RegisterConcrete(&Valset{}, "gravity/Valset", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "gravity/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "gravity/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&MsgERC20DeployedClaim{}, "gravity/MsgERC20DeployedClaim", nil)
	cdc.RegisterConcrete(&MsgLogicCallExecutedClaim{}, "gravity/MsgLogicCallExecutedClaim", nil)
	cdc.RegisterConcrete(&MsgValsetUpdatedClaim{}, "gravity/MsgValsetUpdatedClaim", nil)
	cdc.RegisterConcrete(&BatchTx{}, "gravity/BatchTx", nil)
	cdc.RegisterConcrete(&MsgCancelSendToEthereum{}, "gravity/MsgCancelSendToEthereum", nil)
	cdc.RegisterConcrete(&SendToEthereum{}, "gravity/SendToEthereum", nil)
	cdc.RegisterConcrete(&ERC20Token{}, "gravity/ERC20Token", nil)
	cdc.RegisterConcrete(&IDSet{}, "gravity/IDSet", nil)
	cdc.RegisterConcrete(&EthereumEventVoteRecord{}, "gravity/EthereumEventVoteRecord", nil)
	cdc.RegisterConcrete(&MsgSubmitBadSignatureEvidence{}, "gravity/MsgSubmitBadSignatureEvidence", nil)
}
