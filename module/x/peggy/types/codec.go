package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetEthAddress{}, "peggy/MsgSetEthAddress", nil)
	cdc.RegisterConcrete(MsgValsetRequest{}, "peggy/MsgValsetRequest", nil)
	cdc.RegisterConcrete(MsgValsetConfirm{}, "peggy/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(MsgSendToEth{}, "peggy/MsgSendToEth", nil)
	cdc.RegisterConcrete(MsgRequestBatch{}, "peggy/MsgRequestBatch", nil)
	cdc.RegisterConcrete(MsgConfirmBatch{}, "peggy/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(MsgBridgeSignatureSubmission{}, "peggy/MsgBridgeSignatureSubmission", nil)

	cdc.RegisterConcrete(Valset{}, "peggy/Valset", nil)

	cdc.RegisterConcrete(MsgCreateEthereumClaims{}, "peggy/MsgCreateEthereumClaims", nil)
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(EthereumBridgeDepositClaim{}, "peggy/EthereumBridgeDepositClaim", nil)
	cdc.RegisterConcrete(EthereumBridgeWithdrawalBatchClaim{}, "peggy/EthereumBridgeWithdrawalBatchClaim", nil)
	// cdc.RegisterConcrete(EthereumBridgeMultiSigUpdateClaim{}, "peggy/EthereumBridgeMultiSigUpdateClaim", nil)
	// cdc.RegisterConcrete(EthereumBridgeBootstrappedClaim{}, "peggy/types.EthereumBridgeBootstrappedClaim.", nil)

	cdc.RegisterConcrete(OutgoingTxBatch{}, "peggy/OutgoingTxBatch", nil)
	cdc.RegisterConcrete(OutgoingTransferTx{}, "peggy/OutgoingTransferTx", nil)
	cdc.RegisterConcrete(ERC20Token{}, "peggy/ERC20Token", nil)

	cdc.RegisterConcrete(BridgedDenominator{}, "peggy/BridgedDenominator", nil)
	cdc.RegisterConcrete(IDSet{}, "peggy/IDSet", nil)

	cdc.RegisterConcrete(Attestation{}, "peggy/Attestation", nil)
	cdc.RegisterInterface((*AttestationDetails)(nil), nil)
	cdc.RegisterConcrete(BridgeDeposit{}, "peggy/BridgeDeposit", nil)
	// cdc.RegisterConcrete(SignedCheckpoint{}, "peggy/SignedCheckpoint", nil)
	// cdc.RegisterConcrete(BridgeBootstrap{}, "peggy/BridgeBootstrap", nil)
}
