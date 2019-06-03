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

func RelayEvent(chainId string, cdc *amino.Codec, claim *types.EthBridgeClaim) error {

	cliCtx := context.NewCLIContext().
					WithCodec(cdc).
					WithAccountDecoder(cdc).
					WithFromAddress(claim.Validator).
					WithFromName("validator")

	txBldr := authtxb.NewTxBuilderFromCLI().
					WithTxEncoder(utils.GetTxEncoder(cdc)).
					WithChainID(chainId)

	fmt.Println("\nChecking validator account...")
	err := cliCtx.EnsureAccountExistsFromAddr(claim.Validator)
	if err != nil {
		fmt.Errorf("Validator account error: %s", err)
	}

	msg := ethbridge.NewMsgMakeEthBridgeClaim(*claim)
	fmt.Println("\nMsg successfully constructed!")
	fmt.Printf("Msg information:\n%+v\n", msg)

	err1 := msg.ValidateBasic()
	if err1 != nil {
		fmt.Errorf("Msg validation error: %s", err1)
	}

	cliCtx.PrintResponse = true

	address, keybaseName, _ := context.GetFromFields(string(claim.Validator))
	fmt.Println("address: ", address)
	fmt.Println("keybase: ", keybaseName)

	// Add the witnessing validator to the event
	// claimCount := events.ValidatorMakeClaim(hex.EncodeToString(event.Id[:]), validator)
	// fmt.Println("Total claims on this event: ", claimCount)

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
