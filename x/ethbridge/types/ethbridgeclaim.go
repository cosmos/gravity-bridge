package types

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

type EthBridgeClaim struct {
	Nonce          int
	EthereumSender string
	CosmosReceiver sdk.AccAddress
	Validator      sdk.AccAddress
	Amount         sdk.Coins
}

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(nonce int, ethereumSender string, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) EthBridgeClaim {
	return EthBridgeClaim{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
		CosmosReceiver: cosmosReceiver,
		Validator:      validator,
		Amount:         amount,
	}
}

func CreateOracleClaimFromEthClaim(cdc *codec.Codec, ethClaim EthBridgeClaim) oracle.Claim {
	id := strconv.Itoa(ethClaim.Nonce) + ethClaim.EthereumSender
	claimBytes := cdc.MustMarshalBinaryBare(ethClaim)
	claim := oracle.NewClaim(id, claimBytes)
	return claim
}
