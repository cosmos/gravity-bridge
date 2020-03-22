package relayer

import (
	"context"
	"crypto/ecdsa"
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
	"github.com/ethereum/go-ethereum/ethclient"
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
	validatorAddress := sdk.ValAddress(validatorAccAddress)

	// Confirm that the key is valid
	_, err = authtxb.MakeSignature(nil, validatorName, keys.DefaultKeyPass, authtxb.StdSignMsg{})
	if err != nil {
		return sdk.ValAddress{}, "", err
	}

	return validatorAddress, validatorName, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(appCodec *amino.Codec, validatorAddress sdk.ValAddress, validatorName string,
	rpcURL string, chainID string) sdkContext.CLIContext {
	// Create the new CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	// Confirm that the validator's address exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists((sdk.AccAddress(validatorAddress)))
	if err != nil {
		log.Fatal(err)
	}
	return cliCtx
}

// InitEthereumRelayer : Subscribes to events emitted by the deployed contracts
func InitEthereumRelayer(cdc *codec.Codec, provider string, registryContractAddress common.Address,
	validatorName string, validatorAddress sdk.ValAddress, cliCtx sdkContext.CLIContext,
	txBldr authtypes.TxBuilder, privateKey *ecdsa.PrivateKey) error {
	client, err := SetupWebsocketEthClient(provider)
	if err != nil {
		return err
	}
	log.Println("Started Ethereum websocket with provider:", provider)

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// We will check logs for new events
	logs := make(chan types.Log)

	// Start BridgeBank subscription, prepare contract ABI and LockLog event signature
	bridgeBankAddress, subBridgeBank := startContractEventSub(logs, client, registryContractAddress, txs.BridgeBank)
	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)
	eventLogLockSignature := bridgeBankContractABI.Events[events.LogLock.String()].Id().Hex()

	// Start CosmosBridge subscription, prepare contract ABI and LogNewProphecyClaim event signature
	cosmosBridgeAddress, subCosmosBridge := startContractEventSub(logs, client, registryContractAddress, txs.CosmosBridge)
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
			log.Printf("Witnessed tx %s on block %d\n", vLog.TxHash.Hex(), vLog.BlockNumber)
			var err error
			switch vLog.Topics[0].Hex() {
			case eventLogLockSignature:
				err = handleLogLockEvent(clientChainID, bridgeBankAddress, bridgeBankContractABI,
					events.LogLock.String(), vLog, cdc, validatorAddress, validatorName, cliCtx, txBldr)
			case eventLogNewProphecyClaimSignature:
				err = handleLogNewProphecyClaimEvent(
					cosmosBridgeContractABI, events.LogNewProphecyClaim.String(), vLog, provider,
					cosmosBridgeAddress, privateKey)
			}
			// TODO: Should this be a Fatal err?
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// startContractEventSub : starts an event subscription on the specified Peggy contract
func startContractEventSub(logs chan types.Log, client *ethclient.Client, registryAddress common.Address,
	contractName txs.ContractRegistry) (common.Address, ethereum.Subscription) {
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
	log.Printf("Subscribed to %v contract at address: %s\n", contractName, subContractAddress.Hex())
	return subContractAddress, sub
}

// handleLogLockEvent unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func handleLogLockEvent(clientChainID *big.Int, contractAddress common.Address, contractABI abi.ABI,
	eventName string, log types.Log, cdc *codec.Codec, validatorAddress sdk.ValAddress, validatorName string,
	cliCtx sdkContext.CLIContext, txBldr authtypes.TxBuilder) error {
	event := events.UnpackLogLock(clientChainID, contractAddress.Hex(), contractABI, eventName, log.Data)

	// Add the event to the record
	events.NewEventWrite(log.TxHash.Hex(), event) // TODO: is this at all useful?

	prophecyClaim, err := txs.LogLockToEthBridgeClaim(validatorAddress, &event)
	if err != nil {
		return err
	}
	return txs.RelayLockToCosmos(cdc, validatorName, &prophecyClaim, cliCtx, txBldr)
}

// Unpacks a LogNewProphecyClaim event, builds a new OracleClaim, and relays it to Ethereum
func handleLogNewProphecyClaimEvent(contractABI abi.ABI, eventName string, log types.Log,
	provider string, contractAddress common.Address, privateKey *ecdsa.PrivateKey) error {
	event := events.UnpackLogNewProphecyClaim(contractABI, eventName, log.Data)

	oracleClaim, err := txs.ProphecyClaimToSignedOracleClaim(event, privateKey)
	if err != nil {
		return err
	}
	return txs.RelayOracleClaimToEthereum(provider, contractAddress, events.LogNewProphecyClaim, oracleClaim, privateKey)
}
