# Design Overview

This will walk through all the details of the technical design. [`notes.md`](../notes.md) is probably a better reference
to get an overview. We will attempt to describe the entire technical design here and break out separate documents
for the details message formats, etc.

## Workflow

The high-level workflow is:

* Bootstrap Cosmos SDK chain
* Install Ethereum contract
* Relay original ERC20 tokens from ETH to Cosmos
* Relay pegged tokens from Cosmos to ETH
* Update Cosmos Validator set on ETH

The first two steps are done once, the other 3 repeated many times.

## Definitions

Words matter and we seek clarity in the terminology, so we can have clarity in our thinking and communication.
Key concepts that we mention below will be defined here:

* `Operator` - This is a person (or people) who control a Cosmos SDK validator node. This is also called `valoper` or "Validator Operator" in the Cosmos SDK staking section
* `Full Node` - This is an *Ethereum* Full Node run by an Operator
* `Validator` - This is a Cosmos SDK Validating Node (signing blocks)
* `Orchestrator` (name WIP) - This is a Rust binary controlled by an Operator that holds Cosmos SDK and Ethereum private keys used for signing transactions used to move tokens between the two chains.
* `REST server` - This is the Cosmos SDK "REST Server" that runs on Port 1317, either on the validator node or another Cosmos SDK node controlled by the Operator
* `Ethereum RPC` - This is the JSON-RPC server for the Ethereum Full Node.

The *Operator* is the key unit of trust here. Each operator is responsible for maintaining 3 secure processes:

1. Cosmos SDK Validator - signing blocks
2. Fully synced Ethereum Full Node
3. Orchestrator, which connects to the above as a client

## Bootstraping



## Relaying ETH to Cosmos


## Relaying Cosmos to ETH
