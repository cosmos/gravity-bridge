## Building

On first run:
make proto-update-deps
make proto-tools
To build:
make

## Early MVP

Happy path implementations

### Oracle

#### Assumptions

- An orchestrator may want to submit multiple claims with a msg (withdrawal batch update + MultiSig Set update )
- Nonces are not unique without a context (withdrawal nonce and MultiSig Set update can have same nonce (=height))
- A nonce is unique in it's context and never reused
- Multiple claims by an orchestrator for the same ETH event are forbidden
- We know the ETH event types beforehand (and handle them as ClaimTypes)
- For an **observation** status in Attestation the power AND count thresholds must be exceeded
- Fraction type allows higher precision math than %. For example with 2/3

A good start to follow the process would be the `x/gravity/handler_test.go` file

### Outgoing TX Pool

#### Features

- Unique denominator for gravity vouchers in cosmos (🚧 cut to 15 chars and without a separator due to sdk limitations in v0.38.4)
- Voucher burning 🔥 (minting in test ⛏️ )
- Store/ resolve bridged ETH denominator and contract
- Persistent transaction pool
- Transactions sorted by fees (on a second index)
- Extended test setup

#### Assumptions

- We have only 1 chainID and 1 ETH contract

### Bundle Outgoing TX into Batches

#### Features

- `BatchTx` type with `OutgoingTransferTx` and `TransferCoin`
- Logic to build batch from pending TXs based on fee desc order
- Logic to cancel a batch and revert TXs back to pending pool
- Incremental and unique IDs for batches to be used for `nonces`
- `VoucherDenom` as first class type

## Not covered/ implemented

- [ ] unhappy cases
- [ ] proper unit + integration tests
- [ ] message validation
- [ ] Genesis I/O
- [ ] Parameters
- [ ] authZ: EthereumChainID whitelisted
- [ ] authZ: bridge contract address whitelisted
