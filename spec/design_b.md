# Design B

## Sending tokens from Cosmos to Ethereum

### Cosmos Peg Zone
The peg zone receives an incoming IBC packets and decodes it. It then verifies the correctness of the IBC packet and write an easily parsed by EVM datastructure into the state tree indicating what asset, address, etc the tokens need to be sent to.

### Relayer Process
The relayer process takes a recent block header and a merkle proof to the
transfer and posts it to the set of smart contracts on Ethereum.

### Ethereum Smart Contracts
The contracts receive a block header and verify that it originates from the
correct validator set using `ecrecover`. Then they verify the merkle proof 
and eventually release funds to the destination address.


## Sending tokens from Ethereum to Cosmos
It is the same as design A.


