package types

import (
	"encoding/json"
	"fmt"
)

const PendingStatus = "pending"
const CompleteStatus = "complete"

// DefaultConsensusNeeded is the default fraction of validators needed to make claims on a prophecy in order for it to pass
const DefaultConsensusNeeded = 0.7

// Prophecy is a struct that contains all the metadata of an oracle ritual
type Prophecy struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	MinimumPower int     `json:"minimum_power"` //The minimum number of staked claiming power needed before completion logic is checked
	Claims       []Claim `json:"claims"`
}

func (prophecy Prophecy) String() string {
	prophecyJSON, err := json.Marshal(prophecy)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(id string, status string, minimumPower int, claims []Claim) Prophecy {
	return Prophecy{
		ID:           id,
		Status:       status,
		MinimumPower: minimumPower,
		Claims:       claims,
	}
}

// NewEmptyProphecy returns a blank prophecy, used with errors
func NewEmptyProphecy() Prophecy {
	return NewProphecy("", "", 0, nil)
}
