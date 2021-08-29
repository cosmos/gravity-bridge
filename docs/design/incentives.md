# Incentives

This document covers all the incentivization systems in Gravity bridge.

## Validators

Currently validators in Gravity have only one carrot. The extra activity brought to the chain by a functioning bridge.

There are on the other hand a lot of negative incentives (sticks) that the validators must watch out for. These are outlined in the [slashing spec](/spec/slashing-spec.md).

One negative incentive that is not covered under slashing is the cost of submitting oracle submissions and signatures. Currently these operations are not incentivized, but still cost the validators fees to submit. This isn't an issue considering the low activity on most Cosmos based chains at the moment. But an active Ethereum bridge and dex may change that issue very quickly.

Some positive incentives for correctly participating in the operation of the bridge should be under consideration. In addition to eliminating the fees for mandatory submissions.
