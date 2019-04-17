package ethbridge

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/common"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(oracleKeeper oracle.Keeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMakeEthBridgeClaim:
			return handleMsgMakeEthBridgeClaim(ctx, cdc, oracleKeeper, msg, codespace)
		default:
			errMsg := fmt.Sprintf("Unrecognized ethbridge message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to make a bridge claim
func handleMsgMakeEthBridgeClaim(ctx sdk.Context, cdc *codec.Codec, oracleKeeper oracle.Keeper, msg MsgMakeEthBridgeClaim, codespace sdk.CodespaceType) sdk.Result {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String()).Result()
	}
	if msg.Nonce < 0 {
		return types.ErrInvalidEthNonce(codespace).Result()
	}
	if !common.IsValidEthAddress(msg.EthereumSender) {
		return types.ErrInvalidEthAddress(codespace).Result()
	}
	claim := types.CreateOracleClaimFromEthClaim(cdc, msg.EthBridgeClaim)
	result, err := oracleKeeper.ProcessClaim(ctx, claim)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{Log: result.Status}
}
