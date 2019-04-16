# ETH Bridge Zone

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

# Update dependencies to match the constraints and overrides above
dep ensure -update -v

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

ebcli keys add testuser

ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000tok

# Now its safe to start `ebd`
ebd start

# Then, in another terminal window, send 10 tok tokens from the validator to the testuser
ebcli tx send $(ebcli keys show testuser -a) 10tok --from=validator --chain-id=testing --yes

# Confirm token balances have changed appropriately
ebcli query account $(ebcli keys show validator -a) --trust-node
ebcli query account $(ebcli keys show testuser -a) --trust-node

# Test out the oracle module by submitting a claim for an ethereum prophecy
# (Ethereum prophecies are stored on the blockchain with an identifier created by concatenating the nonce and sender address)
ebcli tx oracle make-claim 0 randomethaddress $(ebcli keys show testuser -a) $(ebcli keys show validator -a) 3eth --from validator --chain-id testing --yes

# Then read the prophecy to confirm it was created with the claim added
ebcli query oracle get-prophecy 0randomethaddress --trust-node
```

## Using the application from rest-server

First, run the cli rest-server

```
ebcli rest-server --trust-node
```

An api collection for Postman (https://www.getpostman.com/) is provided [here](./cosmos-ethereum-bridge.postman_collection.json) which documents some API endpoints and can be used to interact with it.
Note: You will need to change the cosmos addresses in the URLs, params and body to match the addresses you generated that you want to check.