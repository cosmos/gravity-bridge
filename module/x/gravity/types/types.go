package types

import (
	"encoding/binary"
	"math"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// UInt64FromBytes create uint from binary big endian representation
func UInt64FromBytes(s []byte) uint64 {
	return binary.BigEndian.Uint64(s)
}

// UInt64Bytes uses the SDK byte marshaling to encode a uint64
func UInt64Bytes(n uint64) []byte {
	return sdk.Uint64ToBigEndian(n)
}

// UInt64FromString to parse out a uint64 for a nonce
func UInt64FromString(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

//////////////////////////////////////
//      Ethereum Signer(S)         //
//////////////////////////////////////

// ValidateBasic performs stateless checks on validity
func (b *EthereumSigner) ValidateBasic() error {
	if b.Power == 0 {
		return sdkerrors.Wrap(ErrEmpty, "power")
	}
	if err := ValidateEthereumAddress(b.EthereumAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	if b.EthereumAddress == "" {
		return sdkerrors.Wrap(ErrEmpty, "address")
	}
	return nil
}

// EthereumSigners is the sorted set of validator data for Ethereum bridge MultiSig set
type EthereumSigners []EthereumSigner

// Sort sorts the validators by power
func (b EthereumSigners) Sort() {
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
func (b EthereumSigners) PowerDiff(c EthereumSigners) float64 {
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
func (b EthereumSigners) TotalPower() (out uint64) {
	for _, v := range b {
		out += v.Power
	}
	return
}

// HasDuplicates returns true if there are duplicates in the set
func (b EthereumSigners) HasDuplicates() bool {
	m := make(map[string]struct{}, len(b))
	for i := range b {
		m[b[i].EthereumAddress] = struct{}{}
	}
	return len(m) != len(b)
}

// GetPowers returns only the power values for all members
func (b EthereumSigners) GetPowers() []uint64 {
	r := make([]uint64, len(b))
	for i := range b {
		r[i] = b[i].Power
	}
	return r
}

// ValidateBasic performs stateless checks
func (b EthereumSigners) ValidateBasic() error {
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

// NewSignerSetTx returns a new valset
func NewSignerSetTx(nonce, height uint64, members EthereumSigners) *SignerSetTx {
	members.Sort()
	var mem []EthereumSigner
	for _, val := range members {
		mem = append(mem, val)
	}
	return &SignerSetTx{Nonce: nonce, Signers: mem}
}

// WithoutEmptyMembers returns a new Valset without member that have 0 power or an empty Ethereum address.
func (v *SignerSetTx) WithoutEmptyMembers() *SignerSetTx {
	if v == nil {
		return nil
	}
	r := SignerSetTx{Nonce: v.Nonce, Signers: make([]EthereumSigner, 0, len(v.Signers))}
	for i := range v.Signers {
		if err := v.Signers[i].ValidateBasic(); err == nil {
			r.Signers = append(r.Signers, v.Signers[i])
		}
	}
	return &r
}

// SignerSetTxs is a collection of valset
type SignerSetTxs []*SignerSetTx

func (v SignerSetTxs) Len() int {
	return len(v)
}

func (v SignerSetTxs) Less(i, j int) bool {
	return v[i].Nonce > v[j].Nonce
}

func (v SignerSetTxs) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// GetFees returns the total fees contained within a given batch
func (b BatchTx) GetFees() sdk.Int {
	sum := sdk.ZeroInt()
	for _, t := range b.Transactions {
		sum.Add(t.Erc20Fee.Amount)
	}
	return sum
}
