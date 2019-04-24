package types

import (
	"encoding/json"
	"fmt"

	"github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle"
)

// defines the params for the following queries:
// - 'custom/ethbridge/prophecies/'
type QueryEthProphecyParams struct {
	Nonce          int
	EthereumSender string
}

func NewQueryEthProphecyParams(nonce int, ethereumSender string) QueryEthProphecyParams {
	return QueryEthProphecyParams{
		Nonce:          nonce,
		EthereumSender: ethereumSender,
	}
}

// Query Result Payload for an eth prophecy query
type QueryEthProphecyResponse struct {
	ID              string           `json:"id"`
	Status          oracle.Status    `json:"status"`
	EthBridgeClaims []EthBridgeClaim `json:"claims"`
}

func NewQueryEthProphecyResponse(id string, status oracle.Status, claims []EthBridgeClaim) QueryEthProphecyResponse {
	return QueryEthProphecyResponse{
		ID:              id,
		Status:          status,
		EthBridgeClaims: claims,
	}
}

func (response QueryEthProphecyResponse) String() string {
	prophecyJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}
