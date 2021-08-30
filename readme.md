![Gravity Bridge](./gravity-bridge.svg)

Gravity bridge is Cosmos <-> Ethereum bridge designed to run on the [Cosmos Hub](https://github.com/cosmos/gaia) focused on maximum design simplicity and efficiency.

Gravity can transfer ERC20 assets originating on Ethereum to a Cosmos based chain and back to Ethereum.

The ability to transfer assets originating on Cosmos to an ERC20 representation on Ethereum is coming within a few months.

## Status

Gravity bridge is under development and will be undergoing audits soon. Instructions for deployment and use are provided in the hope that they will be useful.

It is your responsibility to understand the financial, legal, and other risks of using this software. There is no guarantee of functionality or safety. You use Gravity entirely at your own risk.

You can keep up with the latest development by watching our [public standups](https://www.youtube.com/playlist?list=PL1MwlVJloJeyeE23-UmXeIx2NSxs_CV4b) feel free to join yourself and ask questions.

- Solidity Contract
  - [x] Multiple ERC20 support
  - [x] Tested with 100+ validators
  - [x] Unit tests for every throw condition
  - [x] Audit for assets originating on Ethereum
  - [ ] Support for issuing Cosmos assets on Ethereum
- Cosmos Module
  - [x] Basic validator set syncing
  - [x] Basic transaction batch generation
  - [x] Ethereum -> Cosmos Token issuing
  - [x] Cosmos -> Ethereum Token issuing
  - [x] Bootstrapping
  - [x] Genesis file save/load
  - [x] Validator set syncing edge cases
  - [x] Slashing
  - [x] Relaying edge cases
  - [ ] Transaction batch edge cases
  - [ ] Support for issuing Cosmos assets on Ethereum
  - [ ] Audit
- Orchestrator / Relayer
  - [x] Validator set update relaying
  - [x] Ethereum -> Cosmos Oracle
  - [x] Transaction batch relaying
  - [ ] Tendermint KMS support
  - [ ] Audit

## The design of Gravity Bridge

- Trust in the integrity of the Gravity bridge is anchored on the Cosmos side. The signing of fraudulent validator set updates and transaction batches meant for the Ethereum contract is punished by slashing on the Cosmos chain. If you trust the Cosmos chain, you can trust the Gravity bridge operated by it, as long as it is operated within certain parameters.
- It is mandatory for peg zone validators to maintain a trusted Ethereum node. This removes all trust and game theory implications that usually arise from independent relayers, once again dramatically simplifying the design.

## Key design Components

- A highly efficient way of mirroring Cosmos validator voting onto Ethereum. The Gravity solidity contract has validator set updates costing ~500,000 gas ($2 @ 20gwei), tested on a snapshot of the Cosmos Hub validator set with 125 validators. Verifying the votes of the validator set is the most expensive on chain operation Gravity has to perform. Our highly optimized Solidity code provides enormous cost savings. Existing bridges incur more than double the gas costs for signature sets as small as 8 signers.
- Transactions from Cosmos to ethereum are batched, batches have a base cost of ~500,000 gas ($2 @ 20gwei). Batches may contain arbitrary numbers of transactions within the limits of ERC20 sends per block, allowing for costs to be heavily amortized on high volume bridges.

## Operational parameters ensuring security

- There must be a validator set update made on the Ethereum contract by calling the `updateValset` method at least once every Cosmos unbonding period (usually 2 weeks). This is because if there has not been an update for longer than the unbonding period, the validator set stored by the Ethereum contract could contain validators who cannot be slashed for misbehavior.
- Cosmos full nodes do not verify events coming from Ethereum. These events are accepted into the Cosmos state based purely on the signatures of the current validator set. It is possible for the validators with >2/3 of the stake to put events into the Cosmos state which never happened on Ethereum. In this case observers of both chains will need to "raise the alarm". We have built this functionality into the relayer.

## Run Gravity bridge right now using docker

We provide a one button integration test that deploys a full arbitrary validator Cosmos chain and testnet Geth chain for both development + validation. We believe having a in depth test environment reflecting the full deployment and production-like use of the code is essential to productive development.

Currently on every commit we send hundreds of transactions, dozens of validator set updates, and several transaction batches in our test environment. This provides a high level of quality assurance for the Gravity bridge.

Because the tests build absolutely everything in this repository they do take a significant amount of time to run. You may wish to simply push to a branch and have Github CI take care of the actual running of the tests.

To run the test simply have docker installed and run.

`bash tests/all-up-test.sh`

There are optional tests for specific features

Valset stress changes the validating power randomly 25 times, in an attempt to break validator set syncing

`bash tests/all-up-test.sh VALSET_STRESS`

Batch stress sends 300 transactions over the bridge and then 3 batches back to Ethereum. This code can do up to 10k transactions but Github Actions does not have the horsepower.

`bash tests/all-up-test.sh BATCH_STRESS`

Validator out tests a validator that is not running the mandatory Ethereum node. This validator will be slashed and the bridge will remain functioning.

`bash tests/all-up-test.sh VALIDATOR_OUT`

# Developer guide

## Solidity Contract

in the `solidity` folder

Run `HUSKY_SKIP_INSTALL=1 npm install`, then `npm run typechain`.

Run `npm run evm` in a separate terminal and then

Run `npm run test` to run tests.

After modifying solidity files, run `npm run typechain` to recompile contract
typedefs.

The Solidity contract is also covered in the Cosmos module tests, where it will be automatically deployed to the Geth test chain inside the development container for a micro testnet every integration test run.

## Cosmos Module

We provide a standard container-based development environment that automatically bootstraps a Cosmos chain and Ethereum chain for testing. We believe standardization of the development environment and ease of development are essential so please file issues if you run into issues with the development flow.

### Go unit tests

These do not run the entire chain but instead test parts of the Go module code in isolation. To run them, go into `/module` and run `make test`

### To hand test your changes quickly

This method is dictinct from the all up test described above. Although it runs the same components it's much faster when editing individual components.

1. run `./tests/build-container.sh`
2. run `./tests/start-chains.sh`
3. switch to a new terminal and run `./tests/run-tests.sh`
4. Or, `docker exec -it gravity_test_instance /bin/bash` should allow you to access a shell inside the test container

Change the code, and when you want to test it again, restart `./tests/start-chains.sh` and run `./tests/run-tests.sh`.

### Explanation:

`./tests/build-container.sh` builds the base container and builds the Gravity test zone for the first time. This results in a Docker container which contains cached Go dependencies (the base container).

`./tests/start-chains.sh` starts a test container based on the base container and copies the current source code (including any changes you have made) into it. It then builds the Gravity test zone, benefiting from the cached Go dependencies. It then starts the Cosmos chain running on your new code. It also starts an Ethereum node. These nodes stay running in the terminal you started it in, and it can be useful to look at the logs. Be aware that this also mounts the Gravity repo folder into the container, meaning changes you make will be reflected there.

`./tests/run-tests.sh` connects to the running test container and runs the integration test found in `./tests/integration-tests.sh`

### Tips for IDEs:

- Launch VS Code in /solidity with the solidity extension enabled to get inline typechecking of the solidity contract
- Launch VS Code in /module/app with the go extension enabled to get inline typechecking of the dummy cosmos chain

### Working inside the container

It can be useful to modify, recompile, and restart the testnet without restarting the container, for example if you are running a text editor in the container and would not like it to exit, or if you are editing dependencies stored in the container's `/go/` folder.

In this workflow, you can use `./tests/reload-code.sh` to recompile and restart the testnet without restarting the container.

For example, you can use VS Code's "Remote-Container" extension to attach to the running container started with `./tests/start-chains.sh`, then edit the code inside the container, restart the testnet with `./tests/reload-code.sh`, and run the tests with `./tests/integration-tests.sh`.

## Debugger

To use a stepping debugger in VS Code, follow the "Working inside the container" instructions above, but set up a one node testnet using `./tests/reload-code.sh 1`. Now kill the node with `pkill gravityd`. Start the debugger from within VS Code, and you will have a 1 node debuggable testnet.
