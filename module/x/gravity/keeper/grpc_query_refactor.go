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
	return &types.ParamsResponse{Params: params}, nil
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

	return &types.SignerSetTxResponse{SignerSet: ss}, nil
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

	batch, ok := otx.(*types.BatchTx)
	if !ok {
		// panic()
	}

	// TODO: handle special case nonce = 0 to find latest by contract address

	return &types.BatchTxResponse{Batch: batch}, nil
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

	return &types.ContractCallTxResponse{LogicCall: cctx}, nil
}

func (k Keeper) SignerSetTxs(c context.Context, req *types.SignerSetTxsRequest) (*types.SignerSetTxsResponse, error) {
	var signers []*types.SignerSetTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		signer, ok := otx.(*types.SignerSetTx)
		if !ok {
			// handle error case
		}
		signers = append(signers, signer)

		return len(signers) < int(req.Count)
	})
	return &types.SignerSetTxsResponse{SignerSets: signers}, nil
}

func (k Keeper) BatchTxs(c context.Context, req *types.BatchTxsRequest) (*types.BatchTxsResponse, error) {
	var batches []*types.BatchTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		batch, ok := otx.(*types.BatchTx)
		if !ok {
			// handle error case
		}
		batches = append(batches, batch)
		return false
	})
	return &types.BatchTxsResponse{Batches: batches}, nil
}

func (k Keeper) ContractCallTxs(c context.Context, req *types.ContractCallTxsRequest) (*types.ContractCallTxsResponse, error) {
	var calls []*types.ContractCallTx
	k.IterateOutgoingTxs(sdk.UnwrapSDKContext(c), types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		call, ok := otx.(*types.ContractCallTx)
		if !ok {
			// handle error case
		}
		calls = append(calls, call)
		return false
	})
	return &types.ContractCallTxsResponse{Calls: calls}, nil
}

func (k Keeper) SignerSetTxEthereumSignatures(c context.Context, req *types.SignerSetTxEthereumSignaturesRequest) (*types.SignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeSignerSetTxKey(req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.SignerSetTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.SignerSetTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) BatchTxEthereumSignatures(c context.Context, req *types.BatchTxEthereumSignaturesRequest) (*types.BatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeBatchTxKey(common.HexToAddress(req.ContractAddress), req.Nonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.BatchTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.BatchTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) ContractCallTxEthereumSignatures(c context.Context, req *types.ContractCallTxEthereumSignaturesRequest) (*types.ContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	key := types.MakeContractCallTxKey(req.InvalidationScope, req.InvalidationNonce)
	if req.Address == "" {
		val, err := k.getSignerValidator(ctx, req.Address)
		if err != nil {
			return nil, err
		}
		return &types.ContractCallTxEthereumSignaturesResponse{Signature: [][]byte{k.GetEthereumSignature(ctx, key, val)}}, nil
	}

	var out [][]byte
	k.IterateEthereumSignatures(ctx, key, func(_ sdk.ValAddress, sig hexutil.Bytes) bool {
		out = append(out, sig)
		return false
	})
	return &types.ContractCallTxEthereumSignaturesResponse{Signature: out}, nil
}

func (k Keeper) PendingSignerSetTxEthereumSignatures(c context.Context, req *types.PendingSignerSetTxEthereumSignaturesRequest) (*types.PendingSignerSetTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var signerSets []*types.SignerSetTx
	k.IterateOutgoingTxs(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			signerSet, ok := otx.(*types.SignerSetTx)
			if !ok {
				// handle error case
			}
			signerSets = append(signerSets, signerSet)
		}
		return false
	})
	return &types.PendingSignerSetTxEthereumSignaturesResponse{signerSets}, nil
}

func (k Keeper) PendingBatchTxEthereumSignatures(c context.Context, req *types.PendingBatchTxEthereumSignaturesRequest) (*types.PendingBatchTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var batches []*types.BatchTx
	k.IterateOutgoingTxs(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			batch, ok := otx.(*types.BatchTx)
			if !ok {
				// handle error case
			}
			batches = append(batches, batch)
		}
		return false
	})
	return &types.PendingBatchTxEthereumSignaturesResponse{batches}, nil
}

func (k Keeper) PendingContractCallTxEthereumSignatures(c context.Context, req *types.PendingContractCallTxEthereumSignaturesRequest) (*types.PendingContractCallTxEthereumSignaturesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := k.getSignerValidator(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	var calls []*types.ContractCallTx
	k.IterateOutgoingTxs(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.GetEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 { // it's pending
			call, ok := otx.(*types.ContractCallTx)
			if !ok {
				// handle error case
			}
			calls = append(calls, call)
		}
		return false
	})
	return &types.PendingContractCallTxEthereumSignaturesResponse{calls}, nil
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
	ctx := sdk.UnwrapSDKContext(c)
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	ethAddr := k.GetValidatorEthereumAddress(ctx, valAddr)
	orchAddr := k.GetEthereumOrchestratorAddress(ctx, ethAddr)
	res := &types.DelegateKeysByValidatorAddressResponse{
		EthAddress:          ethAddr.Hex(),
		OrchestratorAddress: orchAddr.String(),
	}
	return res, nil
}

func (k Keeper) DelegateKeysByEthereumSigner(c context.Context, req *types.DelegateKeysByEthereumSignerRequest) (*types.DelegateKeysByEthereumSignerResponse, error) {
	return &types.DelegateKeysByEthereumSignerResponse{}, nil
}

func (k Keeper) DelegateKeysByOrchestrator(c context.Context, req *types.DelegateKeysByOrchestratorAddress) (*types.DelegateKeysByOrchestratorAddressResponse, error) {
	return &types.DelegateKeysByOrchestratorAddressResponse{}, nil
}
