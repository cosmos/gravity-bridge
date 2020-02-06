package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/nftbridge/types"
	"github.com/cosmos/peggy/x/oracle"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	nftKeeper    types.NFTKeeper
	oracleKeeper ethbridge.OracleKeeper
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(cdc *codec.Codec, nftKeeper types.NFTKeeper, oracleKeeper ethbridge.OracleKeeper) Keeper {
	return Keeper{
		cdc:          cdc,
		nftKeeper:    nftKeeper,
		oracleKeeper: oracleKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(ctx sdk.Context, claim types.BridgeClaim) (oracle.Status, error) {
	oracleClaim, err := types.CreateOracleClaimFromNFTClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, err
	}

	return k.oracleKeeper.ProcessClaim(ctx, oracleClaim)
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(ctx sdk.Context, claim string) error {
	oracleClaim, err := types.CreateOracleNFTClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	// moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	receiverAddress := oracleClaim.CosmosReceiver
	newNFT := nft.NewBaseNFT(oracleClaim.ID, receiverAddress, "")
	switch oracleClaim.ClaimType {
	case ethbridge.LockText:
		err = k.nftKeeper.MintNFT(ctx, oracleClaim.Denom, &newNFT)
	default:
		err = types.ErrInvalidClaimType
	}

	if err != nil {
		return err
	}

	return nil
}

// ProcessBurn processes the burn of bridged NFT from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.Address, denom, id string) error {
	cosmosNFT, err := k.nftKeeper.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}
	if !cosmosNFT.GetOwner().Equals(cosmosSender) {
		return types.ErrInvalidTokenAddress
	}
	if err := k.nftKeeper.DeleteNFT(ctx, denom, id); err != nil {
		return err
	}

	return nil
}

// ProcessLock processes the lockup of cosmos nft from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.Address, denom, id string) error {
	cosmosNFT, err := k.nftKeeper.GetNFT(ctx, denom, id)
	if err != nil {
		return err
	}
	if !cosmosNFT.GetOwner().Equals(cosmosSender) {
		return types.ErrInvalidTokenAddress
	}
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	cosmosNFT.SetOwner(moduleAcct)

	return k.nftKeeper.UpdateNFT(ctx, denom, cosmosNFT)
}
