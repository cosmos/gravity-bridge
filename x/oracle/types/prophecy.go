package types

const PendingStatus = "pending"
const SuccessStatus = "success"
const FailedStatus = "failed"

// Prophecy is a struct that contains all the metadata of an oracle ritual
type Prophecy struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Claims []Claim `json:"claims"`
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(id string, status string, claims []Claim) Prophecy {
	return Prophecy{
		ID:     id,
		Status: status,
		Claims: claims,
	}
}

// NewEmptyProphecy returns a blank prophecy, used with errors
func NewEmptyProphecy() Prophecy {
	return NewProphecy("", "", nil)
}
