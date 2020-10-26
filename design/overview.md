# Design Overview

This will walk through all the details of the technical design. [`notes.md`](../notes.md) is probably a better reference
to get an overview. We will attempt to describe the entire technical design here and break out separate documents
for the details message formats, etc.

## Workflow

The high-level workflow is:

Activation Steps:

- Bootstrap Cosmos SDK chain
- Install Ethereum contract

Token Transfer Steps:

- Transfer original ERC20 tokens from ETH to Cosmos
- Transfer pegged tokens from Cosmos to ETH
- Update Cosmos Validator set on ETH

The first two steps are done once, the other 3 repeated many times.

## Definitions

Words matter and we seek clarity in the terminology, so we can have clarity in our thinking and communication.
Key concepts that we mention below will be defined here:

- `Operator` - This is a person (or people) who control a Cosmos SDK validator node. This is also called `valoper` or "Validator Operator" in the Cosmos SDK staking section
- `Full Node` - This is an _Ethereum_ Full Node run by an Operator
- `Validator` - This is a Cosmos SDK Validating Node (signing blocks)
- `Eth Signer` (name WIP) - This is a Rust binary controlled by an Operator that holds Cosmos SDK and Ethereum private keys used for signing transactions used to move tokens between the two chains.
- `Relayer` - This is a type of node that submits updates to the Peggy contract on Ethereum. It earns fees from the transactions in a batch.
- `REST server` - This is the Cosmos SDK "REST Server" that runs on Port 1317, either on the validator node or another Cosmos SDK node controlled by the Operator
- `Ethereum RPC` - This is the JSON-RPC server for the Ethereum Full Node.
- `Validator Set` - The set of validators on the Cosmos SDK chain, along with their respective voting power. These are ed25519 public keys used to sign tendermint blocks.
- `Peggy Tx pool` - Is a transaction pool that exists in the chain store of Cosmos -> Ethereum transactions waiting to be placed into a transaction batch
- `Transaction batch` - A transaction batch is a set of Ethereum transactions to be sent from the Peggy Ethereum contract at the same time. This helps reduce the costs of submitting a batch. Batches have a maximum size (currently around 100 transactions) and are only involved in the Cosmos -> Ethereum flow
- `Peggy Batch pool` - Is a transaction pool like strucutre that exists in the chains to store, seperate from the `Pegg Tx pool` it stores transactions that have been placed in batches that are in the process of being signed or being submitted by the `Orchestrator Set`
- `EthBlockDelay` - Is a agreed upon number of Ethereum blocks all oracle attestations are delayed by. No `Orchestrator` will attest to have seen an event occur on Ethereum until this number of blocks has elapsed as denoted by their trusted Ethereum full node. This should prevent short forks form causing disagreements on the Cosmos side. The current value being consdiered is 50 blocks.
- `Observed` - events on Ethereum are considered `Observed` when the `Eth Signers` of 66% of the active Cosmos validator set during a given block has submitted an oracle message attesting to seeing the event.
- `Validator set delta` - This is a term for the difference between the validator set currently in the Peggy Ethereum contract and the actual validator set on the Cosmos chain. Since the validator set may change every single block there is essentially guaranteed to be some nonzero `Validator set delta` at any given time.
- `Allowed validator set delta` - This is the maximum allowed `Validator set delta` this parameter is used to determine if the Peggy contract in MsgProposePeggyContract has a validator set 'close enough' to accept. It is also used to determine when validator set updates need to be sent. This is decided by a governance vote _before_ MsgProposePeggyContract can be sent.
- `Peggy ID` - This is a random 32 byte value required to be included in all Peggy signatures for a particular contract instance. It is passed into the contract constructor on Ethereum and used to prevent signature reuse when contracts may share a validator set or subsets of a validator set. This is also set by a governance vote _before_ MsgProposePeggyContract can be sent.
- `Peggy contract code hash` - This is the code hash of a known good version of the Peggy contract solidity code. It will be used to verify exactly which version of the bridge will be deployed.
- `Start Threshold` - This is the percentage of total voting power that must be online and participating in Peggy operations before a bridge can start operating.
- `Claim` - an Ethereum event signed and submitted to cosmos by a single `Orchestrator` instance 
- `Attestation` - aggregate of claims that eventually becomes `observed` by all orchestrators
- `Voucher` - represent a bridged ETH token on the Cosmos side. Their denom is has a `peggy` prefix and a hash that is build from contract address and contract token. The denom is considered unique within the system.
- `Counterpart` - to a `Voucher` is the locked ETH token in the contract
  
The _Operator_ is the key unit of trust here. Each operator is responsible for maintaining 3 secure processes:

1. Cosmos SDK Validator - signing blocks
1. Fully synced Ethereum Full Node
1. `Eth Signer`, which signs things with the `Operator's` Eth keys

## Security Concerns

The **Validator Set** is the actual set of keys with stake behind them, which are slashed for double-signs or other
misbehavior. We typically consider the security of a chain to be the security of a _Validator Set_. This varies on
each chain, but is our gold standard. Even IBC offers no more security than the minimum of both involved Validator Sets.

The **Eth Signer** is a binary run alongside the main Cosmos daemon (gaiad or equivalent) by the validator set. It exists purely as a matter of code organization and is in charge of signing Ethereum transactions, as well as observing events on Ethereum and bringing them into the Cosmos state. It signs transactions bound for Ethereum with an Ethereum key, and signs over events coming from Ethereum with a Cosmos SDK key. We can add slashing conditions to any mis-signed message by any _Eth Signer_ run by the _Validator Set_ and be able to provide the same security as the _Valiator Set_, just a different module detecting evidence of malice and deciding how much to slash. If we can prove a transaction signed by any _Eth Signer_ of the _Validator Set_ was illegal or malicious, then we can slash on the Cosmos chain side and potentially provide 100% of the security of the _Validator Set_. Note that this also has access to the 3 week unbonding
period to allow evidence to slash even if they immediately unbond.


