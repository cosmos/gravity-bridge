<!--
order: 2
-->

# State

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the staking module.

- Params: `Paramsspace("gravity") -> legacy_amino(params)`

+++ <https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/genesis.proto#L72-L104>

### OutgoingTxBatch

Stored in two possible ways, first with a height and second without (unsafe). Unsafe is used for testing and export and import of state.

| key                                                                | Value                            | Type                    | Encoding         |
| ------------------------------------------------------------------ | -------------------------------- | ----------------------- | ---------------- |
| `[]byte{0xa} + []byte(tokenContract) + nonce (big endian encoded)` | A batch of outgoing transactions | `types.OutgoingTxBatch` | Protobuf encoded |

```
message OutgoingTxBatch {
  // The batch_nonce is an incrementing nonce which is assigned to a batch on creation.
  // The Gravity.sol Ethereum contract stores the last executed batch nonce for each token type
  // and it will only execute batches with a lower nonce. Note that the nonce sequence is
  // PER TOKEN, i.e. GRT tokens could have a last executed nonce of 3002 while DAI tokens had a nonce of 4556
  // This property is important for creating batches that are profitable to submit, which is covered in greater
  // detail in the [state transitions spec](03_state_transitions.md)
  uint64                      batch_nonce    = 1;
  // The batch_timeout is an Ethereum block at which this batch will no longer be executed by Gravity.sol. This
  // allows us to cancel batches that we know have timed out, releasing their transactions to be included in a new batch
  // or cancelled by their sender.
  uint64                      batch_timeout  = 2;
  // These are the transactions sending tokens to destinations on Ethereum.
  repeated OutgoingTransferTx transactions   = 3;
  // This is the token contract of the tokens that are being sent in this batch.
  string                      token_contract = 4;
  // The Cosmos block height that this batch was created. This is used in slashing.
  uint64                      block          = 5;
}
```

### ValidatorSet

This is the validator set of the bridge.

Stored in two possible ways, first with a height and second without (unsafe). Unsafe is used for testing and export and import of state.

| key                                        | Value         | Type           | Encoding         |
| ------------------------------------------ | ------------- | -------------- | ---------------- |
| `[]byte{0x2} + nonce (big endian encoded)` | Validator set | `types.Valset` | Protobuf encoded |

### ValsetNonce

The latest validator set nonce, this value is updated on every write.

| key            | Value | Type     | Encoding               |
| -------------- | ----- | -------- | ---------------------- |
| `[]byte{0xf6}` | Nonce | `uint64` | encoded via big endian |

### SlashedValeSetNonce

The latest validator set slash nonce. This is used to track which validator set needs to be slashed and which already has been.

| Key            | Value | Type   | Encoding               |
| -------------- | ----- | ------ | ---------------------- |
| `[]byte{0xf5}` | Nonce | uint64 | encoded via big endian |

### Validator Set Confirmation

When a validator signs over a validator set this is considered a `valSetConfirmation`, these are saved via the current nonce and the orchestrator address.

| Key                                         | Value                  | Type                     | Encoding         |
| ------------------------------------------- | ---------------------- | ------------------------ | ---------------- |
| `[]byte{0x3} + (nonce + []byte(AccAddress)` | Validator Confirmation | `types.MsgValsetConfirm` | Protobuf encoded |

### ConfirmBatch

When a validator confirms a batch it is added to the confirm batch store. It is stored using the orchestrator, token contract and nonce as the key.

| Key                                                                 | Value                        | Type                    | Encoding         |
| ------------------------------------------------------------------- | ---------------------------- | ----------------------- | ---------------- |
| `[]byte{0xe1} + []byte(tokenContract) + nonce + []byte(AccAddress)` | Validator Batch Confirmation | `types.MsgConfirmBatch` | Protobuf encoded |

### OrchestratorValidator

When a validator would like to delegate their voting power to another key. The value is stored using the orchestrator address as the key

| Key                                 | Value                                        | Type     | Encoding         |
| ----------------------------------- | -------------------------------------------- | -------- | ---------------- |
| `[]byte{0xe8} + []byte(AccAddress)` | Orchestrator address assigned by a validator | `[]byte` | Protobuf encoded |

### EthAddress

A validator has an associated Ethereum address.

| Key                                | Value                                    | Type     | Encoding         |
| ---------------------------------- | ---------------------------------------- | -------- | ---------------- |
| `[]byte{0x1} + []byte(ValAddress)` | Ethereum address assigned by a validator | `[]byte` | Protobuf encoded |

### OutgoingLogicCall

```
message OutgoingLogicCall {
  // This is the address of the logic contract that Gravity.sol will call
  string              logic_contract_address = 3;
  // This is the content of the function call on the logic contract. It is formatted
  // as an Ethereum function call which can be passed to .call() and contains the function
  // name and arguments.
  bytes               payload                = 4;
  // The timeout is an Ethereum block at which this logic call will no longer be executed by Gravity.sol. This
  // allows the calling module to cancel logic calls that we know have timed out.
  uint64              timeout                = 5;
  // These are ERC20 transfers to the logic contract that take place before the logic call is made. This is useful
  // if the logic contract implements logic that deals with tokens.
  repeated ERC20Token transfers              = 1;
  // These are fees that go to the relayer of the logic call.
  repeated ERC20Token fees                   = 2;
  // The invalidation_id and invalidation_nonce provide a way for the calling module to implement a variety of
  // replay protection/invalidation strategies. The rules are simple: When a logic call is submitted to the
  // Gravity.sol contract, it will not be executed if a previous logic call with the same invalidation_id
  // and an equal or higher invalidation_nonce was executed previously. To use a strategy where a submitted logic
  // call invalidates all earlier unsubmitted logic calls, the calling module would simply keep the invalidation_id
  // the same and increment the invalidation_nonce. To implement a strategy where logic calls do not invalidate each other
  // at all, and are only invalidated by timing out, the calling module would increment the invalidation_id with each call.
  // To implement a strategy identical to the one used by this module for transaction batches, the calling module would set the
  // invalidation_id to the token contract, and increment the invalidation_nonce.
  bytes               invalidation_id        = 6;
  uint64              invalidation_nonce     = 7;
}
```

