package relayer

// -----------------------------------------------------
//      Ethereum relayer
//
//      Initializes the relayer service, which parses,
//      encodes, and packages named events on an Ethereum
//      Smart Contract for validator's to sign and send
//      to the Cosmos bridge.
// -----------------------------------------------------

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/cosmos/peggy/cmd/ebrelayer/contract"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// InitEthereumRelayer : Subscribes to events emitted by the deployed contracts
func InitEthereumRelayer(
	cdc *codec.Codec,
	provider string,
	registryContractAddress common.Address,
	validatorName string,
	validatorAddress sdk.ValAddress,
	cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
	privateKey *ecdsa.PrivateKey,
) error {
	// Start client with infura ropsten provider
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		return err
	}

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// We will check logs for new events
	logs := make(chan types.Log)

	// Start BridgeBank contract subscription
	bridgeBankAddress, subBridgeBank := startContractEventSub(logs, client, registryContractAddress, txs.BridgeBank)

	// Start CosmosBridge contract subscription
	cosmosBridgeAddress, subCosmosBridge := startContractEventSub(logs, client, registryContractAddress, txs.CosmosBridge)

	// Load BridgeBank contract ABI and LogLock event signature
	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)
	eventLogLockSignature := bridgeBankContractABI.Events[events.LogLock.String()].Id().Hex()

	// Load CosmosBridge contract ABI and LogNewProphecyClaim event signature
	cosmosBridgeContractABI := contract.LoadABI(txs.CosmosBridge)
	eventLogNewProphecyClaimSignature := cosmosBridgeContractABI.Events[events.LogNewProphecyClaim.String()].Id().Hex()

	for {
		select {
		// Handle any errors
		case err := <-subBridgeBank.Err():
			log.Fatal(err)
		case err := <-subCosmosBridge.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			fmt.Println("\nWitnessed new Tx...")
			fmt.Println("Block number:", vLog.BlockNumber)
			fmt.Println("Tx hash:", vLog.TxHash.Hex())

			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogLockSignature:
				err = handleLogLockEvent(
					clientChainID, bridgeBankAddress, bridgeBankContractABI, events.LogLock.String(),
					vLog, cdc, validatorAddress, validatorName, cliCtx, txBldr,
				)
			case eventLogNewProphecyClaimSignature:
				err = handleLogNewProphecyClaimEvent(
					cosmosBridgeContractABI, events.LogNewProphecyClaim.String(), vLog, provider, cosmosBridgeAddress, privateKey,
				)
			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// startContractEventSub : starts an event subscription on the specified Peggy contract
func startContractEventSub(
	logs chan types.Log,
	client *ethclient.Client,
	registryAddress common.Address,
	contractName txs.ContractRegistry,
) (common.Address, ethereum.Subscription) {
	// Get the contract address for this subscription
	subContractAddress, err := txs.GetAddressFromBridgeRegistry(client, registryAddress, contractName)
	if err != nil {
		log.Fatal(err)
	}

	// We need the address in []bytes for the query
	subQuery := ethereum.FilterQuery{
		Addresses: []common.Address{subContractAddress},
	}

	// Start the contract subscription
	sub, err := client.SubscribeFilterLogs(context.Background(), subQuery, logs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nSubscribed to %v contract at address: %s\n", contractName, subContractAddress.Hex())

	return subContractAddress, sub
}

// handleLogLockEvent : unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func handleLogLockEvent(
	clientChainID *big.Int,
	contractAddress common.Address,
	contractABI abi.ABI,
	eventName string,
	log types.Log,
	cdc *codec.Codec,
	validatorAddress sdk.ValAddress,
	validatorName string,
	cliContext sdkContext.CLIContext,
	txBldr authtypes.TxBuilder,
) error {
	// Unpack the LogLock event using its unique event signature from the contract's ABI
	event := events.UnpackLogLock(clientChainID, contractAddress.Hex(), contractABI, eventName, log.Data)

	// Add the event to the record
	events.NewEventWrite(log.TxHash.Hex(), event)

	// Parse the LogLock event's payload into a struct
	prophecyClaim, err := txs.LogLockToEthBridgeClaim(validatorAddress, &event)
	if err != nil {
		return err
	}

	// Initiate the relay
	return txs.RelayLockToCosmos(
		cdc, validatorName, &prophecyClaim, cliContext, txBldr,
	)
}

// handleLogNewProphecyClaimEvent unpacks a LogNewProphecyClaim event,
// converts it to a OracleClaim, and relays a tx to Ethereum
func handleLogNewProphecyClaimEvent(
	contractABI abi.ABI,
	eventName string,
	log types.Log,
	provider string,
	contractAddress common.Address,
	privateKey *ecdsa.PrivateKey,
) error {
	// Unpack the LogNewProphecyClaim event using its unique event signature from the contract's ABI
	event := events.UnpackLogNewProphecyClaim(contractABI, eventName, log.Data)

	// Parse ProphecyClaim's data into an OracleClaim
	oracleClaim, err := txs.ProphecyClaimToSignedOracleClaim(event, privateKey)
	if err != nil {
		return err
	}

	// Initiate the relay
	return txs.RelayOracleClaimToEthereum(provider, contractAddress, events.LogNewProphecyClaim, oracleClaim, privateKey)
}
