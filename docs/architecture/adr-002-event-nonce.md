# ADR 02: In-contract Event Nonce

## Changelog

- 9/20: Original decision
- 5/10/2021: This retroactive ADR

## Status

Accepted

## Abstract

When an event is emitted by the Gravity.sol contract, it includes a field called the "event nonce". This is incremented with every event and helps keep the events in a strict ordering, but also uses more Ethereum gas, since it means that every event fired causes a storage write.

## Context

It is very important that Ethereum events being brought over to Cosmos stay in strict order, and have no gaps. If this invariant is violated, it could be a direct double spend vulnerability. Relatedly, this invariant is also used to make sure that validators are not able to vote for tow different events at the same event nonce. There are `blockNumber` and `transactionIndex` fields on Ethereum event logs, which allows one to order all Ethereum events. However, since events emitted by Gravity.sol are a tiny fraction of the total events on Ethereum, their `transactionIndex` sequence has gaps, and cannot be used to by the Gravity module to check that there are no gaps. Another option is to put the event nonce on events ourselves, which uses a storage write and cost some gas. This makes the `SendToCosmos` method use 38% more gas. The effect on other events is negligable.

## Decision

We chose to set our own `event_nonce` by incrementing a counter stored in Gravity.sol. Using the `transactionIndex` to order events is cheaper, and arguably the optimal solution. However, setting our own `event_nonce` has several advantages which we judged would enable shipping a higher security bridge more quickly:

- The `event_nonce` sequence is free of gaps. This makes it much easier to ensure in the Gravity Cosmos module that no later event is applied when an earlier event has not yet been applied.
- Applying the `event_nonce` in the Gravity.sol contract, and reading it in the Gravity module means that all code concerned with relaying events over (the orchestrator, the code run by validators to put Ethereum events onto the Cosmos chain) is non security critical. All security checks are done in on-chain code. This was a property we were trying to preserve during the design phase.

## Consequences

### Positive

This feature was simpler to implement securely than the alternative.

### Negative

The price we pay is the extra 38% gas fee on `SendToCosmos`, and smaller (proportionally) extra gas fees on other methods.

## Further Discussions

While writing this and revisiting the decision, I realized that we may be able to eliminate the extra gas fee without changing Gravity Cosmos module code. This solution would break the principle of not putting security critical code in the off chain codebase, but would be easy and limited:

- We would stop adding the `event_nonce` to events emitted by Gravity.sol, saving on the storage write and the extra gas fee.
- Instead, we would have the orchestrator add the `event_nonce` after receiving events. This would keep the sequence contigous with no gap, unlike relying on the `transactionIndex`.
- Cosmos module code would be unchanged.
