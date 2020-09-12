package types

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestValsetConfirmSig(t *testing.T) {
	powers := [3]int64{3333, 3333, 3333}
	ethAddresses := [3]string{"0xc783df8a850f42e7F7e57013759C285caa701eB6", "0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4", "0xE5904695748fe4A84b40b3fc79De2277660BD1D3"}
	var v = Valset{
		Nonce:        0,
		Powers:       powers[:],
		EthAddresses: ethAddresses[:],
	}
	hash := v.GetCheckpoint()
	hexHash := hex.EncodeToString(hash)
	correctHash := "88165860d955aee7dc3e83d9d1156a5864b708841965585d206dbef6e9e1a499"
	if correctHash != hexHash {
		panic(fmt.Sprintf("%s does not match correct hash %s\n", hexHash, correctHash))
	}

}
