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

// Size is a custom size implementation to ensure that the Amount is encoded correctly
func (m *ERC20Token) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Amount.Size()
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = len(m.TokenContractAddress)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	return n
}

// Marshal is a custom marshal impl to ensure that the custom Size is used
func (m *ERC20Token) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

// MarshalTo is a custom marshal impl to ensure that the custom Size is used
func (m *ERC20Token) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

// MarshalToSizedBuffer is a custom marshal impl to ensure that the custom Size is used
func (m *ERC20Token) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenContractAddress) > 0 {
		i -= len(m.TokenContractAddress)
		copy(dAtA[i:], m.TokenContractAddress)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.TokenContractAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x12
	}
	if m.Amount.Size() > 0 {
		i -= m.Amount.Size()
		bz, err := m.Amount.Marshal()
		if err != nil {
			return 0, err
		}
		copy(dAtA[i:], bz)
		i = encodeVarintAttestation(dAtA, i, uint64(m.Amount.Size()))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
