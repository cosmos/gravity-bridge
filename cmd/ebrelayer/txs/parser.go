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

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	ethbridgeTypes "github.com/cosmos/peggy/x/ethbridge/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ETH : ETH constant specifies a token type of Ethereum
const (
	ETH string = "eth"
	// TODO: Get ERC20 token symbol from witnessClaim struct
	ERC string = "erc"
)

// ParsePayload : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func ParsePayload(valAddr sdk.ValAddress, event *events.LockEvent) (ethbridgeTypes.EthBridgeClaim, error) {

	witnessClaim := ethbridgeTypes.EthBridgeClaim{}

	// Nonce type casting (*big.Int -> int)
	nonce := int(event.Nonce.Int64())

	// Sender type casting (address.common -> string)
	sender := ethbridgeTypes.NewEthereumAddress(event.From.Hex())

	// Recipient type casting ([]bytes -> sdk.AccAddress)
	recipient, err := sdk.AccAddressFromBech32(string(event.To))
	if err != nil {
		return witnessClaim, err
	}
	if recipient.Empty() {
		return witnessClaim, errors.New("empty recipient address")
	}

	// TODO: Replace conditional with 'event.TokenSymbol'
	var symbol string
	if event.Token == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		symbol = ETH
	} else {
		symbol = ERC
	}

	// Amount type casting (*big.Int -> sdk.Coins)
	coins := sdk.Coins{sdk.NewInt64Coin(symbol, event.Value.Int64())}

	// Package the information in a unique EthBridgeClaim
	witnessClaim.Nonce = nonce
	witnessClaim.EthereumSender = sender
	witnessClaim.ValidatorAddress = valAddr
	witnessClaim.CosmosReceiver = recipient
	witnessClaim.Amount = coins

	return witnessClaim, nil
}
