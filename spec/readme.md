# 2-way permissionless pegzones

## Introduction

This document describes the process of building 2-way pegzones between two
separate blockchains. Pegzones are needed if one and/or both chains only have
probabalistic finality. An example of probalistic finality is Bitcoin or
Ethereum. If both chains have true finality transferring value between both is
achieved by implementing the "Inter Blockchain Communiciation Protocol" or IBC for
short. Blockchain engines such as [Tendermint Core](https://github.com/tendermint/tendermint)
support true finality due to the use of [Tendermint Consensus]().

The reason why pegzones are needed in the former case is that in order to 
separate the global state into two blockchains instead of just one we need a 
finality guarantee after which neither of the chains is allowed to revert any
transactions. The pegzones main duty is to guarantee this finality even though
the underlying chain does not offer it.

Even though this document explains a general concept it choses to use Ethereum
and Cosmos as practical examples and shows how to build a pegzone between them.

In the following paragraphs we describe how to build a pegzone between Ethereum
and Cosmos, which is a Tendermint-based chain. It provides pegged assets on
either side that have the same security properties as the pegzone. Any failure
in the pegzone results in massive issues, which have to resolved through 
governance. However this design guarantees that no assets are locked forever
should the consensus of the pegzone fail.

Concretely this specification describes a pegzone which allows the movement of
ERC20 and Ether tokens from Ethereum to Cosmos and of all Cosmos native assets
to Ethereum, where they are represented as ERC20 tokens.

This design does not describe how tokens move into the Cosmos pegzone. It
assumes that the pegzone itself establishes connectivity over IBC to a hub or
other chains.


## Structure

* [Overview of the design and its components](##overview-of-the-design-and-its-components)
* [Ethereum smart contract component](##smart-contract-component)
* [Witness component](##witness-component)
* [Pegzone blockchain](##pegzone-blockchain)
* [Signer component](##signer-component)
* [Relayer component](##relayer-component)
* Conclusion


## Overview of the design and its components

The pegzone is split into five logical components, the Ethereum smart contract,
the witness, the pegzone, the signer and the relayer. The last four can all be
implemented within the some software and are only split for convenience of 
understanding.

The Ethereum smart contract acts as the custodian of assets, whose origin is
Ethereum and as issuer of assets, whose origin is on Cosmos.

The witness attests to the state of Ethereum and implements the finality
threshold over the non-finality chain. It follows and validates the Ethereum
consensus and attests to events happening within Ethereum by posting a 
message to the pegzone.

The pegzone is a blockchain that allows users to store balances of assets and
interact with them via transactions. It allows users to send assets to other
users, receive and send assets from/to Ethereum and receive/send from/to other
blockchains via IBC.

The signer signs messages in a form that the Ethereum smart contract can
validate. Those messages contain information in order to release tokens on
Ethereum.

The relayer is responsible for taking signed messages and posting them to the
Ethereum smart contract.

### Flow

Now we will describe the typical flow of an asset transfer that originates on 
Ethereum to Cosmos and back. Furthermore we will describe the transfer of an 
assets that originates on Cosmos to Ethereum and back.

#### Ethereum native token
1. Alice sends a transaction with 10 Ether and a destination address to the
Ethereum smart contract. The destination address is Alice's on the Pegzone 
blockchain.

2. The witness component sees the transaction on the Ethereum smart contract
and after the finality threshold attests to the fact by submitting a WitnessTx
to the pegzone. 

3. The pegzone receives the WitnessTx and credits the destination address from
step 1 with 10 Ether.

4. Alice receives the 10 CEther on the pegzone. 

5. Alice now sends 4 CEther to Bob on the pegzone.

6. Bob now sends those 4 CEther to Ethereum by submitting a RedeemTx. 

7. The signer component sees the RedeemTx and generates a signature for it. It
then posts that signature to the pegzone via a SignTx.

8. The relayer component sees the SignTxs. After >2/3 of the validator by 
voting power have submitted a SignTx it posts all SignTx and the original
message to the Ethereum smart contract. 

Alice started with 10 Ether on Ethereum. She sent all of them to Cosmos. There
she sent 4 CEther to Bob. Bob then redeemed those 4 CEther back to Ethereum.
The end result is that Alice holds 6 CEther and 0 Ether and Bob holds 4 Ether 
and 0 CEther.


#### Cosmos native token
1. Alice owns 10 Photons. She sends a RedeemTx with 10 Photons and a destination 
address to the pegzone. The pegzone takes custody of those 10 Photons and Alice
still holds 0 Photons.

2. The signer sees that transaction and signs a message that can be interpreted
by the Ethereum smart contract. It then posts a SignTx to the pegzone.

3. The relayer sees the SignTxs and once >2/3 of the validators have posted
them invokes a function on the Ethereum smart contract.

4. The Ethereum smart contract validates that the data was correctly signed
by a super-majority of the validator set. It then generates an ERC20 contract
for Photons unless it already has generated an ERC20 contract for Photons
previously. It then credits the destination address, which Alice controls,
with 6 Photons.

5. Alice then sends 10 EPhotons to Bob on Ethereum. 

6. Bob transfers 4 EPhotons to the Ethereum smart contract with a destination
address. The contract burns those tokens and raises an event.

7. The witness process attests to the burning of 4 EPhotons and sends a 
WitnessTx to the pegzone.

8. The pegzone receives the WitnessTx and credits Bob with 4 Photons.

Alice started with 10 Photons on Cosmos. She sent all of them to herself on
Ethereum. On Ethereum she sent 4 EPhotons to Bob. Bob sent those EPhotons back
to Cosmos. The end result is that Alice holds 6 EPhotons and 0 Photons and Bob
holds 4 Photons and 0 EPhotons.


## Ethereum smart contract component

The smart contracts verify updates coming from the pegzone
using the known keys of the signing apps. The smart contracts
track updates to the set of signer components, and their associated
signatures. The smart contracts supports 6 functions:

* `lock` ETH or ERC20 tokens for use in Cosmos
* `unlock` previously-locked (encumbered) ETH or ERC20 tokens
* `update` signing app set signatures
* `mint` ERC20 tokens for encumbered denominations
* `burn` ERC20 tokens for encumbered denominations
* `register` denomination


## Witness component

The witness component runs a full Ethereum node. When it sees an event from
the smart contract that it tracks it sends a WitnessTx to the pegzone.

## Pegzone blockchain

The pegzone is a normal blockchain that keeps user accounts. It allows for
users to send assets to each other. 

It allows querying of transactions in these ways:

* query all transactions
* query all transactions >= a specific block height
* query all state, including signatures, for a particular transaction

## Signer component

The signing apps sign transactions using secp256k1 such that the
Ethereum smart contracts can verify them. The signing apps also have an
ethereum address, because they have an identity in the Ethereum
contract. They watch for new Ethereum-bound transactions using
the ABCI app's query functionality, and submit their signatures
back to it for replication.

## Relayer component

The relayer process is responsible for communication
of state changes between Tendermint and Ethereum.
It is stateless, and has at-least-once delivery semantics 
from one chain to another. Every update it delivers to 
either chain is idempotent.

Generally anyone that wants the peg zone to be successful
has an incentive to run the relayer process.



# Transfer Protocols

## Sending Ethereum tokens from Ethereum to Cosmos

![Ethereum to Cosmos](./ether-to-pegzone.jpg)

1. The contract receives a `lock` transaction with a `ERC20` token and a destination address
on the Cosmos side. It locks the received funds to the consensus of the peg
zone, logging an event that notifies the relayers.
1. The relayers process connected via RPC to an Ethereum full node, listening for `Lock` event.
1. Once the node receives a deposit to the smart contract it waits for 100 blocks (finality threshold) and then generates and signs a `SignWitnessMsg` that attests witness to the event
to which the Cosmos peg zone is listening.
1. The peg zone receives witness transactions until a super-majority of the voting power has witnessed an event. Every BeginBlock invocation the peg zone checks whether any incoming Ethereum transfers have reached a super-majority of confirmations.
1. The node then updates the state with an internal transaction to reflect that someone wants to send tokens from Ethereum and generates `IBCWitness` to mint the tokens to specified destination chain.

## Sending Ethereum tokens from Cosmos to Ethereum

![Cosmos to Ethereum](./pegzone-to-ether.jpg)

1. The ABCI app receives an `IBCSignature` that requests for burning Ethereum tokens and handles it according to the IBC specification. The ABCI app generates a valid Ethereum transaction containing {address, token address, amount, nonce}, and writes it to its state.
1. Each signing app is watching for new transactions in the ABCI state, and detects the new transaction. 
1. Each signing app signs the transaction using secp256k1 using a key that is known to the Ethereum smart contracts.
1. Each signing app submits their signatures back to the ABCI app as `SignSignatureMsg` for replication.
1. The relayer processes, which periodically query the ABCI app's transactions,
   see that the transaction has reached the required signature threshold. 
1. One of the relayers send the transaction to the smart contract by calling the `unlock` function
1. The smart contracts use `ecrecover` to check that it was signed by a super-majority of the validator set corresponding to the height of the transaction (this may have been updated). The smart contracts release the token as specified in the transaction making it available to the destination address.

## Sending Cosmos tokens from Cosmos to Ethereum

![Cosmos to Ethereum](./pegzone-to-ether.jpg)

1. the ABCI app receives an `IBCSignature` from the hub that requests for locking Cosmos tokens and handles it according to the IBC specification. The ABCI app generates a valid Ethereum transaction containing {address, denomination, amount, nonce}, and writes it to its state. 
1. Each signing app is watching for new transactions in the ABCI state,
   and detects the new transaction. 
1. Each signing app signs the transaction using secp256k1 using a key that is known to the Ethereum smart contracts.
1. Each signing app submits their signatures back to the ABCI app as `SignSignatureMsg` for replication.
1. The relayer processes, which periodically query the ABCI app's transactions,
   see that the transaction has reached the required signature threshold.
1. One of the relayers send the transaction to the smart contract by calling the `mint` function.
1. The smart contracts use `ecrecover` to check that it was signed by a super-majority of the validator set corresponding to the height of the transaction (this may have been updated). The smart contracts make newly minted `CosmosERC20` tokens available to the specified address in the transaction.

## Sending Cosmos tokens from Ethereum to Cosmos

![Ethereum to Cosmos](./ether-to-pegzone.jpg)

1. The contract receives a `burn` transaction with a `CosmosERC20` token and a destination address on the Cosmos side. It burns the received funds, logging an event that notifies the relayers.
1. The relayers process conttected via RPC to an Ethereum full node, listening for `Burn` event.  
1. Once the node receives a deposit to the smart contract it waits for 100 blocks (finality threshold) and then generates and signs a `SignWitnessMsg` that attests witness to the event
1. The peg zone receives witness transactions until a super-majority of the voting power has witnessed an event. Every BeginBlock invocation the peg zone checks whether any incoming Ethereum transfers have reached a super-majority of confirmations.
1. The node then updates the state with an internal transaction to reflect that someone wants to send tokens from Ethereum and generates `IBCWitness` to release the tokens to specified destination chain.

# API

## ABCI app

### Common Types

#### Witness{Nonce() (uint64), Chain() ([]byte)}
Interface that Lock{} and Burn{} implements. All Witnesses are generated by validators multisig and generates IBC packet that goes to Chain().

#### Signature{}
Interface that Unlock{}, Mint{}, Register{} and Update{} implements. All signature packets are generated by IBC packets and stored in ABCI app storage with nonce.

#### Update{Validators []Validator}
#### Lock{To []byte, Value uint64, Token common.Address, Chain []byte, Nonce uint64}
#### Unlock{To common.Address, Value uint64, Token common.Address, Chain []byte}
#### Mint{To common.Address, Value uint64, Token []byte, Chain []byte}
#### Burn{To []byte, Value uint64, Token []byte, Chain []byte, Nonce uint64}
#### Register{Denom string}

### IBC Packet Types

#### IBCWitness{Witness}

The zones that uses the pegzone will receives and handles IBCWitness packet.

#### IBCSignature{Signature}

The zones that uses the pegzone will sends IBCSignature packet.

### Msg Types

#### SignWitnessMsg{Witness} 

Used for voting on witness packets.

#### SignSignatureMsg{Signature, Nonce uint64, Sig []byte}

Used for add signs on signature packets which will be submitted on the Ethereum contract later.

* `Sig` must be length of 65 and concatenated value of `v`, `r`, and `s` of the sender's signature. 

### Querying Functions

#### SignaturePacketNonce() (uint64)

Returns last signature packet nonce. 

#### SignaturePacketByNonce(nonce uint64) (Signature)

Returns signature packet by its nonce. Returns nil if it dosen't exist.

#### SignatureSignByNonce(nonce uint64) ([]uint16, [][]byte)

Returns signature packet's validator signatures by its nonce. (Signed validators by their indexes, Array of signatures).

#### SignatureSignSatisfied(nonce uint64) (bool)

True if the signature packet has enough signatures.

#### WitnessPacketNonce() (uint64)

Returns last witness packet nonce.

#### WitnessSignByNonce(nonce uint64) ([]uint16)

Returns the indexes of the validators who signed on the witness.

## Signing app

## Ethereum Smart Contracts

### External Entry Points

#### update(address[] newAddress, uint64[] newPower, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Updates validator set. Called by the relayers.

* hash value for `ecrecover` is calculated as: 
```
byte(0) + newAddress.length.PutUint256() + newAddress[0].Bytes() + ... + newPower[0].PutUint64() + ...
```

#### lock(bytes to, uint64 value, address token, bytes chain) payable

Locks Ethereum user's ethers/ERC20s in the contract and loggs an event. Called by the users.

* `token` being `0x0` means ethereum; in this case `msg.value` must be same with `value`
* `event Lock(bytes to, uint64 value, address token, bytes chain, uint64 nonce)` is logged, seen by the relayers

#### unlock(address to, uint64 value, address token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Unlocks Ethereum tokens according to the information from the pegzone. Called by the relayers.

* transfer tokens to `to`
* hash value for `ecrecover` is calculated as:
```
byte(1) + to.Bytes() + value.PutUint64() + chain.length.PutUint256() + chain
```
#### mint(address to, uint64 value, bytes token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Mints 1:1 backed credit for atoms/photons. Called by the relayers.

* `token` has to be `register`ed before the call
* transfer minted tokens to `to`
* hash value for `ecrecover` is calculated as:
```
byte(2) + to.Bytes() + value.PutUint64() + token.length.PutUint256() + token + chain.length.PutUint256() + chain
```

#### burn(bytes to, uint64 value, bytes token, bytes chain)

Burns credit for atoms/photons and loggs an event. Called by the users.

* `event Burn(bytes to, uint64 value, bytes token, bytes chain, uint64 nonce)` is logged, seen by the relayers

#### register(string name, address token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s)

Registers new Cosmos token name with its CosmosERC20 address. Called by the relayers.

* deploys new CosmosERC20 contract and stores it in a mapping

## Relayer Process

### Rejected alternate design

The team reccomends this design over an alternate design we called Design B.
Design B minimizes the role of the signing apps and places more functionality in
tendermint core and the abci app.

Here the ABCI app would contain a ethereum light client implementation and the relayer 
would send light client proofs and block headers from the peg zone contract.

The pegzone contract would release funds based on light client proofs from tendermint.

The biggest changes we realized we would need.

1. Tendermint header serialization that is easy to for solidity to parse. Most likely bitcoin style fixed byte structure.
1. Secp256k1 signatures in Tendermint consensus.

This design seems cleaner in some ways but more difficult to MVP.


