<!--
order: 4
-->

# Messages

In this section we describe the processing of the pegggy messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](./02_state_transitions.md) section.

### MsgSetOrchestratorAddress

Allows validators to delegate their voting responsibilities to a given key. This Key can be used to authenticate oracle claims. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L38-L40

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L56-60

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

### MsgValsetConfirm

When the peggy daemon witnesses a complete validator set within the peggy module, the validator submits a signature of a message containing the entire validator set. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L79-84

This message is expected to fail if:

- If the validator set is not present.
- The signature is encoded incorrectly.
- Signature verification of the ethereum key fails.
- If the signature submitted has already been submitted previously.
- The validator address is incorrect. 
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails


### MsgSendToEth

When a user wants to bridge an asset to an EVM. If the token has originated from the cosmos chain it will be held in a module account. If the token is originally from ethereum it will be burned on the cosmos side.

> Note: this message will later be removed when it is included in a batch.


+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L100-109

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

When enough transactions have been added into a batch, a user or validator can call send this message in order to send a batch of transactions across the bridge. 

+++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L122-125

This message will fail if:

- The denom is not supported.
- Failure to build a batch of transactions.
- If the orchestrator address is not present in the validator set

### MsgConfirmBatch

When a `MsgRequestBatch` is observed, validators 

### MsgConfirmLogicCall

### MsgDepositClaim

### MsgWithdrawClaim

### MsgERC20DeployedClaim

### MsgLogicCallExecutedClaim
