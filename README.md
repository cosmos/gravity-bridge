# ETH Bridge Zone

[![CircleCI](https://circleci.com/gh/swishlabsco/cosmos-ethereum-bridge/tree/master.svg?style=svg)](https://circleci.com/gh/swishlabsco/cosmos-ethereum-bridge/tree/master)

This repository contains the source code of the ethereum bridge zone. 

## Installing and building the application

```
# Clone the repository
mkdir -p $GOPATH/src/github.com/swishlabsco
cd $GOPATH/src/github.com/swishlabsco
git clone https://github.com/swishlabsco/cosmos-ethereum-bridge
cd cosmos-ethereum-bridge && git checkout master

# Install dep, as well as your dependencies
make get_tools
dep ensure -v
go get -u github.com/kardianos/govendor

# Fetch the C file dependencies (this is a manual hack, as dep does not support pulling non-golang files used in dependencies)
govendor fetch -tree github.com/ethereum/go-ethereum/crypto/secp256k1

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands, confirming the build is successful:
ebd help
ebcli help
```

## Running and using the application

First, initialize a chain and create accounts to test sending of a random token.

```
# Initialize the genesis.json file that will help you to bootstrap the network
ebd init --chain-id=testing

# Create a key to hold your validator account and for another test account
ebcli keys add validator
# Enter password

ebcli keys add testuser
# Enter password

ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom

# Now its safe to start `ebd`
ebd start

# Then, wait 10 seconds and in another terminal window, test things are ok by sending 10 tok tokens from the validator to the testuser
ebcli tx send $(ebcli keys show testuser -a) 10stake --from=validator --chain-id=testing --yes

# Confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a) --trust-node
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Next, setup the staking module prerequisites
# First, create a validator and stake
ebcli tx staking create-validator \
  --amount=100000000stake \
  --pubkey=$(ebd tendermint show-validator) \
  --moniker="test_moniker" \
  --chain-id=testing \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas=200000 \
  --gas-prices="0.001stake" \
  --from=validator

# Then wait 10 seconds then confirm your validator was created correctly, and has become Bonded status
ebcli query staking validators --trust-node 

# See the help for the ethbridge make claim function
ebcli tx ethbridge make-claim --help

# Now you can test out the ethbridge module by submitting a claim for an ethereum prophecy
# Make a bridge claim (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
ebcli tx ethbridge make-claim 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 $(ebcli keys show testuser -a) $(ebcli keys show validator -a) 3eth --from validator --chain-id testing --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query ethbridge get-prophecy 0 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 --trust-node

# And finally, confirm that the prophecy was successfully processed and that new eth was minted to the testuser address
ebcli query account $(ebcli keys show testuser -a) --trust-node 
```

## Using the application from rest-server

First, run the cli rest-server

```
ebcli rest-server --trust-node
```

An api collection for Postman (https://www.getpostman.com/) is provided [here](./cosmos-ethereum-bridge.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: For checking account details/balance, you will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.