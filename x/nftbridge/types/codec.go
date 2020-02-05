package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateNFTBridgeClaim{}, "nftbridge/MsgCreateNFTBridgeClaim", nil)
	cdc.RegisterConcrete(MsgBurnNFT{}, "nftbridge/MsgBurnNFT", nil)
	cdc.RegisterConcrete(MsgLockNFT{}, "nftbridge/MsgLockNFT", nil)
}
