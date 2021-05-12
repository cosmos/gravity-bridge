package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the params of the gravity module
func (k Keeper) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.ParamsResponse{Params: params}, nil

}

func (k Keeper) CurrentSignerSetTx(
	c context.Context,
	req *types.CurrentSignerSetTxRequest) (*types.CurrentSignerSetTxResponse, error) {
	return &types.CurrentSignerSetTxResponse{SignerSetTx: k.GetCurrentSignerSetTx(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) SignerSetTx(
	c context.Context,
	req *types.SignerSetTxRequest) (*types.SignerSetTxResponse, error) {
	return &types.SignerSetTxResponse{SignerSetTx: k.GetSignerSetTx(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

func (k Keeper) SignerSetTxSignature(
	c context.Context,
	req *types.SignerSetTxSignatureRequest) (*types.SignerSetTxSignatureResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}
	return &types.SignerSetTxSignatureResponse{SignatureMsg: k.GetSignerSetTxSignature(sdk.UnwrapSDKContext(c), req.Nonce, addr)}, nil
}

func (k Keeper) SignerSetTxSignaturesByNonce(
	c context.Context,
	req *types.SignerSetTxSignaturesByNonceRequest) (*types.SignerSetTxSignaturesByNonceResponse, error) {
	var confirms []*types.MsgSignerSetTxSignature
	k.IterateSignerSetTxSignatureByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(_ []byte, c types.MsgSignerSetTxSignature) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.SignerSetTxSignaturesByNonceResponse{SignatureMsgs: confirms}, nil
}

func (k Keeper) LastSignerSetTxs(
	c context.Context,
	req *types.LastSignerSetTxsRequest) (*types.LastSignerSetTxsResponse, error) {
	valReq := k.GetSignerSetTxs(sdk.UnwrapSDKContext(c))
	valReqLen := len(valReq)
	retLen := 0
	if valReqLen < maxSignerSetTxRequestsReturned {
		retLen = valReqLen
	} else {
		retLen = maxSignerSetTxRequestsReturned
	}
	return &types.LastSignerSetTxsResponse{SignerSetTxs: valReq[0:retLen]}, nil
}

func (k Keeper) LastPendingSignerSetTxByAddr(
	c context.Context,
	req *types.LastPendingSignerSetTxByAddrRequest) (*types.LastPendingSignerSetTxByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingSignerSetTxReq []*types.SignerSetTx
	k.IterateSignerSetTxs(sdk.UnwrapSDKContext(c), func(_ []byte, val *types.SignerSetTx) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := k.GetSignerSetTxSignature(sdk.UnwrapSDKContext(c), val.Nonce, addr) != nil
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
	return &types.LastPendingSignerSetTxByAddrResponse{SignerSetTxs: pendingSignerSetTxReq}, nil
}

// BatchFees queries the batch fees from unbatched pool
func (k Keeper) BatchFees(
	c context.Context,
	req *types.BatchFeeRequest) (*types.BatchFeeResponse, error) {
	return &types.BatchFeeResponse{BatchFees: k.GetAllBatchFees(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) LastPendingBatchTxByAddr(
	c context.Context,
	req *types.LastPendingBatchTxByAddrRequest) (*types.LastPendingBatchTxByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingBatchReq *types.BatchTx
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.BatchTx) bool {
		foundConfirm := k.GetBatchConfirm(sdk.UnwrapSDKContext(c), batch.BatchNonce, batch.TokenContract, addr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})

	return &types.LastPendingBatchTxByAddrResponse{Batch: pendingBatchReq}, nil
}

func (k Keeper) LastPendingContractCallTxByAddr(
	c context.Context,
	req *types.LastPendingContractCallTxByAddrRequest) (*types.LastPendingContractCallTxByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	var pendingLogicReq *types.ContractCallTx
	k.IterateContractCallTxs(sdk.UnwrapSDKContext(c), func(_ []byte, logic *types.ContractCallTx) bool {
		foundConfirm := k.GetLogicCallConfirm(sdk.UnwrapSDKContext(c),
			logic.InvalidationId, logic.InvalidationNonce, addr) != nil
		if !foundConfirm {
			pendingLogicReq = logic
			return true
		}
		return false
	})
	return &types.LastPendingContractCallTxByAddrResponse{Call: pendingLogicReq}, nil
}

func (k Keeper) BatchTxs(
	c context.Context,
	req *types.BatchTxsRequest) (*types.BatchTxsResponse, error) {
	var batches []*types.BatchTx
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.BatchTx) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	return &types.BatchTxsResponse{Batches: batches}, nil
}

func (k Keeper) ContractCallTxs(
	c context.Context,
	req *types.ContractCallTxsRequest) (*types.ContractCallTxsResponse, error) {
	var calls []*types.ContractCallTx
	k.IterateContractCallTxs(sdk.UnwrapSDKContext(c), func(_ []byte, call *types.ContractCallTx) bool {
		calls = append(calls, call)
		return len(calls) == MaxResults
	})
	return &types.ContractCallTxsResponse{Calls: calls}, nil
}

func (k Keeper) BatchTxByNonce(
	c context.Context,
	req *types.BatchTxByNonceRequest) (*types.BatchTxByNonceResponse, error) {
	if err := types.ValidateEthAddress(req.ContractAddress); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(c), req.ContractAddress, req.Nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	return &types.BatchTxByNonceResponse{Batch: foundBatch}, nil
}

func (k Keeper) BatchTxSignatures(
	c context.Context,
	req *types.BatchTxSignaturesRequest) (*types.BatchTxSignaturesResponse, error) {
	var confirms []*types.MsgBatchTxSignature
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c),
		req.Nonce, req.ContractAddress, func(_ []byte, c types.MsgBatchTxSignature) bool {
			confirms = append(confirms, &c)
			return false
		})
	return &types.BatchTxSignaturesResponse{SignatureMsgs: confirms}, nil
}

func (k Keeper) ContractCallTxSignatures(
	c context.Context,
	req *types.ContractCallTxSignaturesRequest) (*types.ContractCallTxSignaturesResponse, error) {
	var confirms []*types.MsgContractCallTxSignature
	k.IterateLogicConfirmByInvalidationIDAndNonce(sdk.UnwrapSDKContext(c), req.InvalidationId,
		req.InvalidationNonce, func(_ []byte, c *types.MsgContractCallTxSignature) bool {
			confirms = append(confirms, c)
			return false
		})

	return &types.ContractCallTxSignaturesResponse{Confirms: confirms}, nil
}

// LastEventNonceByAddr returns the last event nonce for the given validator address,
// this allows eth oracles to figure out where they left off
func (k Keeper) LastEventNonceByAddr(
	c context.Context,
	req *types.LastEventNonceByAddrRequest) (*types.LastEventNonceByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var ret types.LastEventNonceByAddrResponse
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, req.Address)
	}
	validator := k.GetOrchestratorValidator(ctx, addr)
	if validator == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "address")
	}
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, validator)
	ret.EventNonce = lastEventNonce
	return &ret, nil
}

