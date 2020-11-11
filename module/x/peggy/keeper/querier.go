package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// USED BY RUST:
// "custom/%s/valsetRequest/%s"
// "custom/%s/valsetConfirms/%s"
// "custom/%s/lastValsetRequests"
// "custom/%s/lastPendingValsetRequest/%s"

const (
	/* USED BY RUST */ QueryValsetRequest = "valsetRequest"
	/* USED BY RUST */ QueryValsetConfirmsByNonce = "valsetConfirms"
	/* USED BY RUST */ QueryLastValsetRequests = "lastValsetRequests"
	/* USED BY RUST */ QueryLastPendingValsetRequestByAddr = "lastPendingValsetRequest"
	QueryCurrentValset                                     = "currentValset"
	QueryValsetConfirm                                     = "valsetConfirm"
	QueryLastPendingBatchRequestByAddr                     = "lastPendingBatchRequest"
	QueryOutgoingTxBatches                                 = "allBatches"
	QueryBridgedDenominators                               = "allBridgedDenominators"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryCurrentValset:
			return queryCurrentValset(ctx, keeper)
		case QueryValsetRequest:
			return queryValsetRequest(ctx, path[1:], keeper)
		case QueryValsetConfirm:
			return queryValsetConfirm(ctx, path[1:], keeper)
		case QueryValsetConfirmsByNonce:
			return allValsetConfirmsByNonce(ctx, path[1], keeper)
		case QueryLastValsetRequests:
			return lastValsetRequests(ctx, keeper)
		case QueryLastPendingValsetRequestByAddr:
			return lastPendingValsetRequest(ctx, path[1], keeper)
		case QueryLastPendingBatchRequestByAddr:
			return lastPendingBatchRequest(ctx, path[1], keeper)
		case QueryOutgoingTxBatches:
			return allBatchesRequest(ctx, keeper)
		case QueryBridgedDenominators:
			return queryBridgedDenominators(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

/* USED BY RUST */
func queryValsetRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := parseNonce(path[0])
	if err != nil {
		return nil, err
	}

	valset := keeper.GetValsetRequest(ctx, nonce)
	if valset == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, *valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

/* USED BY RUST */
// allValsetConfirmsByNonce returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func allValsetConfirmsByNonce(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := parseNonce(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var confirms []types.MsgValsetConfirm
	keeper.IterateValsetConfirmByNonce(ctx, nonce, func(_ []byte, c types.MsgValsetConfirm) bool {
		confirms = append(confirms, c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

const maxValsetRequestsReturned = 5

/* USED BY RUST */
func lastValsetRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var counter int
	var valReq []types.Valset
	keeper.IterateValsetRequest(ctx, func(_ []byte, val types.Valset) bool {
		valReq = append(valReq, val)
		counter++
		return counter >= maxValsetRequestsReturned
	})
	if len(valReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, valReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

/* USED BY RUST */
func lastPendingValsetRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingValsetReq *types.Valset
	keeper.IterateValsetRequest(ctx, func(_ []byte, val types.Valset) bool {
		found := keeper.GetValsetConfirm(ctx, val.Nonce, addr)
		if found == nil {
			pendingValsetReq = &val
		}
		return true
	})
	if pendingValsetReq == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, pendingValsetReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCurrentValset(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentValset(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryValsetConfirm returns the confirm msg for single orchestrator address and nonce
// When nothing found a nil value is returned
func queryValsetConfirm(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	nonce, err := parseNonce(path[0])
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
	res, err := codec.MarshalJSONIndent(keeper.cdc, *valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

type MultiSigUpdateResponse struct {
	Valset     types.Valset `json:"valset"`
	Signatures [][]byte     `json:"signatures,omitempty"`
}

func lastPendingBatchRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	// todo: find validator address by operator key
	validatorAddr := sdk.ValAddress(addr)

	var pendingBatchReq *types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch types.OutgoingTxBatch) bool {
		found := keeper.HasBatchApprovalSignature(ctx, batch.TokenContract, batch.Nonce, validatorAddr)
		if !found {
			pendingBatchReq = &batch
		}
		return true
	})
	if pendingBatchReq == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, pendingBatchReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

type ApprovedOutgoingTxBatchResponse struct {
	Batch      types.OutgoingTxBatch `json:"batch"`
	Signatures [][]byte              `json:"signatures,omitempty"`
}

const MaxResults = 100 // todo: impl pagination

func allBatchesRequest(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var batches []types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	if len(batches) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, batches)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func parseNonce(nonceArg string) (types.UInt64Nonce, error) {
	return types.UInt64NonceFromString(nonceArg)
}

func queryBridgedDenominators(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var r []types.BridgedDenominator
	keeper.IterateCounterpartDenominators(ctx, func(rawKey []byte, denominator types.BridgedDenominator) bool {
		r = append(r, denominator)
		return false
	})
	if len(r) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, r)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
