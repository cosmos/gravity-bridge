package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	bridgeBank "github.com/cosmos/peggy/cmd/ebrelayer/generated/bridgebank"
	bridgeRegistry "github.com/cosmos/peggy/cmd/ebrelayer/generated/bridgeregistry"
	cosmosBridge "github.com/cosmos/peggy/cmd/ebrelayer/generated/cosmosbridge"
)

const (
	// GasLimit : the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit = uint64(300000)
)

// RelayToEthereum : relays the provided transaction data to a smart contract deployed on Ethereum
func RelayToEthereum(provider string, cosmosBridgeContractAddress common.Address, rawPrivateKey string, eventData *events.MsgEvent) error {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(rawPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Parse public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Set up CallOpts auth
	callOptsAuth := bind.CallOpts{
		Pending:     true,
		From:        fromAddress,
		BlockNumber: header.Number,
		Context:     context.Background(),
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(privateKey)
	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	// Initialize BridgeRegistry instance
	bridgeRegistryInstance, err := bridgeRegistry.NewBridgeRegistry(cosmosBridgeContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Get the specific contract's address (Valset, Oracle, CosmosBridge, or BridgeBank) based on EventName
	address, err := getAddressFromBridgeRegistry(bridgeRegistryInstance, &callOptsAuth, eventData.EventName)
	if err != nil {
		log.Fatal(err)
	}

	var txHash common.Hash

	// Relay tx to appropriate contract depending on the event type
	switch eventData.EventName {
	case events.Burn:
		fmt.Println("\nFetching CosmosBridge contract...")
		// Initialize CosmosBridge instance
		cosmosBridgeInstance, err := cosmosBridge.NewCosmosBridge(address, client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sending tx to CosmosBridge...")
		tx, err := cosmosBridgeInstance.NewProphecyClaim(transactOptsAuth, 0, eventData.CosmosSender, eventData.EthereumReceiver, eventData.TokenContractAddress, eventData.Symbol, eventData.Amount)
		if err != nil {
			log.Fatal(err)
		}

		// Set tx hash
		txHash = tx.Hash()

		fmt.Println("\nNewProphecyClaim tx hash:", txHash.Hex())
	case events.Lock:
		fmt.Println("\nFetching BridgeBank contract...")
		// Initialize BridgeBank instance
		bridgeBankInstance, err := bridgeBank.NewBridgeBank(address, client)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Send lock related transaction to appropriate contract method
		fmt.Println(bridgeBankInstance)
	}

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	// Report tx status
	switch receipt.Status {
	case 0:
		fmt.Println("Status: 0 - Failed")
	case 1:
		fmt.Println("Status: 1 - Successful")
	}

	return nil
}

func getAddressFromBridgeRegistry(instance *bridgeRegistry.BridgeRegistry, auth *bind.CallOpts, eventName events.EventType) (common.Address, error) {

	var contractAddress common.Address

	switch eventName {
	case events.Burn:
		cosmosBridgeAddress, err := instance.CosmosBridge(auth)
		if err != nil {
			log.Fatal(err)
		}
		contractAddress = cosmosBridgeAddress
	case events.Lock:
		bridgeBankAddress, err := instance.BridgeBank(auth)
		if err != nil {
			log.Fatal(err)
		}
		contractAddress = bridgeBankAddress
	}

	return contractAddress, nil
}
