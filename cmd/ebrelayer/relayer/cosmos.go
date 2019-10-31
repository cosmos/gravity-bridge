package relayer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitCosmosRelayer : initializes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(tendermintProvider string, web3Provider string, bridgeContractAddress common.Address, rawPrivateKey string) error {
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))
	client := tmclient.NewHTTP(tendermintProvider, "/websocket")

	client.SetLogger(logger)

	err := client.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}

	defer client.Stop()

	// Subscribe to all tendermint transactions
	query := "tm.event = 'Tx'"

	out, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-out:
			tx, ok := result.Data.(tmtypes.EventDataTx)
			if !ok {
				logger.Error("Type casting failed while extracting event data from new tx")
			}

			logger.Info("\t New transaction witnessed")

			txRes := tx.Result
			for i := 1; i < len(txRes.Events); i++ {
				event := txRes.Events[i]

				var claimType events.Event

				// Parse event type from event name
				if string(event.Type) == "burn" {
					claimType = events.MsgBurn
				} else if string(event.Type) == "lock" {
					claimType = events.MsgLock
				} else {
					claimType = events.Unsupported
				}

				switch claimType {
				case events.MsgBurn, events.MsgLock:
					// Package the data into an array for proper parsing
					cosmosSender := string(event.Attributes[0].Value)
					ethereumReceiver := string(event.Attributes[1].Value)
					coin := string(event.Attributes[3].Value)
					eventData := [3]string{cosmosSender, ethereumReceiver, coin}

					// Parse the eventData into a new CosmosMsg
					cosmosMsg := events.NewCosmosMsg(claimType, eventData)

					// TODO: Data mapping
					prophecyClaim := txs.CosmosMsgToProphecyClaim(cosmosMsg)

					// TODO: Need some sort of delay on this so validators aren't all submitting at the same time
					// Relay the CosmosMsg to the Ethereum network
					err = txs.RelayProphecyClaimToEthereum(web3Provider, bridgeContractAddress, claimType, prophecyClaim)
					if err != nil {
						return err
					}
				case events.Unsupported:
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
