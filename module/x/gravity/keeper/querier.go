package keeper

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

const (

	// SignerSetTxs

	// This retrieves a specific validator set by it's nonce
	// used to compare what's on Ethereum with what's in Cosmos
	// to perform slashing / validation of system consistency
	QuerySignerSetTxRequest = "valsetRequest"
	// Gets all the confirmation signatures for a given validator
	// set, used by the relayer to package the validator set and
	// it's signatures into an Ethereum transaction
	QuerySignerSetTxSignaturesByNonce = "valsetConfirms"
	// Gets the last N (where N is currently 5) validator sets that
	// have been produced by the chain. Useful to see if any recently
	// signed requests can be submitted.
	QueryLastSignerSetTxRequests = "lastSignerSetTxRequests"
	// Gets a list of unsigned valsets for a given validators delegate
	// orchestrator address. Up to 100 are sent at a time
	QueryLastPendingSignerSetTxRequestByAddr = "lastPendingSignerSetTxRequest"

	QueryCurrentSignerSetTx = "currentSignerSetTx"
	// TODO remove this, it's not used, getting one confirm at a time
	// is mostly useless
	QuerySignerSetTxSignature = "valsetConfirm"

	// used by the contract deployer script. GravityID is set in the Genesis
	// file, then read by the contract deployer and deployed to Ethereum
	// a unique GravityID ensures that even if the same validator set with
	// the same keys is running on two chains these chains can have independent
	// bridges
	QueryGravityID = "gravityID"

	// Batches
	// note the current logic here constrains batch throughput to one
	// batch (of any type) per Cosmos block.

	// This retrieves a specific batch by it's nonce and token contract
	// or in the case of a Cosmos originated address it's denom
	QueryBatch = "batch"
	// Get the last unsigned batch (of any denom) for the validators
	// orchestrator to sign
	QueryLastPendingBatchRequestByAddr = "lastPendingBatchRequest"
	// gets the last 100 outgoing batches, regardless of denom, useful
	// for a relayer to see what is available to relay
	QueryBatchTxs = "lastBatches"
	// Used by the relayer to package a batch with signatures required
	// to submit to Ethereum
	QueryBatchConfirms = "batchConfirms"
	// Used to query all pending SendToEthereum transactions and fees available for each
	// token type, a relayer can then estimate their potential profit when requesting
	// a batch
	QueryBatchFees = "batchFees"

	// Logic calls
	// note the current logic here constrains logic call throughput to one
	// call (of any type) per Cosmos block.

	// This retrieves a specific logic call by it's nonce and token contract
	// or in the case of a Cosmos originated address it's denom
	QueryLogicCall = "logicCall"
	// Get the last unsigned logic call for the validators orchestrator
	// to sign
	QueryLastPendingLogicCallByAddr = "lastPendingLogicCall"
	// gets the last 5 outgoing logic calls, regardless of denom, useful
	// for a relayer to see what is available to relay
	QueryContractCallTxs = "lastLogicCalls"
	// Used by the relayer to package a logic call with signatures required
	// to submit to Ethereum
	QueryLogicCallConfirms = "logicCallConfirms"

	// Token mapping
	// This retrieves the denom which is represented by a given ERC20 contract
	QueryERC20ToDenom = "ERC20ToDenom"
	// This retrieves the ERC20 contract which represents a given denom
	QueryDenomToERC20 = "DenomToERC20"

	// Query pending transactions
	PendingSendToEthereumRequest = "PendingSendToEthereum"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {

		// SignerSetTxs
		case QueryCurrentSignerSetTx:
			return queryCurrentSignerSetTx(ctx, keeper)
		case QuerySignerSetTxRequest:
			return querySignerSetTxRequest(ctx, path[1:], keeper)
		case QuerySignerSetTxSignature:
			return querySignerSetTxSignature(ctx, path[1:], keeper)
		case QuerySignerSetTxSignaturesByNonce:
			return queryAllSignerSetTxSignatures(ctx, path[1], keeper)
		case QueryLastSignerSetTxRequests:
			return lastSignerSetTxRequests(ctx, keeper)
		case QueryLastPendingSignerSetTxRequestByAddr:
			return lastPendingSignerSetTxRequest(ctx, path[1], keeper)

		// Batches
		case QueryBatch:
			return queryBatch(ctx, path[1], path[2], keeper)
		case QueryBatchConfirms:
			return queryAllBatchConfirms(ctx, path[1], path[2], keeper)
		case QueryLastPendingBatchRequestByAddr:
			return lastPendingBatchRequest(ctx, path[1], keeper)
		case QueryBatchTxs:
			return lastBatchesRequest(ctx, keeper)
		case QueryBatchFees:
			return queryBatchFees(ctx, keeper)

		// Logic calls
		case QueryLogicCall:
			return queryLogicCall(ctx, path[1], path[2], keeper)
		case QueryLogicCallConfirms:
			return queryAllLogicCallConfirms(ctx, path[1], path[2], keeper)
		case QueryLastPendingLogicCallByAddr:
			return lastPendingLogicCallRequest(ctx, path[1], keeper)
		case QueryContractCallTxs:
			return lastLogicCallRequests(ctx, keeper)

		case QueryGravityID:
			return queryGravityID(ctx, keeper)

		// Token mappings
		case QueryDenomToERC20:
			return queryDenomToERC20(ctx, path[1], keeper)
		case QueryERC20ToDenom:
			return queryERC20ToDenom(ctx, path[1], keeper)

		// Pending transactions
		case PendingSendToEthereumRequest:
			return queryPendingSendToEthereum(ctx, path[1], keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func querySignerSetTxRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(path[0])
	if err != nil {
		return nil, err
	}

	valset := keeper.GetSignerSetTx(ctx, nonce)
	if valset == nil {
		return nil, nil
	}
	// TODO: replace these with the GRPC response types
	// TODO: fix the use of module codec here
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allSignerSetTxSignaturesByNonce returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllSignerSetTxSignatures(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []*types.MsgSignerSetTxSignature
	keeper.IterateSignerSetTxSignatureByNonce(ctx, nonce, func(_ []byte, c types.MsgSignerSetTxSignature) bool {
		confirms = append(confirms, &c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allBatchConfirms returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllBatchConfirms(ctx sdk.Context, nonceStr string, tokenContract string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []types.MsgBatchTxSignature
	keeper.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, c types.MsgBatchTxSignature) bool {
		confirms = append(confirms, c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

const maxSignerSetTxRequestsReturned = 5

// lastSignerSetTxRequests returns up to maxSignerSetTxRequestsReturned valsets from the store
func lastSignerSetTxRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var counter int
	var valReq []*types.SignerSetTx
	keeper.IterateSignerSetTxs(ctx, func(_ []byte, val *types.SignerSetTx) bool {
		valReq = append(valReq, val)
		counter++
		return counter >= maxSignerSetTxRequestsReturned
	})
	if len(valReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, valReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// lastPendingSignerSetTxRequest gets a list of validator sets that this validator has not signed
// limited by 100 sets per request.
func lastPendingSignerSetTxRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingSignerSetTxReq []*types.SignerSetTx
	keeper.IterateSignerSetTxs(ctx, func(_ []byte, val *types.SignerSetTx) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := keeper.GetSignerSetTxSignature(ctx, val.Nonce, addr) != nil
		// if this valset has NOT been signed by operatorAddr, store it in pendingSignerSetTxReq
		// and exit the loop
		if !foundConfirm {
			pendingSignerSetTxReq = append(pendingSignerSetTxReq, val)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, TODO pagination
		if len(pendingSignerSetTxReq) > 100 {
			return true
		}
		// return false to continue the loop
		return false
	})
	if len(pendingSignerSetTxReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingSignerSetTxReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCurrentSignerSetTx(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentSignerSetTx(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// querySignerSetTxSignature returns the confirm msg for single orchestrator address and nonce
// When nothing found a nil value is returned
func querySignerSetTxSignature(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	accAddress, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	valset := keeper.GetSignerSetTxSignature(ctx, nonce, accAddress)
	if valset == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, *valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

type MultiSigUpdateResponse struct {
	SignerSetTx types.SignerSetTx `json:"valset"`
	Signatures  [][]byte          `json:"signatures,omitempty"`
}

// lastPendingBatchRequest gets the latest batch that has NOT been signed by operatorAddr
func lastPendingBatchRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingBatchReq *types.BatchTx
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.BatchTx) bool {
		foundConfirm := keeper.GetBatchConfirm(ctx, batch.BatchNonce, batch.TokenContract, addr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	if pendingBatchReq == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingBatchReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

const MaxResults = 100 // todo: impl pagination

// Gets MaxResults batches from store. Does not select by token type or anything
func lastBatchesRequest(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var batches []*types.BatchTx
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.BatchTx) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	if len(batches) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, batches)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryBatchFees(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	val := types.BatchFeeResponse{BatchFees: keeper.GetAllBatchFees(ctx)}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, val)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// Gets MaxResults logic calls from store.
func lastLogicCallRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var calls []*types.ContractCallTx
	keeper.IterateContractCallTxs(ctx, func(_ []byte, call *types.ContractCallTx) bool {
		calls = append(calls, call)
		return len(calls) == MaxResults
	})
	if len(calls) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, calls)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryBatch gets a batch by tokenContract and nonce
func queryBatch(ctx sdk.Context, nonce string, tokenContract string, keeper Keeper) ([]byte, error) {
	parsedNonce, err := types.UInt64FromString(nonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	if types.ValidateEthAddress(tokenContract) != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	foundBatch := keeper.GetOutgoingTXBatch(ctx, tokenContract, parsedNonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, foundBatch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	return res, nil
}

// lastPendingLogicCallRequest gets the latest call that has NOT been signed by operatorAddr
func lastPendingLogicCallRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingLogicCalls *types.ContractCallTx
	keeper.IterateContractCallTxs(ctx, func(_ []byte, call *types.ContractCallTx) bool {
		foundConfirm := keeper.GetLogicCallConfirm(ctx, call.InvalidationId, call.InvalidationNonce, addr) != nil
		if !foundConfirm {
			pendingLogicCalls = call
			return true
		}
		return false
	})
	if pendingLogicCalls == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingLogicCalls)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryLogicCall gets a logic call by nonce and invalidation id
func queryLogicCall(ctx sdk.Context, invalidationId string, invalidationNonce string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(invalidationNonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	foundCall := keeper.GetContractCallTx(ctx, []byte(invalidationId), nonce)
	if foundCall == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find logic call")
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, foundCall)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	return res, nil
}

// allLogicCallConfirms returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllLogicCallConfirms(ctx sdk.Context, invalidationId string, invalidationNonce string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(invalidationNonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	invalidationIdBytes, err := hex.DecodeString(invalidationId)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []*types.MsgContractCallTxSignature
	keeper.IterateLogicConfirmByInvalidationIDAndNonce(ctx, invalidationIdBytes, nonce, func(_ []byte, c *types.MsgContractCallTxSignature) bool {
		confirms = append(confirms, c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGravityID(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	gravityID := keeper.GetGravityID(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, gravityID)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return res, nil
	}
}

func queryDenomToERC20(ctx sdk.Context, denom string, keeper Keeper) ([]byte, error) {
	cosmos_originated, erc20, err := keeper.DenomToERC20Lookup(ctx, denom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	var response types.DenomToERC20Response
	response.CosmosOriginated = cosmos_originated
	response.Erc20 = erc20
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryERC20ToDenom(ctx sdk.Context, ERC20 string, keeper Keeper) ([]byte, error) {
	cosmos_originated, denom := keeper.ERC20ToDenomLookup(ctx, ERC20)
	var response types.ERC20ToDenomResponse
	response.CosmosOriginated = cosmos_originated
	response.Denom = denom
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, response)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}

func queryPendingSendToEthereum(ctx sdk.Context, senderAddr string, k Keeper) ([]byte, error) {
	batches := k.GetBatchTxs(ctx)
	unbatched_tx := k.GetPoolTransactions(ctx)
	sender_address := senderAddr
	res := types.PendingSendToEthereumResponse{}
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			if tx.Sender == sender_address {
				res.TransfersInBatches = append(res.TransfersInBatches, tx)
			}
		}
	}
	for _, tx := range unbatched_tx {
		if tx.Sender == sender_address {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
	}
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return bytes, nil
	}
}
