## Ethereum to Cosmos asset transfers

With the application set up, you can now use Peggy by sending a lock transaction to the smart contract.

### Start local Ethereum blockchain (terminal 1)

```bash
# Open a new terminal window, this will be terminal 1

# Download dependencies
yarn

# Start local blockchain
yarn develop
```

### Deploy Peggy contracts (terminal 2)

Next, compile and deploy Peggy's contracts to the local Ethereum blockchain.

```bash
# Open a new terminal window, this will be terminal 2

# Deploy contract to local blockchain
yarn migrate

# Activate the contracts
yarn peggy:setup

# Copy contract ABI to Relayer it can subscribe to deployed contracts
yarn peggy:abi

# Get the address of Peggy's registry contract
yarn peggy:address
```

### Start Bridge blockchain (terminal 3)

If you've already started the Bridge blockchain, you can skip this step.

```bash
# Open a new terminal window, this will be terminal 3

# Build the Ethereum Bridge application
make install

# Start the Bridge's blockchain
ebd start
```

### Start the Relayer service which watches Ethereum (terminal 4)

For automated relaying, validators can run a relayer service which will automatically watch for relevant events on the Ethereum network and relay them to the Bridge. Note that your local web socket and registry contract address may vary.

```bash
# Open a new terminal window, this will be terminal 4

# Check ebrelayer connection to ebd
ebrelayer status

# Start ebrelayer
# Use the contract address returned by `yarn peggy:address` for [PEGGY_DEPLOYED_ADDRESS]
# Note that you may need to update your websocket's localhost/port
ebrelayer init ethereum ws://127.0.0.1:7545/ [PEGGY_DEPLOYED_ADDRESS] validator --chain-id=peggy

# Enter password and press enter

# You should see a message like:
#   Started Ethereum websocket with provider: ws://127.0.0.1:7545/
#   Subscribed to bridgebank contract at address: 0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2

# The Relayer will now listen to the deployed contracts and create a claim whenever it detects a new lock event
```

### Lock Ethereum assets on contracts (use terminal 2)

```bash
# Default parameter values:
# [COSMOS_RECIPIENT_ADDRESS] = cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
# [TOKEN_CONTRACT_ADDRESS] = eth (Ethereum has no token contract and is denoted by 'eth')
# [WEI_AMOUNT] = 10

# Send lock transaction with default parameters
yarn peggy:lock --default

# Send lock transaction with custom parameters
yarn peggy:lock [COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]

```

`yarn peggy:lock --default` expected output in Relayer console (terminal 4):

```bash
Witnessed new event: LogLock
Block number: 19
Tx hash: 0xd90225d86aa59ff8fc3b59eb61b622c040c7f81a4f75dde32bcedb95494ccf12

Chain ID: 5777
Bridge contract address: 0x0823eFE0D0c6bd134a48cBd562fE4460aBE6e92c
Token symbol: ETH
Token contract address: 0x0000000000000000000000000000000000000000
Sender: 0x115F6e2004D7b4ccd6b9D5ab34e30909e0F612CD
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 10
Nonce: 5

height: 0
txhash: C1835DA4533BB9F9CD69DB80049CE8BF1576A6480D26D161FF04851CAAF305F6
code: 0
data: ""
rawlog: '[{"msg_index":0,"success":true,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"create_bridge_claim"}]}]}]'
```

## Using the Bridge with the Ropsten testnet

### Setup

Before you can use the Bridge with the Ropsten testnet, you'll need to add two environment variables to the configuration file at `testnet-contracts/.env`. Add MNEMONIC from your MetaMask account, this will allow you to deploy the contracts to the Ropsten testnet. Add INFURA_PROJECT_ID from your Infura account, this will allow you to start a Relayer service which listens for events on the Ropsten testnet.

Further reading:

- [MetaMask Mnemonic](https://metamask.zendesk.com/hc/en-us/articles/360015290032-How-to-Reveal-Your-Seed-Phrase)
- [Infura Project ID](https://blog.infura.io/introducing-the-infura-dashboard-8969b7ab94e7)


### Usage

```bash

# Deploy the contracts to the Ropsten network with the --network flag
yarn migrate --network ropsten

# Get the Registry contract's address on the Ropsten network with the --network flag
yarn peggy:address --network ropsten

# Restart Relayer with Infura's Ropsten network websocket
ebrelayer init ethereum wss://ropsten.infura.io/ [PEGGY_DEPLOYED_ADDRESS] validator --chain-id=peggy

# Send funds to the deployed contracts on the Ropsten testnet
# Note: [TOKEN_CONTRACT_ADDRESS] is 'eth' for Ethereum
yarn peggy:lock --network ropsten [COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]

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
# Note: ERC20 token locking requires 3 custom params and does not support the --default flag
yarn peggy:lock [COSMOS_RECIPIENT_ADDRESS] [TEST_TOKEN_CONTRACT_ADDRESS] [TOKEN_AMOUNT]

```

`yarn peggy:lock` ERC20 expected output in ebrelayer console (with a TOKEN_AMOUNT of 3):

`yarn peggy:lock` ERC20 expected output in ebrelayer console (with a TOKEN_AMOUNT of 11):

```bash
Witnessed new event: LogLock
Block number: 28
Tx hash: 0xab84de6d2f6bde3f2249cc1c31e23901432fa75b83a5b5b52c19e99479a797f1

Chain ID: 5777
Bridge contract address: 0x0823eFE0D0c6bd134a48cBd562fE4460aBE6e92c
Token symbol: TEST
Token contract address: 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB
Sender: 0x115F6e2004D7b4ccd6b9D5ab34e30909e0F612CD
Recipient: cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh
Value: 11
Nonce: 12

height: 0
txhash: 013B79C59828872BA477FC8C2B98C155A0F8D520C42693363B7156F56B6C0A32
code: 0
data: ""
rawlog: '[{"msg_index":0,"success":true,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"create_bridge_claim"}]}]}]'
```