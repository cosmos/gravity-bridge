package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/common"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

type EthBridgeClaim struct {
	Nonce            int                    `json:"nonce"`
	EthereumSender   common.EthereumAddress `json:"ethereum_sender"`
	CosmosReceiver   sdk.AccAddress         `json:"cosmos_receiver"`
	ValidatorAddress sdk.AccAddress         `json:"validator_address"`
	Amount           sdk.Coins              `json:"amount"`
}

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(nonce int, ethereumSender common.EthereumAddress, cosmosReceiver sdk.AccAddress, validator sdk.AccAddress, amount sdk.Coins) EthBridgeClaim {
	return EthBridgeClaim{
		Nonce:            nonce,
		EthereumSender:   ethereumSender,
		CosmosReceiver:   cosmosReceiver,
		ValidatorAddress: validator,
		Amount:           amount,
	}
}

//OracleClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
type OracleClaimContent struct {
	CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
	Amount         sdk.Coins      `json:"amount"`
}

// NewOracleClaim is a constructor function for OracleClaim
func NewOracleClaimContent(cosmosReceiver sdk.AccAddress, amount sdk.Coins) OracleClaimContent {
	return OracleClaimContent{
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
	}
}

func CreateOracleClaimFromEthClaim(cdc *codec.Codec, ethClaim EthBridgeClaim) oracle.Claim {
	oracleId := strconv.Itoa(ethClaim.Nonce) + string(ethClaim.EthereumSender)
	claimContent := NewOracleClaimContent(ethClaim.CosmosReceiver, ethClaim.Amount)
	claimBytes, _ := json.Marshal(claimContent)
	claimString := string(claimBytes)
	validator := sdk.ValAddress(ethClaim.ValidatorAddress)
	claim := oracle.NewClaim(oracleId, validator, claimString)
	return claim
}

func CreateEthClaimFromOracleString(nonce int, ethereumSender string, validator sdk.ValAddress, oracleClaimString string) (EthBridgeClaim, sdk.Error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return EthBridgeClaim{}, err
	}

	valAccAddress := sdk.AccAddress(validator)
	return NewEthBridgeClaim(
		nonce,
		common.EthereumAddress(ethereumSender),
		oracleClaim.CosmosReceiver,
		valAccAddress,
		oracleClaim.Amount,
	), nil
}

func CreateOracleClaimFromOracleString(oracleClaimString string) (OracleClaimContent, sdk.Error) {
	var oracleClaimContent OracleClaimContent

	stringBytes := []byte(oracleClaimString)
	errRes := json.Unmarshal(stringBytes, &oracleClaimContent)
	if errRes != nil {
		return OracleClaimContent{}, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", errRes))
	}

	return oracleClaimContent, nil
}
