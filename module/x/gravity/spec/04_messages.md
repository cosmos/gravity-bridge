<!--
order: 4
-->

# Messages

In this section we describe the processing of the gravity messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](./02_state_transitions.md) section.

### MsgSetOrchestratorAddress

Allows validators to delegate their voting responsibilities to a given key. This Key can be used to authenticate oracle claims.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L38-L40

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L56-60

This message is expected to fail if:

- The validator address is incorrect.
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails
- The orchestrator address is incorrect.
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails
- The ethereum address is incorrect.
  - The address is empty (`""`)
  - Not a length of 42
  - Does not start with 0x
- The validator is not present in the validator set.

### MsgSendToEth

When a user wants to bridge an asset to an EVM. If the token has originated from the cosmos chain it will be held in a module account. If the token is originally from ethereum it will be burned on the cosmos side.

> Note: this message will later be removed when it is included in a batch.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L100-109

This message will fail if:

- The sender address is incorrect.
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails
- The denom is not supported.
- If the token is cosmos originated
  - The sending of the token to the module account fails
- If the token is non-cosmos-originated.
  - If sending to the module account fails
  - If burning of the token fails

### MsgRequestBatch

Anyone can send this message to trigger [creation](03_state_transitions.md#batch-creation) of an `OutgoingTxBatch`.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L122-125

This message will fail if:

- The denom is not supported.
- Failure to build a batch of transactions.
- If the orchestrator address is not present in the validator set

### MsgConfirmBatch

Validators sign `OutgoingTxBatch`'s with their Ethereum keys, and send the signatures to the Gravity module using this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L137-143

This message will fail if:

- The batch does not exist
- If checkpoint generation fails
- If a none validator address or delegated address
- If the counter chain address is empty or incorrect.
- If counter chain address fails signature validation
- If the signature was already presented in a previous message

### MsgConfirmLogicCall

Validators sign `OutgoingLogicCall`'s with their Ethereum keys, and send the signatures to the Gravity module using this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L155-161

This message will fail if:

- The id encoding is incorrect
- The outgoing logic call which is confirmed can not be found
- Invalid checkpoint generation
- Signature decoding failed
- The address calling this function is not a validator or its delegated key
- The counter chain address is incorrect or empty
- Counter party signature verification failed
- A duplicate signature is observed

### MsgValsetConfirm

When a `Valset` is created by the Gravity module, validators sign it with their Ethereum keys, and send the signatures to the Gravity module using this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L79-84

This message is expected to fail if:

- If the validator set is not present.
- The signature is encoded incorrectly.
- Signature verification of the ethereum key fails.
- If the signature submitted has already been submitted previously.
- The validator address is incorrect.
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails

### MsgDepositClaim

When a user deposits funds into the Gravity.sol Ethereum contract an event is emitted by the contract. When validators observe this event, they send this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L170-181

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgWithdrawClaim

When a transaction batch is executed by the Gravity.sol Ethereum contract, an event is emitted by the contract. When validators observe this event, they send this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L187-193

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgERC20DeployedClaim

When a transaction batch is executed by the Gravity.sol Ethereum contract, an event is emitted by the contract. When validators observe this event, they send this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L200-209

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgLogicCallExecutedClaim

When a logic call is executed by the Gravity.sol Ethereum contract, an event is emitted by the contract. When validators observe this event, they send this message.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L215-221

This message will fail if:

- The validator submitting the claim is unknown
- The validator is not in the active set
- Creation of attestation has failed.
