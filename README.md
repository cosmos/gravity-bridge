# ETGate

Send ethereum contract logs to tendermint blockchain.

~~This code is not yet completed, lacking features such as uploading ethereum logs and create packets towards to ethereum. Currently, it supports only uploading ethereum headers.~~

**2017/08/03** Now etgate supports log uploading.

## Usage

### Start etgate

1. Install [golang](https://golang.org/dl), [geth](https://github.com/ethereum/go-ethereum), and [basecoin](https://github.com/tendermint/basecoin).
2. Git clone this repository.
3. cd to cmd/etgate and go build.
4. Run `./init.sh`.
5. Run `geth --testnet --fast` on the other window and wait until sync is completed.
6. Run `./etgate gate start --testnet --nodeaddr=tcp://localhost:12347 ../../static/example.json`

### Deposit ethers

1. Run `basecli --home ~/.etgate/client keys list`, and copy the address named "money".
2. On Ropsten testnet, call function `deposit` of contract [`0xe991802b4d8a6a544c303b623fdf02ecc13d26ae`](https://ropsten.etherscan.io/address/0xe991802b4d8a6a544c303b623fdf02ecc13d26ae) with your "money" address as `to` argument.
3. Wait until etgate submits the log. It takes a few minutes.
4. Run `basecli --home ~/.etgate/client query account *(your "money" address)*` to check your balance.
