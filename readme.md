# Althea Peggy

Althea Peggy is a simplified vision of a Cosmos <-> Ethereum 'peg zone' focused on maximum design simplicity and efficiency.

For now Althea Peggy is focused on only one of the major functions of an Ethereum peg zone. Bidirectional transfer of ERC20 assets originating on Ethereum to a Cosmos based chain.

An expansion of this feature set is expected, but will only be performed once basic transfers are in production. This gives us the opportunity to develop solid foundations and technology useful for the larger peg zone vision without getting bogged down by a larger feature surface.

## Status

Althea Peggy is under development and unaudited. Instructions for deployment and use are provided in the hope that they will be useful.

It is your responsibility to understand the financial, legal, and other risks of using this software. There is no guarantee of functionality or safety. You use Peggy entirely at your own risk.

You can keep up with the latest development by watching our [public standups](https://www.youtube.com/playlist?list=PL1MwlVJloJeyeE23-UmXeIx2NSxs_CV4b) feel free to join yourself and ask questions.

- Solidity Contract
  - [x] Multiple ERC20 support
  - [x] Tested with 100+ validators
  - [X] Unit tests for every throw condition
  - [x] Audit
- Cosmos Module
  - [x] Basic validator set syncing
  - [x] Basic transaction batch generation
  - [x] Ethereum -> Cosmos Token issuing
  - [x] Cosmos -> Ethereum Token issuing
  - [X] Bootstrapping
  - [ ] Genesis file save/load
  - [ ] Validator set syncing edge cases
  - [ ] Transaction batch edge cases
  - [ ] Relaying edge cases
  - [ ] Audit
- Orchestrator / Relayer
  - [x] Validator set update relaying
  - [x] Ethereum -> Cosmos Oracle
  - [x] Transaction batch relaying
  - [ ] Tendermint KMS support
  - [ ] Audit

## Design simplifications from the larger peg zone vision

- Validators are fully trusted to manage the bridge. Validator powers and votes are replicated on the Ethereum side so trust in bridge assets depends entirely on trust in the validator set of the peg zone chain. This has known problems where the assets in the bridge exceed the market cap of the native token. We accept these known issues in exchange for the dramatic design simplification combined with acceptable decentralization this design provides.
- The Althea Peggy Ethereum contract only supports ERC20 transfers and not arbitrary data. This helps keep the contract simple enough to optimize heavily and reach production quality quickly.
- The Relayer as a discrete binary only exists to facilitate the bridge fee market. All chain oracle functions will be integrated directly into the Gaiad binary. This makes it mandatory for peg zone validators to maintain a trusted Ethereum node and removes all trust and game theory implications that usually arise from independent relayers.

## Key Components you can run today

- A highly efficient way of mirroring Cosmos validator voting onto Ethereum. The Althea-Peggy solidity contract has validator set updates costing ~500,000 gas ($2 @ 20gwei) and transaction batches have a base cost of ~500,000 gas ($2 @ 20gwei). This is tested using a snapshot of the Cosmos Hub validator set, with 100+ unique validators. We hope to further reduce these gas costs, see `solidity/possible_optimizations.md` for more details. Batches may contain arbitrary numbers of transactions within the limits of ERC20 sends per block. Allowing for costs to be heavily amortized on high volume bridges. This code will likely be re-used in any iteration of Peggy.
- All up integration tests, we provide a one button integration test that deploys a full arbitrary validator Cosmos chain and testnet Geth chain for both development and integration test validation. We believe having a in depth test environment reflecting the full deployment and production-like use of the code is essential to productive development.

## Running the all up tests

These tests cover everything that's working in this repo (or not) currently they build the Ethereum contract, run the happy path tests and then build + deploy a small Cosmos blockchain in docker running the code in the `module` folder. This code doesn't do much at the moment as we're working on getting it to generate validator updates to relay to the Ethereum contract.

To run the test simply have docker installed and run.

`bash tests/all-up-test.sh`

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
4. Or, `docker exec -it peggy_test_instance /bin/bash` should allow you to access a shell inside the test container

Change the code, and when you want to test it again, restart `./tests/start-chains.sh` and run `./tests/run-tests.sh`.

### Explanation:

`./tests/build-container.sh` builds the base container and builds the peggy test zone for the first time. This results in a Docker container which contains cached Go dependencies (the base container).

`./tests/start-chains.sh` starts a test container based on the base container and copies the current source code (including any changes you have made) into it. It then builds the peggy test zone, benefiting from the cached Go dependencies. It then starts the Cosmos chain running on your new code. It also starts an Ethereum node. These nodes stay running in the terminal you started it in, and it can be useful to look at the logs. Be aware that this also mounts the peggy folder into the container, meaning changes you make will be reflected there.

`./tests/run-tests.sh` connects to the running test container and runs the integration test found in `./tests/integration-tests.sh`

### Tips for IDEs:

- Launch VS Code in /solidity with the solidity extension enabled to get inline typechecking of the solidity contract
- Launch VS Code in /module/app with the go extension enabled to get inline typechecking of the dummy cosmos chain

### Working inside the container

It can be useful to modify, recompile, and restart the testnet without restarting the container, for example if you are running a text editor in the container and would not like it to exit, or if you are editing dependencies stored in the container's `/go/` folder.

In this workflow, you can use `./tests/reload-code.sh` to recompile and restart the testnet without restarting the container.

For example, you can use VS Code's "Remote-Container" extension to attach to the running container started with `./tests/start-chains.sh`, then edit the code inside the container, restart the testnet with `./tests/reload-code.sh`, and run the tests with `./tests/integration-tests.sh`.

## Debugger

To use a stepping debugger in VS Code, follow the "Working inside the container" instructions above, but set up a one node testnet using `./tests/reload-code.sh 1`. Now kill the node with `pkill peggyd`. Start the debugger from within VS Code, and you will have a 1 node debuggable testnet.
