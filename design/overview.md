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
- `Orchestrator` (name WIP) - This is a Rust binary controlled by an Operator that holds Cosmos SDK and Ethereum private keys used for signing transactions used to move tokens between the two chains. (This was called peggy module in other docs, which is confusing as the module runs on the node... I would love a better name than Orchestrator)
- `REST server` - This is the Cosmos SDK "REST Server" that runs on Port 1317, either on the validator node or another Cosmos SDK node controlled by the Operator
- `Ethereum RPC` - This is the JSON-RPC server for the Ethereum Full Node.
- `Validator Set` - The set of validators on the Cosmos SDK chain, along with their respective voting power. These are ed25519 public keys used to sign tendermint blocks.
- `Orchestrator Set` - The validator set mapped over to Cosmos SDK Orchestrator keys. This is used on the Cosmos SDK chain to authorize messages from Ethereum. Note that we can use the current validator set, but we use Orchestrator keys (not tendermint consensus keys)
- `MultiSig Set` - The set of Ethereum keys along with respective voting power. This is based on the validator set and mapped over the registered keys. However, as this is a different chain, this is mirrored with a delay.
- `Peggy Bridge Tx pool` - Is a transaction pool that exists in the chain store of Cosmos -> Ethereum transactions waiting to be placed into a transaction batch
- `Transaction batch` - A transaction batch is a set of Ethereum transactions to be sent from the Peggy Ethereum contract at the same time. This helps reduce the costs of submitting a batch. Batches have a maximum size (currently around 100 transactions) and are only involved in the Cosmos -> Ethereum flow
- `Peggy Bridge Batch pool` - Is a transaction pool like strucutre that exists in the chains to store, seperate from the `Peggy Bridge Tx pool` it stores transactions that have been placed in batches that are in the process of being signed or being submitted by the `Orchestrator Set`
- `EthBlockDelay` - Is a agreed upon number of Ethereum blocks all oracle attestations are delayed by. No `Orchestrator` will attest to have seen an event occur on Ethereum until this number of blocks has elapsed as denoted by their trusted Ethereum full node. This should prevent short forks form causing disagreements on the Cosmos side. The current value being consdiered is 50 blocks.
- `Last ValSet` - The `Last Valset` is the last validator set known to have been accepted by the Ethereum chain. When a new ValSet update is submitted it is first `Confirmed` making it possible to submit and then `Observed` updating the value of `Last ValSet` in the store.
- `Confirmed` - Many actions in Peggy require RLP encoded signatures to be submitted by the `Orchestrators`, an operation is `Confirmed` when it is possible to execute it on Ethereum by collecting and submitting these RLP encoded signatures representing 66% of all voting power encoded in the `Last Valset`.
- `Observed` - events on Ethereum are considered `Observed` when 66% of the active Cosmos validator set during a given block has submitted an oracle message attesting to seeing the event. Note this is distinct from `Confirmed`! At any given time the set of who can `Observe` an event and who can `Confirm` and event will be slightly different. Since `Confirmed` is based off of the `Last ValSet` and `Observed` is always up to date with the latest validator set.
- `Validator set delta` - This is a term for the difference between the validator set currently in the Peggy Ethereum contract and the actual validator set on the Cosmos chain. Since the validator set may change every single block there is essentially guaranteed to be some nonzero `Validator set delta` at any given time.
- `Allowed validator set delta` - This is the maximum allowed `Validator set delta` this parameter is used to determine if the Peggy contract in MsgProposePeggyContract has a validator set 'close enough' to accept. It is also used to determine when validator set updates need to be sent. This is decided by a governance vote _before_ MsgProposePeggyContract can be sent.
- `Peggy ID` - This is a random 32 byte value required to be included in all Peggy signatures for a particular contract instance. It is passed into the contract constructor on Ethereum and used to prevent signature reuse when contracts may share a validator set or subsets of a validator set. This is also set by a governance vote _before_ MsgProposePeggyContract can be sent.
- `Peggy contract code hash` - This is the code hash of a known good version of the Peggy contract solidity code. It will be used to verify exactly which version of the bridge will be deployed.
- `Start Threshold` - This is the percentage of total voting power that must be online and participating in Peggy operations before a bridge can start operating.

