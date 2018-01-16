# ETGate specs

This is the documentation for current implementation of ETGate. Actual implementation will be Design A/B

## Sending tokens from Cosmos to Ethereum

**ETGate does not support sending Cosmos native tokens(atoms/photons) to Ethereum**

### Cosmos Peg Zone

When the pegzone receives the IBC packet for withdrawal, it stores the packet in the state.

### Relayer Process

When the relayers see a new WithdrawTx, they submit it to the contract with its IAVL proof. If the contract did not receive the header that is needed for merkle proving, they also submit it with validators' signatures.

### Ethereum Smart Contracts

Once the contract receives the withdrawal data with its IAVL proof, it proves the data with its already stored pegzone header. The code is just a translation of original tendermint/iavl. Then, the contract releases its locked contract 

## Sending tokens from Ethereum to Cosmos

### Ethereum Smart Contracts

When the contract receives funds, it generates Deposit() event.

### Relayer Process

The validators submit Ethereum headers to the pegzone. When a header gets more than +2/3 of the validators' sign, the DELAYth ancestor of the header will be finalized. The relayers listen for new finalization and takes all Deposit() event contained in the finalized from ethclient. Then they rlp-encode the event and submit it to the pegzone with its merkle proof. 

### Cosmos Peg Zone

The pegzone verifies DepositTx with its merkle proof with corresponding Ethereum header that is already finalized. If it isn't finalized the pegzone rejects the tx. Then it the pegzone generates IBC packet that goes to the destination chain.
