package types

import (
	"encoding/json"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var claimTypeToNames = map[ClaimType]string{
	CLAIM_TYPE_ETHEREUM_BRIDGE_DEPOSIT:          "bridge_deposit",
	CLAIM_TYPE_ETHEREUM_BRIDGE_WITHDRAWAL_BATCH: "bridge_withdrawal_batch",
}

// AllOracleClaimTypes types that are observed and submitted by the current orchestrator set
var AllOracleClaimTypes = []ClaimType{CLAIM_TYPE_ETHEREUM_BRIDGE_DEPOSIT, CLAIM_TYPE_ETHEREUM_BRIDGE_WITHDRAWAL_BATCH}

// ClaimTypeFromName given a string, returns the claim of teh given type if it exists
func ClaimTypeFromName(s string) (ClaimType, bool) {
	for _, v := range AllOracleClaimTypes {
		name, ok := claimTypeToNames[v]
		if ok && name == s {
			return v, true
		}
	}
	return CLAIM_TYPE_UNKNOWN, false
}

// ToClaimTypeNames returns the string representation of the claim type
func ToClaimTypeNames(s ...ClaimType) []string {
	r := make([]string, len(s))
	for i := range s {
		r[i] = s[i].String()
	}
	return r
}

func (claim ClaimType) String() string {
	return claimTypeToNames[claim]
}

// Bytes implements bytes
func (claim ClaimType) Bytes() []byte {
	return []byte{byte(claim)}
}

// MarshalJSON implements proto.Message
func (claim ClaimType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", claim.String())), nil
}

// UnmarshalJSON implements proto.Message
func (claim ClaimType) UnmarshalJSON(input []byte) error {
	if string(input) == `""` {
		return nil
	}
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	c, exists := ClaimTypeFromName(s)
	if !exists {
		return sdkerrors.Wrap(ErrUnknown, "claim type")
	}
	claim = c
	return nil
}

var (
	_ AttestationDetails = &BridgeDeposit{}
	_ AttestationDetails = &WithdrawalBatch{}
)

// AttestationDetails is the payload of an attestation.
type AttestationDetails interface {
	// Hash creates hash of the object that is supposed to be unique during the live time of the block chain.
	// purpose of the hash is to very that orchestrators submit the same payload data and not only the nonce.
	Hash() []byte
}

// Hash implements hash
func (b *WithdrawalBatch) Hash() []byte {
	path := fmt.Sprintf("%s/%d/", b.Erc20Token, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// Hash implements Hash
func (b *BridgeDeposit) Hash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", b.Erc20Token.String(), string(b.EthereumSender), b.CosmosReceiver)
	return tmhash.Sum([]byte(path))
}
