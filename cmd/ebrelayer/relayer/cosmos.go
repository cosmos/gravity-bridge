package relayer

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitCosmosRelayer : initializes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(tendermintProvider string, web3Provider string, peggyContractAddress common.Address, rawPrivateKey string) error {

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
			// txIndex := result.Data.(tmtypes.EventDataTx).Index
			logger.Info("\t New transaction witnessed")

			// fmt.Println(tx)

			runes := "[99] [111] [115] [109] [111] [115]" // 95 115 101 110 100 101 114]
			outString := ""
			for _, v := range runes {
				outString += string(v)
			}
			fmt.Println(outString)

			txRes := tx.Result
			for i := 1; i < len(txRes.Events); i++ {
				switch txRes.Events[i].Type {
				case "burn":
					logger.Info("\tMsgBurn")
					fmt.Printf("%v", txRes.Events[i])
					// TODO: Parse event attributes and pass them to txs.relayToEthereum
					err = txs.RelayToEthereum(web3Provider, peggyContractAddress, rawPrivateKey)
					if err != nil {
						return err
					}
				case "create_claim":
					logger.Info("\tMsgCreateClaim")
					err = txs.RelayToEthereum(web3Provider, peggyContractAddress, rawPrivateKey)
					if err != nil {
						return err
					}
				case "create_prophecy":
					logger.Info("\tMsgCreateProphecy")
					err = txs.RelayToEthereum(web3Provider, peggyContractAddress, rawPrivateKey)
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
