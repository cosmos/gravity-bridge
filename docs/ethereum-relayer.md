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

# Activate the contracts (required)
yarn peggy:setup

# Get contract's address
yarn peggy:address
```

### Start Bridge blockchain (terminal 3)

TODO: Add generation step

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

# Start ebrelayer on the contract's deployed address with [LOCAL_WEB_SOCKET] and [REGISTRY_DEPLOYED_ADDRESS]
# [LOCAL_WEB_SOCKET] should be similar to 'ws://127.0.0.1:7545/'
# [REGISTRY_DEPLOYED_ADDRESS] should be similar to '0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB'
ebrelayer init ethereum [LOCAL_WEB_SOCKET] [REGISTRY_DEPLOYED_ADDRESS] validator --chain-id=peggy

# Enter password and press enter

# You should see a message like: Started Ethereum websocket with provider: [LOCAL_WEB_SOCKET] \ Subscribed to contract events on address: [PEGGY_DEPLOYED_ADDRESS]

# The Relayer will now listen to the deployed contracts and create a claim whenever it detects a new lock event
```

### Lock Ethereum assets on contracts (use terminal 2)

```bash
# Default parameter values:
# [HASHED_COSMOS_RECIPIENT_ADDRESS] = '0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668'
# [TOKEN_CONTRACT_ADDRESS] = '0x0000000000000000000000000000000000000000' (null address denotes Ethereum)
# [WEI_AMOUNT] = '10'

# Send lock transaction with default parameters
yarn peggy:lock --default

# Send lock transaction with custom parameters
yarn peggy:lock [HASHED_COSMOS_RECIPIENT_ADDRESS] [TOKEN_CONTRACT_ADDRESS] [WEI_AMOUNT]

```

`yarn peggy:lock --default` expected output in Relayer console (terminal 4):

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

## Using the Bridge with the Ropsten testnet

### Setup

Before you can use the Bridge with the Ropsten testnet, you'll need to add two environment variables to the configuration file at `testnet-contracts/.env`. Add MNEMONIC from your MetaMask account, this will allow you to deploy the contracts to the Ropsten testnet. Add INFURA_PROJECT_ID from your Infura account, this will allow you to start a Relayer service which listens for events on the Ropsten testnet.


### Usage

```bash

# Deploy the contracts to the Ropsten network with the --network flag
yarn migrate --network ropsten

# Get the Registry contract's address on the Ropsten network with the --network flag
yarn peggy:address --network ropsten

# Restart Relayer with Infura's Ropsten network websocket
ebrelayer init ethereum wss://ropsten.infura.io/ [REGISTRY_DEPLOYED_ADDRESS] validator --chain-id=peggy

# Send funds to the deployed contracts on the Ropsten testnet
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