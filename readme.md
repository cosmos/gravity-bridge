To run tests, 

1. run `./tests/build-container.sh`
2. run `./tests/start-chains.sh`
3. switch to a new terminal and run `./tests/run-tests.sh`
4. Or, `docker exec -it peggy_test_instance /bin/bash` should allow you to access a shell inside the test container

Change the code, and when you want to test it again, restart `./tests/start-chains.sh` and run `./tests/run-tests.sh`. 

Explanation:

`./tests/build-container.sh` builds the base container and builds the peggy test zone for the first time. This results in a Docker container which contains cached Go dependencies (the base container).

`./tests/start-chains.sh` starts a test container based on the base container and copies the current source code (including any changes you have made) into it. It then builds the peggy test zone, benefiting from the cached Go dependencies. It then starts the Cosmos chain running on your new code. It also starts an Ethereum node. These nodes stay running in the terminal you started it in, and it can be useful to look at the logs.

`./tests/run-tests.sh` connects to the running test container and runs the integration test found in `./tests/integration-tests.sh`

Tips for IDEs:

- Launch VS Code in /solidity with the solidity extension enabled to get inline typechecking of the solidity contract
- Launch VS Code in /module/app with the go extension enabled to get inline typechecking of the dummy cosmos chain


