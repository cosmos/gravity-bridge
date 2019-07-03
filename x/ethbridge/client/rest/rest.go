package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/gorilla/mux"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/peggy/x/ethbridge"
	"github.com/swishlabsco/peggy/x/ethbridge/querier"
	"github.com/swishlabsco/peggy/x/ethbridge/types"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	restNonce          = "nonce"
	restEthereumSender = "ethereumSender"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, queryRoute string) {
	r.HandleFunc(fmt.Sprintf("/%s/prophecies", queryRoute), createClaimHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prophecies/{%s}/{%s}", queryRoute, restNonce, restEthereumSender), getProphecyHandler(cdc, cliCtx, queryRoute)).Methods("GET")
}

type createEthClaimReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Nonce          int          `json:"nonce"`
	EthereumSender string       `json:"ethereum_sender"`
	CosmosReceiver string       `json:"cosmos_receiver"`
	Validator      string       `json:"validator"`
	Amount         string       `json:"amount"`
}

func createClaimHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createEthClaimReq

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		ethereumSender := gethCommon.HexToAddress(req.EthereumSender)
		cosmosReceiver, err := sdk.AccAddressFromBech32(req.CosmosReceiver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		validator, err := sdk.ValAddressFromBech32(req.Validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		ethBridgeClaim := types.NewEthBridgeClaim(req.Nonce, ethereumSender, cosmosReceiver, validator, amount)
		msg := ethbridge.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func getProphecyHandler(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[restNonce]

		nonceString, err := strconv.Atoi(nonce)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		ethereumSender := vars[restEthereumSender]

		bz, err := cdc.MarshalJSON(ethbridge.NewQueryEthProphecyParams(nonceString, ethereumSender))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", queryRoute, querier.QueryEthProphecy)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
