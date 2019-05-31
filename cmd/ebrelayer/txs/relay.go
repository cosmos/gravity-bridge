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

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	ethbridge "github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge"

	"github.com/cosmos/cosmos-sdk/codec"
)

func RelayEvent(
	cdc *codec.Codec,
	cosmosRecipient sdk.AccAddress,
	validator sdk.AccAddress,
	nonce int,
	ethereumAddress string,
	amount sdk.Coins) error {

	fmt.Printf("\relayEvent() received:\n")
	fmt.Printf("\n Cosmos Recipient: %s, \n Nonce: %d,\n Ethereum Address: %s,\n Amount: %s\n\n",
		cosmosRecipient, nonce, ethereumAddress, amount) //\n Validator: %s, validator

	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithAccountDecoder(cdc)

	txBldr := authtxb.NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc))

	err := cliCtx.EnsureAccountExists()
	if err != nil {
		return err
	}

	ethBridgeClaim := ethbridge.NewEthBridgeClaim(nonce, ethereumAddress, cosmosRecipient, validator, amount)
	msg := ethbridge.NewMsgMakeEthBridgeClaim(ethBridgeClaim)

	err1 := msg.ValidateBasic()
	if err1 != nil {
		return err1
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})

}