The _Operator_ is the key unit of trust here. Each operator is responsible for maintaining 3 secure processes:

1. Cosmos SDK Validator - signing blocks
1. Fully synced Ethereum Full Node
1. Orchestrator, which connects to the above as a client

## Security Concerns

The **Validator Set** is the actual set of keys with stake behind them, which are slashed for double-signs or other
misbehavior. We typically consider the security of a chain to be the security of a _Validator Set_. This varies on
each chain, but is our gold standard. Even IBC offers no more security than the minimum of both involved Validator Sets.

The **Orchestrator Set** is another Cosmos SDK client key associated with the same validator set. We can add slashing
conditions to any mis-signed message by the Orchestrator Set and be able to provide the same security as the
_Valiator Set_, just a different module detecting evidence of malice and deciding how much to slash. If we can prove a
transaction signed by any member of the _Orchestrator Set_ was illegal or malicious, then we can slash on the Cosmos chain
side an potentially provide 100% of the security of the Validator Set. Note that this also has access to the 3 week unbonding
period to allow evidence to slash even if they immediately unbond.

The **MultiSig Set** is a (possibly aged) mirror of the _Validator Set_ but with Ethereum keys, and stored on the Ethereum
contract. If we ensure the _MultiSig Set_ is updated much more often than the unbonding period (eg at least once per week),
then we can guarantee that all members of the _MultiSig Set_ have slashable atoms for misbehavior. However, in some extreme
cases of stake shifting, the _MultiSig Set_ and _Validator/Orchestrator Set_ could get quite far apart, meaning there is
many of the members in the _MultiSig Set_ are no longer active validators and may not bother to transfer Eth messages.
Thus, to avoid censorship attacks/inactivity, we should also update this everytime there is a significant change
in the Validator Set (eg. > 3-5%). If we maintain those two conditions, the MultiSig Set should offer a similar level of
security as the Validator Set.

There are now 3 conditions that can be slashed for any validator: Double-signing a block with the tendermint key from the
**Validator Set**, signing an invalid/malicious message with the Cosmos SDK key from the **Orchestrator Set**, or
signing an invalid/malicious Ethereum message with the key from the **MultiSig Set**. If all conditions of misbehavior can
be attributed to a signature from one of these sets, and proven **on the Cosmos chain**, then we can argue that Peggy offers
a security level equal to the minimum of the Peg-Zone Validator Set, or reorganizing the Ethereum Chain 50 blocks.
And provide a security equivalent to or greater than IBC.

## Bootstrapping

We assume the act of upgrading the Cosmos-based binary to have peggy module is already complete,
as approaches to that are discussed in many other places. Here we focus on the _activation_ step.

1. Each `Operator` generates an Ethereum and Cosmos private key for the `Orchestrator`. These addresses are signed and submitted by the Operators valoper key in a MsgRegisterOrchestrator. The `Orchestrator` is now free to use these delegated keys for all Peggy messages.
1. A governance vote is held on bridge parameters including `Peggy ID`, `Allowed validator set delta`, `start threshold`, and `Peggy contract code hash`
1. Anyone deploys a Peggy contract using a known codehash and the current validator set of the Cosmos zone to an Ethereum compatible blockchain.
1. Each `Operator` may or may not configure their `Orchestrator` with the above Peggy contract address
1. If configured with an address the `Orchestrator` checks the provided address. If the contract passes validation the `Orchestrator` signs and submits a MsgProposePeggyContract. Validation is defined as finding the correct `Peggy contract code hash` and a validator set matching the current set within `Allowed validator set delta`.
1. A contract address is considered adopted when voting power exceeding the `start threshold` has sent a MsgProposePeggyContract with the same Ethereum address.
1. Because validator sets change quickly `Orchestrators` not configured with a contract address observe the Cosmos blockchain for submissions. When an address is submitted they validate it and approve it themselves if it passes. This results in a workflow where once a valid contract is proposed it will be ratified in a matter of a few seconds.
1. It is possible for the adoption process to fail if a race condition is intentionally created resulting in less than 66% of the validator power approving more than one valid Peggy Ethereum contract. In this case the Orchestrator will check the contract address with the majority of the power (or at random in the case of a perfect tie) and switch it's vote. This leaves only the possible edge case of >33% of `Operators` intentionally selecting a different contract address. This would be a consensus failure and the bridge can not progress.
1. The bridge ratification process is complete, the contract address is now placed in the store to be referenced and other operations are allowed to move forward.

