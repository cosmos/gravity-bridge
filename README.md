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

# Then wait 10 seconds then confirm your validator was created correctly, and has become Bonded status
ebcli query staking validators --trust-node

# See the help for the ethbridge create claim function
ebcli tx ethbridge create-claim --help

# Now you can test out the ethbridge module by submitting a claim for an ethereum prophecy
# Create a bridge lock claim (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
# ebcli tx ethbridge create-claim [ethereum-chain-id] [bridge-contract] [nonce] [symbol] [token-contract] [ethereum-sender-address] [cosmos-receiver-address] [validator-address] [amount]",
ebcli tx ethbridge create-claim 3 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB 0 eth 0x0000000000000000000000000000000000000000 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 3eth lock --from=validator --chain-id=peggy --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query ethbridge prophecy 3 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB 0 eth 0x0000000000000000000000000000000000000000 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 --trust-node

# Confirm that the prophecy was successfully processed and that new eth was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out burning 1 of the eth for the return trip
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 1eth --from=testuser --chain-id=peggy --yes

# Confirm that the eth was successfully burned
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out locking up a cosmos stake coin for relaying over to ethereum
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 1stake --from=testuser --chain-id=peggy --yes

# Confirm that the eth was successfully locked
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out creating a bridge burn claim for the return trip back
ebcli tx ethbridge create-claim 1 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 1stake burn --from=validator --chain-id=peggy --yes

# Confirm that the prophecy was successfully processed and that stake coin was returned to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node

```

## Using the application from rest-server

First, run the cli rest-server

```bash
ebcli rest-server --trust-node
```

An api collection for [Postman](https://www.getpostman.com/) is provided [here](./docs/peggy.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.

## Running the bridge locally

With the application set up, you can now use Peggy by sending a lock transaction to the smart contract.

### Set-up

```bash
cd testnet-contracts/

# Create .env with sample environment variables
cp .env.example .env
```

For running the bridge locally, you'll only need the `LOCAL_PROVIDER` environment variables. Environment variables `MNEMONIC` and `INFURA_PROJECT_ID` are required for using the Ropsten testnet.

Further reading:

- [MetaMask Mnemonic](https://metamask.zendesk.com/hc/en-us/articles/360015290032-How-to-Reveal-Your-Seed-Phrase)
- [Infura Project ID](https://blog.infura.io/introducing-the-infura-dashboard-8969b7ab94e7)

### Terminal 1: Start local blockchain

```bash
# Download dependencies
yarn

# Start local blockchain
yarn develop
```

### Terminal 2: Compile and deploy Peggy contract

```bash
# Deploy contract to local blockchain
yarn migrate

# Copy contract ABI to go modules:
yarn peggy:abi

# Get contract's address
yarn peggy:address
```

### Terminal 3: Build and start Ethereum Bridge

```bash
# Build the Ethereum Bridge application
make install

# Start the Bridge's blockchain
ebd start
```

### Terminal 4: Start the Relayer service

For automated relaying, there is a relayer service that can be run that will automatically watch and relay events (local web socket and deployed address parameters may vary).

```bash
# Check ebrelayer connection to ebd
ebrelayer status

# Start ebrelayer on the contract's deployed address with [LOCAL_WEB_SOCKET] and [PEGGY_DEPLOYED_ADDRESS]
# Example [LOCAL_WEB_SOCKET]: ws://127.0.0.1:7545/
# Example [PEGGY_DEPLOYED_ADDRESS]: 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB

ebrelayer init [LOCAL_WEB_SOCKET] [PEGGY_DEPLOYED_ADDRESS] LogLock\(address,bytes,address,string,uint256,uint256\) validator --chain-id=peggy

# Enter password and press enter
# You should see a message like: Started ethereum websocket with provider: [LOCAL_WEB_SOCKET] \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]

