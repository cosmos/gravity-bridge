package txs

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// TODO: Refactor solc generation of Peggy.go for compatibility with Truffle
	peggy "github.com/cosmos/peggy/cmd/ebrelayer/peggy"
)

// TODO: Remove these test constants once event data is passed in params
const (
	// CosmosSender : hashed address "cosmos1gn8409qq9hnrxde37kuxwx5hrxpfpv8426szuv"
	CosmosSender = "0x636F736D6F7331676E38343039717139686E7278646533376B75787778356872787066707638343236737A7576"
	// EthereumRecipient : intended recipient of token transfer
	EthereumRecipient = "0x115f6e2004d7b4ccd6b9d5ab34e30909e0f612cd"
	// WeiAmount : transfer amount in wei
	WeiAmount = 100
	// ItemID : unique identifier of the tokens stored on the bridge contract
	ItemID = "2064e17083eed31b4a77fc929bedb9e97a1508b484b3a60185413ba58fd36b6d"
)

// RelayToEthereum : relays the provided transaction data to a peggy smart contract deployed on Ethereum
func RelayToEthereum(provider string, peggyContractAddress common.Address, rawPrivateKey string, eventName string, eventData []string) error {

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

	// Initialize Peggy contract instance
	// TODO: Update this to 'CosmosBridge'
	instance, err := peggy.NewPeggy(peggyContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Tokens get locked on Cosmos ->
	// 	Tokens get minted on Ethereum ->
	// 		Validators reach consensus about the state ->
	// 			Tokens get burned on Cosmos ->
	// 				Tokens get unlocked on Ethereum

	// Send transaction to the instance's specified method
	switch eventName {
	case "burn":
		// Parse Cosmos sender
		cosmosSender := [32]byte{}
		copy(cosmosSender[:], []byte(eventData[0]))

		// Parse Ethereum receiver
		if !common.IsHexAddress(eventData[1]) {
			return fmt.Errorf("Invalid recipient address: %v", eventData[1])
		}
		ethereumReceiver := common.HexToAddress(eventData[1])

		// TODO: Parse symbol, amount from sdk.Coin coin
		// coin := eventData[2]
		symbol := "eth"
		amount := 3

		// TODO: Get token address from chain
		tokenAddressString := "0x0000000000000000000000000000000000000000"
		// Parse Ethereum receiver
		if !common.IsHexAddress(tokenAddressString) {
			return fmt.Errorf("Invalid token address: %v", tokenAddressString)
		}
		tokenAddress := common.HexToAddress(tokenAddressString)

		// TODO: REMOVE THIS PRINT
		noncePrint := nonce
		cosmosSenderPrint := hex.EncodeToString(cosmosSender[:])
		ethereumReceiverPrint := ethereumReceiver.Hex()
		tokenAddressPrint := tokenAddress.Hex()
		symbolPrint := symbol
		amountPrint := amount
		// id := hex.EncodeToString(event.Id[:])

		fmt.Printf("\nNonce: %v\nCosmos Sender: %v\nEthereum Recipient: %v\nToken Address: %v\nSymbol: %v\nAmount: %v\n\n",
			noncePrint, cosmosSenderPrint, ethereumReceiverPrint, tokenAddressPrint, symbolPrint, amountPrint)

		// TODO: Remove nonce from function on contract
		tx, err := instance.newProphecyClaim(auth, nonce, cosmosSender, ethereumReceiver, tokenAddress, symbol, amount)
		if err != nil {
			log.Fatal(err)
		}

		// TODO: Deal with this
		// case "create_claim":
		// case "create_prophecy":

		// tx, err := instance.Unlock(auth, itemID)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	// Get the transaction receipt
	// receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("\nTx relayed to Ethereum")
	// fmt.Println("Tx hash:", tx.Hash().Hex())
	// fmt.Println("Status:", receipt.Status)

	return nil
}
