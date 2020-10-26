package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	hexUtil "github.com/ethereum/go-ethereum/common/hexutil"
)

type updateEthAddressReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	EthSig  string       `json:"ethSig"`
}

// accepts a sig proving that the given Cosmos address is owned by a given ethereum key
func updateEthAddressHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateEthAddressReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		cosmosAddr := cliCtx.GetFromAddress()
		// the signed message should be the hash of the presented CosmosAddr
		ethHash := ethCrypto.Keccak256Hash(cosmosAddr)

		ethSig, err := hexUtil.Decode(req.EthSig)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// we recover the address and public key from the sig
		ethPubkey, err := ethCrypto.SigToPub(ethHash.Bytes(), ethSig)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ethPubkeyBytes := ethCrypto.FromECDSAPub(ethPubkey)
		ethAddr := ethCrypto.PubkeyToAddress(*ethPubkey)
		correct := ethCrypto.VerifySignature(ethPubkeyBytes, ethHash.Bytes(), ethSig)
		if correct == false {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Make the message, we convert the recovered address into a string
		// so at this point we have verified that this address signed this
		// cosmos address
		msg := types.NewMsgSetEthAddress(types.EthereumAddress(ethAddr), cosmosAddr, hex.EncodeToString(ethSig))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type createValsetReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func createValsetRequestHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createValsetReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		// this can be the sender, since we don't really care who's name is on this
		cosmosAddr := cliCtx.GetFromAddress()
		// Make the message
		msg := types.NewMsgValsetRequest(cosmosAddr)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type valsetConfirmReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	EthAddress string       `json:"eth_address"`
	Nonce      string       `json:"nonce"`
	EthSig     string       `json:"ethSig"`
}

// check the ethereum sig on a particular valset and broadcast a transaction containing
// it if correct. The nonce / block height is used to determine what valset to look up
// locally and verify
func createValsetConfirmHandler(cliCtx context.CLIContext, storeKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req valsetConfirmReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		ethAddress := types.NewEthereumAddress(req.EthAddress)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, req.Nonce), nil)
		if err != nil {
			fmt.Printf("could not get valset")
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		var valset types.Valset
		cliCtx.Codec.MustUnmarshalJSON(res, &valset)
		checkpoint := valset.GetCheckpoint()

		// the signed message should be the hash of the checkpoint at the given nonce
		ethHash := ethCrypto.Keccak256Hash(checkpoint)

		ethSig, err := hexUtil.Decode(req.EthSig)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ethPubkey, err := ethCrypto.SigToPub(ethHash.Bytes(), ethSig)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ethPubkeyBytes := ethCrypto.FromECDSAPub(ethPubkey)

		correct := ethCrypto.VerifySignature(ethPubkeyBytes, ethHash.Bytes(), ethSig)
		if correct == false {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cosmosAddr := cliCtx.GetFromAddress()
		msg := types.NewMsgValsetConfirm(valset.Nonce, ethAddress, cosmosAddr, req.EthSig)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type bootstrapConfirmReq struct {
	BaseReq               rest.BaseReq           `json:"base_req"`
	Orchestrator          sdk.AccAddress         `json:"orchestrator"`
	EthereumChainID       uint64                 `json:"ethereum_chain_id"`
	BridgeContractAddress types.EthereumAddress  `json:"bridge_contract_address"`
	Block                 string                 `json:"block"`
	BridgeValidators      types.BridgeValidators `json:"bridge_validators"`
	PeggyID               string                 `json:"peggy_id"`
	StartThreshold        uint64                 `json:"start_threshold"`
}

func bootstrapConfirmHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req bootstrapConfirmReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		blockNumber, err := strconv.ParseUint(req.Block, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse block")
			return
		}
		claims := []types.EthereumClaim{
			types.EthereumBridgeBootstrappedClaim{
				Block:            blockNumber,
				BridgeValidators: req.BridgeValidators,
				PeggyID:          req.PeggyID,
				StartThreshold:   req.StartThreshold,
			},
		}
		msg := types.NewMsgCreateEthereumClaims(req.EthereumChainID, req.BridgeContractAddress, req.Orchestrator, claims)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid data: %s", err))
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type BridgeApprovalSignatureReq struct {
	BaseReq           rest.BaseReq   `json:"base_req"`
	Orchestrator      sdk.AccAddress `json:"orchestrator"`
	EthereumSignature []byte         `json:"ethereum_signature"`
}

func BridgeApprovalSignatureHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BridgeApprovalSignatureReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		vars := mux.Vars(r)
		signType, exists := types.SignTypeFromName(vars[claimType])
		if !exists {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "unknown sign type")
			return
		}
		nonce, err := types.UInt64NonceFromString(vars[nonce])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid nonce: %s", err))
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := types.MsgBridgeSignatureSubmission{
			Nonce:             nonce,
			SignType:          signType,
			Orchestrator:      req.Orchestrator,
			EthereumSignature: req.EthereumSignature,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid data: %s", err))
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
