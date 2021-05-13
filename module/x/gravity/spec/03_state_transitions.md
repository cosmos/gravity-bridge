<!--
order: 3
-->

# State Transitions

This document describes the state transition operations pertaining to:

## EthereumEventVoteRecord

### First vote

The first time any validator sees a given Ethereum event on the Ethereum blockchain, and calls `SendToCosmosEvent`, or one of the other endpoints for other types of ethereum events (claims):

- We check that the event nonce of the submitted event is exactly one higher than that validator's last submitted event. This keeps validators from voting on different events at the same event nonce, which makes tallying votes easier later.
- An EthereumEventVoteRecord is created for that event at that event nonce. Event nonces are created by the Gravity.sol Ethereum contract, and increment every time it fires an event. It is possible for validators to disagree about what event happened at a given event nonce, but only in the case of an attempted attack by Cosmos validators, or in the case of serious issues with Ethereum (like a hard fork).
- That validator's address is added to the votes array.
- The observed field is initialized to false.
- The height field is filled with the current Cosmos block height.

### Subsequent votes

When other validators see the same event at the same event nonce, and call `SendToCosmosEvent`, or one of the other endpoints for other types of ethereum events:

- We check that the event nonce of the submitted event is exactly one higher than that validator's last submitted event. This keeps validators from voting on different events at the same event nonce, which makes tallying votes easier later.
- We look up the event's EthereumEventVoteRecord.
- The validator's address is added to the votes array.

### Counting EthereumEventVoteRecord votes

Every endblock, the module attempts to tally up the votes for un-Accepted ethereumEventVoteRecords. Which ethereumEventVoteRecords it chooses to tally is covered in the [end blocker spec](05_end_block.md).

When tallying the votes a given ethereumEventVoteRecord, we follow this algorithm:

- First get `LastTotalPower` from the StakingKeeper
- `requiredPower` = `EthereumEventVoteRecordPowerThreshold` \* `LastTotalPower` / 100
  - This effectively calculates `EthereumEventVoteRecordPowerThreshold` percent (usually 66%) of `LastTotalPower`, truncating all decimal points.
- Set `ethereumEventVoteRecordPower` = 0

- For every validator in the ethereumEventVoteRecord's votes field:
  - Add the validators current power to `ethereumEventVoteRecordPower`.
  - Check if the `ethereumEventVoteRecordPower` is greater than or equal to `requiredPower`
    - If so, we first check if the `eventNonce` of the ethereumEventVoteRecord's event is exactly one greater than the global `LastAcceptedEventNonce`. If it is not, something is very wrong and we panic (this could only be caused by programmer error elsewhere in the module).
    - We set the `observed` field to true, set the global `LastAcceptedEventNonce` to the ethereumEventVoteRecord's event's `event_nonce`. This will only ever result in incrementing the `LastAcceptedEventNonce` by one, given the preceding conditions.
    - We set the `LatestEthereumBlockHeight` to the Ethereum block height from the ethereumEventVoteRecord's event. This is used later when we need a recent Ethereum block height, for example to calculate batch timeouts.

Now we are ready to apply the ethereumEventVoteRecord's event to the Cosmos state. This is different depending on which event we are dealing with, see state transtions for the individual events.

## MsgSendToCosmosEvent

### On event observed:

- Check if deposited token is Ethereum or Cosmos originated, and get it's Cosmos denom, using the `MsgSendToCosmosEvent`'s `token_contract` field.
- If it is Cosmos originated:
  - Send the number of coins in the `amount` field to the Cosmos address in the `cosmos_receiver` field, from the Gravity module's wallet. This works because any Cosmos originated tokens that are circulating on Ethereum must have been created by depositing into the Gravity module at some point in the past.
- If it is Ethereum originated:
  - Mint the number of coins in the `amount` field and send to the Cosmos address in the `cosmos_receiver` field.

## MsgBatchExecutedEvent

This event is fired when a `BatchTx` is executed on Ethereum, sending the tokens in that `OutgoingTXBatch` to their destinations on Ethereum.

### On event observed:

- Delete all the transactions in the batch from the `SendToEthereumPool`, since they have been spent on Ethereum.
- For all batches with a `BatchNonce` lower than this one, put their transactions back into the `UnbatchedTXIndex`, which allows them to either be put into a new batch, or canceled by their sender using `MsgCancelSendToEthereum`. This is because the Gravity.sol Ethereum contract does not allow batches to be executed with a lower nonce than the last executed batch, meaning that the transactions in these batches can never be spent, making it safe to cancel them or put them in a new batch.

## MsgERC20DeployedEvent

Cosmos originated assets are represented by ERC20 contracts deployed on Ethereum by the Gravity.sol contract. This deployment can cost over $100, and somebody needs to pay for the gas. Gravity allows anybody to pay for this, as long as they deploy the contract with the correct parameters. Once this happens, the `MsgERC20DeployedEvent` event is fired and picked up by the Gravity module.

