package keeper

import (
	"encoding/base64"
	"strconv"

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
	QueryLastObservedNonce              = "lastObservedNonce"
	QueryLastObservedNonces             = "lastObservedNonces"
	QueryLastPendingBatchRequestByAddr  = "lastPendingBatchRequest"
	QueryOutgoingTxBatches              = "allBatches"
	// last valset that was updated on the ETH side successfully
	QueryLastObservedValset = "lastObservedMultiSigUpdate"
	// last valset with enough signatures
	QueryLastApprovedValset      = "lastApprovedMultiSigUpdate"
	QueryAttestationsByClaimType = "allAttestations"
	QueryAttestation             = "attestation"
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
		case QueryLastObservedNonce:
			return lastObservedNonce(ctx, path[1], keeper)
		case QueryLastObservedNonces:
			return lastObservedNonces(ctx, keeper)
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
		case QueryAttestation:
			return queryAttestation(ctx, path[1], path[2], keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
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

	valset := keeper.GetValsetRequest(ctx, int64(nonce.Uint64()))
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

	valset := keeper.GetValsetConfirm(ctx, int64(nonce.Uint64()), accAddress)
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
	keeper.IterateValsetConfirmByNonce(ctx, int64(nonce.Uint64()), func(_ []byte, c types.MsgValsetConfirm) bool {
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
	keeper.IterateValsetRequest(ctx, func(key []byte, val types.Valset) bool {
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
	validatorAddr := addr

	var pendingValsetReq *types.Valset
	keeper.IterateValsetRequest(ctx, func(key []byte, val types.Valset) bool {
		found := keeper.HasValsetConfirm(ctx, val.Nonce, validatorAddr)
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

// lastObservedNonce returns as single nonce value or nil
func lastObservedNonce(ctx sdk.Context, claimType string, keeper Keeper) ([]byte, error) {
	if !types.IsClaimType(claimType) {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}
	att := keeper.GetLastObservedAttestation(ctx, types.ClaimType(claimType))
	if att == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, att.Nonce)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// lastObservedNonce returns a list of nonces. One for each claim type if exists
func lastObservedNonces(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	result := make(map[string]types.Nonce, len(types.AllOracleClaimTypes))
	for _, v := range types.AllOracleClaimTypes {
		att := keeper.GetLastObservedAttestation(ctx, v)
		if att != nil {
			result[v.String()] = att.Nonce
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
	Valset     types.Valset
	Signatures []string
	Checkpoint []byte
}

func lastApprovedMultiSigUpdate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	nonce := keeper.GetLastValsetApprovedNonce(ctx)
	return fetchMultiSigUpdateData(ctx, nonce, keeper)
}

func lastObservedMultiSigUpdate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	nonce := keeper.GetLastValsetObservedNonce(ctx)
	return fetchMultiSigUpdateData(ctx, nonce, keeper)
}

func fetchMultiSigUpdateData(ctx sdk.Context, nonce types.Nonce, keeper Keeper) ([]byte, error) {
	if nonce.IsEmpty() {
		return nil, nil
	}

	valset := keeper.GetValsetRequest(ctx, int64(nonce.Uint64())) // todo: revisit nonce type
	if valset == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "no valset found for nonce")
	}

	result := MultiSigUpdateResponse{
		Checkpoint: valset.GetCheckpoint(),
		Valset:     *valset,
	}

	// todo: revisit nonce type
	keeper.IterateValsetConfirmByNonce(ctx, int64(nonce.Uint64()), func(_ []byte, confirm types.MsgValsetConfirm) bool {
		result.Signatures = append(result.Signatures, confirm.Signature)
		return false
	})
	if len(result.Signatures) == 0 {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "no signatures found for nonce")
	}

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
	keeper.IterateOutgoingTXBatches(ctx, func(batchID uint64, batch types.OutgoingTxBatch) bool {
		found := keeper.HasOutgoingTXBatchConfirm(ctx, batchID, validatorAddr)
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

const MaxResults = 100 // todo: impl pagination

func allBatchesRequest(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var batches []types.OutgoingTxBatch
	keeper.IterateOutgoingTXBatches(ctx, func(batchID uint64, batch types.OutgoingTxBatch) bool {
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

func allAttestations(ctx sdk.Context, claimType string, keeper Keeper) ([]byte, error) {
	if !types.IsClaimType(claimType) {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}
	var attestations []types.Attestation
	keeper.IterateAttestationByClaimTypeDesc(ctx, types.ClaimType(claimType), func(_ []byte, att types.Attestation) bool {
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

func queryAttestation(ctx sdk.Context, claimType, nonceStr string, keeper Keeper) ([]byte, error) {
	if !types.IsClaimType(claimType) {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "claim type")
	}
	nonce, err := parseNonce(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "nonce")
	}
	attestation := keeper.GetAttestation(ctx, types.ClaimType(claimType), nonce)
	if attestation == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, *attestation)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// todo: we mix nonces as int64 and base64 bytes at the moment
func parseNonce(nonceArg string) (types.Nonce, error) {
	if len(nonceArg) != base64.StdEncoding.EncodedLen(8) {
		// not a byte nonce byte representation
		v, err := strconv.ParseUint(nonceArg, 10, 64)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "nonce")
		}
		return types.NonceFromUint64(v), nil
	}
	return base64.URLEncoding.DecodeString(nonceArg)
}
