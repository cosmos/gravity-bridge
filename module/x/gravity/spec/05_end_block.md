<!--
order: 5
-->

# End-Block

Each abci end block call, the operations to update queues and validator set
changes are specified to execute.

## Slashing

Slashing groups multiple types of slashing (validator set, batch and claim slashing). We will cover how these work in the following sections.

### Valset Slashing

This slashing condition is triggered when a validator does not sign a validator set update or transaction batch which is produced by the Gravity Cosmos module. This prevents two bad scenarios-

1. A validator simply does not bother to keep the correct binaries running on their system,
2. A cartel of >1/3 validators unbond and then refuse to sign updates, preventing any validator set updates from getting enough signatures to be submitted to the Gravity Ethereum contract. If they prevent validator set updates for longer than the Cosmos unbonding period, they can no longer be punished for submitting fake validator set updates and tx batches.

To deal with scenario 2, we also need to slash validators who are no longer validating, but are still in the unbonding period for up to `UnbondSlashingValsetsWindow` blocks. This means that when a validator leaves the validator set, they will need to keep running their equipment for at least `UnbondSlashingValsetsWindow` blocks. This is unusual for a Cosmos chain, and may not be accepted by the validators.

The current value of `UnbondSlashingValsetsWindow` is 10,000 blocks, or about 12-14 hours. We have determined this to be a safe value based on the following logic. So long as every validator leaving hte validator set signs at least one validator set update that they are not contained in then it is guaranteed to be possible for a relayer to produce a chain of validator set updates to transform the current state on the chain into the present state.

It should be noted that this slashing requirement could be eliminated with no loss of security if it where possible to perform the Ethereum signatures inside the consensus code. This is a pretty limited feature addition to Tendermint.

### Batch Slashing

This slashing condition is triggered when a validator does not sign a transaction batch which is produced by the Gravity Cosmos module. This prevents two bad scenarios-

1. A validator simply does not bother to keep the correct binaries running on their system,
2. A cartel of >1/3 validators unbond and then refuse to sign updates, preventing any batches from getting enough signatures to be submitted to the Gravity Ethereum contract.

## Attestation

Iterates through all attestations currently being voted on. Once an attestation nonce one higher than the previous one, we stop searching for an attestation and call `TryAttestation`. Once an attestation at a specific nonce has enough votes all the other attestations will be skipped and the `lastObservedEventNonce` incremented.

## Cleanup

Cleanup loops through batches and logic calls in order to clean up the timed out transactions.

### Batches

When a batch of transactions are created they have a specified height of the opposing chain for when the batch becomes invalid. When this happens we must remove them from the store. At the end of every block, we loop through the store of logic calls checking the the timeout heights.

### Logic Calls

When a logic call is created it consists of a timeout height. This height is used to know when the logic call becomes invalid. At the end of every block, we loop through the store of logic calls checking the the timeout heights.
