package ethbridge

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/common"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

// NewHandler returns a handler for "ethbridge" type messages.
func NewHandler(oracleKeeper oracle.Keeper, bankKeeper bank.Keeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMakeEthBridgeClaim:
			return handleMsgMakeEthBridgeClaim(ctx, cdc, oracleKeeper, bankKeeper, msg, codespace)
		default:
			errMsg := fmt.Sprintf("Unrecognized ethbridge message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to make a bridge claim
func handleMsgMakeEthBridgeClaim(ctx sdk.Context, cdc *codec.Codec, oracleKeeper oracle.Keeper, bankKeeper bank.Keeper, msg MsgMakeEthBridgeClaim, codespace sdk.CodespaceType) sdk.Result {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String()).Result()
	}
	if msg.Nonce < 0 {
		return types.ErrInvalidEthNonce(codespace).Result()
	}
	if !common.IsValidEthAddress(msg.EthereumSender) {
		return types.ErrInvalidEthAddress(codespace).Result()
	}
	oracleId, validator, claimString := types.CreateOracleClaimFromEthClaim(cdc, msg.EthBridgeClaim)
	status, err := oracleKeeper.ProcessClaim(ctx, oracleId, validator, claimString)
	if err != nil {
		return err.Result()
	}
	if status.StatusText == oracle.SuccessStatus {
		err = processSuccessfulClaim(ctx, bankKeeper, status.FinalClaim)
		if err != nil {
			return err.Result()
		}
	}
	return sdk.Result{Log: status.StatusText}
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
