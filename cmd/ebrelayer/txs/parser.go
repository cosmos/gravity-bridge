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
  // "log"
  // "encoding/hex"
  // "fmt"
  // "math/big"

  // "github.com/ethereum/go-ethereum/common"
  "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
  sdk "github.com/cosmos/cosmos-sdk/types"
  // "github.com/cosmos/cosmos-sdk/codec"
)

// Witness claim builds a Cosmos transaction
type WitnessClaim struct {
  Nonce          int            `json:"nonce"`
  EthereumSender string         `json:"ethereum_sender"`
  CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
  Validator      sdk.AccAddress `json:"validator"`
  Amount         sdk.Coins      `json:"amount"`
}

func ParsePayloadAndRelay(validator sdk.AccAddress, event *events.LockEvent) string { //cdc *codec.Codec, 
  
  var witnessClaim WitnessClaim

  witnessClaim.EthereumSender = event.From.Hex() // address.common to string

  witnessClaim.Validator = validator

  // witnessClaim.Nonce = (event.Nonce).Int64()

  // recipient, err := sdk.AccAddressFromHex(string(event.To[:]).Hex())
  // if err != nil {
  //   log.Fatal(err)
  // }
  // witnessClaim.CosmosReceiver = recipient

  // // Correct for wei 10**18. Does not currently support erc20.
  // weiAmount, err = sdk.ParseCoins(strings, Join(strconv.Itoa(amount/(Pow(10.0, 18))), "ethereum"))
  // if err3 != nil {
  //     fmt.Errorf("%s", err3)
  // }
  // witnessClaim.Amount = weiAmount
 

  // err := RelayEvent(cdc,
  //                   witnessClaim.CosmosReceiver,
  //                   witnessClaim.Validator,
  //                   witnessClaim.Nonce,
  //                   witnessClaim.EthereumSender,
  //                   witnessClaim.Amount)

  return "No error"
}
