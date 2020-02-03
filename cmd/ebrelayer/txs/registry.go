package txs

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	bridgeRegistry "github.com/cosmos/peggy/cmd/ebrelayer/generated/bridgeregistry"
)

// ContractRegistry :
type ContractRegistry byte

const (
	// Valset valset contract
	Valset ContractRegistry = iota + 1
	// Oracle oracle contract
	Oracle
	// BridgeBank bridgeBank contract
	BridgeBank
	// CosmosBridge cosmosBridge contract
	CosmosBridge
)

// String returns the event type as a string
func (d ContractRegistry) String() string {
	return [...]string{"valset", "oracle", "bridgebank", "cosmosbridge"}[d-1]
}

// GetAddressFromBridgeRegistry utility method which queries the requested contract address from the BridgeRegistry
func GetAddressFromBridgeRegistry(
	client *ethclient.Client, registry common.Address, target ContractRegistry,
) (common.Address, error) {
	// Load the sender's address
	sender, err := LoadSender()
	if err != nil {
		log.Fatal(err)
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set up CallOpts auth
	auth := bind.CallOpts{
		Pending:     true,
		From:        sender,
		BlockNumber: header.Number,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	registryInstance, err := bridgeRegistry.NewBridgeRegistry(registry, client)
	if err != nil {
		log.Fatal(err)
	}

	var address common.Address
	switch target {
	case Valset:
		address, err = registryInstance.Valset(&auth)
	case Oracle:
		address, err = registryInstance.Oracle(&auth)
	case BridgeBank:
		address, err = registryInstance.BridgeBank(&auth)
	case CosmosBridge:
		address, err = registryInstance.CosmosBridge(&auth)
	default:
		panic("invalid target contract address")
	}

	if err != nil {
		log.Fatal(err)
	}

	return address, nil
}
