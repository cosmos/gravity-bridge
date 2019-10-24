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

	// TODO: Add a contract abigen generator via solc to testnet-contracts/scripts
	cosmosBridge "github.com/cosmos/peggy/cmd/ebrelayer/cosmosbridge"
	"github.com/cosmos/peggy/cmd/ebrelayer/events"
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

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Set up tx signature authorization
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // 300,000 Gwei in units
	auth.GasPrice = gasPrice

	// // Initialize BridgeRegistry instance
	// instance, err := cosmosBridge.NewBridgeRegistry(cosmosBridgeContractAddress, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// TODO: Get cosmosBridgeAddress from BridgeRegistry
	// Initialize CosmosBridge contract instance
	instance, err := cosmosBridge.NewCosmosBridge(cosmosBridgeContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	var txHash common.Hash
	switch eventData.EventName {
	case "burn":
		// TODO: Delete nonce (2nd param) once it is removed from NewProphecyClaim() params
		tx, err := instance.NewProphecyClaim(auth, big.NewInt(2), eventData.CosmosSender, eventData.EthereumReceiver, eventData.TokenContractAddress, eventData.Symbol, eventData.Amount)
		if err != nil {
			log.Fatal(err)
		}
		txHash = tx.Hash()
		fmt.Println("\nNewProphecyClaim tx:", txHash.Hex())

	case "lock":
		// TODO: Integrate with MsgLock feature
	}

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	switch receipt.Status {
	case 0:
		fmt.Println("Status: 0 - Failed")
	case 1:
		fmt.Println("Status: 1 - Successful")
	}

	return nil
}
