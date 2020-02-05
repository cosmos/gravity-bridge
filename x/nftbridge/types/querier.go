package types

import (
	"encoding/json"
	"fmt"

	ethbridge "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/cosmos/peggy/x/oracle"
)

// query endpoints supported by the oracle Querier
const (
	QueryNFTProphecy = "prophecies"
)

// QueryNFTProphecyParams defines the params for the following queries:
// - 'custom/nftbridge/prophecies/'
type QueryNFTProphecyParams struct {
	EthereumChainID       int                       `json:"ethereum_chain_id"`
	BridgeContractAddress ethbridge.EthereumAddress `json:"bridge_contract_address"`
	Nonce                 int                       `json:"nonce"`
	Symbol                string                    `json:"symbol"`
	TokenContractAddress  ethbridge.EthereumAddress `json:"token_contract_address"`
	EthereumSender        ethbridge.EthereumAddress `json:"ethereum_sender"`
}

// NewQueryNFTProphecyParams creates a new QueryNFTProphecyParams
func NewQueryNFTProphecyParams(
	ethereumChainID int, bridgeContractAddress ethbridge.EthereumAddress, nonce int, symbol string,
	tokenContractAddress ethbridge.EthereumAddress, ethereumSender ethbridge.EthereumAddress,
) QueryNFTProphecyParams {
	return QueryNFTProphecyParams{
		EthereumChainID:       ethereumChainID,
		BridgeContractAddress: bridgeContractAddress,
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContractAddress,
		EthereumSender:        ethereumSender,
	}
}

// QueryNFTProphecyResponse defines the result payload for an nft prophecy query
type QueryNFTProphecyResponse struct {
	ID     string           `json:"id"`
	Status oracle.Status    `json:"status"`
	Claims []NFTBridgeClaim `json:"claims"`
}

// NewQueryNFTProphecyResponse creates a new QueryNFTProphecyResponse instance
func NewQueryNFTProphecyResponse(
	id string, status oracle.Status, claims []NFTBridgeClaim,
) QueryNFTProphecyResponse {
	return QueryNFTProphecyResponse{
		ID:     id,
		Status: status,
		Claims: claims,
	}
}

// String implements fmt.Stringer interface
func (response QueryNFTProphecyResponse) String() string {
	prophecyJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}
