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

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"

	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
)

func RelayEvent(chainId string, cdc *amino.Codec, passphrase string, claim *types.EthBridgeClaim) error {

	cliCtx := context.NewCLIContext().
		WithCodec(cdc).
		WithAccountDecoder(cdc).
		WithFromAddress(claim.Validator).
		WithFromName("validator")

	cliCtx.SkipConfirm = true

	txBldr := authtxb.NewTxBuilderFromCLI().
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainId)

	err := cliCtx.EnsureAccountExistsFromAddr(claim.Validator)
	if err != nil {
		fmt.Errorf("Validator account error: %s", err)
	}

	msg := ethbridge.NewMsgMakeEthBridgeClaim(*claim)

	err1 := msg.ValidateBasic()
	if err1 != nil {
		fmt.Errorf("Msg validation error: %s", err1)
	}

	cliCtx.PrintResponse = true

	//prepare tx
	txBldr, err = utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign("validator", passphrase, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	cliCtx.PrintOutput(res)
	return err
}
