package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/gorilla/mux"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	restNonce = "nonce"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), makeClaimHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prophecies/{%s}", storeName, restNonce), getProphecyHandler(cdc, cliCtx, storeName)).Methods("GET")
}

type buyNameReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Nonce          int          `json:"nonce"`
	EthereumSender string       `json:"ethereum_sender"`
	CosmosReceiver string       `json:"cosmos_receiver"`
	Validator      string       `json:"validator"`
	Amount         string       `json:"amount"`
}

func makeClaimHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyNameReq

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		ethereumSender := req.EthereumSender
		cosmosReceiver, err2 := sdk.AccAddressFromBech32(req.CosmosReceiver)
		if err2 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err2.Error())
			return
		}
		validator, err3 := sdk.AccAddressFromBech32(req.Validator)
		if err3 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err3.Error())
			return
		}

		amount, err4 := sdk.ParseCoins(req.Amount)
		if err4 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err4.Error())
			return
		}

		// create the message
		msg := oracle.NewMsgMakeBridgeClaim(req.Nonce, ethereumSender, cosmosReceiver, validator, amount)
		err5 := msg.ValidateBasic()
		if err5 != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err5.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func getProphecyHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restNonce]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/prophecy/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