// DenomToERC20 queries the Cosmos Denom that maps to an Ethereum ERC20
func (k Keeper) DenomToERC20(
	c context.Context,
	req *types.DenomToERC20Request) (*types.DenomToERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, erc20, err := k.DenomToERC20Lookup(ctx, req.Denom)
	var ret types.DenomToERC20Response
	ret.Erc20 = erc20
	ret.CosmosOriginated = cosmosOriginated

	return &ret, err
}

// ERC20ToDenom queries the ERC20 contract that maps to an Ethereum ERC20 if any
func (k Keeper) ERC20ToDenom(
	c context.Context,
	req *types.ERC20ToDenomRequest) (*types.ERC20ToDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosOriginated, name := k.ERC20ToDenomLookup(ctx, req.Erc20)
	var ret types.ERC20ToDenomResponse
	ret.Denom = name
	ret.CosmosOriginated = cosmosOriginated

	return &ret, nil
}

func (k Keeper) GetDelegateKeyByValidator(
	c context.Context,
	req *types.DelegateKeysByValidatorAddress) (*types.DelegateKeysByValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	reqValidator, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		keyValidator, err := sdk.ValAddressFromBech32(key.Validator)
		// this should be impossible due to the validate basic on the set orchestrator message
		if err != nil {
			panic("Invalid validator addr in store!")
		}
		if reqValidator.Equals(keyValidator) {
			return &types.DelegateKeysByValidatorAddressResponse{EthAddress: key.EthAddress, OrchestratorAddress: key.Orchestrator}, nil
		}
	}

	return nil, sdkerrors.Wrap(types.ErrInvalid, "No validator")
}

func (k Keeper) GetDelegateKeyByOrchestrator(
	c context.Context,
	req *types.DelegateKeysByOrchestratorAddress) (*types.DelegateKeysByOrchestratorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	reqOrchestrator, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		keyOrchestrator, err := sdk.AccAddressFromBech32(key.Orchestrator)
		// this should be impossible due to the validate basic on the set orchestrator message
		if err != nil {
			panic("Invalid orchestrator addr in store!")
		}
		if reqOrchestrator.Equals(keyOrchestrator) {
			return &types.DelegateKeysByOrchestratorAddressResponse{ValidatorAddress: key.Validator, EthAddress: key.EthAddress}, nil
		}

	}
	return nil, sdkerrors.Wrap(types.ErrInvalid, "No validator")
}

func (k Keeper) GetDelegateKeyByEth(
	c context.Context,
	req *types.DelegateKeysByEthAddress) (*types.DelegateKeysByEthAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	if err := types.ValidateEthAddress(req.EthAddress); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid eth address")
	}
	for _, key := range keys {
		if req.EthAddress == key.EthAddress {
			return &types.DelegateKeysByEthAddressResponse{
				ValidatorAddress:    key.Validator,
				OrchestratorAddress: key.Orchestrator,
			}, nil
		}
	}

	return nil, sdkerrors.Wrap(types.ErrInvalid, "No validator")
}

func (k Keeper) GetPendingSendToEthereum(
	c context.Context,
	req *types.PendingSendToEthereumRequest) (*types.PendingSendToEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	batches := k.GetBatchTxs(ctx)
	unbatchedTx := k.GetPoolTransactions(ctx)
	senderAddress := req.SenderAddress
	var res *types.PendingSendToEthereumResponse

	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			if tx.Sender == senderAddress {
				res.TransfersInBatches = append(res.TransfersInBatches, tx)
			}
		}
	}

	for _, tx := range unbatchedTx {
		if tx.Sender == senderAddress {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
	}

	return res, nil
}
