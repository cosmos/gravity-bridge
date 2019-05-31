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
  "log"
  "fmt"

  "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/cosmos/cosmos-sdk/codec"
)

type WitnessClaim struct {
  Nonce          int            `json:"nonce"`
  EthereumSender string         `json:"ethereum_sender"`
  CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
  Validator      sdk.AccAddress `json:"validator"`
  Amount         sdk.Coins      `json:"amount"`
}

func ParsePayloadAndRelay(cdc *codec.Codec, validator sdk.AccAddress, event *events.Event) {
  // Set the witnessClaim's validator
  var witnessClaim WitnessClaim
  witnessClaim.Validator = validator

  // Get the keyset of the payload's fields
  // payloadKeySet := event.EventPayload()["_id"].([]sdk.AccAddress)
  payloadKeySet := event.EventPayload()["Keys"].([]string)

  // Parse each key field individually
  for _, field := range payloadKeySet {
      switch(field) {
          case "_id":
              // Print the unique id of the event.
              fmt.Print(field);
          case "_from":
              ethereumSender, ok := field.Address();
              if !ok {
                  return eventPayload, errors.New("Error while parsing transaction's ethereum sender");
              }
              witnessClaim.EthereumSender = ethereumSender;
          case "_to":
              cosmosReceiver, ok := field.Bytes32();

              bech32CosmosReceiver, err2 := sdk.AccAddressFromBech32(cosmosReceiver)
              if err2 != nil {
                  fmt.Errorf("%s", err2)
              }

              if !ok {
                  return eventPayload, errors.New("Error while parsing transaction's Cosmos recipient");
              }
              witnessClaim.CosmosReceiver = bech32CosmosReceiver;
          case "_token":
              tokenType, ok := field.Bytes32();
              if !ok {
                  return eventPayload, errors.New("Error while parsing the token type");
              }
              witnessClaim.Token = tokenType;               
          case "_value":
              amount, ok := field.BigInt()
              if !ok {
                  return eventPayload, errors.New("Error while parsing transaction's value")
              }
              // Correct for wei 10**18. Does not currently support erc20.
              weiAmount, err = sdk.ParseCoins(strings, Join(strconv.Itoa(amount/(Pow(10.0, 18))), "ethereum"))
              if err3 != nil {
                  fmt.Errorf("%s", err3)
              }
              witnessClaim.Amount = weiAmount
          case "_nonce":
              nonce, ok := field.BigInt()
              if !ok {
                  return eventPayload, errors.New("Error while parsing transaction's nonce")
              }
              witnessClaim.Nonce = nonce
          }
  }

  err := RelayEvent(cdc,
        witnessClaim.CosmosReceiver,
        witnessClaim.Validator,
        witnessClaim.Nonce,
        witnessClaim.EthereumSender,
        witnessClaim.Amount)

  if err != nil {
    log.Fatal(err)
  }
}
