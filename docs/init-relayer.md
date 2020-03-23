## Start Relayer

### Start local Ethereum blockchain and application

You must have both the application and Ethereum blockchain running before you'll be able to relay assets between the two. If you've already started the application + started a local Ethereum blockchain and deployed the contracts, move to the next section.
- To start the application, follow these ([steps](./setup-bridge-chain.md))
- To start a local Ethereum blockchain and deploy the contrats, follow these [steps](./setup-eth-local.md)

### Setup

In order to send transactions to the contracts, the Relayer requires the private key of an active validator. The private key must be set as an environment variable named `ETHEREUM_PRIVATE_KEY` and located in the .env file at the root of the project. If testing locally, can use the private key of accounts[1], which can be found in the truffle console running in terminal 1. If testing on a live network, you'll need to use the private key of your Ethereum address.

### Setup

1. Copy root `.env.example` file:  

```bash
cp .env.example .env
```

In order to send transactions to the contracts, the Relayer requires the private key of an active validator. The private key must be set as an environment variable named `ETHEREUM_PRIVATE_KEY` and located in the .env file at the root of the project. If testing locally, can use the private key of accounts[1], which can be found in the truffle console running in terminal 1. If testing on a live network, you'll need to use the private key of your Ethereum address.

2. Copy contracts `.env.example` file:  

```bash
cd testnet-contracts/
cp .env.example .env
```

### Start the Relayer service

For automated tx relaying, validators can run a Relayer service which will automatically watch for events on both the Ethereum network and Tendermint. Note that Peggy's deployed registry contract address will vary. Peggy's deployed registry contract address is the address returned by `yarn peggy:address`.

```bash
# In a new terminal window, check ebrelayer connection to ebd
ebrelayer status
# Start relayer
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ [REGISTRY_CONTRACT_ADDRESS] validator --chain-id=peggy
```

Expected terminal output:

```bash
Password to sign with 'validator':
I[2020-03-22|11:23:33.920] Starting WSEvents                            impl=WSEvents
I[2020-03-22|11:23:33.920] Starting WSClient                            impl="WSClient{localhost:26657 (/websocket)}"
I[2020-03-22|11:23:33.922] sent a request                               req="RPCRequest{0 subscribe/7B227175657279223A22746D2E6576656E74203D2027547827227D}"
I[2020-03-22|11:23:33.922] got response                                 id=0 result=7B7D
I[2020-03-22|11:23:33.922] Started Ethereum websocket with provider:    ws://localhost:7545/=(MISSING)
I[2020-03-22|11:23:33.955] Subscribed to bridgebank contract at address: 0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 
I[2020-03-22|11:23:33.980] Subscribed to cosmosbridge contract at address: 0x8f0483125FCb9aaAEFA9209D8E9d7b9C8B9Fb90F
```

The relayer is now subscribed to specific events and messages in transactions on the Ethereum and Tendermint chains, respectively. The relayer will create a prophecy claim whenever it detects a burn or lock event/message on either chain.