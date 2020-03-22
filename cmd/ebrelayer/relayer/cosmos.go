package relayer

import (
	"context"
	"crypto/ecdsa"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
	"github.com/cosmos/peggy/cmd/ebrelayer/types"
)

// InitCosmosRelayer initializes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(tendermintProvider string, web3Provider string, contractAddress common.Address,
	key *ecdsa.PrivateKey) error {
	// TODO: move logger to main
	logger := tmLog.NewTMLogger(tmLog.NewSyncWriter(os.Stdout))
	client, err := tmClient.NewHTTP(tendermintProvider, "/websocket")
	if err != nil {
		return err
	}

	client.SetLogger(logger)

	if err := client.Start(); err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}

	defer client.Stop() //nolint:errcheck

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
			tx, ok := result.Data.(tmTypes.EventDataTx)
			if !ok {
				logger.Error("Type casting failed while extracting event data from new tx")
			}
			logger.Info("New transaction witnessed")

			// Iterate over each event in the transaction
			for _, event := range tx.Result.Events {
				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case types.MsgBurn, types.MsgLock:
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					err := handleBurnLockMsg(event.GetAttributes(), claimType, web3Provider, contractAddress, key)
					if err != nil {
						return err
					}
				case types.Unsupported:
				}
			}
		case <-quit:
			os.Exit(0)
		}
	}
}

// getOracleClaimType sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) types.Event {
	var claimType types.Event
	switch eventType {
	case types.MsgBurn.String():
		claimType = types.MsgBurn
	case types.MsgLock.String():
		claimType = types.MsgLock
	default:
		claimType = types.Unsupported
	}
	return claimType
}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func handleBurnLockMsg(attributes []tmKv.Pair, claimType types.Event, web3Provider string,
	contractAddress common.Address, key *ecdsa.PrivateKey) error {
	cosmosMsg := txs.BurnLockEventToCosmosMsg(claimType, attributes)
	log.Println(cosmosMsg.String()) // TODO: use logger here

	// TODO: Ideally one validator should relay the prophecy and other validators make oracle claims upon that prophecy
	prophecyClaim := txs.CosmosMsgToProphecyClaim(cosmosMsg)
	err := txs.RelayProphecyClaimToEthereum(web3Provider, contractAddress, claimType, prophecyClaim, key)
	if err != nil {
		return err
	}
	return nil
}
