<!--
order: 2
-->

# State

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the staking module.

- Params: `Paramsspace("gravity") -> legacy_amino(params)`

+++ <https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/gravity/v1/genesis.proto#L72-L104>


### BatchTx

Stored in two possible ways, first with a height and second without (unsafe). Unsafe is used for testing and export and import of state.

| key          | Value | Type   | Encoding               |
|--------------|-------|--------|------------------------|
| `[]byte{0xa} + common.HexToAddress(tokenContract).Bytes() + nonce (big endian encoded)` | A batch of outgoing transactions | `types.BatchTx` | Protobuf encoded |

### ValidatorSet

This is the validator set of the bridge.

Stored in two possible ways, first with a height and second without (unsafe). Unsafe is used for testing and export and import of state.

| key          | Value | Type   | Encoding               |
|--------------|-------|--------|------------------------|
| `[]byte{0x2} + nonce (big endian encoded)` | Validator set | `types.Valset` | Protobuf encoded |

### ValsetNonce

The latest validator set nonce, this value is updated on every write. 

| key          | Value | Type   | Encoding               |
|--------------|-------|--------|------------------------|
| `[]byte{0xf6}` | Nonce | `uint64` | encoded via big endian |

### SlashedValeSetNonce

The latest validator set slash nonce. This is used to track which validator set needs to be slashed and which already has been. 

| Key            | Value | Type   | Encoding               |
|----------------|-------|--------|------------------------|
| `[]byte{0xf5}` | Nonce | uint64 | encoded via big endian |

### Validator Set Confirmation

When a validator signs over a validator set this is considered a `valSetConfirmation`, these are saved via the current nonce and the orchestrator address. 


| Key                                         | Value                  | Type                     | Encoding         |
|---------------------------------------------|------------------------|--------------------------|------------------|
| `[]byte{0x3} + (nonce + []byte(AccAddress)` | Validator Confirmation | `types.MsgSubmitEthereumTxConfirmation` | Protobuf encoded |

### ConfirmBatch

When a validator confirms a batch it is added to the confirm batch store. It is stored using the orchestrator, token contract and nonce as the key. 

| Key                                                                 | Value                        | Type                    | Encoding         |
|---------------------------------------------------------------------|------------------------------|-------------------------|------------------|
| `[]byte{0xe1} + common.HexToAddress(tokenContract).Bytes() + nonce + []byte(AccAddress)` | Validator Batch Confirmation | `types.MsgConfirmBatch` | Protobuf encoded |

### OrchestratorValidator

When a validator would like to delegate their voting power to another key. The value is stored using the orchestrator address as the key

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xe8} + []byte(AccAddress)` | Orchestrator address assigned by a validator | `[]byte` | Protobuf encoded |

### EthAddress

A validator has an associated counter chain address. 

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0x1} + []byte(ValAddress)` | Ethereum address assigned by a validator | `[]byte` | Protobuf encoded |


### ContractCallTx

When a user requests a logic call to be executed on an opposing chain it is stored in a store within the gravity module.

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xde} + []byte(invalidationId) + nonce (big endian encoded)` | A user created logic call to be sent to the counter chain | `types.ContractCallTx` | Protobuf encoded |

### ConfirmLogicCall

When a logic call is executed validators confirm the execution. 

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
|`[]byte{0xae} + []byte(invalidationId) + nonce (big endian encoded) + []byte(AccAddress)` | Confirmation of execution of the logic call | `types.MsgConfirmLogicCall` | Protobuf encoded |

### OutgoingTx

Sets an outgoing transactions into the applications transaction pool to be included into a batch. 

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0x6} + id (big endian encoded)` | User created transaction to be included in a batch | `types.OutgoingTx` | Protobuf encoded |

### IDS

### SlashedBlockHeight

Represents the latest slashed block height. There is always only a singe value stored. 

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xf7}` | Latest height a batch slashing occurred | `uint64` | Big endian encoded |

### TokenContract & Denom

A denom that is originally from a counter chain will be from a contract. The toke contract and denom are stored in two ways. First, the denom is used as the key and the value is the token contract. Second, the contract is used as the key, the value is the denom the token contract represents. 

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xf3} + []byte(denom)` | Token contract address | `[]byte` | stored in byte format |

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xf4} + common.HexToAddress(tokenContract).Bytes()` | Latest height a batch slashing occurred | `[]byte` | stored in byte format |

### LastEventNonce

The last observed event nonce. This is set when `TryAttestation()` is called. There is always only a single value held in this store.

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xf2}` | Last observed event nonce| `uint64` | Big endian encoded |

### LastObservedEthereumHeight 

This is the last observed height on ethereum. There will always only be a single value stored in this store.

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0xf9}` | Last observed Ethereum Height| `uint64` | Protobuf encoded |

### Attestation

| Key                                 | Value                                        | Type     | Encoding         |
|-------------------------------------|----------------------------------------------|----------|------------------|
| `[]byte{0x5} + evenNonce (big endian encoded) + []byte(claimHash)` | Attestation of occurred events/claims| `types.Attestation` | Protobuf encoded |
