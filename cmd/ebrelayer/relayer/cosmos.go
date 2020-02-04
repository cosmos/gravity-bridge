package relayer

import (
	"context"
	"crypto/ecdsa"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"

	tmKv "github.com/tendermint/tendermint/libs/kv"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitCosmosRelayer initializes a relayer which witnesses events on the Cosmos network and relays them to Ethereum
func InitCosmosRelayer(
	tendermintProvider string,
	web3Provider string,
	contractAddress common.Address,
	key *ecdsa.PrivateKey,
) error {
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

			// Iterate over each event inside of the transaction
			for _, event := range tx.Result.Events {
				// Get type of OracleClaim based on the event's type
				claimType := getOracleClaimType(event.GetType())

				switch claimType {
				case events.MsgBurn, events.MsgLock:
					// Parse event data, then package it as a ProphecyClaim and relay to the Ethereum Network
					err := handleBurnLockMsg(event.GetAttributes(), claimType, web3Provider, contractAddress, key)
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

// getOracleClaimType sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) events.Event {
	var claimType events.Event

	switch eventType {
	case events.MsgBurn.String():
		claimType = events.MsgBurn
	case events.MsgLock.String():
		claimType = events.MsgLock
	default:
		claimType = events.Unsupported
	}

	return claimType
}

// handleBurnLockMsg parse event data as a CosmosMsg,
// package it into a ProphecyClaim, then relay tx to the Ethereum Network
func handleBurnLockMsg(
	attributes []tmKv.Pair,
	claimType events.Event,
	web3Provider string,
	contractAddress common.Address,
	key *ecdsa.PrivateKey,
) error {
	// Parse the witnessed event's data into a new CosmosMsg
	cosmosMsg := txs.BurnLockEventToCosmosMsg(claimType, attributes)

	// Parse the CosmosMsg into a ProphecyClaim for relay to Ethereum
	prophecyClaim := txs.CosmosMsgToProphecyClaim(cosmosMsg)

	// TODO: Need some sort of delay on this so validators aren't all submitting at the same time
	// Relay the CosmosMsg to the Ethereum network
	err := txs.RelayProphecyClaimToEthereum(web3Provider, contractAddress, claimType, prophecyClaim, key)
	if err != nil {
		return err
	}

	return nil
}
