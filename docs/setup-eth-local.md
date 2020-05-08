## Setup Peggy locally

Peggy uses Truffle for running a local Ethereum blockchain which you can deploy the contracts to for testing.

Note: Truffle is currently incompatible with node v14 because of a bug in ganache - see [here](trufflesuite/ganache-cli#732). Until the issue is resolved, we recommend any prior stable version of node (such as v10.16.3).

### Setup

In order for Peggy to process cross-chain asset transfers, the Relayer service must be run by a set of validators. Before validators participate in asset transfers, they must set up the appropriate configuration files with the following commands:

```bash
cd testnet-contracts/

# Create .env with environment variables required for contract deployment
cp .env.example .env
```

### Start local blockchain

```bash
# Open a new terminal window

# Download dependencies
yarn # or npm i

# Start local blockchain
yarn develop # or npm run develop
```

### Set up Peggy contracts

Next, compile and deploy Peggy's contracts to the Ethereum blockchain:

```bash
# Open a new terminal window

# Deploy and set up contracts, then mint ERC20 TEST tokens and approve some to bank contract
yarn peggy:all

# Take note of Peggy's BridgeRegistry contract address and the ERC20 TEST contract address,
# you'll need them in the next step.
```

To set up the relayer, go to [the next step](./setup-relayer.md).
