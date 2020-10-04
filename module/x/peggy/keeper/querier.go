package keeper

import (
	"github.com/althea-net/peggy/module/x/peggy/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryCurrentValset                  = "currentValset"
	QueryValsetRequest                  = "valsetRequest"
	QueryValsetConfirm                  = "valsetConfirm"
	QueryValsetConfirmsByNonce          = "valsetConfirms"
	QueryLastValsetRequests             = "lastValsetRequests"
	QueryLastPendingValsetRequestByAddr = "lastPendingValsetRequest"
	QueryLastProcessedNonce             = "lastNonce"
	QueryLastObservedNonces             = "lastObservedNonces"
	QueryLastApprovedNonces             = "lastApprovedNonces"
	QueryLastPendingBatchRequestByAddr  = "lastPendingBatchRequest"
	QueryBatchRequestByNonce            = "batch"
	QueryOutgoingTxBatches              = "allBatches"
	// batches that are approved but not observed
	QueryInflightOutgoingTxBatches = "inflightBatches"
	// last valset that was updated on the ETH side successfully
	QueryLastObservedValset = "lastObservedMultiSigUpdate"
	// last valset with enough signatures
	QueryLastApprovedValset      = "lastApprovedMultiSigUpdate"
	QueryAttestationsByClaimType = "allAttestations"
	QueryAttestation             = "attestation"
	QueryBridgedDenominators     = "allBridgedDenominators"
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
		case QueryLastProcessedNonce:
			return lastProcessedNonce(ctx, path[1], keeper)
		case QueryLastObservedNonces:
			return lastObservedNonces(ctx, keeper)
		case QueryLastApprovedNonces:
			return lastApprovedNonces(ctx, keeper)
		case QueryLastPendingBatchRequestByAddr:
			return lastPendingBatchRequest(ctx, path[1], keeper)
		case QueryOutgoingTxBatches:
			return allBatchesRequest(ctx, keeper)
		case QueryLastApprovedValset:
			return lastApprovedMultiSigUpdate(ctx, keeper)
		case QueryLastObservedValset:
			return lastObservedMultiSigUpdate(ctx, keeper)
		case QueryAttestationsByClaimType:
			return allAttestations(ctx, path[1], keeper)
		case QueryBatchRequestByNonce:
			return queryBatchByNonce(ctx, path[1], keeper)
		case QueryAttestation:
			return queryAttestation(ctx, path[1], path[2], keeper)
		case QueryBridgedDenominators:
			return queryBridgedDenominators(ctx, keeper)
		case QueryInflightOutgoingTxBatches:
			return queryInflightBatches(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", types.ModuleName)
		}
	}
}

