package txs

import (
	"context"
	"fmt"

	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/cosmos/peggy/cmd/ebrelayer/events"
	cosmosBridge "github.com/cosmos/peggy/cmd/ebrelayer/generated/cosmosbridge"
	oracle "github.com/cosmos/peggy/cmd/ebrelayer/generated/oracle"
)

const (
	// GasLimit : the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit = uint64(600000)
)

// RelayProphecyClaimToEthereum : relays the provided ProphecyClaim to CosmosBridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(provider string, contractAddress common.Address, event events.Event, claim ProphecyClaim) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := initRelayConfig(provider, contractAddress, event)

	// Initialize CosmosBridge instance
	fmt.Println("\nFetching CosmosBridge contract...")
	cosmosBridgeInstance, err := cosmosBridge.NewCosmosBridge(target, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	fmt.Println("Sending new ProphecyClaim to CosmosBridge...")
	tx, err := cosmosBridgeInstance.NewProphecyClaim(auth, uint8(claim.ClaimType), claim.CosmosSender, claim.EthereumReceiver, claim.TokenContractAddress, claim.Symbol, claim.Amount)
	if err != nil {
		log.Fatal(err)
	}

	// Get the transaction receipt
	fmt.Println("NewProphecyClaim tx hash:", tx.Hash().Hex())
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// Report tx status
	switch receipt.Status {
	case 0:
		fmt.Println("Tx Status: 0 - Failed")
	case 1:
		fmt.Println("Tx Status: 1 - Successful")
	}

	return nil
}

// RelayOracleClaimToEthereum : relays the provided OracleClaim to Oracle contract on the Ethereum network
func RelayOracleClaimToEthereum(provider string, contractAddress common.Address, event events.Event, claim OracleClaim) error {
	// Initialize client service, validator's tx auth, and target contract address
	client, auth, target := initRelayConfig(provider, contractAddress, event)

	// Initialize Oracle instance
	fmt.Println("\nFetching Oracle contract...")
	oracleInstance, err := oracle.NewOracle(target, client)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	fmt.Println("Sending new OracleClaim to Oracle...")
	tx, err := oracleInstance.NewOracleClaim(auth, claim.ProphecyID, claim.Message, claim.Signature)
	if err != nil {
		log.Fatal(err)
	}

	// Get the transaction receipt
	fmt.Println("NewOracleClaim tx hash:", tx.Hash().Hex())
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// Report tx status
	switch receipt.Status {
	case 0:
		fmt.Println("Tx Status: 0 - Failed")
	case 1:
		fmt.Println("Tx Status: 1 - Successful")
	}

	return nil
}

// initRelayConfig : set up Ethereum client, validator's transaction auth, and the target contract's address
func initRelayConfig(
	provider string,
	registry common.Address,
	event events.Event,
) (*ethclient.Client,
	*bind.TransactOpts,
	common.Address,
) {
	// Start Ethereum client
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	// Load the validator's private key
	key, err := LoadPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Load the validator's address
	sender, err := LoadSender()
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Set up TransactOpts auth's tx signature authorization
	transactOptsAuth := bind.NewKeyedTransactor(key)
	transactOptsAuth.Nonce = big.NewInt(int64(nonce))
	transactOptsAuth.Value = big.NewInt(0) // in wei
	transactOptsAuth.GasLimit = GasLimit
	transactOptsAuth.GasPrice = gasPrice

	var targetContract ContractRegistry

	switch event {
	// ProphecyClaims are sent to the CosmosBridge contract
	case events.MsgBurn, events.MsgLock:
		targetContract = CosmosBridge
	// OracleClaims are sent to the Oracle contract
	case events.LogNewProphecyClaim:
		targetContract = Oracle
	default:
		panic("Invalid target contract address")
	}

	// Get the specific contract's address
	target, err := GetAddressFromBridgeRegistry(client, registry, targetContract)
	if err != nil {
		log.Fatal(err)
	}

	return client, transactOptsAuth, target
}
