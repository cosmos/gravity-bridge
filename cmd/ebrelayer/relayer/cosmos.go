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
func InitCosmosRelayer(tendermintProvider string, web3Provider string, cosmosBridgeContractAddress common.Address, rawPrivateKey string) error {

	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))

	client := tmclient.NewHTTP(tendermintProvider, "/websocket")
	client.SetLogger(logger)
	err := client.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer client.Stop()

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
			tx := result.Data.(tmtypes.EventDataTx)
			logger.Info("\t New transaction witnessed")

			txRes := tx.Result
			for i := 1; i < len(txRes.Events); i++ {
				event := txRes.Events[i]
				eventName := string(txRes.Events[i].Type)
				switch eventName {
				case "burn", "lock":
					// Package the data into an array for proper parsing
					cosmosSender := string(event.Attributes[0].Value)
					ethereumReceiver := string(event.Attributes[1].Value)
					coin := string(event.Attributes[3].Value)
					eventData := [3]string{cosmosSender, ethereumReceiver, coin}

					// Parse the eventData into a new MsgEvent
					msgEvent := events.NewMsgEvent(eventName, eventData)

					// Relay the MsgEvent to the Ethereum network
					err = txs.RelayToEthereum(web3Provider, cosmosBridgeContractAddress, rawPrivateKey, &msgEvent)
					if err != nil {
						return err
					}
				default:
					// do nothing
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
