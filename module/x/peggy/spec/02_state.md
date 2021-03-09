<!--
order: 2
-->

# State

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the staking module.

- Params: `Paramsspace("peggy") -> legacy_amino(params)`

+++ <https://github.com/althea-net/cosmos-gravity-bridge/blob/main/module/proto/peggy/v1/genesis.proto#L72-L104>


### OutgoingTxBatch

Stored in two possible ways, first with a height and second without (unsafe)

### ValidatorSet

This is the validator set of the bridge.

Stored in two possible ways, first with a height and second without (unsafe)

### ValsetNonce

### SlashableValeSetNonce

### ValsetConfirmation

<!-- MsgValsetConfirm -->

### ConfirmBatch

<!-- MsgConfirmBatch -->

### OrchestratorValidator

<!-- sets orchestrator address of a validator -->

### EthValidator

<!-- sets eth address of a validator -->

### OutgoingLogicCall

### ConfirmLogicCall

### OutgoingTx

### IDS

### SlashedBlockHeight

### TokenContract

<!-- GetDenomToERC20Key(denom), []byte(tokenContract)) -->
### Denom

<!-- GetERC20ToDenomKey -->

### LastEventNonce

<!-- set by a validator and/or by anyone -->

### LastObservedHeight 

This is the last observed height on ethereum.

### Attestation
