# ETGate

Send ethereum tokens to tendermint zones.

## Usage

### Start etgate

1. Install [golang](https://golang.org/dl), [geth](https://github.com/ethereum/go-ethereum), and [basecoin](https://github.com/tendermint/basecoin).
2. Git clone this repository.
3. cd to cmd/etgate and go build.
4. Run `geth --testnet --fast` on the other window and wait until sync is completed.
5. Run `./init.sh`
6. Run `./etgate gate start --testnet --nodeaddr=tcp://localhost:12347`

### Deposit/Withdraw ethers

1. cd to server/ and go build
2. Run `./server`
3. Open [http://localhost:12349](http://localhost:12349) on your web browser
4. Paste your deployed contract's address and name of the key(money)
5. Type the amount of ethers you want to deposit and press Deposit
6. Type the destination address, the amount of ethers, password of your tendermint account, and press Withdraw

## Demo

[![Demo](https://img.youtube.com/vi/2vtTLzYZE-o/0.jpg)](https://www.youtube.com/watch?v=2vtTLzYZE-o)

## Features

### Deposit

Each time when users send deposit message to the contract, it will generate an event. The relayers, on the other hand, will consistently upload ethereum headers to the tendermint zone. Since an ethereum header contains the merkle root of the receipts, the events could be proven with proving merkle path. This scheme is called [sidechaining](http://www.rsk.co/blog/sidechains-drivechains-and-rsk-2-way-peg-design). The zone will mint new coins right after the relayers uploaded the events.

### Withdraw

NOTE: As [merkleeyes](https://github.com/tendermint/merkleeyes/tree/master/iavl) uses ripemd160 instead of keccak-256, it is extremely expensive to prove data in a tendermint zone with solidity(nearly 0.1ether). For now, ETGate uses 2/3+ validator's multisig for each withdrawal to be validated.

The relayers also upload tendermint headers to the contract. The contract will confirm the header when 2/3+ validators signed. After the corresponding header is confirmed, the users can withdraw their tokens via submitting the necessary data(destination, value, etc) and its merkle proof. The contract will verify the proof, and release the tokens.
