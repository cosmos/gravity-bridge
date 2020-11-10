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

// Here are the routes that are actually queried by the rust
// "peggy/valset_request/{}"
// "peggy/pending_valset_requests/{}"
// "peggy/valset_requests"
// "peggy/valset_confirm/{}"
// "peggy/pending_batch_requests/{}" UNIMPLEMENTED???
// "peggy/transaction_batches/" UNIMPLEMENTED???
// "peggy/signed_batches" UNIMPLEMENTED???

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/current_valset", storeName), currentValsetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/valset_request/{%s}", storeName, nonce), getValsetRequestHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/query_valset_confirm", storeName), getValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/update_ethaddr", storeName), updateEthAddressHandler(cliCtx)).Methods("POST")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/valset_request", storeName), createValsetRequestHandler(cliCtx)).Methods("POST")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/valset_confirm", storeName), createValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm/{%s}", storeName, nonce), allValsetConfirmsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/valset_requests", storeName), lastValsetRequestsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pending_valset_requests/{%s}", storeName, bech32ValidatorAddress), lastValsetRequestsByAddressHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/lastNonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonces", storeName), lastObservedNoncesHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/last_observed_valset", storeName), lastObservedValsetHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/last_approved_valset", storeName), lastApprovedValsetHandler(cliCtx, storeName)).Methods("GET")
	/* UNUSED */ r.HandleFunc(fmt.Sprintf("/%s/attestation/{%s}/{%s}", storeName, claimType, nonce), queryAttestation(cliCtx, storeName)).Methods("GET")

}
