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
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {

	// validator set endpoints

	// Provides the current validator set with powers and eth addresses, useful to check the current validator state
	// but not actually used to operate the bridge
	r.HandleFunc(fmt.Sprintf("/%s/current_valset", storeName), currentValsetHandler(cliCtx, storeName)).Methods("GET")
	// gets valset request by nonce, mostly unused in favor of /valset_requests which provides the most recent validator sets
	// since the validators only really care about upgrading to the latest available validator set. Might be useful for a relayer
	// who wants to relay historical validator sets
	r.HandleFunc(fmt.Sprintf("/%s/valset_request/{%s}", storeName, nonce), getValsetRequestHandler(cliCtx, storeName)).Methods("GET")
	// gets the latest 5 validator set requests, used heavily by the relayer. Which hits this endpoint before checking which
	// of these last 5 have sufficient signatures to relay
	r.HandleFunc(fmt.Sprintf("/%s/valset_requests", storeName), lastValsetRequestsHandler(cliCtx, storeName)).Methods("GET")
	// This endpoint gets all of the validator set confirmations for a given nonce. In order to determine if a valset is complete
	// the relayer queries the latest valsets and then compares the number of members they show versus the length of this endpoints output
	// if they match every validator has submitted a signature and we can go forward with relaying that validator set update.
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm/{%s}", storeName, nonce), allValsetConfirmsHandler(cliCtx, storeName)).Methods("GET")
	// The Ethereum signer queries this endpoint and signs whatever it returns once per loop iteration
	r.HandleFunc(fmt.Sprintf("/%s/pending_valset_requests/{%s}", storeName, bech32ValidatorAddress), lastValsetRequestsByAddressHandler(cliCtx, storeName)).Methods("GET")

	// batch endpoints

	// The Ethereum signer queries this endpoint and signs whatever it returns once per loop iteration
	r.HandleFunc(fmt.Sprintf("/%s/pending_batch_requests/{%s}", storeName, bech32ValidatorAddress), lastBatchRequestsByAddressHandler(cliCtx, storeName)).Methods("GET")
	// Gets all outgoing batches in the batch queue, up to 100
	r.HandleFunc(fmt.Sprintf("/%s/transaction_batches", storeName), lastBatchRequestsHandler(cliCtx, storeName)).Methods("GET")
	// Gets all outgoing batches in the batch queue, up to 100
	r.HandleFunc(fmt.Sprintf("/%s/transaction_batch/{%s}", storeName, nonce), lastBatchRequestByNonceHandler(cliCtx, storeName)).Methods("GET")
	// Gets all outgoing batches that have signatures up to 100, these may or may not be possible to submit it's up to the relayer
	// to attempt that.
	r.HandleFunc(fmt.Sprintf("/%s/signed_batches", storeName), lastSignedBatchesHandler(cliCtx, storeName)).Methods("GET")

	// Ethereum oracle endpoints

	// This is the nonce of the last processed attestation nonce of the provided claimType from Ethereum. For example if you requested the last
	// eth deposit claim you'd get the last eht deposit nonce
	r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")
	// Returns the last nonce of each Ethereum oracle claim that exists.
	r.HandleFunc(fmt.Sprintf("/%s/last_observed_nonces", storeName), lastObservedNoncesHandler(cliCtx, storeName)).Methods("GET")
	// Returns the last validator set nonce the Ethereum oracle process has observed. This is somewhat redundant with the above endpoint
	r.HandleFunc(fmt.Sprintf("/%s/last_observed_valset", storeName), lastObservedValsetHandler(cliCtx, storeName)).Methods("GET")
	// A generic endpoint to the actual attestation message for any given type and nonce
	r.HandleFunc(fmt.Sprintf("/%s/attestation/{%s}/{%s}", storeName, claimType, nonce), queryAttestation(cliCtx, storeName)).Methods("GET")

	// Transaction generation endpoints

	// Helper endpoint for generating an update ethaddr transaction, not used in production but useful for testing and debugging
	r.HandleFunc(fmt.Sprintf("/%s/update_ethaddr", storeName), updateEthAddressHandler(cliCtx)).Methods("POST")
	// Helper endpoint for generating a valset request transaction, not used in production but useful for testing and debugging
	r.HandleFunc(fmt.Sprintf("/%s/valset_request", storeName), createValsetRequestHandler(cliCtx)).Methods("POST")
	// Helper endpoint for generating a valset confirm transaction, not used in production but useful for testing and debugging
	r.HandleFunc(fmt.Sprintf("/%s/valset_confirm", storeName), createValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
	// Helper endpoint for generating either a transaction batch signature or a validator set signature transaction, useful for testing and debugging
	r.HandleFunc(fmt.Sprintf("/%s/sign_bridge_request/{%s}/{%s}", storeName, claimType, nonce), BridgeApprovalSignatureHandler(cliCtx)).Methods("POST")
	// Helper endpoint for generating a bootstrap transaction, useful for testing or debugging
	r.HandleFunc(fmt.Sprintf("/%s/bootstrap", storeName), bootstrapConfirmHandler(cliCtx)).Methods("POST")

	// Cruft

	// duplicate of last_observed_nonce
	r.HandleFunc(fmt.Sprintf("/%s/lastNonce/{%s}", storeName, claimType), lastObservedNonceHandler(cliCtx, storeName)).Methods("GET")
	// this endpoint will be removed, since this is now the same as the last valset with any signatures.
	r.HandleFunc(fmt.Sprintf("/%s/last_approved_valset", storeName), lastApprovedValsetHandler(cliCtx, storeName)).Methods("GET")
	// this endpoint is unused it requires the Cosmos address of the other validators, not the eth address. Since no endpoint provides
	// a public mapping from Cosmos to Ethereum address for other validators it can't be practically used without a lot of extra
	// queries.
	r.HandleFunc(fmt.Sprintf("/%s/query_valset_confirm", storeName), getValsetConfirmHandler(cliCtx, storeName)).Methods("POST")
}
