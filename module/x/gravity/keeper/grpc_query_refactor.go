package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.ParamsResponse{params}, nil
}

func (k Keeper) SignerSetTx(c context.Context, req *types.SignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	// TODO: audit once we finalize storage
	storeIndex := sdk.Uint64ToBigEndian(req.Nonce)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		// handle not found case
	}

	ss, ok := otx.(*types.SignerSetTx)
	if !ok {
		// panic("this shouldn't happen")
	}

	// TODO: special case nonce = 0 to find latest
	// TODO: ensure that latest signer set tx nonce index is set properly
	// TODO: ensure nonce sequence starts at one

	return &types.SignerSetTxResponse{ss}, nil
}

func (k Keeper) BatchTx(c context.Context, req *types.BatchTxRequest) (*types.BatchTxResponse, error) {
	if !common.IsHexAddress(req.ContractAddress) {
		// return err
	}

	// TODO: audit once we finalize storage
	storeIndex := append(sdk.Uint64ToBigEndian(req.Nonce), common.Hex2Bytes(req.ContractAddress)...)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		// handle not found case
	}

	btx, ok := otx.(*types.BatchTx)
	if !ok {
		// panic()
	}

	// TODO: handle special case nonce = 0 to find latest by contract address

	return &types.BatchTxResponse{btx}, nil
}

func (k Keeper) ContractCallTx(c context.Context, req *types.ContractCallTxRequest) (*types.ContractCallTxResponse, error) {
	storeIndex := append(sdk.Uint64ToBigEndian(req.InvalidationNonce), req.InvalidationScope...)
	otx := k.GetOutgoingTx(sdk.UnwrapSDKContext(c), types.GetOutgoingTxKey(storeIndex))
	if otx == nil {
		// handle not found case
	}

	cctx, ok := otx.(*types.ContractCallTx)
	if !ok {
		// panic()
	}

	// TODO: figure out how to call latest

	return &types.ContractCallTxResponse{cctx}, nil
}

func (k Keeper) SignerSetTxs(c context.Context, req *types.SignerSetTxsRequest) (*types.SignerSetTxsResponse, error) {
	var out []*types.SignerSetTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sstx, ok := otx.(*types.SignerSetTx)
		if !ok {
			// handle error case
		}
		out = append(out, sstx)
		if len(out) < int(req.Count) {
			return true
		}
		return false
	})
	return &types.SignerSetTxsResponse{out}, nil
}

func (k Keeper) BatchTxs(c context.Context, req *types.BatchTxsRequest) (*types.BatchTxsResponse, error) {
	var out []*types.BatchTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sstx, ok := otx.(*types.BatchTx)
		if !ok {
			// handle error case
		}
		out = append(out, sstx)
		return false
	})
	return &types.BatchTxsResponse{}, nil
}

func (k Keeper) ContractCallTxs(c context.Context, req *types.ContractCallTxsRequest) (*types.ContractCallTxsResponse, error) {
	var out []*types.ContractCallTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sstx, ok := otx.(*types.ContractCallTx)
		if !ok {
			// handle error case
		}
		out = append(out, sstx)
		return false
	})
	return &types.ContractCallTxsResponse{}, nil
}

func (k Keeper) SignerSetTxEthereumSignatures(c context.Context, req *types.SignerSetTxEthereumSignaturesRequest) (*types.SignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeSignerSetTxKey(req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.SignerSetTxEthereumSignaturesResponse{[][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.SignerSetTxEthereumSignaturesResponse{out}, nil
}

func (k Keeper) BatchTxEthereumSignatures(c context.Context, req *types.BatchTxEthereumSignaturesRequest) (*types.BatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeBatchTxKey(common.HexToAddress(req.ContractAddress), req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.BatchTxEthereumSignaturesResponse{[][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.BatchTxEthereumSignaturesResponse{out}, nil
}

func (k Keeper) ContractCallTxEthereumSignatures(c context.Context, req *types.ContractCallTxEthereumSignaturesRequest) (*types.ContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.ContractCallTxEthereumSignaturesResponse{[][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.ContractCallTxEthereumSignaturesResponse{out}, nil
}

func (k Keeper) PendingSignerSetTxEthereumSignatures(c context.Context, req *types.PendingSignerSetTxEthereumSignaturesRequest) (*types.PendingSignerSetTxEthereumSignaturesResponse, error) {
	return &types.PendingSignerSetTxEthereumSignaturesResponse{}, nil
}

func (k Keeper) PendingBatchTxEthereumSignatures(c context.Context, req *types.PendingBatchTxEthereumSignaturesRequest) (*types.PendingBatchTxEthereumSignaturesResponse, error) {
	return &types.PendingBatchTxEthereumSignaturesResponse{}, nil
}
func (k Keeper) PendingContractCallTxEthereumSignatures(c context.Context, req *types.PendingContractCallTxEthereumSignaturesRequest) (*types.PendingContractCallTxEthereumSignaturesResponse, error) {
	return &types.PendingContractCallTxEthereumSignaturesResponse{}, nil
}
func (k Keeper) LastSubmittedEthereumEvent(c context.Context, req *types.LastSubmittedEthereumEventRequest) (*types.LastSubmittedEthereumEventResponse, error) {
	return &types.LastSubmittedEthereumEventResponse{}, nil
}
func (k Keeper) BatchTxFees(c context.Context, req *types.BatchTxFeesRequest) (*types.BatchTxFeesResponse, error) {
	return &types.BatchTxFeesResponse{}, nil
}
func (k Keeper) ERC20ToDenom(c context.Context, req *types.ERC20ToDenomRequest) (*types.ERC20ToDenomResponse, error) {
	return &types.ERC20ToDenomResponse{}, nil
}
func (k Keeper) DenomToERC20(c context.Context, req *types.DenomToERC20Request) (*types.DenomToERC20Response, error) {
	return &types.DenomToERC20Response{}, nil
}
func (k Keeper) PendingSendToEthereums(c context.Context, req *types.PendingSendToEthereumsRequest) (*types.PendingSendToEthereumsResponse, error) {
	return &types.PendingSendToEthereumsResponse{}, nil
}
func (k Keeper) DelegateKeysByValidator(c context.Context, req *types.DelegateKeysByValidatorAddress) (*types.DelegateKeysByValidatorAddressResponse, error) {
	return &types.DelegateKeysByValidatorAddressResponse{}, nil
}
func (k Keeper) DelegateKeysByEthereumSigner(c context.Context, req *types.DelegateKeysByEthereumSignerRequest) (*types.DelegateKeysByEthereumSignerResponse, error) {
	return &types.DelegateKeysByEthereumSignerResponse{}, nil
}
func (k Keeper) DelegateKeysByOrchestrator(c context.Context, req *types.DelegateKeysByOrchestratorAddress) (*types.DelegateKeysByOrchestratorAddressResponse, error) {
	return &types.DelegateKeysByOrchestratorAddressResponse{}, nil
}
