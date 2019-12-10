## Cosmos to Ethereum asset transfers

### Start local Ethereum blockchain and application

Before you can transfer Cosmos assets to Ethereum, you'll need to have a local Ethereum blockchain with the Peggy contracts deployed to it as described [here](./local-ethereum-usage.md). If you've already started a local blockchain and deployed the contracts, you can skip this step.

You'll also need to start the Bridge blockchain if it's not already running. To do so, follow these ([steps](./initialization.md)).

### Setup

In order to send transactions to the contracts, the Cosmos Relayer requires the private key of an active validator. The private key must be set as an environment variable named `ETHEREUM_PRIVATE_KEY` and located in the .env file at the root of the project. If testing locally, can use the private key of accounts[1], which can be found in the truffle console running in terminal 1. If testing on a live network, you'll need to use the private key of your Ethereum address.

### Start the Relayer service

```bash
# Open a new terminal window

# Check ebrelayer connection to ebd
ebrelayer status

# Start Cosmos relayer
# Note: ports for the tendermint node (tcp://localhost:) and web3 provider (http://localhost:) may vary
# Note: Use the address from 'yarn peggy:address` for [PEGGY_CONTRACT_ADDRESS]
ebrelayer init cosmos tcp://localhost:26657 http://localhost:7545 [PEGGY_CONTRACT_ADDRESS]

# You should see a message like:
# [2019-10-24|19:02:21.888] Starting WSEvents         impl=WSEvents

# The relayer will now watch the Cosmos network and create a prophecy claim whenever it detects a burn or lock event
```

### Start the Oracle Claim Relayer

To make an Oracle Claim on every Prophecy Claim witnessed, start an Ethereum relayer with flag `--make-claims=true`

Note: For now, close any other active Ethereum Relayers currently running.

```bash
# Open a new terminal window

# Start ebrelayer on the contract's deployed address with [PEGGY_DEPLOYED_ADDRESS]
ebrelayer init ethereum ws://127.0.0.1:7545/ [PEGGY_DEPLOYED_ADDRESS] validator --make-claims=true --chain-id=peggy

# Enter password and press enter

# The relayer will now watch the contract on Ropsten and create a new oracle claim whenever it detects a new prophecy claim event
```

### Sending Cosmos assets to Ethereum via Lock

To send Cosmos assets an EVM based chain, you'll use a transaction containing a lock message:

```bash
# Open a new terminal window

# Send tokens to the testuser (10stake tokens)
ebcli tx send validator $(ebcli keys show testuser -a) 10stake --chain-id=peggy --yes

# Send lock transaction (1stake token)
ebcli tx ethbridge lock $(ebcli keys show testuser -a) [RECIPIENT_ETHEREUM_ADDRESS] 1stake --from testuser --chain-id peggy --ethereum-chain-id 3 --token-contract-address [TOKEN_CONTRACT_ADDRESS]
# Note: --token-contract-address will be '0x0000000000000000000000000000000000000000' for Ethereum

# Enter 'y' to confirm the transaction

# Enter testuser's password

# You should see the transaction output in this terminal with 'success:true' in the 'rawlog' field:
# rawlog: '[{"msg_index":0,"success":true,"log":""}]'
```

Expected output in the Cosmos Relayer console

```bash
[2019-10-24|19:07:01.714]       New transaction witnessed

Msg Type: lock
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

Expected output in Oracle Claim Relayer console

```bash
Witnessed new event: LogNewProphecyClaim
Block number: 43
Tx hash: 0xb14695d7ca229c713c89ab2e78c41549cfac11daed6d09ab4b9755b12b46f17c

Prophecy ID: 2
Claim Type: 2
Sender: cosmos1qwnw2r9ak79536c4dqtrtk2pl2nlzpqh763rls
Recipient 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359
Symbol stake
Token 0xbEDdB076fa4dF04859098A9873591dcE3E9C404d
Amount: 1
Validator: 0xc230f38FF05860753840e0d7cbC66128ad308B67

Using validator account: c230f38ff05860753840e0d7cbc66128ad308b67
Attempting to sign message: "0xb8b701ef59944e115d6ecfd4aa1bd03025d85338d771b0099d4061923bd0a1ed"
Success! Signature: 0x919ca03752269c87c5df9f4af99ba49be84cb2bbc77921db581719379e95c548164b55822e89294b8066f77812695d9575b4827c04592d4daa41dd087ba1ba7f01

Fetching Oracle contract...
Sending new OracleClaim to Oracle...
NewOracleClaim tx hash: 0x89c1c905f65170e799fc17b16406aad61e07c857f3379190829f5fd5f9a157d9
Tx Status: 1 - Successful
```

Congratulations, you've automatically relayed information from the lock transaction on Tendermint to the contracts deployed on the Ethereum network as a new prophecy claim, witnessed the new prophecy claim, and signed its information to create an oracle claim. When enough validators submit oracle claims for the prophecy claim, it will be processed. When a prophecy claim is successfully processed, the amount of tokens specified will be minted by the contracts to the intended recipient on the Ethereum network.

### Returning Cosmos assets originally based on Ethereum via Burn

In the `Ethereum to Cosmos asset transfers` section, you sent assets to a Cosmos-SDK enabled chain. In order to return these assets to Ethereum and unlock the funds currently locked on the deployed contracts, you'll need to use a second type of transaction - burn. It's simple, just replace the ebcli `lock` command with `burn`:

```bash
# Send burn transaction (1stake token)
ebcli tx ethbridge burn $(ebcli keys show testuser -a) [RECIPIENT_ETHEREUM_ADDRESS] 1stake --from testuser --chain-id peggy --ethereum-chain-id 3 --token-contract-address [TOKEN_CONTRACT_ADDRESS]
```

## Prophecy claim processing

You are able to check the status of active prophecy claims. Prophecy claims reach the signed power threshold when the weighted signed power surpasses the weighted total power, where *weighted total power* = (total power * 2) and *weighted signed power* = (signed power * 3).

```bash
# Check prophecy claim status
yarn peggy:check [PROPHECY_CLAIM_ID]
```

Expected output (for a prophecy claim with an ID of 2)

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


Once the prophecy claim has reached the signed power threshold, anyone may initiate its processing. Any attempts to process prophecy claims under the signed power threshold will be rejected by the contracts.   


```bash
# Process the prophecy claim
yarn peggy:process [PROPHECY_CLAIM_ID]
```