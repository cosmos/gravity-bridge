# ETGate

Send ethereum contract logs to tendermint blockchain.

This code is not yet completed, lacking features such as uploading ethereum logs and create packets towards to ethereum. Currently, it supports only uploading ethereum headers.

## Usage

1. Install [golang](https://golang.org/dl), [geth](https://github.com/ethereum/go-ethereum), and [basecoin](https://github.com/tendermint/basecoin).
2. Git clone & go build this repository.
3. cd to cmd/etgate
4. Run `./init.sh`.
5. Run `geth --testnet --fast` on the other window.
6. Run `./etgate start &> etgate.log &`
7. Run `./etgate gate start --testnet --nodeaddr=tcp://localhost:12347 ../../static/example.json`
