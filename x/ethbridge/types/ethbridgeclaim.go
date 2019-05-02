package types

import (
	"encoding/json"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

type EthBridgeClaim struct {
	Nonce          int            `json:"nonce"`
	EthereumSender string         `json:"ethereum_sender"`
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Validator      sdk.AccAddress `json:"validator"`
	Amount         sdk.Coins      `json:"amount"`
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
	claimBytes, _ := json.Marshal(ethClaim)
	claim := oracle.NewClaim(id, claimBytes)
	return claim
}
