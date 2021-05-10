# ADR 001: No light clients

## Changelog

- 6/2020: Initial decision
- 5/6/2021: First draft in ADR format

## Abstract

Most cross-chain communication software relies on in-consensus light clients: each chain has a process in consensus, that is able to understand the headers of the other chain and verify that an event or transaction on that chain actually happened. Gravity instead uses a multi-sig with the signatures of the current Cosmos validator set on the Ethereum side, and an oracle powered by the validators on the Cosmos side.

## Context


Light clients are relatively complicated and computationally expensive compared to oracles and multi-sig systems. They need to verify blockchain headers using a specialized procedure, instead of just verifying some signatures. Headers are also very specific to the underlying blockchain technology, which means that you'd need different light clients to connect Gravity to Ethereum vs Polygon, for example.

## Decision


PoS chains, including Cosmos, generally have a small number of validators which make up the 2/3s majority needed to reach quorum. Cosmos also allows us great freedom to design the software that validators run. These two attributes of Cosmos allowed us to implement a design with no in-consensus light clients. In general, our strategy is to lean on the flexibility of Cosmos to make things easier and cheaper on Ethereum.

**Ethereum -> Cosmos**

The mechanism that Ethereum provides to monitor the chain is events, which can be emitted by any smart contract. To read events from Ethereum into the Cosmos state, we rely on a simple oracle within the Gravity module. All validators in the validator set watch events on the Ethereum blockchain. There is no requirement on where they get these events, but the simplest option is to run an Ethereum full node. When validators see an event, they send a message to the Cosmos blockchain. When this is received by the Gravity module:

- It counts up the voting power of the validators who have seen an event.
- When an event passes the consensus quorum threshold (usually >2/3s), it is applied to the Cosmos state.
- There is also a mechanism that ensures that events are only applied to the state in the order they were emitted on Ethereum.

**Cosmos -> Ethereum**

Gravity communicates from Cosmos to Ethereum with what are basically remote procedure calls. There is a specialized type of transaaction used for transferring tokens called a transfer batch, but Gravity also allows other Cosmos modules to call arbitrary Ethereum contracts using this mechanism. The Gravity.sol Ethereum contract only executes one of these transactions if it is signed by a quorum of the current Cosmos validator set. It does this by verifying their signatures.

The Gravity.sol Ethereum contract must for this reason be kept updated with the latest validator set, something which introduces many subtle security considerations which we will discuss below.

Gravity gathers the signatures of the validators over transactions destined for Ethereum by having the validators send them in to the Cosmos chain with messages. Once a transaction has been signed by a quorum of the validators, it can be submitted to Ethereum by a relayer.

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Positive

- The Gravity Cosmos module is extremely portable between different EVM chains. It does not care at all how the blockchain works. It can be deployed on Polygon, xDAI, Ethereum, BSC, Ethermint, and any other blockchain with an EVM with no additional code or configuration.
- Verifying signatures can be cheaper than running a light client <TODO: work out some hard numbers here>
- We did not need to get an Ethereum light client working in Cosmos consensus, which made development quicker.

### Negative

- The current signer set stored by the Ethereum contract must be kept updated, otherwise security is broken. We have adddressed this by requiring validators departing the validator set to keep signing signer set updates. Currently this means that they need to keep their hardware running for one hour after departing the validator set, which is not usually a requirement for Cosmos validators.
- It is possible for validators to sign fake transactions that never were produced by Cosmos and submit them to the Ethereum contract. This is somewhat circular, since to be accepted, these fake transactions would need to be signed by a quorum of validators, who could corrupt the Cosmos chain in more straightforward ways. However, we have addressed this with evidence-based slashing, where validators are slashed on Cosmos for signing fake transactions.
- It is possible for validators to sign fake events which never existed on Ethereum. These would have to be signed by a quorum of validators, who again, could corrupt the Cosmos chain in a simpler manner. The difference is that if a quorum of validators broke the rules of Cosmos consensus, this would be detected by Cosmos full nodes who could "raise the alarm". We have addressed this by building functionality into the Gravity relayer which can similarly detect fake Ethereum events being signed, and raise the alarm. These fake events can also be detected by anyone who is watching both the Ethereum and the Cosmos chains.

The audit performed by Informal addresses all of these points in much greater detail.

### Neutral

- The Gravity.sol Ethereum contract is somewhat portable between PoS chains. This is not an intended use case, but it will work with any chain that can put together signatures over transactions from a quorum of validators.

## Further Discussions

While an ADR is in the DRAFT or PROPOSED stage, this section should contain a summary of issues to be solved in future iterations (usually referencing comments from a pull-request discussion).
Later, this section can optionally list ideas or improvements the author or reviewers found during the analysis of this ADR.
