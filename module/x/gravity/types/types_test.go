package types

import (
	"bytes"
	"encoding/hex"
	mrand "math/rand"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestValsetConfirmHash(t *testing.T) {
	powers := []uint64{3333, 3333, 3333}
	ethAddresses := []string{
		"0xc783df8a850f42e7F7e57013759C285caa701eB6",
		"0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4",
		"0xE5904695748fe4A84b40b3fc79De2277660BD1D3",
	}
	members := make(EthereumSigners, len(powers))
	for i := range powers {
		members[i] = &EthereumSigner{
			Power:           powers[i],
			EthereumAddress: ethAddresses[i],
		}
	}

	var mem []*EthereumSigner
	for _, m := range members {
		mem = append(mem, m)
	}
	v := SignerSetTx{Signers: mem}
	// TODO: this is hardcoded to foo, replace?
	hash := v.GetCheckpoint("foo")
	hexHash := hex.EncodeToString(hash)
	correctHash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"
	assert.Equal(t, correctHash, hexHash)
}

func TestEthereumSigners_PowerDiff(t *testing.T) {
	specs := map[string]struct {
		start EthereumSigners
		diff  EthereumSigners
		exp   float64
	}{
		"no diff": {
			start: EthereumSigners{
				{Power: 1, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 2, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 3, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: EthereumSigners{
				{Power: 1, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 2, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 3, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: 0.0,
		},
		"one fifth": {
			start: EthereumSigners{
				{Power: 1073741823, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 1073741823, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 2147483646, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: EthereumSigners{
				{Power: 858993459, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 858993459, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 2576980377, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: 0.2,
		},
		"real world": {
			start: EthereumSigners{
				{Power: 678509841, EthereumAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 685294939, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 671724742, EthereumAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, EthereumAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 617443955, EthereumAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 6785098, EthereumAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
				{Power: 291759231, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			diff: EthereumSigners{
				{Power: 642345266, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 678509841, EthereumAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, EthereumAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, EthereumAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 671724742, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 617443955, EthereumAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 291759231, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
				{Power: 6785098, EthereumAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
			},
			exp: 0.010000000011641532,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			assert.Equal(t, spec.exp, spec.start.PowerDiff(spec.diff))
		})
	}
}

func TestValsetSort(t *testing.T) {
	specs := map[string]struct {
		src EthereumSigners
		exp EthereumSigners
	}{
		"by power desc": {
			src: EthereumSigners{
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
				{Power: 2, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 3, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
			},
			exp: EthereumSigners{
				{Power: 3, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 2, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
		},
		"by eth addr on same power": {
			src: EthereumSigners{
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
			exp: EthereumSigners{
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String()},
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(2)}, 20)).String()},
				{Power: 1, EthereumAddress: gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(3)}, 20)).String()},
			},
		},
		// if you're thinking about changing this due to a change in the sorting algorithm
		// you MUST go change this in gravity_utils/types.rs as well. You will also break all
		// bridges in production when they try to migrate so use extreme caution!
		"real world": {
			src: EthereumSigners{
				{Power: 678509841, EthereumAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 685294939, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 671724742, EthereumAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, EthereumAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 617443955, EthereumAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 6785098, EthereumAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
				{Power: 291759231, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
			},
			exp: EthereumSigners{
				{Power: 685294939, EthereumAddress: "0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D"},
				{Power: 678509841, EthereumAddress: "0x6db48cBBCeD754bDc760720e38E456144e83269b"},
				{Power: 671724742, EthereumAddress: "0x0A7254b318dd742A3086882321C27779B4B642a6"},
				{Power: 671724742, EthereumAddress: "0x454330deAaB759468065d08F2b3B0562caBe1dD1"},
				{Power: 671724742, EthereumAddress: "0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140"},
				{Power: 617443955, EthereumAddress: "0x3511A211A6759d48d107898302042d1301187BA9"},
				{Power: 291759231, EthereumAddress: "0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326"},
				{Power: 6785098, EthereumAddress: "0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8"},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			spec.src.Sort()
			assert.Equal(t, spec.src, spec.exp)
			shuffled := shuffled(spec.src)
			shuffled.Sort()
			assert.Equal(t, shuffled, spec.exp)
		})
	}
}

func shuffled(v EthereumSigners) EthereumSigners {
	mrand.Shuffle(len(v), func(i, j int) {
		v[i], v[j] = v[j], v[i]
	})
	return v
}
