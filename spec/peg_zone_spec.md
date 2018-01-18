# Specifcation for 2-way peg between a Tendermint chain and an Ethereum chain

#### Terminology
* The *Cosmos Peg Zone* is the blanket term for the four components
involved in Ethereum <-> Tendermint asset transfer.

## Overview
The goal of the Peg Zone is to enable the movement of assets between a Tendermint
chain and an Ethereum chain. It is designed to allow for secure and cheap
transfers of all Ethereum tokens (Ether and ERC20) as well as all Cosmos
tokens.

The Cosmos peg zone accepts and sends IBC packets. When it receives an IBC
packet it processes it and then affects a change on the Ethereum state. When
the app is informed of a state change on Ethereum it generates and sends an IBC
packet.

### Cosmos Peg Zone Components
1. a *Cosmos ABCI app*
1. a set of *signing apps* 
1. a set of Ethereum *smart contracts* 
1. a set of *relayer* processes

#### Cosmos ABCI App
The ABCI app serves as the interface to the peg zone. It communicates
using IBC packets with the hub.

It allows querying of transactions in these ways:

1. query all transactions
1. query all transactions >= a specific block height
1. query all state, including signatures, for a particular transaction

#### Signing Apps
The signing apps sign transactions using secp256k1 such that the
Ethereum smart contracts can verify them. The signing apps also have an
ethereum address, because they have an identity in the Ethereum
contract. They watch for new Ethereum-bound transactions using
the ABCI app's query functionality, and submit their signatures
back to it for replication.

#### Ethereum Smart Contracts
The smart contracts verify updates coming from the ABCI app
using the known keys of the signing apps. The smart contracts
track updates to the set of signing apps, and their associated
signatures. The smart contracts support 6 functions:

1. `lock` ETH or ERC20 tokens for use in Cosmos
1. `unlock` previously-locked (encumbered) ETH or ERC20 tokens
1. `update` signing app set signatures
1. `mint` ERC20 tokens for encumbered denominations
1. `burn` ERC20 tokens for encumbered denominations
1. `register` denomination

#### Relayer Process
The relayer process is responsible for communication
of state changes between Tendermint and Ethereum.
It is stateless, and has at-least-once delivery semantics 
from one chain to another. Every update it delivers to 
either chain is idempotent.

Generally anyone that wants the peg zone to be successful
has an incentive to run the relayer process.

It follows updates to the Ethereum chain by communicating
with a node-local Ethereum node.
When it detects locked or burned updates by the smart contracts,
it sends a signed message to the ABCI app.

# Transfer Protocols

## Sending Ethereum tokens from Ethereum to Cosmos

![Ethereum to Cosmos](./ether-to-pegzone.jpg)

1. The contract receives a `lock()` transaction with a `ERC20` token and a destination address
on the Cosmos side. It locks the received funds to the consensus of the peg
zone, logging an event that notifies the relayers.
1. The relayers process connected via RPC to an Ethereum full node, listening for `Lock` event.
1. Once the node receives a deposit to the smart contract it waits for 100 blocks (finality threshold) and then generates and signs a transactions that attests witness to the event
to which the Cosmos peg zone is listening.
1. The peg zone receives witness transactions. Until a super-majority of the voting power has witnessed an event.
1. The node then updates the state with an internal transaction to reflect that someone wants to send tokens from Ethereum. Every subsequent node adds another confirmation to the peg zone state. Every BeginBlock invocation the peg zone checks whether any incoming Ethereum transfers have reached a super-majority of confirmations and if so creates an IBC packet.

## Sending Ethereum tokens from Cosmos to Ethereum

![Cosmos to Ethereum](./pegzone-to-ether.jpg)

1. The peg zone receives an IBC packet that requests burning Ethereum tokens on the Cosmos side. The ABCI app burns the tokens and stores its serialized data, containing {address, token address, amount, nonce}(see hashing method in design_a.md). 
1. The relayers keep querying the peg zone every seconds to find a new burning event. 
1. The relayers will sign on the data using their signing app.
1. The relayers submit the signature. The peg zone waits until the signs of a super-majority of the validator set are collected.
1. The zone will notify the relayers to submit the data to the contract with its signatures. 
1. One of the relayers calls `unlock()` function on the contract.
1. The contract will release the tokens as specified, after verifying the signatures with `ecrecover()`.

## Sending Cosmos tokens from Cosmos to Ethereum

![Cosmos to Ethereum](./pegzone-to-ether.jpg)

1. the ABCI app receives an IBC packet from the hub that requests for locking Cosmos tokens and handles it according to the IBC specification. The ABCI app generates a valid Ethereum transaction containing {address, denomination, amount, nonce}, and writes it to its state. 
1. Each signing app is watching for new transactions in the ABCI state,
   and detects the new transaction. 
1. Each signing app signs the transaction using secp256k1 using a key that is known to the Ethereum smart contracts.
1. Each signing app submits their signatures back to the ABCI app for replication.
1. The relayer processes, which periodically query the ABCI app's transactions,
   see that the transaction has reached the required signature threshold.
1. The transaction is sent to the smart contracts by calling the `mint` function.
1. The smart contracts use `ecrecover` to check that it was signed by a super-majority of the validator set corresponding to the height of the transaction (this may have been updated). The smart contracts make newly minted `CosmosERC20` tokens available to the specified address in the transaction.

## Sending Cosmos tokens from Ethereum to Cosmos

![Ethereum to Cosmos](./ether-to-pegzone.jpg)

1. The contract receives a `burn()` transaction with a `CosmosERC20` token and a destination address on the Cosmos side. It burns the received funds, logging an event that notifies the relayers.
1. The relayers runs an Ethereum full node, listening for `Burn` event.  
1. The relayers waits for 100 blocks (finality threshold) after they receives the event and sign on it using the signing app.
1. The peg zone receives witness transactions. 
1. When a super-majority of the voting power has witnessed an event, the peg zone will update the internal balance state. It generates IBC packet to release the tokens to specified destination chain.
