package oracle

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/common"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/types"
)

// NewHandler returns a handler for "oracle" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMakeBridgeEthClaim:
			return handleMsgMakeBridgeEthClaim(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized oracle message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to make a bridge claim
func handleMsgMakeBridgeEthClaim(ctx sdk.Context, keeper Keeper, msg MsgMakeBridgeEthClaim) sdk.Result {
	if msg.CosmosReceiver.Empty() {
		return sdk.ErrInvalidAddress(msg.CosmosReceiver.String()).Result()
	}
	if msg.Nonce < 0 {
		return types.ErrInvalidEthereumNonce(keeper.Codespace()).Result()
	}
	if !common.IsValidEthereumAddress(msg.EthereumSender) {
		return types.ErrInvalidEthereumAddress(keeper.Codespace()).Result()
	}
	id := strconv.Itoa(msg.Nonce) + msg.EthereumSender
	_, err := keeper.GetProphecy(ctx, id)
	if err == nil {
		return sdk.ErrInternal("Not yet implemented").Result()
		//check if prophecy exists or not
		//if exist
		//	//get it and continue checks
		//	//check if claim for this validator exists or not
		//	//if does
		//		//return error
		//	//else
		//		//add claim to list
		//	//check if claimthreshold is passed
		//	//if does
		//		//check enough claims match and are valid
		//		//update prophecy to be successful
		//		//trigger minting
		//		//save finalized prophecy to db
		//		//return
		//	//if doesnt
		//		//save updated prophecy to db
		//		//return
	} else {
		if err.Code() != types.CodeNotFound {
			return err.Result()
		}
		bridgeClaim := NewBridgeClaim(id, msg.CosmosReceiver, msg.Validator, msg.Amount)
		bridgeClaims := []BridgeClaim{bridgeClaim}
		newProphecy := NewBridgeProphecy(id, PendingStatus, getPowerThreshold(), bridgeClaims)
		err := keeper.CreateProphecy(ctx, newProphecy)
		if err != nil {
			return err.Result()
		}
		return sdk.Result{}
	}
}

func getPowerThreshold() int {
	minimumPower := float64(getTotalPower()) * DefaultConsensusNeeded
	return int(math.Ceil(minimumPower))

}

func getTotalPower() int {
	//TODO: Get from Tendermint/last block/staking module?
	return 10
}
