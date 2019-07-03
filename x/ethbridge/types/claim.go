package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/peggy/x/oracle"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

type EthBridgeClaim struct {
	Nonce            int                `json:"nonce"`
	EthereumSender   gethCommon.Address `json:"ethereum_sender"`
	CosmosReceiver   sdk.AccAddress     `json:"cosmos_receiver"`
	ValidatorAddress sdk.ValAddress     `json:"validator_address"`
	Amount           sdk.Coins          `json:"amount"`
}

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(nonce int, ethereumSender gethCommon.Address, cosmosReceiver sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coins) EthBridgeClaim {
	return EthBridgeClaim{
		Nonce:            nonce,
		EthereumSender:   ethereumSender,
		CosmosReceiver:   cosmosReceiver,
		ValidatorAddress: validator,
		Amount:           amount,
	}
}

// OracleClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
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

func CreateOracleClaimFromEthClaim(cdc *codec.Codec, ethClaim EthBridgeClaim) (oracle.Claim, error) {
	oracleId := strconv.Itoa(ethClaim.Nonce) + ethClaim.EthereumSender.String()
	claimContent := NewOracleClaimContent(ethClaim.CosmosReceiver, ethClaim.Amount)
	claimBytes, err := json.Marshal(claimContent)
	if err != nil {
		return oracle.Claim{}, err
	}
	claimString := string(claimBytes)
	validator := sdk.ValAddress(ethClaim.ValidatorAddress)
	claim := oracle.NewClaim(oracleId, validator, claimString)
	return claim, nil
}

func CreateEthClaimFromOracleString(nonce int, ethereumAddress gethCommon.Address, validator sdk.ValAddress, oracleClaimString string) (EthBridgeClaim, sdk.Error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return EthBridgeClaim{}, err
	}

	valAccAddress := sdk.ValAddress(validator)
	return NewEthBridgeClaim(
		nonce,
		ethereumAddress,
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
