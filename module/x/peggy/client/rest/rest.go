package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	nonce                  = "nonce"
	bech32ValidatorAddress = "bech32ValidatorAddress"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/current_valset", storeName), currentValsetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/valset_request/{%s}", storeName, nonce), getValsetRequestHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/query_valset_confirm", storeName), getValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/update_ethaddr", storeName), updateEthAddressHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_request", storeName), createValsetRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm", storeName), createValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
}
