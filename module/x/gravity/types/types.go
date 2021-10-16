package types

import (
	"bytes"
	"crypto/sha256"
	"math"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

//////////////////////////////////////
//      Ethereum Signer(S)         //
//////////////////////////////////////

// ValidateBasic performs stateless checks on validity
func (b *EthereumSigner) ValidateBasic() error {
	if !common.IsHexAddress(b.EthereumAddress) {
		return sdkerrors.Wrap(ErrInvalid, "ethereum address")
	}
	return nil
}

// EthereumSigners is the sorted set of validator data for Ethereum bridge MultiSig set
type EthereumSigners []*EthereumSigner

// Sort sorts the validators by power
func (b EthereumSigners) Sort() {
	sort.Slice(b, func(i, j int) bool {
		if b[i].Power == b[j].Power {
			// Secondary sort on eth address in case powers are equal
			return EthereumAddrLessThan(b[i].EthereumAddress, b[j].EthereumAddress)
		}
		return b[i].Power > b[j].Power
	})
}

// Hash takes the sha256sum of a representation of the signer set
func (b EthereumSigners) Hash() []byte {
	b.Sort()
	var out bytes.Buffer
	for _, s := range b {
		out.Write(append(common.HexToAddress(s.EthereumAddress).Bytes(), sdk.Uint64ToBigEndian(s.Power)...))
	}
	hash := sha256.Sum256(out.Bytes())
	return hash[:]
}

// PowerDiff returns the difference in power between two bridge validator sets
// note this is Gravity bridge power *not* Cosmos voting power. Cosmos voting
// power is based on the absolute number of tokens in the staking pool at any given
// time Gravity bridge power is normalized using the equation.
//
// validators cosmos voting power / total cosmos voting power in this block = gravity bridge power / u32_max
//
// As an example if someone has 52% of the Cosmos voting power when a validator set is created their Gravity
// bridge voting power is u32_max * .52
//
// Normalized voting power dramatically reduces how often we have to produce new validator set updates. For example
// if the total on chain voting power increases by 1% due to inflation, we shouldn't have to generate a new validator
// set, after all the validators retained their relative percentages during inflation and normalized Gravity bridge power
// shows no difference.
func (b EthereumSigners) PowerDiff(c EthereumSigners) float64 {
	// loop over b and initialize the map with their powers
	powers := map[string]int64{}
	for _, bv := range b {
		powers[bv.EthereumAddress] = int64(bv.Power)
	}

	// subtract c powers from powers in the map, initializing
	// uninitialized keys with negative numbers
	for _, es := range c {
		if val, ok := powers[es.EthereumAddress]; ok {
			powers[es.EthereumAddress] = val - int64(es.Power)
		} else {
			powers[es.EthereumAddress] = -int64(es.Power)
		}
	}

	var delta int64
	for _, v := range powers {
		// NOTE: we care about the absolute value of the changes
		delta += absInt(v)
	}

	return math.Abs(float64(delta) / float64(math.MaxUint32))
}

func absInt(x int64) int64 {
	if x < 0 {
		x = -x
	}
	return x
}

// TotalPower returns the total power in the bridge validator set
func (b EthereumSigners) TotalPower() (out uint64) {
	for _, v := range b {
		out += v.Power
	}
	return
}

// GetPowers returns only the power values for all members
func (b EthereumSigners) GetPowers() []uint64 {
	r := make([]uint64, len(b))
	for i := range b {
		r[i] = b[i].Power
	}
	return r
}

// NewSignerSetTx returns a new valset
func NewSignerSetTx(nonce, height uint64, members EthereumSigners) *SignerSetTx {
	members.Sort()
	var mem []*EthereumSigner
	for _, val := range members {
		mem = append(mem, val)
	}
	return &SignerSetTx{Nonce: nonce, Height: height, Signers: mem}
}

// GetFees returns the total fees contained within a given batch
func (b BatchTx) GetFees() sdk.Int {
	sum := sdk.ZeroInt()
	for _, t := range b.Transactions {
		sum.Add(t.Erc20Fee.Amount)
	}
	return sum
}
