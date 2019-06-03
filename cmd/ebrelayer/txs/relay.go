package txs

// ------------------------------------------------------------
//      Relay
//
//      Builds and encodes EthBridgeClaim Msgs with the
//      specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the Cosmos Bridge.
// ------------------------------------------------------------

import (
	"fmt"

  amino "github.com/tendermint/go-amino"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"

  "github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
  "github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge"
)

func RelayEvent(cdc *amino.Codec, claim *types.EthBridgeClaim) error {

	cliCtx := context.NewCLIContext().
					WithCodec(cdc).
					WithAccountDecoder(cdc)

	txBldr := authtxb.NewTxBuilderFromCLI().
					WithTxEncoder(utils.GetTxEncoder(cdc))

	fmt.Println("\nChecking validator account...")
	err := cliCtx.EnsureAccountExistsFromAddr(claim.Validator)
	if err != nil {
		return err
	}

	fmt.Println("\nChecking recipient account...")
	errRecipient := cliCtx.EnsureAccountExistsFromAddr(claim.CosmosReceiver)
	if errRecipient != nil {
		return errRecipient
	}

	msg := ethbridge.NewMsgMakeEthBridgeClaim(*claim)
	fmt.Println("\nMsg successfully constructed!")
	fmt.Printf("Msg information:\n%+v\n", msg)

	err1 := msg.ValidateBasic()
	if err1 != nil {
		fmt.Println("ERROR validation")
		return err1
	}

	// Add the witnessing validator to the event
	// claimCount := events.ValidatorMakeClaim(hex.EncodeToString(event.Id[:]), validator)
	// fmt.Println("Total claims on this event: ", claimCount)

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
