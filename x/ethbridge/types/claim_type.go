package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ClaimType is an enum used to represent the type of claim
type ClaimType int

const (
	LockText = ClaimType(iota)
	BurnText
)

var ClaimTypeToString = [...]string{"lock", "burn"}
var StringToClaimType = map[string]ClaimType{
	"lock": LockText,
	"burn": BurnText,
}

func (text ClaimType) String() string {
	return ClaimTypeToString[text]
}

func (text ClaimType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", text.String())), nil
}

func (text *ClaimType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	stringKey, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	if value, ok := StringToClaimType[stringKey]; ok {
		*text = value
		return nil
	} else {
		return ErrInvalidClaimType()
	}
}