# The relayer will now watch the contract on Ropsten and create a claim whenever it detects a lock event.
```

### Using Terminal 2: Send lock transaction to contract

```bash
# Default parameter values:
# [HASHED_COSMOS_RECIPIENT_ADDRESS] = 0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668
# [TOKEN_CONTRACT_ADDRESS] = 0x0000000000000000000000000000000000000000
# [WEI_AMOUNT] = 10

# Send lock transaction with default parameters
yarn peggy:lock --default

# Send lock transaction with custom parameters
yarn peggy:lock [HASHED_COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]

```

`yarn peggy:lock --default` expected output in ebrelayer console:

```bash
New Lock Transaction:
Tx hash: 0x83e6ee88c20178616e68fee2477d21e84f16dcf6bac892b18b52c000345864c0
Block number: 5
Event ID: cc10955295e555130c865949fb1fd48dba592d607ae582b43a2f3f0addce83f2
Token: 0x0000000000000000000000000000000000000000
Sender: 0xc230f38FF05860753840e0d7cbC66128ad308B67
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 10
Nonce: 1

Response:
Height: 48
TxHash: AD842C51B4347F0F610CB524529C2D8A875DACF12C8FE4B308931D266FEAD067
Logs: [{"msg_index":0,"success":true,"log":"success"}]
GasWanted: 200000
GasUsed: 42112
Tags: - action = create_bridge_claim
```

## Running the bridge on the Ropsten testnet

To run the Ethereum Bridge on the Ropsten testnet, repeat the steps for running locally with the following changes:

```bash
# Add environment variable MNEMONIC from your MetaMask account

# Add environment variable INFURA_PROJECT_ID from your Infura account.

# Specify the Ropsten network via a --network flag for the following commands...

yarn migrate --network ropsten
yarn peggy:address --network ropsten

# Make sure to start ebrelayer with Ropsten network websocket
ebrelayer init wss://ropsten.infura.io/ [PEGGY_DEPLOYED_ADDRESS] LogLock\(bytes32,address,bytes,address,string,uint256,uint256\) validator --chain-id=peggy

# Send lock transaction on Ropsten testnet

yarn peggy:lock --network ropsten [HASHED_COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]

```

## Testing ERC20 token support

The bridge supports the transfer of ERC20 token assets. A sample TEST token is deployed upon migration and can be used to locally test the feature.

### Local

```bash
# Mint 1,000 TEST tokens to your account for local use
yarn token:mint

# Approve 100 TEST tokens to the Bridge contract
yarn token:approve --default

# You can also approve a custom amount of TEST tokens to the Bridge contract:
yarn token:approve 3

# Get deployed TEST token contract address
yarn token:address

# Lock TEST tokens on the Bridge contract
yarn peggy:lock [HASHED_COSMOS_RECIPIENT_ADDRESS] [TEST_TOKEN_CONTRACT_ADDRESS] [TOKEN_AMOUNT]

```

`yarn peggy:lock` ERC20 expected output in ebrelayer console (with a TOKEN_AMOUNT of 3):

```bash
New Lock Transaction:
Tx hash: 0xce7b219427c613c8927f7cafe123af4145016a490cd9fef6e3796d468f72e09f
Event ID: bb1c4798aaf4a1236f4f0235495f54a135733446f6c401c1bb86b690f3f35e60
Token Symbol: TEST
Token Address: 0x5040BA3Cf968de7273201d7C119bB8D8F03BDcBc
Sender: 0xc230f38FF05860753840e0d7cbC66128ad308B67
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 3
Nonce: 2

Response:
  height: 0
  txhash: DF1F55D2B8F4277671772D9A72188D0E4E15097AD28272E31116FF4B5D832B08
  code: 0
  data: ""
  rawlog: '[{"msg_index":0,"success":true,"log":""}]'
  logs:
  - msgindex: 0
    success: true
    log: ""
```

## Using the modules in other projects

The ethbridge and oracle modules can be used in other cosmos-sdk applications by copying them into your application's modules folders and including them in the same way as in the example application. Each module may be moved to its own repo or integrated into the core Cosmos-SDK in future, for easier usage.
