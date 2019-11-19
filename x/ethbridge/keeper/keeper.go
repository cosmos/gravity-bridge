package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/peggy/x/ethbridge/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/supply"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	supplyKeeper supply.Keeper
	codespace    sdk.CodespaceType
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc *codec.Codec, supplyKeeper supply.Keeper, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
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
