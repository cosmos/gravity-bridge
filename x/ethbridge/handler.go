package ethbridge

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/swishlabsco/peggy/x/ethbridge/types"
	"github.com/swishlabsco/peggy/x/oracle"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(oracleKeeper oracle.Keeper, bankKeeper bank.Keeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateEthBridgeClaim:
			return handleMsgCreateEthBridgeClaim(ctx, cdc, oracleKeeper, bankKeeper, msg, codespace)
		default:
			errMsg := fmt.Sprintf("unrecognized ethbridge message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create a bridge claim
func handleMsgCreateEthBridgeClaim(ctx sdk.Context, cdc *codec.Codec, oracleKeeper oracle.Keeper, bankKeeper bank.Keeper, msg MsgCreateEthBridgeClaim, codespace sdk.CodespaceType) sdk.Result {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String()).Result()
	}
	if msg.Nonce < 0 {
		return types.ErrInvalidEthNonce(codespace).Result()
	}
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(cdc, types.EthBridgeClaim(msg))
	if err != nil {
		return types.ErrJSONMarshalling(codespace).Result()
	}
	status, sdkErr := oracleKeeper.ProcessClaim(ctx, oracleClaim)
	if sdkErr != nil {
		return sdkErr.Result()
	}
	if status.Text == oracle.SuccessStatus {
		sdkErr = processSuccessfulClaim(ctx, bankKeeper, status.FinalClaim)
		if sdkErr != nil {
			return sdkErr.Result()
		}
	}
	return sdk.Result{Log: status.Text}
}

func processSuccessfulClaim(ctx sdk.Context, bankKeeper bank.Keeper, claim string) sdk.Error {
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return err
	}
	receiverAddress := oracleClaim.CosmosReceiver
	_, _, err = bankKeeper.AddCoins(ctx, receiverAddress, oracleClaim.Amount)
	if err != nil {
		return err
	}
	return nil
}
