package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"

	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	restEthereumChainID = "ethereumChainID"
	restBridgeContract  = "bridgeContract"
	restNonce           = "nonce"
	restSymbol          = "symbol"
	restTokenContract   = "tokenContract"
	restEthereumSender  = "ethereumSender"
)

type createNFTClaimReq struct {
	BaseReq               rest.BaseReq `json:"base_req"`
	EthereumChainID       int          `json:"ethereum_chain_id"`
	BridgeContractAddress string       `json:"bridge_contract_address"`
	Nonce                 int          `json:"nonce"`
	Symbol                string       `json:"symbol"`
	TokenContractAddress  string       `json:"token_contract_address"`
	EthereumSender        string       `json:"ethereum_sender"`
	CosmosReceiver        string       `json:"cosmos_receiver"`
	Validator             string       `json:"validator"`
	Denom                 string       `json:"denom"`
	ID                    string       `json:"id"`
	ClaimType             string       `json:"claim_type"`
}

type burnOrLockNFTReq struct {
	BaseReq          rest.BaseReq `json:"base_req"`
	EthereumChainID  string       `json:"ethereum_chain_id"`
	TokenContract    string       `json:"token_contract_address"`
	CosmosSender     string       `json:"cosmos_sender"`
	EthereumReceiver string       `json:"ethereum_receiver"`
	Denom            string       `json:"denom"`
	ID               string       `json:"id"`
}

// RegisterRESTRoutes - Central function to define routes that get registered by the main application
func RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prophecies", storeName), createClaimHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prophecies/{%s}/{%s}/{%s}/{%s}/{%s}/{%s}", storeName, restEthereumChainID, restBridgeContract, restNonce, restSymbol, restTokenContract, restEthereumSender), getProphecyHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/burn", storeName), burnOrLockHandler(cliCtx, "burn")).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/lock", storeName), burnOrLockHandler(cliCtx, "lock")).Methods("POST")
}

func createClaimHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createNFTClaimReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		bridgeContractAddress := ethbridge.NewEthereumAddress(req.BridgeContractAddress)

		tokenContractAddress := ethbridge.NewEthereumAddress(req.TokenContractAddress)

		ethereumSender := ethbridge.NewEthereumAddress(req.EthereumSender)

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
		denom := req.Denom
		id := req.ID

		claimType, err := ethbridge.StringToClaimType(req.ClaimType)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, types.ErrInvalidClaimType.Error())
			return
		}

		// create the message
		nftBridgeClaim := types.NewNFTBridgeClaim(req.EthereumChainID, bridgeContractAddress, req.Nonce, req.Symbol, tokenContractAddress, ethereumSender, cosmosReceiver, validator, denom, id, claimType)
		msg := types.NewMsgCreateNFTBridgeClaim(nftBridgeClaim)
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

		ethereumChainID := vars[restEthereumChainID]
		ethereumChainIDString, err := strconv.Atoi(ethereumChainID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		bridgeContract := ethbridge.NewEthereumAddress(vars[restBridgeContract])

		nonce := vars[restNonce]
		nonceString, err := strconv.Atoi(nonce)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tokenContract := ethbridge.NewEthereumAddress(vars[restTokenContract])

		symbol := vars[restSymbol]
		if strings.TrimSpace(symbol) == "" {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		ethereumSender := ethbridge.NewEthereumAddress(vars[restEthereumSender])

		bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryNFTProphecyParams(ethereumChainIDString, bridgeContract, nonceString, symbol, tokenContract, ethereumSender))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryNFTProphecy)
		res, _, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func burnOrLockHandler(cliCtx context.CLIContext, lockOrBurn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnOrLockNFTReq

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

		tokenContract := ethbridge.NewEthereumAddress(req.TokenContract)

		cosmosSender, err := sdk.AccAddressFromBech32(req.CosmosSender)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ethereumReceiver := ethbridge.NewEthereumAddress(req.EthereumReceiver)

		denom := req.Denom
		id := req.ID

		// create the message
		var msg sdk.Msg
		switch lockOrBurn {
		case "lock":
			msg = types.NewMsgLockNFT(ethereumChainID, tokenContract, cosmosSender, ethereumReceiver, denom, id)
		case "burn":
			msg = types.NewMsgBurnNFT(ethereumChainID, tokenContract, cosmosSender, ethereumReceiver, denom, id)
		}
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
