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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
)

const (
	ETH string = "eth"
)

func ParsePayload(validator sdk.AccAddress, event *events.LockEvent) (types.EthBridgeClaim, error) {

	witnessClaim := types.EthBridgeClaim{}

	// Nonce type casting (*big.Int -> int)
	nonce := Int64(event.Nonce)
	witnessClaim.Nonce = nonce

	// EthereumSender type casting (address.common -> string)
	witnessClaim.EthereumSender = event.From.Hex()

	// CosmosReceiver type casting (bytes[] -> sdk.AccAddress)
	recipient := sdk.AccAddress(event.To)
	if recipient.Empty() {
		return
	}
	witnessClaim.CosmosReceiver = recipient

	// Validator is already the correct type (sdk.AccAddress)
	witnessClaim.ValidatorAddress = validator

	// Amount type casting (*big.Int -> sdk.Coins)
	weiAmount := sdk.NewCoins(sdk.NewInt64Coin(ETH, event.Value.Int64))
	witnessClaim.Amount = weiAmount

	return witnessClaim, nil
}
