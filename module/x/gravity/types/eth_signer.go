package types

import (
	"math"
	"sort"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic performs stateless checks on validity
func (b *BridgeValidator) ValidateBasic() error {
	if b.Power == 0 {
		return sdkerrors.Wrap(ErrEmpty, "power")
	}
	if err := ValidateEthAddress(b.EthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	if b.EthereumAddress == "" {
		return sdkerrors.Wrap(ErrEmpty, "address")
	}
	return nil
}

// BridgeValidators is the sorted set of validator data for Ethereum bridge MultiSig set
type BridgeValidators []*BridgeValidator

// Sort sorts the validators by power
func (b BridgeValidators) Sort() {
	sort.Slice(b, func(i, j int) bool {
		if b[i].Power == b[j].Power {
			// Secondary sort on eth address in case powers are equal
			return EthAddrLessThan(b[i].EthereumAddress, b[j].EthereumAddress)
		}
		return b[i].Power > b[j].Power
	})
}

// PowerDiff returns the difference in power between two bridge validator sets
// TODO: this needs to be potentially refactored
func (b BridgeValidators) PowerDiff(c BridgeValidators) float64 {
	powers := map[string]int64{}
	var totalB int64
	// loop over b and initialize the map with their powers
	for _, bv := range b {
		powers[bv.EthereumAddress] = int64(bv.Power)
		totalB += int64(bv.Power)
	}

	// subtract c powers from powers in the map, initializing
	// uninitialized keys with negative numbers
	for _, bv := range c {
		if val, ok := powers[bv.EthereumAddress]; ok {
			powers[bv.EthereumAddress] = val - int64(bv.Power)
		} else {
			powers[bv.EthereumAddress] = -int64(bv.Power)
		}
	}

	var delta float64
	for _, v := range powers {
		// NOTE: we care about the absolute value of the changes
		delta += math.Abs(float64(v))
	}

	return math.Abs(delta / float64(totalB))
}

// TotalPower returns the total power in the bridge validator set
func (b BridgeValidators) TotalPower() (out uint64) {
	for _, v := range b {
		out += v.Power
	}
	return
}

// HasDuplicates returns true if there are duplicates in the set
func (b BridgeValidators) HasDuplicates() bool {
	m := make(map[string]struct{}, len(b))
	for i := range b {
		m[b[i].EthereumAddress] = struct{}{}
	}
	return len(m) != len(b)
}

// GetPowers returns only the power values for all members
func (b BridgeValidators) GetPowers() []uint64 {
	r := make([]uint64, len(b))
	for i := range b {
		r[i] = b[i].Power
	}
	return r
}

// ValidateBasic performs stateless checks
func (b BridgeValidators) ValidateBasic() error {
	// TODO: check if the set is sorted here?
	if len(b) == 0 {
		return ErrEmpty
	}
	for i := range b {
		if err := b[i].ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "member %d", i)
		}
	}
	if b.HasDuplicates() {
		return sdkerrors.Wrap(ErrDuplicate, "addresses")
	}
	return nil
}

// NewValset returns a new valset
func NewValset(nonce, height uint64, members BridgeValidators) *Valset {
	members.Sort()
	var mem []*BridgeValidator
	for _, val := range members {
		mem = append(mem, val)
	}
	return &Valset{Nonce: uint64(nonce), Members: mem, Height: height}
}
