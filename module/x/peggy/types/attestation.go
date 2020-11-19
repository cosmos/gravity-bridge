package types

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	_ AttestationDetails = &BridgeDeposit{}
	_ AttestationDetails = &WithdrawalBatch{}
)

// AttestationDetails represents the payload of an attestation.
// Hash creates hash of the object that is unique in the context of the chain
// the hash verifies orchestrator's submited payload data
type AttestationDetails interface {
	Hash() []byte
}

// Hash implements WithdrawBatch.Hash
func (b *WithdrawalBatch) Hash() []byte {
	path := fmt.Sprintf("%s/%d/", b.Erc20Token, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// Hash implements BridgeDeposit.Hash
func (b *BridgeDeposit) Hash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", b.Erc20Token.String(), string(b.EthereumSender), b.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}
