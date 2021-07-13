# Orchestrator folder

### client/

This folder builds a binary that is a client application for the gravity system. It contains the following commands:
- `cosmos-to-eth`
- `eth-to-cosmos`
- `deploy-erc20-representation`

### cosmos_gravity/

This is a library for interacting with the cosmos chain both queries and transactions. It substantally wraps `gravity_proto`.

### ethereum_gravity/

This is a library that contains code for the interactions with the counterparty ethereum chain.

### gravity_proto/

`prost` generated bindings for working with the gravity protobuf objects.

### gravity_utils/

Various utilities for working with the `gravity` code.

### orchestrator/

The package to build the orchestartor binary.

### proto_build/

Run `cargo run` in this folder to build `gravity_proto` also note, this will generate too many files. Only `gravity.v1.rs` is required.

### register_delegate_keys/

This is a sepreate binary for running a command to register delegate keys for a validator. NOTE: this needs to be done in `gentx` now so this is likely no longer needed.

### relayer/

This is to build the relayer logic (i.e. cosmos to ethereum) as a seperate binary. It also contains the library for the relayer.

### scripts/

Supporting bash scripts for this library

### test_runner/

A binary which runs tests against a cosmos chain


## CLI

### CURRENT

```
client cosmos-to-eth --cosmos-phrase=<key> --cosmos-grpc=<url> --cosmos-prefix=<prefix> --cosmos-denom=<denom> --amount=<amount> --eth-destination=<dest> [--no-batch] [--times=<number>]
client eth-to-cosmos --ethereum-key=<key> --ethereum-rpc=<url> --cosmos-prefix=<prefix> --contract-address=<addr> --erc20-address=<addr> --amount=<amount> --cosmos-destination=<dest> [--times=<number>]
client deploy-erc20-representation --cosmos-grpc=<url> --cosmos-prefix=<prefix> --cosmos-denom=<denom> --ethereum-key=<key> --ethereum-rpc=<url> --contract-address=<addr> --erc20-name=<name> --erc20-symbol=<symbol> --erc20-decimals=<decimals>
orchestrator --cosmos-phrase=<key> --ethereum-key=<key> --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --fees=<denom> --contract-address=<addr>
register-delegate-key --validator-phrase=<key> --address-prefix=<prefix> [--cosmos-phrase=<key>] [--ethereum-key=<key>] --cosmos-grpc=<url> --fees=<denom>
relayer --ethereum-key=<key> --cosmos-grpc=<url> --address-prefix=<prefix> --ethereum-rpc=<url> --contract-address=<addr> --cosmos-grpc=<gurl>
test_runner 
```

## PROPOSED

Proposing the name `gorc` for the binary. This is short for `gravity-orchestrator`.

```
gorc
  tx
    eth
      send-to-cosmos [from-eth-key] [to-cosmos-addr] [erc20 conract] [erc20 amount] [[--times=int]]
      send [from-key] [to-addr] [amount] [token-contract]
    cosmos
      send-to-eth [from-cosmos-key] [to-eth-addr] [erc20-coin] [[--times=int]]
      send [from-key] [to-addr] [coin-amount]
  query
    eth
      balance [key-name]
      contract
    cosmos
      balance [key-name]
      gravity-keys [key-name]
  deploy
    cosmos-erc20 [denom] [erc20_name] [erc20_symbol] [erc20_decimals]
  start
    orchestrator [contract-address] [fee-denom]
    relayer
  tests
    runner
  keys
    eth
      add [name]
      import [name] [privkey]
      delete [name]
      update [name] [new-name]
      list
      show [name]
    cosmos 
      add [name]
      import [name] [mnemnoic]
      delete [name]
      update [name] [new-name]
      list
      show [name]
```

```json
[gravity]
	contract = "0x6b175474e89094c44da98b954eedeac495271d0f"
	
[ethereum]
key = "testkey"
rpc = "http://localhost:8545"

[cosmos]
key = "testkey"
grpc = "http://localhost:9090"
prefix = "cosmos"
```
