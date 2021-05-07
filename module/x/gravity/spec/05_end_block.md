<!--
order: 5
-->

# End-Block

Each abci end block call, the operations to update queues and validator set
changes are specified to execute.

## Valset Creation

Every endblock, we run the following procedure to determine whether to make a new `Valset` which will then need to be signed by all validators.

1. If there are no valset requests, create a new one.
2. If there is at least one validator who started unbonding in current block, create a `Valset`. This will make sure the unbonding validator has to provide an attestation to a new Valset that excludes them before they completely Unbond. Otherwise they will be slashed.
3. If power change between validators of CurrentValset and latest valset request is > 5%, create a new `Valset`.

If the above conditions are met, we create a new `Valset` using the procedure described [here](03_state_transitions.md#valset-creation)

## Slashing

Slashing groups multiple types of slashing (validator set, batch and claim slashing). We will cover how these work in the following sections.

### Valset Slashing

This slashing condition is triggered when a validator does not sign a validator set update or transaction batch which is produced by the Gravity Cosmos module. This prevents two bad scenarios-

1. A validator simply does not bother to keep the correct binaries running on their system,
2. A cartel of >1/3 validators unbond and then refuse to sign updates, preventing any validator set updates from getting enough signatures to be submitted to the Gravity Ethereum contract. If they prevent validator set updates for longer than the Cosmos unbonding period, they can no longer be punished for submitting fake validator set updates and tx batches.

To deal with scenario 2, we also need to slash validators who are no longer validating, but are still in the unbonding period for up to `UnbondSlashingValsetsWindow` blocks. This means that when a validator leaves the validator set, they will need to keep running their equipment for at least `UnbondSlashingValsetsWindow` blocks. This is unusual for a Cosmos chain, and may not be accepted by the validators.

The current value of `UnbondSlashingValsetsWindow` is 10,000 blocks, or about 12-14 hours. We have determined this to be a safe value based on the following logic. So long as every validator leaving hte validator set signs at least one validator set update that they are not contained in then it is guaranteed to be possible for a relayer to produce a chain of validator set updates to transform the current state on the chain into the present state.

#### Slashing currently bonded validators:

- We compare the current Cosmos block height with the `SignedValsetsWindow` parameter. If the current block height is less than `SignedValsetsWindow`, this procedure completes, doing nothing.
- Get a list of all the valsets that have not yet been processed by this procedure in earlier blocks, up to `SignedValsetsWindow` blocks ago. We only slash for each valset once.
- For each of these valsets:
  - Get the current set of bonded validators with `StakingKeeper.GetBondedValidatorsByPower`. For each validator:
    - Check that the validator started validating before the valset was created.
    - Check if the validator has signed this valset with a `MsgConfirmBatch`. If not, slash the validator by `SlashFractionValset` and jail them.

#### Slashing recently unbonded validators:

- Get a list of all validators who are currently unbonding. For each of these:
  - Check that the validator started validating before the valset was created.
  - Check that the valset was created before the `UnbondSlashingValsetsWindow` for that validator ended. (validators `UnbondingHeight` + `UnbondSlashingValsetsWindow`)
  - Check if the validator has signed this valset with a `MsgConfirmBatch`. If not, slash the validator by `SlashFractionValset` and jail them.

### Batch Slashing

This slashing condition is triggered when a validator does not sign a transaction batch which is produced by the Gravity Cosmos module. This prevents two bad scenarios-

1. A validator simply does not bother to keep the correct binaries running on their system,
2. A cartel of >1/3 validators unbond and then refuse to sign updates, preventing any batches from getting enough signatures to be submitted to the Gravity Ethereum contract.

Procedure:

- We compare the current Cosmos block height with the `SignedBatchesWindow` parameter. If the current block height is less than `SignedBatchesWindow`, this procedure completes, doing nothing.
- Get a list of all the valsets that have not yet been processed by this procedure in earlier blocks, up to `SignedBatchesWindow` blocks ago. We only slash for each valset once.
- For each of these valsets:
  - Get the current set of bonded validators with `StakingKeeper.GetBondedValidatorsByPower`. For each validator:
    - Check that the validator started validating before the valset was created.
    - Check if the validator has signed this valset with a `MsgValsetConfirm`. If not, slash the validator by `SlashFractionBatch` and jail them.

## Attestation vote counting

This logic counts up votes on `Attestation`s and kicks off the process of bringing Ethereum events into the Cosmos state;

- We retrieve all attestations from storage and order them into a map of event nonces to attestations, sorted by nonce: `map[uint64][]types.Attestation`.
  - Note that the only time one nonce will have more than one attestation is when validators are disagreeing about which event happened at which event nonce.
- We then loop over the nonces:
  - For each attestation, we check that the event nonce is exactly 1 higher than the `LastObservedEventNonce`.
  - If it is, we count up the votes on that attestation using the procedure described [here](03_state_transitions.md#counting-attestation-votes)
  - If the attestation passes the `AttestationVotesPowerThreshold`, we apply it to the Cosmos state, and increment the `LastObservedEventNonce`. As a result of this, any additional attestations at the same nonce do not have their votes counted, but the first attestation at the next nonce will have its votes counted.
  - If the attestation does not pass the `AttestationVotesPowerThreshold`, it is not applied to the Cosmos state, and `LastObservedEventNonce` is not incremented. As a result of this, the next attestation at that nonce will have its votes counted. If no attestations at that nonce pass the `AttestationVotesPowerThreshold`, then all attestations at subsequent nonces will be skipped and this procedure ends.

This procedure has the following attributes:

- Attestations will never be observed and applied to Cosmos state out of order, since to have their votes counted, they must have a nonce exactly one higher than the last observed attestation.
- It is only possible for one attestation at a given nonce to pass the `AttestationVotesPowerThreshold` and become `Observed`, since we have [enforced](03_state_transitions.md#counting-attestation-votes) that validators cannot vote for different attestations at the same height.
- If there is an attestation that has not passed the `AttestationVotesPowerThreshold`, but there are later attestations which have, we do not count the later attestations until the earlier one passes the `AttestationVotesPowerThreshold` and is observed. At this point, all later attestations which have passed the `AttestationVotesPowerThreshold` will also be counted and be applied to the Cosmos state.

## Cleanup

Cleanup loops through batches and logic calls in order to clean up the timed out transactions.

### Batches

When a batch of transactions are created they have a specified Ethereum block height when the batch becomes invalid. When this happens we must remove them from the store. At the end of every block, we loop through the store of batches checking the the timeout heights.

Here is the procedure:

- Get the `LastObservedEthereumBlockHeight`. This is the Ethereum block height from the last observed attestation, and is the most current information we have available about the current Ethereum block height.
- Loop through all `OutgoingTxBatches`. For each batch:
  - If the `BatchTimeout` on that batch is lower than the `LastObservedEthereumBlockHeight`, cancel the batch, freeing its transactions to be included in another batch, or cancelled by the sender.

An important attribute of this procedure is that it is not possible for a batch to be executed on Ethereum AND cleaned up by the above procedure. This would result in double spends. Since the `LastObservedEthereumBlockHeight` can only be updated by observed Ethereum events, and events must be observed in order, it is impossible for the `LastObservedEthereumBlockHeight` to be higher than a given batches `BatchTimeout` AND for that batch to have been executed without the Gravity module having received the batch executed event.

### Logic Calls

Same procedure as batch cleanup above, but with logic calls.
