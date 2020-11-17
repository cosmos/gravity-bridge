package types

import (
	"bytes"
	"encoding/hex"
	mrand "math/rand"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestValsetConfirmHash(t *testing.T) {
	powers := []uint64{3333, 3333, 3333}
	ethAddresses := []EthereumAddress{
		NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"),
		NewEthereumAddress("0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"),
		NewEthereumAddress("0xE5904695748fe4A84b40b3fc79De2277660BD1D3"),
	}
	members := make(BridgeValidators, len(powers))
	for i := range powers {
		members[i] = &BridgeValidator{
			Power:           powers[i],
			EthereumAddress: ethAddresses[i].Bytes(),
		}
	}

	var mem []*BridgeValidator
	for _, m := range members {
		mem = append(mem, m)
	}
	v := Valset{Members: mem}
	hash := v.GetCheckpoint()
	hexHash := hex.EncodeToString(hash)
	correctHash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"
	assert.Equal(t, correctHash, hexHash)
}

func TestValsetCheckpointGold1(t *testing.T) {
	src := NewValset(0xc, BridgeValidators{{
		Power:           0xffffffff,
		EthereumAddress: EthereumAddress{0xb4, 0x62, 0x86, 0x4e, 0x39, 0x5d, 0x88, 0xd6, 0xbc, 0x7c, 0x5d, 0xd5, 0xf3, 0xf5, 0xeb, 0x4c, 0xc2, 0x59, 0x92, 0x55}.Bytes(),
	}})

	ourHash := src.GetCheckpoint()

	// hash from bridge contract
	goldHash := "0xf024ab7404464494d3919e5a7f0d8ac40804fb9bd39ad5d16cdb3e66aa219b64"[2:]
	assert.Equal(t, goldHash, hex.EncodeToString(ourHash))
}

func TestValsetSort(t *testing.T) {
	specs := map[string]struct {
		src BridgeValidators
		exp BridgeValidators
	}{
		"by power desc": {
			src: BridgeValidators{
				{Power: 1, EthereumAddress: createEthAddress(3).Bytes()},
				{Power: 2, EthereumAddress: createEthAddress(1).Bytes()},
				{Power: 3, EthereumAddress: createEthAddress(2).Bytes()},
			},
			exp: BridgeValidators{
				{Power: 3, EthereumAddress: createEthAddress(2).Bytes()},
				{Power: 2, EthereumAddress: createEthAddress(1).Bytes()},
				{Power: 1, EthereumAddress: createEthAddress(3).Bytes()},
			},
		},
		"by eth addr on same power": {
			src: BridgeValidators{
				{Power: 1, EthereumAddress: createEthAddress(2).Bytes()},
				{Power: 1, EthereumAddress: createEthAddress(1).Bytes()},
				{Power: 1, EthereumAddress: createEthAddress(3).Bytes()},
			},
			exp: BridgeValidators{
				{Power: 1, EthereumAddress: createEthAddress(1).Bytes()},
				{Power: 1, EthereumAddress: createEthAddress(2).Bytes()},
				{Power: 1, EthereumAddress: createEthAddress(3).Bytes()},
			},
		},
		// if you're thinking about changing this due to a change in the sorting algorithm
		// you MUST go change this in peggy_utils/types.rs as well. You will also break all
		// bridges in production when they try to migrate so use extreme caution!
		"real world": {
			src: BridgeValidators{
				{Power: 678509841, EthereumAddress: NewEthereumAddress("0x6db48cBBCeD754bDc760720e38E456144e83269b").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140").Bytes()},
				{Power: 685294939, EthereumAddress: NewEthereumAddress("0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x0A7254b318dd742A3086882321C27779B4B642a6").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x454330deAaB759468065d08F2b3B0562caBe1dD1").Bytes()},
				{Power: 617443955, EthereumAddress: NewEthereumAddress("0x3511A211A6759d48d107898302042d1301187BA9").Bytes()},
				{Power: 6785098, EthereumAddress: NewEthereumAddress("0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8").Bytes()},
				{Power: 291759231, EthereumAddress: NewEthereumAddress("0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326").Bytes()},
			},
			exp: BridgeValidators{
				{Power: 685294939, EthereumAddress: NewEthereumAddress("0x479FFc856Cdfa0f5D1AE6Fa61915b01351A7773D").Bytes()},
				{Power: 678509841, EthereumAddress: NewEthereumAddress("0x6db48cBBCeD754bDc760720e38E456144e83269b").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x0A7254b318dd742A3086882321C27779B4B642a6").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x454330deAaB759468065d08F2b3B0562caBe1dD1").Bytes()},
				{Power: 671724742, EthereumAddress: NewEthereumAddress("0x8E91960d704Df3fF24ECAb78AB9df1B5D9144140").Bytes()},
				{Power: 617443955, EthereumAddress: NewEthereumAddress("0x3511A211A6759d48d107898302042d1301187BA9").Bytes()},
				{Power: 291759231, EthereumAddress: NewEthereumAddress("0xF14879a175A2F1cEFC7c616f35b6d9c2b0Fd8326").Bytes()},
				{Power: 6785098, EthereumAddress: NewEthereumAddress("0x37A0603dA2ff6377E5C7f75698dabA8EE4Ba97B8").Bytes()},
			},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			// when
			spec.src.Sort()
			// then
			assert.Equal(t, spec.src, spec.exp)
		})
	}
}

func shuffled(v BridgeValidators) BridgeValidators {
	mrand.Shuffle(len(v), func(i, j int) {
		v[i], v[j] = v[j], v[i]
	})
	return v
}

func createEthAddress(i int) EthereumAddress {
	return EthereumAddress(gethCommon.BytesToAddress(bytes.Repeat([]byte{byte(i)}, 20)))
}
