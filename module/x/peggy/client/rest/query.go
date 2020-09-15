package rest

import (
	"fmt"
	"net/http"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

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

		// why doesn't this cliCtx update it's height on it's own?
		// looks like the sdk uses client context TODO investigate
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
