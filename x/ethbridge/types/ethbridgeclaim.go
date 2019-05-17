package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

//OracleClaim is the details of how the claim for each validator will be stored in the oracle
type OracleClaim struct {
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Amount         sdk.Coins      `json:"amount"`
}

// NewOracleClaim is a constructor function for OracleClaim
func NewOracleClaim(cosmosReceiver sdk.AccAddress, amount sdk.Coins) OracleClaim {
	return OracleClaim{
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
	}
}

func CreateOracleClaimFromEthClaim(cdc *codec.Codec, ethClaim EthBridgeClaim) (string, sdk.ValAddress, string) {
	oracleId := strconv.Itoa(ethClaim.Nonce) + ethClaim.EthereumSender
	claimContent := NewOracleClaim(ethClaim.CosmosReceiver, ethClaim.Amount)
	claimBytes, _ := json.Marshal(claimContent)
	claim := string(claimBytes)
	validator := sdk.ValAddress(ethClaim.Validator)
	return oracleId, validator, claim
}

func CreateEthClaimFromOracleString(nonce int, ethereumSender string, validator sdk.ValAddress, oracleClaimString string) (EthBridgeClaim, sdk.Error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return EthBridgeClaim{}, err
	}

	valAccAddress := sdk.AccAddress(validator)
	return NewEthBridgeClaim(
		nonce,
		ethereumSender,
		oracleClaim.CosmosReceiver,
		valAccAddress,
		oracleClaim.Amount,
	), nil
}

func CreateOracleClaimFromOracleString(oracleClaimString string) (OracleClaim, sdk.Error) {
	var oracleClaim OracleClaim

	stringBytes := []byte(oracleClaimString)
	errRes := json.Unmarshal(stringBytes, &oracleClaim)
	if errRes != nil {
		return OracleClaim{}, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", errRes))
	}

	return oracleClaim, nil
}
