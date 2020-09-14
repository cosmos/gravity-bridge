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

	cdc.RegisterConcrete(Valset{}, "peggy/Valset", nil)

	cdc.RegisterConcrete(MsgCreateEthereumClaims{}, "peggy/MsgCreateEthereumClaims", nil)
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(EthereumBridgeDepositClaim{}, "peggy/EthereumBridgeDepositClaim", nil)
	cdc.RegisterConcrete(EthereumBridgeWithdrawalBatchClaim{}, "peggy/EthereumBridgeWithdrawalBatchClaim", nil)
	cdc.RegisterConcrete(EthereumBridgeMultiSigUpdateClaim{}, "peggy/EthereumBridgeMultiSigUpdateClaim", nil)
}
