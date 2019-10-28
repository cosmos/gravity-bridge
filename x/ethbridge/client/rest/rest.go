package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"

	"github.com/cosmos/peggy/x/ethbridge/querier"
	"github.com/cosmos/peggy/x/ethbridge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	restNonce          = "nonce"
	restEthereumSender = "ethereumSender"
)

type createEthClaimReq struct {
	BaseReq        rest.BaseReq `json:"base_req"`
	Nonce          int          `json:"nonce"`
	EthereumSender string       `json:"ethereum_sender"`
	CosmosReceiver string       `json:"cosmos_receiver"`
	Validator      string       `json:"validator"`
	Amount         string       `json:"amount"`
	ClaimType      string       `json:"claim_type"`
}

type burnEthReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	EthereumChainID  string       `json:"ethereum_chain_id"`
	Token            string       `json:"token"`
	CosmosSender     string       `json:"cosmos_sender"`
	EthereumReceiver string       `json:"ethereum_receiver"`
	Amount           string       `json:"amount"`
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), createClaimHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prophecies/{%s}/{%s}", storeName, restNonce, restEthereumSender), getProphecyHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/burn", storeName), burnHandler(cliCtx)).Methods("POST")
}

func createClaimHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createEthClaimReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		ethereumSender := types.NewEthereumAddress(req.EthereumSender)

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

		var claimType types.ClaimType
		if value, ok := types.StringToClaimType[req.ClaimType]; ok {
			claimType = value
		} else {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidClaimType().Error())
			return
		}

		// create the message
		ethBridgeClaim := types.NewEthBridgeClaim(req.Nonce, ethereumSender, cosmosReceiver, validator, amount, claimType)
		msg := types.NewMsgCreateEthBridgeClaim(ethBridgeClaim)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func getProphecyHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonce := vars[restNonce]

		nonceString, err := strconv.Atoi(nonce)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		ethereumSender := types.NewEthereumAddress(vars[restEthereumSender])

		bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryEthProphecyParams(nonceString, ethereumSender))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, querier.QueryEthProphecy)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func burnHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnEthReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		ethereumChainID, err := strconv.Atoi(req.EthereumChainID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		token := types.NewEthereumAddress(req.Token)

		cosmosSender, err := sdk.AccAddressFromBech32(req.CosmosSender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ethereumReceiver := types.NewEthereumAddress(req.EthereumReceiver)

		amount, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgBurn(ethereumChainID, token, cosmosSender, ethereumReceiver, amount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
