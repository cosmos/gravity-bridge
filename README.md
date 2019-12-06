# Peggy

[![version](https://img.shields.io/github/tag/cosmos/peggy.svg)](https://github.com/cosmos/peggy/releases/latest)
[![CircleCI](https://circleci.com/gh/cosmos/peggy/tree/master.svg?style=svg)](https://circleci.com/gh/cosmos/peggy/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/cosmos/peggy)](https://goreportcard.com/report/github.com/cosmos/peggy)
[![LoC](https://tokei.rs/b1/github/cosmos/peggy)](https://github.com/cosmos/peggy)
[![API Reference](https://godoc.org/github.com/cosmos/peggy?status.svg)](https://godoc.org/github.com/cosmos/peggy)

## Summary

Peggy is the starting point for cross chain value transfers from the Ethereum blockchain to Cosmos-SDK based blockchains as part of the Ethereum Cosmos Bridge project. The system accepts incoming transfers of Ethereum tokens on an Ethereum smart contract, locking them while the transaction is validated and equitable funds issued to the intended recipient on the Cosmos bridge chain.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Disclaimer

This codebase, including all smart contract components, has not been professionally audited and are not intended for use in a production environment. As such, users should NOT trust the system to securely hold mainnet funds. Any developers attempting to use Peggy on the mainnet at this time will need to develop their own smart contracts or find another implementation.

## Architecture

A diagram of the protocol's architecture can be found [here](./docs/architecture.md).

## Requirements

- Go 1.13

## Building the application

These modules can be added to any Cosmos-SDK based chain, but a demo application/blockchain is provided with example code for how to integrate them. It can be installed and built as follows:

```bash
# Clone the repository
mkdir -p $GOPATH/src/github.com/cosmos
cd $GOPATH/src/github.com/cosmos
git clone https://github.com/cosmos/peggy
cd peggy && git checkout master

# Install tools (golangci-lint v1.18)
make tools-clean
make tools

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
ebd help
ebcli help
ebrelayer help
```

## Initializing the application

First, initialize a chain and create accounts.

```bash
# Initialize the genesis.json file that will help you to bootstrap the network
ebd init local --chain-id=peggy

# Create a key to hold your validator account and for another test account
ebcli keys add validator
# Enter password

ebcli keys add testuser
# Enter password

# Initialize the genesis account and transaction
ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom

# Create genesis transaction
ebd gentx --name validator
# Enter password

# Collect genesis transaction
ebd collect-gentxs

# Now its safe to start `ebd`
ebd start

```

To test the initialized application, follow the steps [here](./docs/running.md). It includes sending tokens between accounts, querying accounts, claim creation, token burning, and token locking. Once the Relayer is running, you can submit new burning/locking txs to the chain using these commands.

## Using the application from rest-server

First, run the cli rest-server

```bash
ebcli rest-server --trust-node
```

An api collection for [Postman](https://www.getpostman.com/) is provided [here](./docs/peggy.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.

## Relayer set up

In order for Peggy to process cross-chain asset transfers, the Relayer service must be run by a set of validators. Before validators participate in asset transfers, they must set up the appropriate configuration files with the following commands.

```bash
# Create .env with sample environment variables for the Cosmos relayer
cp .env.example .env

cd testnet-contracts/

# Create .env with sample environment variables for the Ethereum relayer
cp .env.example .env
```

## Ethereum to Cosmos asset transfers

Validators can participate in Ethereum -> Cosmos asset transfers by following the process described [here](./docs/ethereum-relayer.md).

## Cosmos to Ethereum asset transfers

Validators can participate in Cosmos -> Ethereum asset transfers by following the process described [here](./docs/cosmos-relayer.md).

## Using the modules in other projects

The ethbridge and oracle modules can be used in other cosmos-sdk applications by copying them into your application's modules folders and including them in the same way as in the example application. Each module may be moved to its own repo or integrated into the core Cosmos-SDK in future, for easier usage.
