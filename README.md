# ETH Bridge Zone

[![version](https://img.shields.io/github/tag/cosmos/peggy.svg)](https://github.com/cosmos/peggy/releases/latest)
[![CircleCI](https://circleci.com/gh/cosmos/peggy/tree/master.svg?style=svg)](https://circleci.com/gh/cosmos/peggy/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/cosmos/peggy)](https://goreportcard.com/report/github.com/cosmos/peggy)
[![LoC](https://tokei.rs/b1/github/cosmos/peggy)](https://github.com/cosmos/peggy)
[![API Reference](https://godoc.org/github.com/cosmos/peggy?status.svg)](https://godoc.org/github.com/cosmos/peggy)

## Summary

Unidirectional Peggy is the starting point for cross chain value transfers from the Ethereum blockchain to Cosmos-SDK based blockchains as part of the Ethereum Cosmos Bridge project. The system accepts incoming transfers of Ethereum tokens on an Ethereum smart contract, locking them while the transaction is validated and equitable funds issued to the intended recipient on the Cosmos bridge chain.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Disclaimer

This codebase, including all smart contract components, have not been professionally audited and are not intended for use in a production environment. As such, users should NOT trust the system to securely hold mainnet funds. Any developers attempting to use Unidirectional Peggy on the mainnet at this time will need to develop their own smart contracts or find another implementation.

## Architecture

See [here](./docs/architecture.md)

## Requirements
 - Go 1.13

## Example application

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

## Running and testing the application

First, initialize a chain and create accounts to test sending of a random token.

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

# Then, wait 10 seconds and in another terminal window, test things are ok by sending 10 tok tokens from the validator to the testuser
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --chain-id=peggy --yes

# Wait a few seconds for confirmation, then confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a) --trust-node
ebcli query account $(ebcli keys show testuser -a) --trust-node

# See the help for the ethbridge create claim function
ebcli tx ethbridge create-claim --help

# Now you can test out the ethbridge module by submitting a claim for an ethereum prophecy
# Create a bridge claim (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
ebcli tx ethbridge create-claim 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 3eth --from=validator --chain-id=peggy --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query ethbridge prophecy 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 --trust-node

# And finally, confirm that the prophecy was successfully processed and that new eth was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node

```

## Using the application from rest-server

First, run the cli rest-server

```bash
ebcli rest-server --trust-node
```

An api collection for [Postman](https://www.getpostman.com/) is provided [here](./docs/peggy.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.

## Running the relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events.

```bash
# Check ebrelayer connection to ebd
ebrelayer status

# Initialize the Relayer service for automatic claim processing
ebrelayer init wss://ropsten.infura.io/ws ec6df30846baab06fce9b1721608853193913c19 "LogLock\(bytes32,address,bytes,address,uint256,uint256\)" validator --chain-id=peggy

# Enter password and press enter
# You should see a message like:  Started ethereum websocket... and Subscribed to contract events...
```

The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.

## Using the bridge

With the application set up and the relayer running, you can now use Peggy by sending a lock transaction to the smart contract. You can do this from any Ethereum wallet/client that supports smart contract transactions.

The easiest way to do this for now, assuming you have Metamask setup for Ropsten in the browser is to use remix or mycrypto as the frontend, for example:

1. Go to remix.ethereum.org
2. Compile Peggy.sol with solc v0.5.0
3. Set the environment as Injected Web3 Ropsten
4. On 'Run' tab, select Peggy and enter "0xec6df30846baab06fce9b1721608853193913c19" in 'At Address' field
5. Select 'At Address' to load the deployed contract
6. Enter the following for the variables under function lock():

```javascript
 _recipient = [HASHED_COSMOS_RECIPIENT_ADDRESS] *(for testuser cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh, enter "0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668")*
 _token = [DEPLOYED_TOKEN_ADDRESS] *(erc20 not currently supported, enter "0x0000000000000000000000000000000000000000" for ethereum)*
 _amount = [WEI_AMOUNT]
```

7. Enter the same number from \_amount as the transaction's value (in wei)
8. Select "transact" to send the lock() transaction

Then, wait for the transaction to confirm and mine, and for the relayer to pick it up. You should see the successful output in the relayer console. You can also confirm the tokens have been minted by using the CLI again:

```bash
ebcli query account cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh --trust-node
```

## Using the modules in other projects

The ethbridge and oracle modules can be used in other cosmos-sdk applications by copying them into your application's modules folders and including them in the same way as in the example application. Each module may be moved to its own repo or integrated into the core Cosmos-SDK in future, for easier usage.

There are 2 nuances you need to be aware of when using these modules in other Cosmos-SDK projects.

- A specific version of golang.org/x/crypto (ie tendermint/crypto) is needed for compatability with go-ethereum. See the Gopkg.toml for constraint details. There is an open pull request to tendermint/crypto to add compatbility, but until that is merged you need to use the [customized version](https://github.com/tendermint/crypto/pull/1)
- The govendor steps in the application as above are needed

For instructions on building and deploying the smart contracts, see the README in their folder.
