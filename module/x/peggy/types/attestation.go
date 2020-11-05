package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ClaimType is the cosmos type of an event from the counterpart chain that can be handled
type ClaimType byte

const (
	ClaimTypeUnknown                       ClaimType = 0
	ClaimTypeEthereumBridgeDeposit         ClaimType = 1
	ClaimTypeEthereumBridgeWithdrawalBatch ClaimType = 2
)

var claimTypeToNames = map[ClaimType]string{
	ClaimTypeEthereumBridgeDeposit:         "bridge_deposit",
	ClaimTypeEthereumBridgeWithdrawalBatch: "bridge_withdrawal_batch",
}

// AllOracleClaimTypes types that are observed and submitted by the current orchestrator set
var AllOracleClaimTypes = []ClaimType{ClaimTypeEthereumBridgeDeposit, ClaimTypeEthereumBridgeWithdrawalBatch}

func ClaimTypeFromName(s string) (ClaimType, bool) {
	for _, v := range AllOracleClaimTypes {
		name, ok := claimTypeToNames[v]
		if ok && name == s {
			return v, true
		}
	}
	return ClaimTypeUnknown, false
}
func ToClaimTypeNames(s ...ClaimType) []string {
	r := make([]string, len(s))
	for i := range s {
		r[i] = s[i].String()
	}
	return r
}

func (c ClaimType) String() string {
	return claimTypeToNames[c]
}

func (c ClaimType) Bytes() []byte {
	return []byte{byte(c)}
}

func (e ClaimType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", e.String())), nil
}

func (e *ClaimType) UnmarshalJSON(input []byte) error {
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
	*e = c
	return nil
}

// Attestation is an aggregate of `claims` that eventually becomes `observed` by all orchestrators
type Attestation struct {
	ClaimType  ClaimType          `json:"claim_type"`
	EventNonce UInt64Nonce        `json:"event_nonce"`
	Observed   bool               `json:"observed"`
	Votes      []sdk.ValAddress   `json:"votes"`
	Details    AttestationDetails `json:"details,omitempty"`
}

// AttestationDetails is the payload of an attestation.
type AttestationDetails interface {
	// Hash creates hash of the object that is supposed to be unique during the live time of the block chain.
	// purpose of the hash is to very that orchestrators submit the same payload data and not only the nonce.
	Hash() []byte
}

var (
	_ AttestationDetails = BridgeDeposit{}
	_ AttestationDetails = WithdrawalBatch{}
)

// WithdrawalBatch is an attestation detail that marks a batch of outgoing transactions executed and
// frees earlier unexecuted batches
type WithdrawalBatch struct {
	BatchNonce sdk.Int    `json:"batch_nonce"`
	ERC20Token ERC20Token `json:"erc_20_token"`
}

func (b WithdrawalBatch) Hash() []byte {
	path := fmt.Sprintf("%s/%s/", b.ERC20Token, b.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// BridgeDeposit is an attestation detail that adds vouchers to an account when executed
type BridgeDeposit struct {
	ERC20Token     ERC20Token      `json:"erc_20_token"`
	EthereumSender EthereumAddress `json:"ethereum_sender" yaml:"ethereum_sender"`
	CosmosReceiver sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
}

func (b BridgeDeposit) Hash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", b.ERC20Token.String(), b.EthereumSender.String(), b.CosmosReceiver.String())
	return tmhash.Sum([]byte(path))
}
