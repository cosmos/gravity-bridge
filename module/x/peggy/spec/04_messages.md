<!--
order: 4
-->

# Messages

In this section we describe the processing of the pegggy messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](./02_state_transitions.md) section.

### MsgDelegateKeys

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

### MsgSubmitConfirm

When the peggy daemon witnesses an event on the counter chain the validator will submit confirmation of the event. 

<!-- +++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L79-84 -->

```proto
  message MsgSubmitClaim {
    google.protobuf.Any confirm = 1 [
        (cosmos_proto.accepts_interface) = "EthereumClaim"];
    string signer = 2;
  }
```

There are three types of confirmations a validator can submit.

- `CONFIRM_TYPE_VALSET`
  - When the peggy daemon witnesses a complete validator set within the peggy module, the validator submits a signature of a message containing the entire validator set. 
- `CONFIRM_TYPE_LOGIC`
  - When a logic call request has been made, it needs to be confirmed by the bridge validators. Each validator has to submit a confirmation of the logic call being executed.
- `CONFIRM_TYPE_BATCH`
  - When a `MsgRequestBatch` is observed, validators need to sign batch request to signify this is not a maliciously created batch and to avoid getting slashed. 

This message is expected to fail if:

- If the validator set is not present.
- The signature is encoded incorrectly.
- Signature verification of the ethereum key fails.
- If the signature submitted has already been submitted previously.
- The validator address is incorrect. 
  - The address is empty (`""`)
  - Not a length of 20
  - Bech32 decoding fails
  
### MsgSubmitClaim

When a message to deposit funds into the peggy contract is created a event will be omitted and observed a message will be submitted confirming the deposit.

<!-- +++ https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/msgs.proto#L170-181 -->

There are four types of claims:

- `CLAIM_TYPE_DEPOSIT`
- `CLAIM_TYPE_WITHDRAW`
  - When a user requests a withdrawal from the peggy contract a event will omitted by the counter party chain. This event will be observed by a bridge validator and submitted to the gravity module.
  
- `CLAIM_TYPE_ERC20_DEPLOYED`
  - This message allows the cosmos chain to learn information about the denom from the counter party chain.
- `CLAIM_TYPE_LOGIC_CALL_EXECUTED`
  - This informs the chain that a logic call has been executed. This message is submitted by bridge validators when they observe a event containing details around the logic call.

This message will fail if:

- The validator is unknown
- The validator is not in the active set
- If the creation of attestation fails
