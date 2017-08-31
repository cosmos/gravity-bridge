# ETGate

Send ethereum tokens to tendermint zones.

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

## Functionality

### Deposit

ETGate contract will make a log when it receives ethers/erc20s. ETGate chain will create an IBC packet when it receives contract logs by relayers.

### Withdraw

ETGate chain will create an ETGate packet when it receives IBC messages for withdrawal. ETGate contract will release its tokens when it detects the packet.

### Transfer(Planned)

Any zone can send an IBC message to ETGate chain to transfer its tokens to another zone. When ETGate chain receives the message, it will create two packets: one for destination zone, one for ETGate contract.

### Validator Update(Planned)

For maximum security, validator list will be maintained on ETGate contract(not chain). Validator will be determined like as DPOS consensus blockchains. 

### Fraud Proof(Planned)

As ETGate only supports three tx types, it is technically possible to make a fraud proof of invalid blocks. If +2/3 of validators are byzantine and they try to create a invalid block, anybody can upload fraud proof and the validators will be punished. 