### On event observed:

- Check if a contract has already been deployed for this asset. If so, error out.
- Check if the Cosmos denom that the contract was deployed even exists. If not, error out.
- Check if the ERC20 parameters, Name, Symbol, and Decimals match the equivalent attributes in the `DenomMetaData`. If not, error out.
- If the previous checks all passed, associate the ERC20's contract address with the denom using the `CosmosOriginatedDenomToERC20` index

## BatchTx

### Batch creation

To create a new batch for a given token type:

- Check if there is a previous active batch for this token type, if so:
  - Calculate the fees (denominated in the batches token) that the new batch would generate for a relayer once submitted to Ethereum.
  - Calculate the fees that the previous batch would generate for a relayer.
  - If the new batch does not have higher fees than the old batch, error out.

This mechanism ensures smooth functioning of the bridge, by keeping batches from being filled with low value transactions. Consider:

If there were many transactions in the transaction pool with an unprofitably low fee, and a few coming in every block with a high fee, each high fee transaction might end up in a batch with a bunch of unprofitable transactions. These batches would not be profitable to submit, and so the profitable transactions would end up in unprofitable batches, and not be submitted.

By making it so that every new batch must be more profitable than any other batch that is waiting to be submitted, it gives the few profitable transactions that come in every block the chance to build up and form a batch profitable enough to submit.

Moving on with the batch creation process:

- Take the `BatchTxSize` unbatched transactions with the highest fees for the given token type, add them to the batches `transactions` field, and remove the transactions from the `UnbatchedTXIndex`, so they cannot be cancelled or added to another batch.
- Increment the `LastOutgoingBatchID` and set the batches `batch_nonce` field to the incremented value.
- Get the `BatchTimeout`. The batch timeout is an Ethereum block height in the future, after which the batch will no longer be accepted by the Gravity.sol contract. This allows unprofitable batches to time out and free their transactions to be added to a more profitable batch or be cancelled. Gravity has knowledge of the `LatestEthereumBlockHeight` which is brought in on every block, but this knowledge is only as recent as the last observed event. For this reason, we estimate the current Ethereum block height using the following procedure:
  - We estimate how many milliseconds it has been since we recorded the `LatestEthereumBlockHeight` by multiplying the number of blocks since then with the average Cosmos block time.
  - We estimate current Ethereum block height by dividing the product of the multiplication above by the average Ethereum block time, and adding to the `LatestEthereumBlockHeight`
  - We set the `BatchTimeout` by adding the proper number of blocks to the estimated current Ethereum block height.
  - In more compact notation:
    - a: Average Cosmos block time in ms
    - b: Cosmos block height at time of last recorded Ethereum block height
    - c: Current Cosmos block height
    - d: Average Ethereum block time in ms
    - e: Last recorded Ethereum block height
    - f: Target batch timeout in ms
    - `BatchTimeout` = ((((c - b) \* a) / d) + e) + (f / d)
- Store the batch, indexed by the token contract and the batch nonce.

### Batch signing

Once a batch has been created and stored, it is up to the current validators to sign it with their Ethereum keys so that it can be submitted to the Ethereum chain. They do this with a separate process called the "orchestrator", and send the signatures to the Cosmos chain as `MsgBatchTxSignature` messages. The Gravity module then checks that the signature is valid and stores it .

Relayers are then able to get all the signatures for a batch, assemble them into an Ethereum transaction, and send it to the Gravity.sol contract.

## ContractCallTx

### Logic call creation

Another module on the same Cosmos chain can call `SetContractCallTx` to create a logic call. All setting of parameters is left up to the external module.

### Logic call signing

Once a logic call has been created and stored, it is up to the current validators to sign it with their Ethereum keys so that it can be submitted to the Ethereum chain. They do this with a separate process called the "orchestrator", and send the signatures to the Cosmos chain as `MsgContractCallTxSignature` messages. The Gravity module then checks that the signature is valid and stores it.

Relayers are then able to get all the signatures for a logic call, assemble them into an Ethereum transaction, and send it to the Gravity.sol contract.

## SignerSetTx

### SignerSetTx creation

To create valsets:

- We get the all bonded validators using `StakingKeeper.GetBondedValidatorsByPower`.
- We get their Ethereum addresses and powers.
- We normalize their powers by dividing each validator's power by the sum of powers in the whole validator set.

We save this data in a `SignerSetTx`

### SignerSetTx signing

Once a valset has been created and stored, it is up to the current validators to sign it with their Ethereum keys so that it can be submitted to the Ethereum chain. They do this with a separate process called the "orchestrator", and send the signatures to the Cosmos chain as `MsgSignerSetTxSignature` messages. The Gravity module then checks that the signature is valid and stores it.

Relayers are then able to get all the signatures for a valset, assemble them into an Ethereum transaction, and send it to the Gravity.sol contract.
