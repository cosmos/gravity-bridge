package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/peggy/x/oracle"
)

type EthBridgeClaim struct {
	Nonce            int             `json:"nonce"`
	EthereumSender   EthereumAddress `json:"ethereum_sender"`
	CosmosReceiver   sdk.AccAddress  `json:"cosmos_receiver"`
	ValidatorAddress sdk.ValAddress  `json:"validator_address"`
	Amount           sdk.Coins       `json:"amount"`
}

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(nonce int, ethereumSender EthereumAddress, cosmosReceiver sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coins) EthBridgeClaim {
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

// CreateOracleClaimFromEthClaim converts a specific ethereum bridge claim to a general oracle claim to be used by
// the oracle module. The oracle module expects every claim for a particular prophecy to have the same id, so this id
// must be created in a deterministic way that all validators can follow. For this, we use the Nonce an Ethereum Sender provided,
// as all validators will see this same data from the smart contract.
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

// CreateEthClaimFromOracleString converts a string from any generic claim from the oracle module into an ethereum bridge specific claim.
func CreateEthClaimFromOracleString(nonce int, ethereumAddress EthereumAddress, validator sdk.ValAddress, oracleClaimString string) (EthBridgeClaim, sdk.Error) {
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

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module. In general, it is
// expecte that the oracle module will store claims in this JSON format and so this should be used to convert oracle claims.
func CreateOracleClaimFromOracleString(oracleClaimString string) (OracleClaimContent, sdk.Error) {
	var oracleClaimContent OracleClaimContent

	stringBytes := []byte(oracleClaimString)
	errRes := json.Unmarshal(stringBytes, &oracleClaimContent)
	if errRes != nil {
		return OracleClaimContent{}, sdk.ErrInternal(fmt.Sprintf("failed to parse claim: %s", errRes))
	}

	return oracleClaimContent, nil
}
