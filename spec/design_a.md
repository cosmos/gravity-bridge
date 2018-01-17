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

## API

### Ethereum Smart Contract

#### about signature variables(common for update(), unlock(), mint(), and register())

* signatures are flattened as `uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s`
* `idxs` must be ordered increasing and not repeated
* `v[i]`, `r[i]`, and `s[i]` must be `ecrecover()` argument for `idxs[i]`th validator 
* hashing method for data is not decided yet

#### update(address[] newAddress, uint64[] newPower, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Updates validator set. Called by the relayers.

#### lock(bytes to, uint64 value, address token, bytes chain) payable

Locks Ethereum user's ethers/ERC20s in the contract and loggs an event. Called by the users.

* `token` being `0x0` means ethereum; in this case `msg.value` must be same with `value`
* `event Lock(bytes to, uint64 value, address token, bytes chain, uint64 nonce)` is logged, seen by the relayers

#### unlock(address to, uint64 value, address token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Unlocks Ethereum tokens according to the incoming information from the pegzone. Called by the relayers.

* transfer tokens to `to`

#### mint(address to, uint64 value, bytes token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Mints 1:1 backed credit for atoms/photons. Called by the relayers.

* `token` has to be `register`ed before the call
* transfer minted tokens to `to`

#### burn(bytes to, uint64 value, bytes token, bytes chain)

Burns credit for atoms/photons and loggs an event. Called by the users.

* `event Burn(bytes to, uint64 value, bytes token, bytes chain, uint64 nonce)` is logged, seen by the relayers

#### register(string name, address token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Registers new Cosmos token name with its CosmosERC20 address. Called by the relayers.

* deploys new CosmosERC20 contract and stores it in a mapping

### Relayer

### Cosmos Peg Zone
