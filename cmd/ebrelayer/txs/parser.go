package txs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      Cosmos Bridge.
// --------------------------------------------------------

import (
	"errors"
	// "math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/peggy_fork/cmd/ebrelayer/events"
	"github.com/swishlabsco/peggy_fork/x/ethbridge/types"
)

const (
	ETH string = "eth"
)

func ParsePayload(valAddr sdk.ValAddress, event *events.LockEvent) (types.EthBridgeClaim, error) {

	witnessClaim := types.EthBridgeClaim{}

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Uint64())
	witnessClaim.Nonce = nonce

	// EthereumSender type casting (address.common -> string)
	// witnessClaim.EthereumSender = event.From // .Hex()

	// CosmosReceiver type casting (bytes[] -> sdk.AccAddress)
	recipient := sdk.AccAddress(event.To)
	if recipient.Empty() {
		return witnessClaim, errors.New("Invalid recipient address")
	}
	witnessClaim.CosmosReceiver = recipient

	// valAddr is correct type (sdk.ValAddress)
	witnessClaim.ValidatorAddress = valAddr

	// Amount type casting (*big.Int -> sdk.Coins)
	// weiAmount := sdk.NewCoins(sdk.NewInt64Coin(ETH, event.Value.Int64()))
	// witnessClaim.Amount = weiAmount

	return witnessClaim, nil
}
