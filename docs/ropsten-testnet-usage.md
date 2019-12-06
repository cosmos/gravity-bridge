
## Using Peggy with the Ropsten testnet

### Setup

Before you can use the Bridge with the Ropsten testnet, you'll need to add two environment variables to the configuration file at `testnet-contracts/.env`. Add MNEMONIC from your MetaMask account, this will allow you to deploy the contracts to the Ropsten testnet. Add INFURA_PROJECT_ID from your Infura account, this will allow you to start a Relayer service which listens for events on the Ropsten testnet.

Further reading:

- [MetaMask Mnemonic](https://metamask.zendesk.com/hc/en-us/articles/360015290032-How-to-Reveal-Your-Seed-Phrase)
- [Infura Project ID](https://blog.infura.io/introducing-the-infura-dashboard-8969b7ab94e7)

### Deploy contracts to Ropsten testnet

```bash
# Deploy the contracts to the Ropsten network with the --network flag
yarn migrate --network ropsten

# Activate the contracts
yarn peggy:setup

# Copy contract ABI to Relayer it can subscribe to deployed contracts
yarn peggy:abi

# Get the Registry contract's address on the Ropsten network with the --network flag
yarn peggy:address --network ropsten

```

### Start the Relayer service on Ropsten testnet

```bash
ebrelayer init ethereum wss://ropsten.infura.io/ [PEGGY_DEPLOYED_ADDRESS] validator --chain-id=peggy

```

### Lock rEth on contracts

Send funds to the deployed contracts on the Ropsten testnet:   

```bash
yarn peggy:lock --network ropsten [COSMOS_RECIPIENT_ADDRESS] eth [RWEI_AMOUNT]

```