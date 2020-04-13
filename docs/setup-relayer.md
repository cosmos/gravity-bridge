## Start Relayer

### Start local Ethereum blockchain and application

You must have both the Comsos SDK application and the Ethereum blockchain running before you'll be able to relay assets between the two. If you've already started the application + started a local Ethereum blockchain and deployed the contracts, move to the next section. If not follow these two sets of instructions:

- To start the Cosmos SDK application, follow these ([steps](./setup-bridge-chain.md))
- To start a local Ethereum blockchain and deploy the contrats, follow these [steps](./setup-eth-local.md)

### Setup

In order to send transactions to the contracts, the Relayer requires the ethereum private key of an active validator. The private key must be set as an environment variable named `ETHEREUM_PRIVATE_KEY` and located in the .env file at the root of the project. If testing locally, can keep the default private key of accounts[0], which can be found in the truffle console running in your `truffle develop` terminal (`c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3`). If testing on a live network, you'll need to use the private key of your Ethereum address.

### Setup

1. Copy root `.env.example` file:

```bash
cp .env.example .env
```

2. Generate contract bindings

```bash
ebrelayer generate
```

### Start the Relayer service

For automated tx relaying, validators can run a Relayer service which will automatically watch for events on both the Ethereum network and Tendermint. Note that Peggy's deployed BridgeRegistry contract address will vary. Peggy's deployed registry contract address is the address returned by `yarn peggy:address`. If you're using a fresh `truffle develop` EVM chain, the address should be `0x30753E4A8aad7F8597332E813735Def5dD395028`.

```bash
# In a new terminal window, check ebrelayer connection to ebd
ebrelayer status
# Start relayer
# ebrelayer init [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMoniker] [flags]
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator --chain-id=peggy
```

Expected terminal output:

```bash
I[2020-03-22|11:23:33.920] Starting WSEvents                            impl=WSEvents
I[2020-03-22|11:23:33.920] Starting WSClient                            impl="WSClient{localhost:26657 (/websocket)}"
I[2020-03-22|11:23:33.922] sent a request                               req="RPCRequest{0 subscribe/7B227175657279223A22746D2E6576656E74203D2027547827227D}"
I[2020-03-22|11:23:33.922] got response                                 id=0 result=7B7D
I[2020-03-22|11:23:33.922] Started Ethereum websocket with provider:    ws://localhost:7545/=(MISSING)
I[2020-03-22|11:23:33.955] Subscribed to bridgebank contract at address: 0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4
I[2020-03-22|11:23:33.980] Subscribed to cosmosbridge contract at address: 0x8f0483125FCb9aaAEFA9209D8E9d7b9C8B9Fb90F
```

The relayer is now subscribed to specific events and messages in transactions on the Ethereum and Tendermint chains, respectively. The relayer will create a prophecy claim whenever it detects a burn or lock event/message on either chain.

To try transferring assets from the EVM chain to the Cosmos SDK chain see the instructions [here](./ethereum-to-cosmos.md).
