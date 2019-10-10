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
)

// InitCosmosRelayer : initializes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(provider string, peggyContractAddress common.Address, rawPrivateKey string) error {

	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))

	// TODO: Parameterize tmclient provider
	client := tmclient.NewHTTP("tcp://localhost:26657", "/websocket")
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
			logger.Info("Tx witnessed:")

			txEvents := tx.Result.Events
			for i := 1; i < len(txEvents); i++ {
				switch txEvents[i].Type {
				case "burn":
					logger.Info("\tMsgBurn")
					// TODO: Parse event attributes and pass them to txs.relayToEthereum
					// err = txs.relayToEthereum(provider, peggyContractAddress, rawPrivateKey)
					// if err != nil {
					// 	return err
					// }
				case "create_claim":
					logger.Info("\tMsgCreateClaim")
				case "create_prophecy":
					logger.Info("\tMsgCreateProphecy")
				default:
					logger.Info("")
					// do nothing
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}
