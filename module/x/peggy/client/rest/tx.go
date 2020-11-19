package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	hexUtil "github.com/ethereum/go-ethereum/common/hexutil"
)

type updateEthAddressReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	EthSig  string       `json:"ethSig"`
}

// accepts a sig proving that the given Cosmos address is owned by a given ethereum key
func updateEthAddressHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateEthAddressReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
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
		msg := types.NewMsgSetEthAddress(ethAddr.String(), cosmosAddr, hex.EncodeToString(ethSig))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

type createValsetReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
}

func createValsetRequestHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createValsetReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
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

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
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
func createValsetConfirmHandler(cliCtx client.Context, storeKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req valsetConfirmReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, req.Nonce), nil)
		if err != nil {
			fmt.Printf("could not get valset")
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		var valset types.Valset
		cliCtx.JSONMarshaler.MustUnmarshalJSON(res, &valset)
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
		msg := types.NewMsgValsetConfirm(valset.Nonce, req.EthAddress, cosmosAddr, req.EthSig)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
	}
}

type bootstrapConfirmReq struct {
	BaseReq               rest.BaseReq           `json:"base_req"`
	Orchestrator          sdk.AccAddress         `json:"orchestrator"`
	EthereumChainID       uint64                 `json:"ethereum_chain_id"`
	BridgeContractAddress string                 `json:"bridge_contract_address"`
	Block                 string                 `json:"block"`
	BridgeValidators      types.BridgeValidators `json:"bridge_validators"`
	PeggyID               string                 `json:"peggy_id"`
	StartThreshold        uint64                 `json:"start_threshold"`
}
