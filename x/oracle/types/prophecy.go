package types

const PendingStatus = "pending"
const SuccessStatus = "success"
const FailedStatus = "failed"

// Prophecy is a struct that contains all the metadata of an oracle ritual
type Prophecy struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	MinimumPower int     `json:"minimum_power"` //The minimum number of staked claiming power needed before completion logic is checked
	Claims       []Claim `json:"claims"`
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
