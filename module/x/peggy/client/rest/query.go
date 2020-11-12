package rest

import (
	"fmt"
	"net/http"

	"github.com/althea-net/peggy/module/x/peggy/keeper"
	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// USED BY RUST //
func getValsetRequestHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetRequest/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func batchByNonceHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batch/%s/%s", storeName, nonce, denom))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastBatchesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastBatches", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset not found")
			return
		}

		var out []types.OutgoingTxBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST //
// gets all the confirm messages for a given validator set nonce
func allValsetConfirmsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirms/%s", storeName, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		var out []types.MsgValsetConfirm
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// gets all the confirm messages for a given transaction batch
func allBatchConfirmsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[nonce]
		denom := vars[tokenAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/batchConfirms/%s/%s", storeName, nonce, denom))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset confirms not found")
			return
		}

		var out []types.MsgConfirmBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST //
func lastValsetRequestsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastValsetRequests", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "valset requests not found")
			return
		}

		var out []types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// USED BY RUST //
func lastValsetRequestsByAddressHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingValsetRequest/%s", storeName, operatorAddr))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastBatchesByAddressHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		operatorAddr := vars[bech32ValidatorAddress]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastPendingBatchRequest/%s", storeName, operatorAddr))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no pending valset requests found")
			return
		}

		var out types.OutgoingTxBatch
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func currentValsetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/currentValset", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

type valsetConfirmRequest struct {
	Nonce   string `json:"nonce"`
	Address string `json:"address"`
}

func getValsetConfirmHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req valsetConfirmRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/valsetConfirm/%s/%s", storeName, req.Nonce, req.Address))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "confirmation not found")
			return
		}

		var out types.MsgValsetConfirm
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastObservedNonceHandler(cliCtx context.CLIContext, storeName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claimType := vars[claimType]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastNonce/%s", storeName, claimType))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no observed nonce found for this type")
			return
		}

		var out types.UInt64Nonce
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func lastObservedNoncesHandler(cliCtx context.CLIContext, storeName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastObservedNonces", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no observed nonces")
			return
		}

		var out map[string]types.UInt64Nonce
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// The last valset that was observed to be on the Ethereum chain
func lastObservedValsetHandler(cliCtx context.CLIContext, storeName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastObservedMultiSigUpdate", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no observed multisig update")
			return
		}

		var out keeper.MultiSigUpdateResponse
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

// approved = signed by >66% of the validator set
func lastApprovedValsetHandler(cliCtx context.CLIContext, storeName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/lastApprovedMultiSigUpdate", storeName))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no approved multisig update")
			return
		}

		var out keeper.MultiSigUpdateResponse
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

func queryAttestation(cliCtx context.CLIContext, storeName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		claimType := vars[claimType]
		nonce := vars[nonce]

		res, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/attestation/%s/%s", storeName, claimType, nonce))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "no attestation")
			return
		}

		var out types.Attestation
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}
