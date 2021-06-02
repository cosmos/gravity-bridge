<!--
order: 4
-->

# Messages

In this section we describe the processing of the gravity messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](./02_state_transitions.md) section.

### MsgDelegateKeys

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

### MsgSubmitEthereumTxConfirmation

When the gravity daemon witnesses a complete validator set within the gravity module, the validator submits a signature of a message containing the entire validator set. 

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


### MsgSendToEthereum

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

### MsgRequestBatchTx

When enough transactions have been added into a batch, a user or validator can call send this message in order to send a batch of transactions across the bridge. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L122-125

This message will fail if:

- The denom is not supported.
- Failure to build a batch of transactions.
- If the orchestrator address is not present in the validator set

### MsgConfirmBatch

When a `MsgRequestBatchTx` is observed, validators need to sign batch request to signify this is not a maliciously created batch and to avoid getting slashed. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L137-143

This message will fail if:

- The batch does not exist
- If checkpoint generation fails
- If a none validator address or delegated address 
- If the counter chain address is empty or incorrect.
- If counter chain address fails signature validation
- If the signature was already presented in a previous message

### MsgConfirmLogicCall

When a logic call request has been made, it needs to be confirmed by the bridge validators. Each validator has to submit a confirmation of the logic call being executed.

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

### MsgDepositClaim

When a message to deposit funds into the gravity contract is created a event will be omitted and observed a message will be submitted confirming the deposit.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L170-181

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgWithdrawClaim

When a user requests a withdrawal from the gravity contract a event will omitted by the counter party chain. This event will be observed by a bridge validator and submitted to the gravity module.


+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L187-193

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgERC20DeployedClaim

This message allows the cosmos chain to learn information about the denom from the counter party chain.

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L200-209

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails

### MsgLogicCallExecutedClaim

This informs the chain that a logic call has been executed. This message is submitted by bridge validators when they observe a event containing details around the logic call. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/msgs.proto#L215-221

This message will fail if: 

- The validator submitting the claim is unknown
- The validator is not in the active set
- Creation of attestation has failed.
