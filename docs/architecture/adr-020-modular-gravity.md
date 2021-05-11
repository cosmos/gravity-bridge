# ADR 020: Modular Gravity

## Changelog

- 5/11/2021: First draft

## Status

DRAFT

## Abstract

This ADR proposes making Gravity fully modular, by making the core Gravity.sol Ethereum contract and Cosmos module responsible only for updating the signer set. Transaction batches would be moved outside of the core Gravity code, putting this functionality on the same footing as any future Gravity modules developed by outside developers. This would hopefully make it easy for users of Gravity to extend it's functionality without needing to fork core code. Removing the token bridge functionality would also stop Gravity from being tightly wedded to the ERC20 token standard.

## Context

The core of Gravity is the ability for it to securely relay information from Ethereum to Cosmos using Ethereum events, and from Cosmos to Ethereum using Ethereum contract calls executed by a multisig. An important part of is updating the signer set which the Gravity.sol Ethereum contract uses to check that relayed contract calls are valid.

For historical reasons, Gravity also builds a token bridge into its core code. This is because the strongest impetus for the creation of Gravity was bridging tokens. Focusing on the token bridging use case allowed us to avoid over-abstraction and navel gazing. Designing and implementing this functionality lead us to the current Gravity design. Now that the code has settled down and is being finalized, we have an opportunity to make it modular.

### Current Contract Call Code

This functionality extends the "contract call" functionality. This functionality was built by a contributor who needed it for a very specific use case, and whose system already was using a different kind of of oracle. For this reason, it does not include the ability to take arbitrary events from Ethereum to Cosmos, and it has some integration with the token bridge.

Currently, when making a contract call through Gravity, there is the option to send some tokens from the Gravity.sol contract into the contract that is being called. This is so that a user who has previously transferred some tokens into the token bridge can then have the option of having those tokens used by the contract being called. This is convenient, but couples the contract call functionality tightly to the built in token bridge, which this ADR proposes to remove.

### Fees

Currently, fees are paid out to relayers (who pay to submit transactions to Ethereum), by sending ERC20 tokens to `msg.sender` on Ethereum, out of the wallet of the Gravity.sol contract. This happens in token bridge transfers as well as the contract call functionality. But this assumes that the Gravity.sol contract has tokens in its wallet, which only makes sense with a built in token bridge.

#### Option 1: Let called contract handle fees

This seems like the most straightforward option: just let the contract that is being called handle the fees. However, this has several caveats.

- In the called contract, `msg.sender` will be set to address of the Gravity.sol contract, and so fees sent there will not reward the relayer. The called contract will instead need to use `tx.origin`. `tx.origin` has some well-known security pitfalls when using it for authorization, but for this use case, it should be OK.
- This model still punts the question of fees to the called contract. Even if the called contract uses `tx.origin`, it will still need to have tokens in its wallet to reward relayers. This creates extra complexity for the called contract, and the associated Cosmos module, and a lot of duplication. Many Gravity modules will probably opt to either run centralized relayers, or put tokens into the wallets of their Ethereum contracts. This is not desirable.

#### Option 2: Keep the built in token bridge

Earlier in this document, I wrote about removing the token bridge functionality, now I am talking about keeping it. That's why this is a draft. There's an argument to be made that, given the need to pay fees, a built in token bridge is actually core to Gravity. In this model, we would keep the existing contract call functionality exactly as-is, allowing the relayer to be paid out from the Gravity.sol wallet. Some caveats:

- This is not as clean and modular any more.
- This allows modules which are not Gravity to pay fees out of Gravity.sol's wallet. In a lot of Cosmos chains this is not an issue, since modules are trusted code, but it is less desireable than something more trustless.
- It keeps Gravity tightly coupled to the ERC20 standard.

#### Option 3: Pay fees on the Cosmos side

The idea to pay fees on the Cosmos side instead of on Ethereum has come up in other contexts before. The rough idea is that when relaying a transaction, the relayer would pass a Cosmos address into Gravity.sol alongside the transaction data and the signatures. When the transaction was accepted by the Ethereum blockchain, this Cosmos address would be emitted in an event. The Gravity Cosmos module would then send fees to this address on Cosmos.

A module calling Gravity could then send some Cosmos tokens to the Cosmos address of a relayer. In the case of a token bridge module, these tokens would be taken from the fees sent my users to the module along with their transactions. In the case of another type of module, such as an NFT bridge, the fees might be in Atoms.

## Decision

We will make the following modifications to Gravity:

- Remove all code dealing with transaction batch creation etc. from the core Gravity module. Move it into its own token bridge module.
- Do the same for the Gravity.sol contract.
- Modify the contract call functionality in the following way:
  -

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## Further Discussions

While an ADR is in the DRAFT or PROPOSED stage, this section should contain a summary of issues to be solved in future iterations (usually referencing comments from a pull-request discussion).
Later, this section can optionally list ideas or improvements the author or reviewers found during the analysis of this ADR.

## Test Cases [optional]

Test cases for an implementation are mandatory for ADRs that are affecting consensus changes. Other ADRs can choose to include links to test cases if applicable.

## References

- {reference link}
