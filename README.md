# Peggy

[![version](https://img.shields.io/github/tag/cosmos/peggy.svg)](https://github.com/cosmos/peggy/releases/latest)
[![CircleCI](https://circleci.com/gh/cosmos/peggy/tree/master.svg?style=svg)](https://circleci.com/gh/cosmos/peggy/tree/master)
![Go Tests](https://github.com/cosmos/peggy/workflows/test%20coverage/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/cosmos/peggy)](https://goreportcard.com/report/github.com/cosmos/peggy)
[![LoC](https://tokei.rs/b1/github/cosmos/peggy)](https://github.com/cosmos/peggy)
[![API Reference](https://godoc.org/github.com/cosmos/peggy?status.svg)](https://godoc.org/github.com/cosmos/peggy)

> :warning: :warning: **WARNING: This bridge is not production ready. There is at least one known security vulnerability. See: [181](https://github.com/cosmos/peggy/issues/181).**


## Introduction

Peggy is the starting point for cross chain value transfers from the Ethereum blockchain to Cosmos-SDK based blockchains as part of the Ethereum Cosmos Bridge project. The system accepts incoming transfers of Ethereum tokens on an Ethereum smart contract, locking them while the transaction is validated and equitable funds issued to the intended recipient on the Cosmos bridge chain. The system supports value transfers from Cosmos-SDK based blockchains to the Ethereum blockchain as well through a reverse process.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Disclaimer

This codebase, including all smart contract components, has **not** been professionally audited and is not intended for use in a production environment. As such, users should **NOT** trust the system to securely hold mainnet funds. Any developers attempting to use Peggy on the mainnet at this time will need to develop their own smart contracts or find another implementation.

## Installation

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

## Usage Steps

- **Initialization**: setup the Bridge chain, add accounts, start the Bridge chain, and test available commands
- **Setup Peggy locally**: start local Ethereum blockchain, compile and deploy contracts, activate the contracts, and get the registry contract's deployed address
- **Run the Relayer**: start the relayer with a validator account in order to relay lock and burn events between the EVM blockchain and the Cosmos SDK blockchain.
- **Ethereum to Cosmos asset transfers**: start the Relayer service, send lock transaction containing local assets to the contracts, and test ERC20 support
- **Using Peggy with the Ropsten testnet**: setup interfacting with the Ropsten testnet, deploy contracts to Ropsten testnet, start the Relayer service on the Ropsten testnet, and send lock transaction containing Ropsten testnet assets to the contracts
- **Cosmos to Ethereum asset transfers**: setup interfacing with tendermint, start the Relayer service, start the Oracle Claim Relayer service, burn assets on tendermint, create prophecy and oracle claims on Ethereum, and process prophecy claims

## Initialization

In order to facilitate cross chain transfers, the Bridge blockchain must be set up by following these [steps](./docs/setup-bridge-chain.md).

## Setup Peggy locally

To test the transfer of Ethereum based assets, set up and start a local Ethereum chain by following these [steps](./docs/setup-eth-local.md).

## Run the Relayer

To set up and operate the relayer, follow the instructions [here](./setup-relayer.md).

## Ethereum to Cosmos asset transfers

With a local Ethereum blockchain running, you can participate in Ethereum -> Cosmos asset transfers by starting the Relayer service and acting as a validator. Validators witness the locking of Ethereum/ERC20 assets and sign a data package containing information about the lock, which is then relayed to the Cosmos chain and witnessed by the EthBridge module. Once a quorum of validators have confirmed that the transaction's information is valid, the funds are released by the Oracle module and transferred to the intended recipient's address. In this way, Ethereum assets can be transferred to Cosmos-SDK based blockchains. Instructions for trying out the process yourself are [here](./docs/ethereum-to-cosmos.md).

## Using Peggy with the Ropsten testnet

Instead of transferring local Ethereum assets to Cosmos-SDK based blockchains, you can test out transferring rEth from the Ropsten testnet by following these [steps](./docs/setup-eth-ropsten.md).

## Cosmos to Ethereum asset transfers

Cosmos -> Ethereum asset transfers are facilitated by a reverse process where validators witness transactions on tendermint and sign a data package containing the information. Cosmos assets can be locked, resulting in the release of funds held on Ethereum, or burned, resulting in the minting of new ERC20 tokens on Ethereum which represent the burned assets. The data package containing the validator's signature is then relayed to the contracts deployed on the Ethereum blockchain. Once enough other validators have confirmed that the transaction's information is valid, the funds are released/minted to the intended recipient's Ethereum address. In this way, assets on Cosmos-SDK based blockchains can be transferred to Ethereum. The process is described [here](./docs/cosmos-to-ethereum.md).

## Using the application from rest-server

First, run the cli rest-server

```bash
ebcli rest-server --trust-node
```

An api collection for [Postman](https://www.getpostman.com/) is provided [here](./docs/peggy.postman_collection.json) which documents some API endpoints and can be used to interact with it.

Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.

## Using the modules in other projects

The ethbridge and oracle modules can be used in other cosmos-sdk applications by copying them into your application's modules folders and including them in the same way as in the example application. Each module may be moved to its own repo or integrated into the core Cosmos-SDK in future, for easier usage.

## Architecture

A diagram of the protocol's architecture can be found [here](./docs/architecture.md).
