# ETH to Cosmos Oracle

All `Operators` run an `Oracle` binary. This separate process monitors an Ethereum node for new events involving the `Gravity Contract` on the Ethereum chain. Every event that `Oracle` monitors has an event nonce. This nonce is a unique coordinating value for a `Claim`. Since every event that may need to be observed by the `Oracle` has a unique event nonce `Claims` can always refer to a unique event by specifying the event nonce.

- An `Oracle` observes an event on the Ethereum chain, it packages this event into a `Claim` and submits this claim to the cosmos chain as an [Oracle message](/docs/design/messages.md##Oracle-messages)
- Within the Gravity Cosmos module this `Claim` either creates or is added to an existing `Attestation` that matches the details of the `Claim` once more than 66% of the active `Validator` set has made a `Claim` that matches the given `Attestation` the `Attestation` is executed. This may mint tokens, burn tokens, or whatever is appropriate for this particular event.
- In the event that the validators can not agree >66% on a single `Attestation` the oracle is halted. This means no new events will be relayed from Ethereum until some of the validators change their votes. There is no slashing condition for this, with reasoning outlined in the [slashing spec](/spec/slashing-spec.md)
