package oracle

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	//else (if doesnt exist yet)
	id := strconv.Itoa(msg.Nonce) + msg.EthereumSender
	bridgeClaim := NewBridgeClaim(id, msg.CosmosReceiver, msg.Validator, msg.Amount)
	bridgeClaims := []BridgeClaim{bridgeClaim}
	newProphecy := NewBridgeProphecy(id, PendingStatus, getMinimumClaims(), bridgeClaims)
	err := keeper.CreateProphecy(ctx, newProphecy)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func getMinimumClaims() int {
	minimumClaims := float64(getTotalNumberValidators()) * DefaultConsensusNeeded
	return int(math.Ceil(minimumClaims))

}

func getTotalNumberValidators() int {
	//TODO: Get from Tendermint/last block?
	return 10
}
