package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/peggy/x/oracle"
)

// QueryEthProphecyParams defines the params for the following queries:
// - 'custom/ethbridge/prophecies/'
type QueryEthProphecyParams struct {
	ChainID        string          `json:"chain_id"`
	Nonce          int             `json:"nonce"`
	EthereumSender EthereumAddress `json:"ethereum_sender"`
}

// QueryEthProphecyParams creates a new QueryEthProphecyParams
func NewQueryEthProphecyParams(chainID string, nonce int, ethereumSender EthereumAddress) QueryEthProphecyParams {
	return QueryEthProphecyParams{
		ChainID:        chainID,
		Nonce:          nonce,
		EthereumSender: ethereumSender,
	}
}

// Query Result Payload for an eth prophecy query
type QueryEthProphecyResponse struct {
	ID     string           `json:"id"`
	Status oracle.Status    `json:"status"`
	Claims []EthBridgeClaim `json:"claims"`
}

func NewQueryEthProphecyResponse(id string, status oracle.Status, claims []EthBridgeClaim) QueryEthProphecyResponse {
	return QueryEthProphecyResponse{
		ID:     id,
		Status: status,
		Claims: claims,
	}
}

func (response QueryEthProphecyResponse) String() string {
	prophecyJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}