func queryCurrentValset(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	valset := keeper.GetCurrentValset(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, valset)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

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

func lastPendingValsetRequest(ctx sdk.Context, operatorAddr string, keeper Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	// todo: find validator address by operator key
	validatorAddr := sdk.ValAddress(addr)

	var pendingValsetReq *types.Valset
	keeper.IterateValsetRequest(ctx, func(_ []byte, val types.Valset) bool {
		found := keeper.HasBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, val.Nonce, validatorAddr)
		if !found {
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

// lastProcessedNonce returns as single nonce value or nil
func lastProcessedNonce(ctx sdk.Context, claimTypeStr string, keeper Keeper) ([]byte, error) {
	claimType, exists := types.ClaimTypeFromName(claimTypeStr)
	if !exists {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}
	att := keeper.GetLastProcessedAttestation(ctx, claimType)
	if att == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, att.Nonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// lastObservedNonces returns a list of nonces. One for each claim type if exists
func lastObservedNonces(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	result := make(map[string]types.UInt64Nonce, len(types.AllOracleClaimTypes))
	for _, v := range types.AllOracleClaimTypes {
		nonce := keeper.GetLastAttestedNonce(ctx, v)
		if nonce != nil {
			result[v.String()] = *nonce
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return sdk.SortJSON(res)
}

// lastApprovedNonces returns a list of nonces. One for each claim type if exists
func lastApprovedNonces(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	result := make(map[string]types.UInt64Nonce, len(types.AllOracleClaimTypes))
	for _, v := range types.AllSignerApprovalClaimTypes {
		nonce := keeper.GetLastAttestedNonce(ctx, v)
		if nonce != nil {
			result[v.String()] = *nonce
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return sdk.SortJSON(res)
}

type MultiSigUpdateResponse struct {
	Valset     types.Valset `json:"valset"`
	Signatures [][]byte     `json:"signatures,omitempty"`
}

func lastApprovedMultiSigUpdate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	nonce := keeper.GetLastAttestedNonce(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate)
	return fetchMultiSigUpdateData(ctx, nonce, keeper)
}

func lastObservedMultiSigUpdate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	nonce := keeper.GetLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeMultiSigUpdate)
	if nonce == nil || nonce.IsEmpty() {
		nonce = keeper.GetLastAttestedNonce(ctx, types.ClaimTypeEthereumBridgeBootstrap)
	}
	return fetchMultiSigUpdateData(ctx, nonce, keeper)
}

func fetchMultiSigUpdateData(ctx sdk.Context, nonce *types.UInt64Nonce, keeper Keeper) ([]byte, error) {
	if nonce == nil || nonce.IsEmpty() {
		return nil, nil
	}

	valset := keeper.GetValsetRequest(ctx, *nonce)
	if valset == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "no valset found for nonce")
	}

	result := MultiSigUpdateResponse{
		Valset: *valset,
	}

	keeper.IterateBridgeApprovalSignatures(ctx, types.ClaimTypeOrchestratorSignedMultiSigUpdate, *nonce, func(_ []byte, sig []byte) bool {
		// todo: check if signature order matters
		result.Signatures = append(result.Signatures, sig)
		return false
	})

	res, err := codec.MarshalJSONIndent(keeper.cdc, result)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return sdk.SortJSON(res)
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
		found := keeper.HasBridgeApprovalSignature(ctx, types.ClaimTypeOrchestratorSignedWithdrawBatch, batch.Nonce, validatorAddr)
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

func queryInflightBatches(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	// batches that are approved but not observed, yet
	// batch nonce > than last observed nonce
	var lastNonce types.UInt64Nonce
	if att := keeper.GetLastProcessedAttestation(ctx, types.ClaimTypeEthereumBridgeWithdrawalBatch); att != nil {
		lastNonce = att.Nonce
	}

	var batches []types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(_ []byte, batch types.OutgoingTxBatch) bool {
		if !batch.Nonce.GreaterThan(lastNonce) {
			return true
		}
		batches = append(batches, batch)
		return false
	})
	if len(batches) == 0 {
		return nil, nil
	}
	out := make([]ApprovedOutgoingTxBatchResponse, len(batches))
	for i := range batches {
		r := ApprovedOutgoingTxBatchResponse{
			Batch: batches[i],
		}
		keeper.IterateBridgeApprovalSignatures(ctx, types.ClaimTypeOrchestratorSignedWithdrawBatch, batches[i].Nonce, func(_ []byte, sig []byte) bool {
			r.Signatures = append(r.Signatures, sig)
			return false
		})
		out[i] = r
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, out)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
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

func allAttestations(ctx sdk.Context, claimTypeStr string, keeper Keeper) ([]byte, error) {
	claimType, exists := types.ClaimTypeFromName(claimTypeStr)
	if !exists {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}
	var attestations []types.Attestation
	keeper.IterateAttestationByClaimTypeDesc(ctx, claimType, func(_ []byte, att types.Attestation) bool {
		attestations = append(attestations, att)
		return len(attestations) == MaxResults
	})
	if len(attestations) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, attestations)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryAttestation(ctx sdk.Context, claimTypeStr, nonceStr string, keeper Keeper) ([]byte, error) {
	claimType, exists := types.ClaimTypeFromName(claimTypeStr)
	if !exists {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}

	nonce, err := parseNonce(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "nonce")
	}
	attestation := keeper.GetAttestation(ctx, claimType, nonce)
	if attestation == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, *attestation)
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

func queryBatchByNonce(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := parseNonce(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "nonce")
	}
	batch := keeper.GetOutgoingTXBatch(ctx, nonce)
	if batch == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, *batch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
