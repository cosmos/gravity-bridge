package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValsetConfirmHash(t *testing.T) {
	powers := []uint64{3333, 3333, 3333}
	ethAddresses := []EthereumAddress{
		NewEthereumAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"),
		NewEthereumAddress("0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4"),
		NewEthereumAddress("0xE5904695748fe4A84b40b3fc79De2277660BD1D3"),
	}

	var v = Valset{
		Nonce:        0,
		Powers:       powers[:],
		EthAddresses: ethAddresses[:],
	}
	hash := v.GetCheckpoint()
	hexHash := hex.EncodeToString(hash)
	correctHash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"
	assert.Equal(t, correctHash, hexHash)
}

func TestValsetCheckpointGold1(t *testing.T) {
	orchestratorAddr := EthereumAddress{0xb4, 0x62, 0x86, 0x4e, 0x39, 0x5d, 0x88, 0xd6, 0xbc, 0x7c, 0x5d, 0xd5, 0xf3, 0xf5, 0xeb, 0x4c, 0xc2, 0x59, 0x92, 0x55}
	src := Valset{
		Nonce:        0xc,
		Powers:       []uint64{0xffffffff},
		EthAddresses: []EthereumAddress{orchestratorAddr},
	}
	ourHash := src.GetCheckpoint()

	// hash from bridge contract
	goldHash := "0xf024ab7404464494d3919e5a7f0d8ac40804fb9bd39ad5d16cdb3e66aa219b64"[2:]
	assert.Equal(t, goldHash, hex.EncodeToString(ourHash))
}
