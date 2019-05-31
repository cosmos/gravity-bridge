package txs

// ------------------------------------------------------------
//      Root
//
//      Registers REST routes for use by ebrelayer.
// ------------------------------------------------------------

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	// TODO: add '/txs/relay' cmd line support
	// r.HandleFunc("/txs/relay", relayEvent(cdc, cliCtx)).Methods("POST")
}
