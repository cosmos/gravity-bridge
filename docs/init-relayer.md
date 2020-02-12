## Start Relayer

### Start local Ethereum blockchain and application

You must have both the application and Ethereum blockchain running before you'll be able to relay assets between the two. If you've already started the application + started a local Ethereum blockchain and deployed the contracts, move to the next section.
- To start the application, follow these ([steps](./initialization.md))
- To start a local Ethereum blockchain and deploy the contrats, follow these [steps](./local-ethereum-usage.md)

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

For automated tx relaying, validators can run a Relayer service which will automatically watch for events on both the Ethereum network and Tendermint. The port of your Tendermint node, the port of your web3 provider, and Peggy's deployed registry contract address will vary. Peggy's deployed registry contract address is the address returned by `yarn peggy:address`.

```bash
# Open a new terminal window

# Check ebrelayer connection to ebd
ebrelayer status

# Start relayer
# Note: double check both websocket ports
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ [PEGGY_CONTRACT_ADDRESS] validator --chain-id=peggy
```

Expected terminal output:

```bash
Password to sign with 'validator':
I[2019-12-10|13:03:11.784] Starting WSEvents                            impl=WSEvents
Subscribed to bridgebank contract at address: 0xd88159878c50e4B2b03BB701DD436e4A98D6fBe2
Subscribed to cosmosbridge contract at address: 0x8E7da79fd36d89a381CcFA2412D34E057bFFAdDe
```

The relayer is now subscribed to specific events and messages in transactions on the Ethereum and Tendermint chains, respectively. The relayer will create a prophecy claim whenever it detects a burn or lock event/message on either chain.