At this point, we know we have a contract on Ethereum with the proper _MultiSig Set_, that > `start threshold` of the _Orchestrator Set_ is online and agrees with this contract, and that the Cosmos chain has stored this contract address. Only then can we begin to accept transactions to transfer tokens

Note: `start threshold` is some security factor for bootstrapping. 67% is sufficient to release, but we don't want to start until there is a margin of error online (not to fall off with a small change of voting power). This may be 70, 80, 90, or even 95% depending on how much assurances we want that all _Orchestrators_ are operational before starting.

## Relaying ETH to Cosmos

**TODO**

## Relaying Cosmos to ETH

- Someone (permissionlessly) sends MsgSendToEth this creates an entry in the 'peggy bridge tx pool' this pool is in the Cosmos chain store and contains all transactions not yet bridged. Funds of the specified demon are then burned from the users account if that demon is valid to bridge, otherwise the tx is rejected.
- Someone (permissionlessly) sends a MsgRequestBatch, this produces a new `Transaction batch` in the `Peggy Bridge Tx pool`. The creation of this batch occurs in ComsosSDK and is entierly deterministic, ordering transactions from highest to lowest fee and including them in the `Transaction batch` until there are no more transactions or the batch size limit is reached.
- The `Orchestrators` view the `Peggy Bridge Tx pool` via REST endpoints. When a new batch appears they sign the RLP-encoded batch and submit a MsgConfirmBatch signed with their Cosmos Key that contains their Ethereum signature over the batch.
- Once enough signatures are collected and the batch is `Confirmed` any `Orcestrator` can collect these signatures from the REST endpoint, assemble an Ethereum transaction and submit the batch.
- The `Orchestrators` observe the Ethereum blockchain via their trusted Ethereum node. Once the batch has entered the blockchain and `EthBlockDelay` has elapsed an `Orchestrator` will send MsgBatchInChain. Containing a signed attestation that the block has entered the Ethereum chain, once the batch is `Observed` by sending a MsgBatchInChain the Cosmos -> Ethereum bridge process is complete. The batch is removed from the `Peggy Bridge Tx pool` and the MsgSendToEth flow is now fully complete.
- It is possible for a `TransactionBatch` to be created and `Confirmed` that is not profitable to submit. This can occur when gas prices are high or the average fee in the bridge is low. In this case the batch will either wait long enough for the gas situation to change or a new MsgRequestBatch will be Created.
- If an existing batch is `Confirmed` but not `Observed` when a new MsgRequestBatch comes in a new batch will not be produced unless it's ratio of fees to transactions is higher. This means that when there is a `Transaction Batch` backlog batches _must_ become progressively more profitable to submit.
- If there is more than one `Confirmed` batch in the `Peggy Bridge Batch pool` and a later one is successfully `Observed` earlier batches will have their transactions returned to the `Peggy Bridge Tx pool` for inclusion in later batches. This is safe because no earlier batch can be submitted after a later batch.

Open Questions for transaction batch logic

- After some number of blocks transactions the`Peggy Bridge Tx pool` should be removed from the pool and refunded. To be re-submitted with a higher fee.
- TODO: consider a gas price oracle which would allow us to better estimate profitable batches
- TODO: consider allowing users to remove their MsgSendToEth transactions from the `Peggy Bridge Tx pool` at their own request, without waiting for the timeout, this should always be safe because only transactions in the `Peggy Bridge Batch pool` are even possible to submit.
