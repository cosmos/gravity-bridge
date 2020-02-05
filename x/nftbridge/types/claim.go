package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/oracle"
)

// NFTBridgeClaim is a structure that contains all the data for a particular bridge claim
type NFTBridgeClaim struct {
	EthereumChainID       int                       `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	BridgeContractAddress ethbridge.EthereumAddress `json:"bridge_contract_address" yaml:"bridge_contract_address"`
	Nonce                 int                       `json:"nonce" yaml:"nonce"`
	Symbol                string                    `json:"symbol" yaml:"symbol"`
	TokenContractAddress  ethbridge.EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	EthereumSender        ethbridge.EthereumAddress `json:"ethereum_sender" yaml:"ethereum_sender"`
	CosmosReceiver        sdk.AccAddress            `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	ValidatorAddress      sdk.ValAddress            `json:"validator_address" yaml:"validator_address"`
	Denom                 string                    `json:"denom" yaml:"denom"`
	ID                    string                    `json:"id" yaml:"id"`
	ClaimType             ethbridge.ClaimType       `json:"claim_type" yaml:"claim_type"`
}

// NewNFTBridgeClaim is a constructor function for NewNFTBridgeClaim
func NewNFTBridgeClaim(ethereumChainID int, bridgeContract ethbridge.EthereumAddress, nonce int, symbol string, tokenContact ethbridge.EthereumAddress, ethereumSender ethbridge.EthereumAddress, cosmosReceiver sdk.AccAddress, validator sdk.ValAddress, denom, id string, claimType ethbridge.ClaimType) NFTBridgeClaim {
	return NFTBridgeClaim{
		EthereumChainID:       ethereumChainID,
		BridgeContractAddress: bridgeContract,
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContact,
		EthereumSender:        ethereumSender,
		CosmosReceiver:        cosmosReceiver,
		ValidatorAddress:      validator,
		Denom:                 denom,
		ID:                    id,
		ClaimType:             claimType,
	}
}

// OracleNFTClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
type OracleNFTClaimContent struct {
	CosmosReceiver sdk.AccAddress      `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	Denom          string              `json:"string" yaml:"string"`
	ID             string              `json:"id" yaml:"id"`
	ClaimType      ethbridge.ClaimType `json:"claim_type" yaml:"claim_type"`
}

// NewOracleNFTClaimContent is a constructor function for OracleClaim
func NewOracleNFTClaimContent(cosmosReceiver sdk.AccAddress, denom, id string, claimType ethbridge.ClaimType) OracleNFTClaimContent {
	return OracleNFTClaimContent{
		CosmosReceiver: cosmosReceiver,
		Denom:          denom,
		ID:             id,
		ClaimType:      claimType,
	}
}

// CreateOracleClaimFromNFTClaim converts a specific ethereum bridge claim to a general oracle claim to be used by
// the oracle module. The oracle module expects every claim for a particular prophecy to have the same id, so this id
// must be created in a deterministic way that all validators can follow. For this, we use the Nonce an Ethereum Sender provided,
// as all validators will see this same data from the smart contract.
func CreateOracleClaimFromNFTClaim(cdc *codec.Codec, nftClaim NFTBridgeClaim) (oracle.Claim, error) {
	oracleID := strconv.Itoa(nftClaim.EthereumChainID) + strconv.Itoa(nftClaim.Nonce) + nftClaim.EthereumSender.String()
	claimContent := NewOracleNFTClaimContent(nftClaim.CosmosReceiver, nftClaim.Denom, nftClaim.ID, nftClaim.ClaimType)
	claimBytes, err := json.Marshal(claimContent)
	if err != nil {
		return oracle.Claim{}, err
	}
	claimString := string(claimBytes)
	claim := oracle.NewClaim(oracleID, nftClaim.ValidatorAddress, claimString)
	return claim, nil
}

// CreateNFTClaimFromOracleString converts a string from any generic claim from the oracle module into an ethereum bridge specific claim.
func CreateNFTClaimFromOracleString(ethereumChainID int, bridgeContract ethbridge.EthereumAddress, nonce int, symbol string, tokenContract ethbridge.EthereumAddress, ethereumAddress ethbridge.EthereumAddress, validator sdk.ValAddress, oracleClaimString string) (NFTBridgeClaim, error) {
	oracleClaim, err := CreateOracleNFTClaimFromOracleString(oracleClaimString)
	if err != nil {
		return NFTBridgeClaim{}, err
	}

	return NewNFTBridgeClaim(
		ethereumChainID,
		bridgeContract,
		nonce,
		symbol,
		tokenContract,
		ethereumAddress,
		oracleClaim.CosmosReceiver,
		validator,
		oracleClaim.Denom,
		oracleClaim.ID,
		oracleClaim.ClaimType,
	), nil
}

// CreateOracleNFTClaimFromOracleString converts a JSON string into an OracleNFTClaimContent struct used by this module. In general, it is
// expected that the oracle module will store claims in this JSON format and so this should be used to convert oracle claims.
func CreateOracleNFTClaimFromOracleString(oracleClaimString string) (OracleNFTClaimContent, error) {
	var oracleClaimContent OracleNFTClaimContent

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaimContent); err != nil {
		return OracleNFTClaimContent{}, sdkerrors.Wrap(ErrJSONMarshalling, fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaimContent, nil
}
