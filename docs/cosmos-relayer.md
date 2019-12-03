## Cosmos to Ethereum asset transfers

### Start local Ethereum blockchain (terminal 1)

```bash
# Download dependencies
yarn

# Start local blockchain
yarn develop

```

### Set up

In order to send transactions to the contracts, the Cosmos Relayer requires the private key of an active validator. The private key must be set as an environment variable named `ETHEREUM_PRIVATE_KEY` and located in the .env file at the root of the project. If testing locally, can use the private key of accounts[1], which can be found in the truffle console running in terminal 1. If testing on a live network, you'll need to use the private key of your Ethereum address.

### Deploy Peggy contracts (terminal 2)

```bash
# Deploy contract to local blockchain
yarn migrate

# Activate the contracts (required)
yarn peggy:setup

# Get the address of Peggy's registry service (required to start Cosmos relayer)
yarn peggy:address
```

### Start Bridge blockchain (terminal 3)
TODO: Add generation step

```bash
# Build the Bridge application
make install

# Start the Bridge's blockchain
ebd start
```

### Start the Relayer service which watches Tendermint (terminal 4)

```bash
# Check ebrelayer connection to ebd
ebrelayer status

# Start Cosmos relayer
# Example [tendermintNode]: tcp://localhost:26657
# Example [web3Provider]: http://localhost:7545
ebrelayer init cosmos [tendermintNode] [web3Provider] [bridgeRegistryContractAddress]

# You should see a message like:
# [2019-10-24|19:02:21.888] Starting WSEvents         impl=WSEvents

# The relayer will now watch the Cosmos network and create a prophecy claim whenever it detects a burn or lock event
```

### Start the Oracle Claim Relayer (terminal 5)

To make an Oracle Claim on every Prophecy Claim witnessed, start an Ethereum relayer with flag `--make-claims=true`

```bash
# Start ebrelayer on the contract's deployed address with [LOCAL_WEB_SOCKET] and [REGISTRY_DEPLOYED_ADDRESS]
ebrelayer init ethereum [LOCAL_WEB_SOCKET] [REGISTRY_DEPLOYED_ADDRESS] validator --make-claims=true --chain-id=peggy

# Enter password and press enter

# The relayer will now watch the contract on Ropsten and create a new oracle claim whenever it detects a new prophecy claim event
```

### Send burn transaction on Cosmos

```bash
# Send some tokens to the testuser using the process described in section "Running and testing the application"

# Send burn transaction in terminal 2
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 1stake --from testuser --chain-id peggy

# Enter 'y' to confirm the transaction

# Enter testuser's password

# You should see the transaction output in this terminal with 'success:true' in the 'rawlog' field:
# rawlog: '[{"msg_index":0,"success":true,"log":""}]'

```

Expected output in the Cosmos Relayer console (terminal 4)

```bash
[2019-10-24|19:07:01.714]       New transaction witnessed

Msg Type: burn
Cosmos Sender: cosmos1qwnw2r9ak79536c4dqtrtk2pl2nlzpqh763rls
Ethereum Recipient: 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359
Token Address: 0xbEDdB076fa4dF04859098A9873591dcE3E9C404d
Symbol: stake
Amount: 1

Fetching CosmosBridge contract...
Sending tx to CosmosBridge...

NewProphecyClaim tx hash: 0x5544bdb31b90da102c0b7fd959b3106b823805871ddcbe972a7877ad15164631
Status: 1 - Successful
```

Expected output in Oracle Claim Relayer console (terminal 5)

```bash

New "LogNewProphecyClaim":
Tx hash: 0xb14695d7ca229c713c89ab2e78c41549cfac11daed6d09ab4b9755b12b46f17c
Block number: 18
Prophecy ID: 2
Claim Type: 0
Sender: cosmos1qwnw2r9ak79536c4dqtrtk2pl2nlzpqh763rls
Recipient 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359
Symbol eth
Token 0xbEDdB076fa4dF04859098A9873591dcE3E9C404d
Amount: 1
Validator: 0xc230f38FF05860753840e0d7cbC66128ad308B67


Attempting to sign message "0xb8b701ef59944e115d6ecfd4aa1bd03025d85338d771b0099d4061923bd0a1ed" with account "c230f38ff05860753840e0d7cbc66128ad308b67"...
Success! Signature: 0x919ca03752269c87c5df9f4af99ba49be84cb2bbc77921db581719379e95c548164b55822e89294b8066f77812695d9575b4827c04592d4daa41dd087ba1ba7f01
```

Congratulations, you've automatically relayed information from the burn transaction on Tendermint to the contracts deployed on the Ethereum network as a new prophecy claim, witnessed the new prophecy claim, and signed its information to create an oracle claim. When enough validators submit oracle claims for the prophecy claim, it will be processed. When a prophecy claim is successfully processed, the funds are unlocked on the deployed contracts and sent to the intended recipient on the Ethereum network.   


You'll be able to check the status of an active prophecy claim

```bash
# Check prophecy claim status
yarn peggy:check [PROPHECY_CLAIM_ID]
```

Expected output:

```bash

Fetching Oracle contract...
Attempting to send checkBridgeProphecy() tx...

        Prophecy 2 status:
----------------------------------------
Weighted total power:    104
Weighted signed power:   150
Reached threshold:       true
----------------------------------------
```