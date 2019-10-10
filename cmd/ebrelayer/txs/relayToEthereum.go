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

// relayToEthereum : relays the provided transaction data to a peggy smart contract deployed on Ethereum
func relayToEthereum(provider string, peggyContractAddress common.Address, rawPrivateKey string) error {

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
	instance, err := peggy.NewPeggy(peggyContractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Event parameters
	itemID := [32]byte{}
	copy(itemID[:], []byte(ItemID))

	// TODO: Refactor unlock() on smart contract to accept (cosmosSender, ethereumRecipient, amount) instead of itemID
	// cosmosSender := []byte{}
	// copy(cosmosSender[:], []byte(CosmosSender))

	// ethereumRecipientString := EthereumTokenAddress
	// if !common.IsHexAddress(ethereumRecipientString) {
	// 	return fmt.Errorf("Invalid contract-address: %v", ethereumRecipientString)
	// }
	// ethereumRecipient := common.HexToAddress(ethereumRecipientString)

	// amount := big.NewInt(WeiAmount)

	// Send transaction to the instance's specified method
	tx, err := instance.Unlock(auth, itemID) // cosmosSender, ethereumReceiver, amount
	if err != nil {
		log.Fatal(err)
	}

	// Get the transaction receipt
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tx hash:", tx.Hash().Hex())
	fmt.Println("Status:", receipt.Status, "\n")

	return nil
}