The **MultiSig Set** is a (possibly aged) mirror of the _Validator Set_ but with Ethereum keys, and stored on the Ethereum
contract. If we ensure the _MultiSig Set_ is updated much more often than the unbonding period (eg at least once per week),
then we can guarantee that all members of the _MultiSig Set_ have slashable atoms for misbehavior. However, in some extreme
cases of stake shifting, the _MultiSig Set_ and _Validator Set_ could get quite far apart, meaning there is
many of the members in the _MultiSig Set_ are no longer active validators and may not bother to transfer Eth messages.
Thus, to avoid censorship attacks/inactivity, we should also update this everytime there is a significant change
in the Validator Set (eg. > 3-5%). If we maintain those two conditions, the MultiSig Set should offer a similar level of
security as the Validator Set.

There are now 3 conditions that can be slashed for any validator: Double-signing a block with the tendermint key from the
**Validator Set**, signing an invalid/malicious event from Ethereum with the Cosmos SDK key held by its _Eth Signer_, or
signing an invalid/malicious Ethereum transaction with the Ethereum key held by its _Eth Signer_. If all conditions of misbehavior can
be attributed to a signature from one of these sets, and proven **on the Cosmos chain**, then we can argue that Peggy offers
a security level equal to the minimum of the Peg-Zone Validator Set, or reorganizing the Ethereum Chain 50 blocks.
And provide a security equivalent to or greater than IBC.

## Bootstrapping

We assume the act of upgrading the Cosmos-based binary to have peggy module is already complete,
as approaches to that are discussed in many other places. Here we focus on the _activation_ step.

1. Each `Operator` generates an Ethereum and Cosmos private key for their `EthSigner`. These addresses are signed and submitted by the Operators valoper key in a MsgRegisterEthSigner. The `EthSigner` is now free to use these delegated keys for all Peggy messages.
1. A governance vote is held on bridge parameters including `Peggy ID`, `Allowed validator set delta`, `start threshold`, and `Peggy contract code hash`
1. Anyone deploys a Peggy contract using a known codehash and the current validator set of the Cosmos zone to an Ethereum compatible blockchain.
1. Each `Operator` may or may not configure their `Eth Signer` with the above Peggy contract address
1. If configured with an address the `Eth Signer` checks the provided address. If the contract passes validation the `Eth Signer` signs and submits a MsgProposePeggyContract. Validation is defined as finding the correct `Peggy contract code hash` and a validator set matching the current set within `Allowed validator set delta`.
1. A contract address is considered adopted when voting power exceeding the `start threshold` has sent a MsgProposePeggyContract with the same Ethereum address.
1. Because validator sets change quickly, `Eth Signers` not configured with a contract address observe the Cosmos blockchain for submissions. When an address is submitted they validate it and approve it themselves if it passes. This results in a workflow where once a valid contract is proposed it will be ratified in a matter of a few seconds.
1. It is possible for the adoption process to fail if a race condition is intentionally created resulting in less than 66% of the validator power approving more than one valid Peggy Ethereum contract. In this case the Orchestrator will check the contract address with the majority of the power (or at random in the case of a perfect tie) and switch it's vote. This leaves only the possible edge case of >33% of `Operators` intentionally selecting a different contract address. This would be a consensus failure and the bridge can not progress.
1. The bridge ratification process is complete, the contract address is now placed in the store to be referenced and other operations are allowed to move forward.

At this point, we know we have a contract on Ethereum with the proper _MultiSig Set_, that > `start threshold` of the _Orchestrator Set_ is online and agrees with this contract, and that the Cosmos chain has stored this contract address. Only then can we begin to accept transactions to transfer tokens

Note: `start threshold` is some security factor for bootstrapping. 67% is sufficient to release, but we don't want to start until there is a margin of error online (not to fall off with a small change of voting power). This may be 70, 80, 90, or even 95% depending on how much assurances we want that all _Orchestrators_ are operational before starting.

## Relaying ETH to Cosmos

**TODO**

## Relaying Cosmos to ETH

- A user sends a MsgSendToEth when they want to transfer tokens across to Ethereum. This debits the tokens from their account, and places a transaction in the `Peggy Tx Pool`
- Someone (permissionlessly) sends a MsgRequestBatch, this produces a new `Transaction batch` in the `Peggy Batch pool`. The creation of this batch occurs in ComsosSDK and is entirely deterministic, and should create the most profitable batch possible out of transactions in the `Peggy Tx Pool`.
    - The `TransactionBatch` includes a batch nonce.
    - It also includes the latest `Valset`
    - The transactions in this batch are removed from the `Peggy Tx Pool`, and cannot be included in a new batch.
- Batches in the `Peggy Batch Pool` are signed over by the `Validator Set`'s `Eth Signers`.
    - `Relayers` may now attempt to submit these batches to the Peggy contract. If a batch has enough signatures (2/3+1 of the `Multisig Set`), it's submission will succeed. The decision whether or not to attempt a batch submission is entirely up to a given `Relayer`.
- Once a batch is `Observed` to have been successfully submitted to Ethereum (this takes at least as long as the `EthBlockDelay`), any batches in the `Peggy Batch Pool` which have a lower nonce, and have not yet been successfully submitted have their transactions returned to the `Peggy Tx Pool` to be tried in a new batch. This is safe because we know that these batches cannot possibly be submitted any more since their nonces are too low.

- When a new MsgRequestBatch comes in a new batch will not be produced unless it is more profitable than any batch currently in the `Peggy Batch Pool`. This means that when there is a batch backlog batches _must_ become progressively more profitable to submit.