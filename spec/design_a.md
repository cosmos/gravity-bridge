# Design A

## Sending tokens from Cosmos to Ethereum

### Cosmos Peg Zone
The Cosmos peg zone receives an incoming IBC packet. It decodes it according to
the IBC specification and comes to consensus over that IBC packet. Processing
the IBC packet causes a state change, whereby a valid Ethereum transaction is
written into the state. It is represented as `go []byte`. The signing app
is listening to events from the peg zone. Once it sees that an incoming IBC
transaction was processed it takes the transaction bytes and signs them with a
valid Ethereum key that is known in the smart contract. In a second step it 
posts the signed transaction bytes as a transaction to the  Cosmos peg zone. 
The peg zone then comes to consensus over the signed bytes. The signed bytes 
are stored under original transaction bytes.

The signing app signs the transactions using secpk256k1 keys in order for 
Ethereum to be able to run `ecrecover`.

This signing app must trigger on validator set changes. Produce a "Update Validators"
transaction and then sign it with every validators private key.

### Relayer Process
The relayer process takes the signed transaction bytes and posts them to the 
set of Ethereum smart contracts.

### Ethereum Smart Contracts
The contracts receive the signed transaction data in a normal smart contract
call. They seen decode that data and run `ecrecover` to check that it was 
correctly signed by a super-majority of the current validator set. It then 
releases the funds. 


## Sending tokens from Ethereum to Cosmos

### Ethereum Smart Contracts
The contract receives a transaction with a token and a destination address
on the Cosmos side. It locks the received funds to the consensus of the peg
zone.

### Relayer Process
In this case the relayer process connected via RPC to an Ethereum full node. Once the node receives
a deposit to the smart contract it waits for 100 blocks (finality threshold)
and then generates and signs a transactions that attests witness to the event
to which the Cosmos peg zone is listening.

### Cosmos Peg Zone
The peg zone receixves witness transactions. When a super-majority of the voting power has witnessed an event,
the node then updates the state with an internal transaction to
reflect that someone wants to send tokens from Ethereum. Every subsequent
node adds another confirmation to the peg zone state. Every BeginBlock
invocation the peg zone checks whether any incoming Ethereum transfers have 
reached a super-majority of confirmations and if so creates an IBC packet.
