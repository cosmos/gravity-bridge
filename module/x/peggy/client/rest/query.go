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
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/currentValset", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
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

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeName, nonce), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// why doesn't this cliCtx update it's height on it's own?
		// looks like the sdk uses client context TODO investigate
		var out types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx.WithHeight(height), res)
	}
}

type getValsetConfirm struct {
	Nonce   string `json:"nonce"`
	Address string `json:"address"`
}

func getValsetConfirmHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getValsetConfirm

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		nonce := req.Nonce
		bech32ValidatorAddress := req.Address

		// this panics too?
		_, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetConfirm/%s", storeName, nonce), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		return
		// this panics why and how do I guard for it?
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetConfirm/%s/%s", storeName, nonce, bech32ValidatorAddress), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var out []types.MsgValsetConfirm
		cliCtx.Codec.MustUnmarshalJSON(res, &out)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
