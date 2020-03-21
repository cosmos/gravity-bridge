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
	"io"
	"log"
	"math/big"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/peggy/cmd/ebrelayer/contract"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	"github.com/cosmos/peggy/cmd/ebrelayer/txs"
)

// LoadValidatorCredentials : loads validator's credentials (address, moniker, and passphrase)
func LoadValidatorCredentials(validatorFrom string, inBuf io.Reader) (sdk.ValAddress, string, error) {
	// Get the validator's name and account address using their moniker
	validatorAccAddress, validatorName, err := sdkContext.GetFromFields(inBuf, validatorFrom, false)
	if err != nil {
		return sdk.ValAddress{}, "", err
	}

	// Convert the validator's account address into type ValAddress
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Test keys.DefaultKeyPass is correct-
	_, err = authtxb.MakeSignature(nil, validatorName, keys.DefaultKeyPass, authtxb.StdSignMsg{})
	if err != nil {
		return sdk.ValAddress{}, "", err
	}

	return validatorAddress, validatorName, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(
	appCodec *amino.Codec,
	validatorAddress sdk.ValAddress,
	validatorName string,
	rpcURL string,
	chainID string,
) sdkContext.CLIContext {
	// Create the new CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}

	cliCtx.SkipConfirm = true

	accountRetriever := authtypes.NewAccountRetriever(cliCtx)

	// Ensure that the validator's address exists
	err := accountRetriever.EnsureExists((sdk.AccAddress(validatorAddress)))
	if err != nil {
		log.Fatal(err)
	}

	return cliCtx
}

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

	fmt.Printf("\nStarted Ethereum websocket with provider: %s", provider)

	var targetContract txs.ContractRegistry
	var eventName string

	// TODO: load (targetContract, eventName, contractABI) for both CosmosBridge, BridgeBank
	targetContract = txs.BridgeBank     // TODO: txs.CosmosBridge
	eventName = events.LogLock.String() // TODO: events.LogNewProphecyClaim.String()
	contractABI := contract.LoadABI(false)

	// Load unique event signature from the named event contained within the contract's ABI
	eventSignature := contractABI.Events[eventName].Id().Hex()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get the specific contract's address (CosmosBridge or BridgeBank)
	targetAddress, err := txs.GetAddressFromBridgeRegistry(client, registryContractAddress, targetContract)
	if err != nil {
		log.Fatal(err)
	}

	// We need the target address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{targetAddress},
	}

	// We will check logs for new events
	logs := make(chan types.Log)

	// Filter by contract and event, write results to logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	fmt.Printf("\nSubscribed to %v contract at address: %s\n", targetContract, targetAddress.Hex())

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			log.Fatal(err)
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSignature {
				fmt.Println("\nWitnessed new event:", eventName)
				fmt.Println("Block number:", vLog.BlockNumber)
				fmt.Println("Tx hash:", vLog.TxHash.Hex())

				var err error
				switch eventName {
				case events.LogLock.String():
					err = handleLogLockEvent(
						clientChainID, registryContractAddress, contractABI, eventName,
						vLog, cdc, validatorAddress, validatorName, cliCtx, txBldr,
					)
				case events.LogNewProphecyClaim.String():
					err = handleLogNewProphecyClaimEvent(
						contractABI, eventName, vLog, provider, registryContractAddress, privateKey,
					)
				}

				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// handleLogLockEvent unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func handleLogLockEvent(
	clientChainID *big.Int,
	contractAddress common.Address,
	contractABI abi.ABI,
	eventName string,
	log types.Log,
	cdc *codec.Codec,
	validatorAddress sdk.ValAddress,
	validatorName string,
	cliCtx sdkContext.CLIContext,
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
	return txs.RelayLockToCosmos(cdc, validatorName, &prophecyClaim, cliCtx, txBldr)
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
