package txs

import (
	"math/big"

	"github.com/trinhtan/peggy/cmd/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
)

// OracleClaim contains data required to make an OracleClaim
type OracleClaim struct {
	ProphecyID *big.Int
	Message    [32]byte
	Signature  []byte
}

// ProphecyClaim contains data required to make an ProphecyClaim
type ProphecyClaim struct {
	ClaimType            types.Event
	CosmosSender         []byte
	EthereumReceiver     common.Address
	Symbol               string
	Amount               *big.Int
}
