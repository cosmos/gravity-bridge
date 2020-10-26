package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	nonce                  = "nonce"
	bech32ValidatorAddress = "bech32ValidatorAddress"
	claimType              = "claimType"
	signType               = "signType"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/current_valset", storeName), currentValsetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/valset_request/{%s}", storeName, nonce), getValsetRequestHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/query_valset_confirm", storeName), getValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/update_ethaddr", storeName), updateEthAddressHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_request", storeName), createValsetRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm", storeName), createValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm/{%s}", storeName, nonce), allValsetConfirmsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/valset_requests", storeName), lastValsetRequestsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pending_valset_requests/{%s}", storeName, bech32ValidatorAddress), lastValsetRequestsByAddressHandler(cliCtx, storeName)).Methods("GET")
	// deprecated
	r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/lastNonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonces", storeName), lastObservedNoncesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/last_observed_valset", storeName), lastObservedValsetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/last_approved_valset", storeName), lastApprovedValsetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/attestation/{%s}/{%s}", storeName, claimType, nonce), queryAttestation(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/sign_bridge_request/{%s}/{%s}", storeName, signType, nonce), BridgeApprovalSignatureHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/bootstrap", storeName), bootstrapConfirmHandler(cliCtx)).Methods("POST")

}
