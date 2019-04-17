package types

// ProgressUpdate is a struct that contains the progress update as a result of processing a single validator's claim
type ProgressUpdate struct {
	Status     string `json:"status"`
	FinalBytes []byte `json:"final_bytes"`
}

// NewProgressUpdate returns a new ProgressUpdate with the given data contained
func NewProgressUpdate(status string, finalBytes []byte) ProgressUpdate {
	return ProgressUpdate{
		Status:     status,
		FinalBytes: finalBytes,
	}
}
