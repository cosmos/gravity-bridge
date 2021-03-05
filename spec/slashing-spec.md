This file names and documents the various slashing conditions we use in Gravity.

## GRAVSLASH-01: Signing fake validator set or tx batch evidence

This slashing condition is intended to stop validators from signing over a validator set and nonce that has never existed on Cosmos. It works via an evidence mechanism, where anyone can submit a message containing the signature of a validator over a fake validator set. This is intended to produce the effect that if a cartel of validators is formed with the intention of submitting a fake validator set, one defector can cause them all to be slashed.

**Implementation considerations:**

The trickiest part of this slashing condition is determining that a validator set has never existed on Cosmos. To save space, we will need to clean up old validator sets. We could keep a mapping of validator set hash to true in the KV store, and use that to check if a validator set has ever existed. This is more efficient than storing the whole validator set, but its growth is still unbounded. It might be possible to use other cryptographic methods to cut down on the size of this mapping. It might be OK to prune very old entries from this mapping, but any pruning reduces the deterrence of this slashing condition.

## GRAVSLASH-02: Failure to sign validator set update or tx batch

This slashing condition is triggered when a validator does not sign a validator set update or transaction batch which is produced by the Gravity Cosmos module. This prevents two bad scenarios- 

1. A validator simply does not bother to keep the correct binaries running on their system,
2. A cartel of >1/3 validators unbond and then refuse to sign updates, preventing any validator set updates from getting enough signatures to be submitted to the Gravity Ethereum contract. If they prevent validator set updates for longer than the Cosmos unbonding period, they can no longer be punished for submitting fake validator set updates and tx batches (GRAVSLASH-01 and GRAVSLASH-02). 

To deal with scenario 2, GRAVSLASH-02 will also need to slash validators who are no longer validating, but are still in the unbonding period. This means that when a validator leaves the validator set, they will need to keep running their equipment for 2 weeks. This is unusual for a Cosmos chain, and may not be accepted by the validators. Research is ongoing for ways to allow validators to stop signing before the unbonding period is fully over.

## GRAVSLASH-03: Submitting incorrect Eth oracle claim - INTENTIONALLY NOT IMPLEMENTED

The Ethereum oracle code (currently mostly contained in attestation.go), is a key part of Gravity. It allows the Gravity module to have knowledge of events that have occurred on Ethereum, such as deposits and executed batches. GRAVSLASH-03 is intended to punish validators who submit a claim for an event that never happened on Ethereum.

**Implementation considerations**

The only way we know whether an event has happened on Ethereum is through the Ethereum event oracle itself. So to implement this slashing condition, we slash validators who have submitted claims for a different event at the same nonce as an event that was observed by >2/3s of validators.

Although well-intentioned, this slashing condition is likely not advisable for most applications of Gravity. This is because it ties the functioning of the Cosmos chain which it is installed on to the correct functioning of the Ethereum chain. If there is a serious fork of the Ethereum chain, different validators behaving honestly may see different events at the same event nonce and be slashed through no fault of their own. Widespread unfair slashing would be very disruptive to the social structure of the Cosmos chain.

Maybe GRAVSLASH-03 is not necessary at all:

The real utility of this slashing condition is to make it so that, if >2/3 of the validators form a cartel to all submit a fake event at a certain nonce, some number of them can defect from the cartel and submit the real event at that nonce. If there are enough defecting cartel members that the real event becomes observed, then the remaining cartel members will be slashed by this condition. However, this would require >1/2 of the cartel members to defect in most conditions. 

If not enough of the cartel defects, then neither event will be observed, and the Ethereum oracle will just halt. This is a much more likely scenario than one in which GRAVSLASH-03 is actually triggered.

Also, GRAVSLASH-03 will be triggered against the honest validators in the case of a successful cartel. This could act to make it easier for a forming cartel to threaten validators who do not want to join.

## GRAVSLASH-04: Failure to submit Eth oracle claims

This is similar to GRAVSLASH-03, but it is triggered against validators who do not submit an oracle claim that has been observed. In contrast to GRAVSLASH-03, GRAVSLASH-04 is intended to punish validators who stop participating in the oracle completely. 

**Implementation considerations**

Unfortunately, GRAVSLASH-04 has the same downsides as GRAVSLASH-03 in that it ties the correct operation of the Cosmos chain to the Ethereum chain. Also, it likely does not incentivize much in the way of correct behavior. To avoid triggering GRAVSLASH-04, a validator simply needs to copy claims which are close to becoming observed. This copying of claims could be prevented by a commit-reveal scheme, but it would still be easy for a "lazy validator" to simply use a public Ethereum full node or block explorer, with similar effects on security. Therefore, the real usefulness of GRAVSLASH-04 is likely minimal

Without GRAVSLASH-03 and GRAVSLASH-04, the Ethereum event oracle only continues to function if >2/3 of the validators voluntarily submit correct claims. Although the arguments against GRAVSLASH-03 and GRAVSLASH-04 are convincing, we must decide whether we are comfortable with this fact. We should probably make it possible to enable or disable GRAVSLASH-03 and GRAVSLASH-04 in the chain's parameters.