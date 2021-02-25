package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (

	// Valsets

	// This retrieves a specific validator set by it's nonce
	// used to compare what's on Ethereum with what's in Cosmos
	// to perform slashing / validation of system consistency
	QueryValsetRequest = "valsetRequest"
	// Gets all the confirmation signatures for a given validator
	// set, used by the relayer to package the validator set and
	// it's signatures into an Ethereum transaction
	QueryValsetConfirmsByNonce = "valsetConfirms"
	// Gets the last N (where N is currently 5) validator sets that
	// have been produced by the chain. Useful to see if any recently
	// signed requests can be submitted.
	QueryLastValsetRequests = "lastValsetRequests"
	// Gets a list of unsigned valsets for a given validators delegate
	// orchestrator address. Up to 100 are sent at a time
	QueryLastPendingValsetRequestByAddr = "lastPendingValsetRequest"

	QueryCurrentValset = "currentValset"
	// TODO remove this, it's not used, getting one confirm at a time
	// is mostly useless
	QueryValsetConfirm = "valsetConfirm"

	// used by the contract deployer script. PeggyID is set in the Genesis
	// file, then read by the contract deployer and deployed to Ethereum
	// a unique PeggyID ensures that even if the same validator set with
	// the same keys is running on two chains these chains can have independent
	// bridges
	QueryPeggyID = "peggyID"

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
	QueryOutgoingTxBatches = "lastBatches"
	// Used by the relayer to package a batch with signatures required
	// to submit to Ethereum
	QueryBatchConfirms = "batchConfirms"

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
	QueryOutgoingLogicCalls = "lastLogicCalls"
	// Used by the relayer to package a logic call with signatures required
	// to submit to Ethereum
	QueryLogicCallConfirms = "logicCallConfirms"

	// Cosmos originated assets
	// This retrieves the denom which is represented by a given ERC20 contract
	QueryERC20ToDenom = "ERC20ToDenom"
	// This retrieves the ERC20 contract which represents a given denom
	QueryDenomToERC20 = "DenomToERC20"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {

		// Valsets
		case QueryCurrentValset:
			return queryCurrentValset(ctx, keeper)
		case QueryValsetRequest:
			return queryValsetRequest(ctx, path[1:], keeper)
		case QueryValsetConfirm:
			return queryValsetConfirm(ctx, path[1:], keeper)
		case QueryValsetConfirmsByNonce:
			return queryAllValsetConfirms(ctx, path[1], keeper)
		case QueryLastValsetRequests:
			return lastValsetRequests(ctx, keeper)
		case QueryLastPendingValsetRequestByAddr:
			return lastPendingValsetRequest(ctx, path[1], keeper)

		// Batches
		case QueryBatch:
			return queryBatch(ctx, path[1], path[2], keeper)
		case QueryBatchConfirms:
			return queryAllBatchConfirms(ctx, path[1], path[2], keeper)
		case QueryLastPendingBatchRequestByAddr:
			return lastPendingBatchRequest(ctx, path[1], keeper)
		case QueryOutgoingTxBatches:
			return lastBatchesRequest(ctx, keeper)

		// Logic calls
		case QueryLogicCall:
			return queryLogicCall(ctx, path[1], path[2], keeper)
		case QueryLogicCallConfirms:
			return queryAllLogicCallConfirms(ctx, path[1], path[2], keeper)
		case QueryLastPendingLogicCallByAddr:
			return lastPendingLogicCallRequest(ctx, path[1], keeper)
		case QueryOutgoingLogicCalls:
			return lastLogicCallRequests(ctx, keeper)

		case QueryPeggyID:
			return queryPeggyID(ctx, keeper)

		// Cosmos originated assets
		case QueryDenomToERC20:
			return queryDenomToERC20(ctx, path[1], keeper)
		case QueryERC20ToDenom:
			return queryERC20ToDenom(ctx, path[1], keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryValsetRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(path[0])
	if err != nil {
		return nil, err
	}

	valset := keeper.GetValset(ctx, nonce)
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

// allValsetConfirmsByNonce returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllValsetConfirms(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []*types.MsgValsetConfirm
	keeper.IterateValsetConfirmByNonce(ctx, nonce, func(_ []byte, c types.MsgValsetConfirm) bool {
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

	var confirms []types.MsgConfirmBatch
	keeper.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
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

const maxValsetRequestsReturned = 5

// lastValsetRequests returns up to maxValsetRequestsReturned valsets from the store
func lastValsetRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var counter int
	var valReq []*types.Valset
	keeper.IterateValsets(ctx, func(_ []byte, val *types.Valset) bool {
		valReq = append(valReq, val)
		counter++
		return counter >= maxValsetRequestsReturned
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

// lastPendingValsetRequest gets a list of validator sets that this validator has not signed
// limited by 100 sets per request.
func lastPendingValsetRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingValsetReq []*types.Valset
	keeper.IterateValsets(ctx, func(_ []byte, val *types.Valset) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := keeper.GetValsetConfirm(ctx, val.Nonce, addr) != nil
		// if this valset has NOT been signed by operatorAddr, store it in pendingValsetReq
		// and exit the loop
		if !foundConfirm {
			pendingValsetReq = append(pendingValsetReq, val)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, TODO pagination
		if len(pendingValsetReq) > 100 {
			return true
		}
		// return false to continue the loop
		return false
	})
	if len(pendingValsetReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingValsetReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCurrentValset(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentValset(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryValsetConfirm returns the confirm msg for single orchestrator address and nonce
// When nothing found a nil value is returned
func queryValsetConfirm(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := types.UInt64FromString(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	accAddress, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	valset := keeper.GetValsetConfirm(ctx, nonce, accAddress)
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
	Valset     types.Valset `json:"valset"`
	Signatures [][]byte     `json:"signatures,omitempty"`
}

// lastPendingBatchRequest gets the latest batch that has NOT been signed by operatorAddr
func lastPendingBatchRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingBatchReq *types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
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
	var batches []*types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
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

// Gets MaxResults logic calls from store.
func lastLogicCallRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var calls []*types.OutgoingLogicCall
	keeper.IterateOutgoingLogicCalls(ctx, func(_ []byte, call *types.OutgoingLogicCall) bool {
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

	var pendingLogicCalls *types.OutgoingLogicCall
	keeper.IterateOutgoingLogicCalls(ctx, func(_ []byte, call *types.OutgoingLogicCall) bool {
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

	foundCall := keeper.GetOutgoingLogicCall(ctx, []byte(invalidationId), nonce)
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

	var confirms []*types.MsgConfirmLogicCall
	keeper.IterateLogicConfirmByInvalidationIdAndNonce(ctx, []byte(invalidationId), nonce, func(_ []byte, c *types.MsgConfirmLogicCall) bool {
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

func queryPeggyID(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	peggyID := keeper.GetPeggyID(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, peggyID)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return res, nil
	}
}

func queryDenomToERC20(ctx sdk.Context, denom string, keeper Keeper) ([]byte, error) {
	_, erc20, err := keeper.DenomToERC20(ctx, denom)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	} else {
		return []byte(erc20), nil
	}
}

func queryERC20ToDenom(ctx sdk.Context, ERC20 string, keeper Keeper) ([]byte, error) {
	_, denom := keeper.ERC20ToDenom(ctx, ERC20)
	return []byte(denom), nil
}
