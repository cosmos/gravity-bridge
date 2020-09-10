# Design Overview

This will walk through all the details of the technical design. [`notes.md`](../notes.md) is probably a better reference
to get an overview. We will attempt to describe the entire technical design here and break out separate documents
for the details message formats, etc.

## Workflow

The high-level workflow is:

Activation Steps:

* Bootstrap Cosmos SDK chain
* Install Ethereum contract

Token Transfer Steps:

* Transfer original ERC20 tokens from ETH to Cosmos
* Transfer pegged tokens from Cosmos to ETH
* Update Cosmos Validator set on ETH

The first two steps are done once, the other 3 repeated many times.

## Definitions

Words matter and we seek clarity in the terminology, so we can have clarity in our thinking and communication.
Key concepts that we mention below will be defined here:

* `Operator` - This is a person (or people) who control a Cosmos SDK validator node. This is also called `valoper` or "Validator Operator" in the Cosmos SDK staking section
* `Full Node` - This is an *Ethereum* Full Node run by an Operator
* `Validator` - This is a Cosmos SDK Validating Node (signing blocks)
* `Orchestrator` (name WIP) - This is a Rust binary controlled by an Operator that holds Cosmos SDK and Ethereum private keys used for signing transactions used to move tokens between the two chains. (This was called peggy module in other docs, which is confusing as the module runs on the node... I would love a better name than Orchestrator)
* `REST server` - This is the Cosmos SDK "REST Server" that runs on Port 1317, either on the validator node or another Cosmos SDK node controlled by the Operator
* `Ethereum RPC` - This is the JSON-RPC server for the Ethereum Full Node.
* `Validator Set` - The set of validators on the Cosmos SDK chain, along with their respective voting power. These are ed25519 public keys used to sign tendermint blocks.
* `Orchestrator Set` - The validator set mapped over to Cosmos SDK Orchestrator keys. This is used on the Cosmos SDK chain to authorize messages from Ethereum. Note that we can use the current validator set, but we use Orchestrator keys (not tendermint consensus keys)
* `MultiSig Set` - The set of Ethereum keys along with respective voting power. This is based on the validator set and mapped over the registered keys. However, as this is a different chain, this is mirrored with a delay.

The *Operator* is the key unit of trust here. Each operator is responsible for maintaining 3 secure processes:

1. Cosmos SDK Validator - signing blocks
1. Fully synced Ethereum Full Node
1. Orchestrator, which connects to the above as a client

## Security Concerns

The **Validator Set** is the actual set of keys with stake behind them, which are slashed for double-signs or other 
misbehavior. We typically consider the security of a chain to be the security of a *Validator Set*. This varies on
each chain, but is our gold standard. Even IBC offers no more security than the minimum of both involved Validator Sets.

The **Orchestrator Set** is another Cosmos SDK client key associated with the same validator set. We can add slashing
conditions to any mis-signed message by the Orchestrator Set and be able to provide the same security as the
*Valiator Set*, just a different module detecting evidence of malice and deciding how much to slash. If we can prove a
transaction signed by any member of the *Orchestrator Set* was illegal or malicious, then we can slash on the Cosmos chain
side an potentially provide 100% of the security of the Validator Set. Note that this also has access to the 3 week unbonding
period to allow evidence to slash even if they immediately unbond.

The **MultiSig Set** is a (possibly aged) mirror of the *Validator Set* but with Ethereum keys, and stored on the Ethereum
contract. If we ensure the *MultiSig Set* is updated much more often than the unbonding period (eg at least once per week),
then we can guarantee that all members of the *MultiSig Set* have slashable atoms for misbehavior. However, in some extreme
cases of stake shifting, the *MultiSig Set* and *Validator/Orchestrator Set* could get quite far apart, meaning there is
many of the members in the *MultiSig Set* are no longer active validators and may not bother to transfer Eth messages. 
Thus, to avoid censorship attacks/inactivity, we should also update this everytime there is a significant change
in the Validator Set (eg. > 3-5%). If we maintain those two conditions, the MultiSig Set should offer a similar level of
security as the Validator Set.

There are now 3 conditions that can be slashed for any validator: Double-signing a block with the tendermint key from the
**Validator Set**, signing an invalid/malicious message with the Cosmos SDK key from the **Orchestrator Set**, or
signing an invalid/malicious Ethereum message with the key from the **MultiSig Set**. If all conditions of misbehavior can
be attributed to a signature from one of these sets, and proven **on the Cosmos chain**, then we can argue that Peggy offers
a security level equal to the minimum of the Peg-Zone Validator Set, or reorganizing the Ethereum Chain 50 blocks.
And provide a security equivalent to or greater than IBC.

## Bootstraping

This is based on [Installing Peggy on a live cosmos chain](notes.md#installing-peggy-on-a-live-cosmos-chain).
We also assume the act of upgrading the Cosmos-based binary to have peggy module is already complete,
as approaches to that are discussed in many other places. Here we focus on the *activation* step.

Set up and register orchestrators:

1. Every *Operator* must initialize an *Orchestrator* on a secure computer
1. Each *Operator* signs a `MsgRegisterOrchestrator` that ties the orchestrator's Cosmos and Ethereum keys to the validator (and makes them liable for slashing conditions)
1. A `StartPeggy` message is triggered, either as a governance proposal or by any validator. If over X% (70? 80? 90?) of the validator power has registered, we create the original *MultiSig Set* and store that in the Cosmos chain.

Upload and connect Ethereum contract:

1. Someone (anyone?) uploads a Peggy contract on Ethereum referencing the *MultiSig Set* stored on the Cosmos chain
1. This Peggy address is proposed is created on the SDK chain (governance vote? other vote?)
1. All orchestrators check the contract has official ethereum bytecode and proper *MultiSig Set*.
1. If the orchestrator approves the contract, it will submit a message to the Cosmos Chain approving this address and signing an "activate peggy" ethereum message
1. Once X% of the *Orchestrator Set* has signed SDK messages approving the peggy, the Cosmos Chain will store this as the official Peggy address. If there are multiple proposals, the first one to hit 70% is stored and other proposals will not succeed.
1. Once the official Peggy address has been stored on the Cosmos Chain, someone (anyone?) can activate the Peggy contract by submitting all the Ethereum signed messages that contain > X% of the *MultiSig Set* contained in the Peggy contract

At this point, we know we have a contract on Ethereum with the proper *MultiSig Set*, that > X% of the *Orchestrator Set* is online and agrees with this contract, and that the Cosmos chain has stored this contract address. Only then can we begin to accept transactions to transfer tokens

Note: X% is some security factor for bootstrapping. 67% is sufficient to release, but we don't want to start until there is a margin of error online (not to fall off with a small change of voting power). This may be 70, 80, 90, or even 95% depending on how much assurances we want that all *Orchestrators* are operational before starting.

## Relaying ETH to Cosmos

**TODO**

## Relaying Cosmos to ETH

**TODO**