When another module requests a logic call to be executed on Ethereum it is stored in a store within the gravity module.

| Key                                                                  | Value                                                | Type                      | Encoding         |
| -------------------------------------------------------------------- | ---------------------------------------------------- | ------------------------- | ---------------- |
| `[]byte{0xde} + []byte(invalidationId) + nonce (big endian encoded)` | A user created logic call to be sent to the Ethereum | `types.OutgoingLogicCall` | Protobuf encoded |

### ConfirmLogicCall

When a logic call is executed validators confirm the execution.

| Key                                                                                       | Value                                       | Type                        | Encoding         |
| ----------------------------------------------------------------------------------------- | ------------------------------------------- | --------------------------- | ---------------- |
| `[]byte{0xae} + []byte(invalidationId) + nonce (big endian encoded) + []byte(AccAddress)` | Confirmation of execution of the logic call | `types.MsgConfirmLogicCall` | Protobuf encoded |

### OutgoingTx

Sets an outgoing transactions into the applications transaction pool to be included into a batch.

| Key                                     | Value                                              | Type               | Encoding         |
| --------------------------------------- | -------------------------------------------------- | ------------------ | ---------------- |
| `[]byte{0x6} + id (big endian encoded)` | User created transaction to be included in a batch | `types.OutgoingTx` | Protobuf encoded |

### IDS

### SlashedBlockHeight

Represents the latest slashed block height. There is always only a singe value stored.

| Key            | Value                                   | Type     | Encoding           |
| -------------- | --------------------------------------- | -------- | ------------------ |
| `[]byte{0xf7}` | Latest height a batch slashing occurred | `uint64` | Big endian encoded |

### Cosmos originated ERC20 - TokenContract & Denom

This is how we associate the ERC20 contracts representing Cosmos originated assets with the asset's denom on Cosmos. First, the denom is used as the key and the value is the token contract. Second, the contract is used as the key, the value is the denom the token contract represents.

| Key                            | Value                  | Type     | Encoding              |
| ------------------------------ | ---------------------- | -------- | --------------------- |
| `[]byte{0xf3} + []byte(denom)` | Token contract address | `[]byte` | stored in byte format |

| Key                                    | Value                                   | Type     | Encoding              |
| -------------------------------------- | --------------------------------------- | -------- | --------------------- |
| `[]byte{0xf4} + []byte(tokenContract)` | Latest height a batch slashing occurred | `[]byte` | stored in byte format |

### LastEventNonce

The last observed event nonce. This is set when `TryAttestation()` is called. There is always only a single value held in this store.

| Key            | Value                     | Type     | Encoding           |
| -------------- | ------------------------- | -------- | ------------------ |
| `[]byte{0xf2}` | Last observed event nonce | `uint64` | Big endian encoded |

### LastObservedEthereumHeight

This is the last observed height on ethereum. There will always only be a single value stored in this store.

| Key            | Value                         | Type     | Encoding         |
| -------------- | ----------------------------- | -------- | ---------------- |
| `[]byte{0xf9}` | Last observed Ethereum Height | `uint64` | Protobuf encoded |

### Attestation

This is a record of all the votes for a given claim (Ethereum event).

| Key                                                                 | Value                                 | Type                | Encoding         |
| ------------------------------------------------------------------- | ------------------------------------- | ------------------- | ---------------- |
| `[]byte{0x5} + eventNonce (big endian encoded) + []byte(claimHash)` | Attestation of occurred events/claims | `types.Attestation` | Protobuf encoded |

```
message Attestation {
  // This field stores whether the Attestation has had its event applied to the Cosmos state. This happens when
  // enough (usually >2/3s) of the validator power votes that they saw the event on Ethereum.
  // For example, once a DepositClaim has modified the token balance of the account that it was deposited to,
  // this boolean will be set to true.
  bool observed = 1;
  // This is an array of the addresses of the validators which have voted that they saw the event on Ethereum.
  repeated string votes = 2;
  // This is the Cosmos block height that this event was first observed by a validator.
  uint64 height = 3;
  // The claim is the Ethereum event that this attestation is recording votes for.
  google.protobuf.Any claim = 4;
}
```

+++ <https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/attestation.proto#L38-L43>

### Valset

This is a record of the Cosmos validator set at a given moment. Can be sent to the Gravity.sol contract to update the signer set.

| Key                                 | Value                      | Type           | Encoding         |
| ----------------------------------- | -------------------------- | -------------- | ---------------- |
| `[]byte{0x5} + uint64 valset nonce` | Validator set for Ethereum | `types.Valset` | Protobuf encoded |

```
message Valset {
  // This nonce is incremented for each subsequent valset produced by the Gravity module.
  // The Gravity.sol contract will only accept a valset with a higher nonce than the last
  // executed Valset.
  uint64 nonce = 1;
  // The validators in the valset.
  repeated BridgeValidator members = 2;
  // TODO: what is this for?
  uint64 height = 3;
}
```
