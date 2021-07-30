package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type GravityHooks interface {
	AfterContractCallExecutedEvent(ctx sdk.Context, event ContractCallExecutedEvent)
	AfterERC20DeployedEvent(ctx sdk.Context, event ERC20DeployedEvent)
	AfterSignerSetExecutedEvent(ctx sdk.Context, event SignerSetTxExecutedEvent)
	AfterBatchExecutedEvent(ctx sdk.Context, event BatchExecutedEvent)
	AfterSendToCosmosEvent(ctx sdk.Context, event SendToCosmosEvent)
}

type MultiGravityHooks []GravityHooks

func NewMultiGravityHooks(hooks ...GravityHooks) MultiGravityHooks {
	return hooks
}

func (mghs MultiGravityHooks) AfterContractCallExecutedEvent(ctx sdk.Context, event ContractCallExecutedEvent) {
	for i := range mghs {
		mghs[i].AfterContractCallExecutedEvent(ctx, event)
	}
}

func (mghs MultiGravityHooks) AfterERC20DeployedEvent(ctx sdk.Context, event ERC20DeployedEvent) {
	for i := range mghs {
		mghs[i].AfterERC20DeployedEvent(ctx, event)
	}
}

func (mghs MultiGravityHooks) AfterSignerSetExecutedEvent(ctx sdk.Context, event SignerSetTxExecutedEvent) {
	for i := range mghs {
		mghs[i].AfterSignerSetExecutedEvent(ctx, event)
	}
}

func (mghs MultiGravityHooks) AfterBatchExecutedEvent(ctx sdk.Context, event BatchExecutedEvent) {
	for i := range mghs {
		mghs[i].AfterBatchExecutedEvent(ctx, event)
	}
}

func (mghs MultiGravityHooks) AfterSendToCosmosEvent(ctx sdk.Context, event SendToCosmosEvent) {
	for i := range mghs {
		mghs[i].AfterSendToCosmosEvent(ctx, event)
	}
}
