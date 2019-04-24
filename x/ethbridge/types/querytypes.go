package types

import (
	"encoding/json"
	"fmt"
)

// defines the params for the following queries:
// - 'custom/ethbridge/prophecies/'
type QueryEthProphecyParams struct {
	ID string
}

func NewQueryEthProphecyParams(id string) QueryEthProphecyParams {
	return QueryEthProphecyParams{
		ID: id,
	}
}

// Query Result Payload for an eth prophecy query
type QueryEthProphecyResponse struct {
	ID              string           `json:"id"`
	Status          string           `json:"status"`
	EthBridgeClaims []EthBridgeClaim `json:"claims"`
}

func NewQueryEthProphecyResponse(id string, status string, claims []EthBridgeClaim) QueryEthProphecyResponse {
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
