package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trinhtan/peggy/x/oracle"
)

// EthBridgeClaim is a structure that contains all the data for a particular bridge claim
type EthBridgeClaim struct {
	EthereumChainID       int             `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	BridgeContractAddress EthereumAddress `json:"bridge_registry_contract_address" yaml:"bridge_registry_contract_address"`
	Nonce                 int             `json:"nonce" yaml:"nonce"`
	Symbol                string          `json:"symbol" yaml:"symbol"`
	TokenContractAddress  EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	EthereumSender        EthereumAddress `json:"ethereum_sender" yaml:"ethereum_sender"`
	CosmosReceiver        sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	ValidatorAddress      sdk.ValAddress  `json:"validator_address" yaml:"validator_address"`
	Amount                int64           `json:"amount" yaml:"amount"`
	ClaimType             ClaimType       `json:"claim_type" yaml:"claim_type"`
}

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(ethereumChainID int, bridgeContract EthereumAddress,
	nonce int, symbol string, tokenContact EthereumAddress, ethereumSender EthereumAddress,
	cosmosReceiver sdk.AccAddress, validator sdk.ValAddress, amount int64, claimType ClaimType,
) EthBridgeClaim {
	return EthBridgeClaim{
		EthereumChainID:       ethereumChainID,
		BridgeContractAddress: bridgeContract,
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContact,
		EthereumSender:        ethereumSender,
		CosmosReceiver:        cosmosReceiver,
		ValidatorAddress:      validator,
		Amount:                amount,
		ClaimType:             claimType,
	}
}

// OracleClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
type OracleClaimContent struct {
	CosmosReceiver       sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	Amount               int64           `json:"amount" yaml:"amount"`
	Symbol               string          `json:"symbol" yaml:"symbol"`
	TokenContractAddress EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	ClaimType            ClaimType       `json:"claim_type" yaml:"claim_type"`
}

// NewOracleClaimContent is a constructor function for OracleClaim
func NewOracleClaimContent(
	cosmosReceiver sdk.AccAddress, amount int64, symbol string, tokenContractAddress EthereumAddress, claimType ClaimType,
) OracleClaimContent {
	return OracleClaimContent{
		CosmosReceiver:       cosmosReceiver,
		Amount:               amount,
		Symbol:               symbol,
		TokenContractAddress: tokenContractAddress,
		ClaimType:            claimType,
	}
}

// CreateOracleClaimFromEthClaim converts a specific ethereum bridge claim to a general oracle claim to be used by
// the oracle module. The oracle module expects every claim for a particular prophecy to have the same id, so this id
// must be created in a deterministic way that all validators can follow.
// For this, we use the Nonce an Ethereum Sender provided,
// as all validators will see this same data from the smart contract.
func CreateOracleClaimFromEthClaim(cdc *codec.Codec, ethClaim EthBridgeClaim) (oracle.Claim, error) {
	oracleID := strconv.Itoa(ethClaim.EthereumChainID) + strconv.Itoa(ethClaim.Nonce) + ethClaim.EthereumSender.String()
	claimContent := NewOracleClaimContent(ethClaim.CosmosReceiver, ethClaim.Amount,
		ethClaim.Symbol, ethClaim.TokenContractAddress, ethClaim.ClaimType)
	claimBytes, err := json.Marshal(claimContent)
	if err != nil {
		return oracle.Claim{}, err
	}
	claimString := string(claimBytes)
	claim := oracle.NewClaim(oracleID, ethClaim.ValidatorAddress, claimString)
	return claim, nil
}

// CreateEthClaimFromOracleString converts a string
// from any generic claim from the oracle module into an ethereum bridge specific claim.
func CreateEthClaimFromOracleString(
	ethereumChainID int, bridgeContract EthereumAddress, nonce int,
	ethereumAddress EthereumAddress, validator sdk.ValAddress, oracleClaimString string,
) (EthBridgeClaim, error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return EthBridgeClaim{}, err
	}

	return NewEthBridgeClaim(
		ethereumChainID,
		bridgeContract,
		nonce,
		oracleClaim.Symbol,
		oracleClaim.TokenContractAddress,
		ethereumAddress,
		oracleClaim.CosmosReceiver,
		validator,
		oracleClaim.Amount,
		oracleClaim.ClaimType,
	), nil
}

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module.
// In general, it is expected that the oracle module will store claims in this JSON format
// and so this should be used to convert oracle claims.
func CreateOracleClaimFromOracleString(oracleClaimString string) (OracleClaimContent, error) {
	var oracleClaimContent OracleClaimContent

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaimContent); err != nil {
		return OracleClaimContent{}, sdkerrors.Wrap(ErrJSONMarshalling, fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaimContent, nil
}
