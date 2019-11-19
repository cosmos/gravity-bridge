package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/oracle"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	supplyKeeper types.SupplyKeeper
	oracleKeeper types.OracleKeeper
	codespace    sdk.CodespaceType
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc *codec.Codec, supplyKeeper types.SupplyKeeper, oracleKeeper types.OracleKeeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		oracleKeeper: oracleKeeper,
		codespace:    codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.EthBridgeClaim) (oracle.Status, sdk.Error) {
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, types.ErrJSONMarshalling(k.Codespace())
	}

	status, sdkErr := k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
	if sdkErr != nil {
		return oracle.Status{}, sdkErr
	}
	return status, nil
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim string) sdk.Error {
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	receiverAddress := oracleClaim.CosmosReceiver

	if oracleClaim.ClaimType == types.LockText {
		err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, oracleClaim.Amount)
		if err != nil {
			return err
		}
	}
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, oracleClaim.Amount)
	if err != nil {
		panic(err)
	}
	return nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
	if err != nil {
		return err
	}
	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount)
	if err != nil {
		panic(err)
	}
	return nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
	if err != nil {
		return err
	}
	return nil
}
