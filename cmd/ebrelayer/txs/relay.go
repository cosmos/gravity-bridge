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

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
