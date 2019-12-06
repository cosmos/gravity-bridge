## Running Ethereum Testnet

Peggy uses Truffle for running a local Ethereum blockchain which you can deploy the contracts to for testing.

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
yarn

# Start local blockchain
yarn develop
```

### Deploy Peggy contracts to local blockchain

Next, compile and deploy Peggy's contracts to the Ethereum blockchain:

```bash
# Open a new terminal window

# Deploy contracts to local blockchain
yarn migrate

# Activate the contracts
yarn peggy:setup

# Copy contract ABI to Relayer it can subscribe to deployed contracts
yarn peggy:abi

# Get the address of Peggy's registry contract
yarn peggy:address

```
