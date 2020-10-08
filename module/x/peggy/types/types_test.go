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
		members[i] = BridgeValidator{
			Power:           powers[i],
			EthereumAddress: ethAddresses[i],
		}
	}

	v := Valset{Members: members}
	hash := v.GetCheckpoint()
	hexHash := hex.EncodeToString(hash)
	correctHash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"
	assert.Equal(t, correctHash, hexHash)
}

func TestValsetCheckpointGold1(t *testing.T) {
	src := NewValset(0xc, BridgeValidators{{
		Power:           0xffffffff,
		EthereumAddress: EthereumAddress{0xb4, 0x62, 0x86, 0x4e, 0x39, 0x5d, 0x88, 0xd6, 0xbc, 0x7c, 0x5d, 0xd5, 0xf3, 0xf5, 0xeb, 0x4c, 0xc2, 0x59, 0x92, 0x55},
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
				{Power: 1, EthereumAddress: createEthAddress(3)},
				{Power: 2, EthereumAddress: createEthAddress(1)},
				{Power: 3, EthereumAddress: createEthAddress(2)},
			},
			exp: BridgeValidators{
				{Power: 3, EthereumAddress: createEthAddress(2)},
				{Power: 2, EthereumAddress: createEthAddress(1)},
				{Power: 1, EthereumAddress: createEthAddress(3)},
			},
		},
		"by eth addr on same power": {
			src: BridgeValidators{
				{Power: 1, EthereumAddress: createEthAddress(2)},
				{Power: 1, EthereumAddress: createEthAddress(1)},
				{Power: 1, EthereumAddress: createEthAddress(3)},
			},
			exp: BridgeValidators{
				{Power: 1, EthereumAddress: createEthAddress(1)},
				{Power: 1, EthereumAddress: createEthAddress(2)},
				{Power: 1, EthereumAddress: createEthAddress(3)},
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
