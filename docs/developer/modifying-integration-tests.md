# Modifying the integration tests

This document is a short guide on how to add new integration tests.

Before starting, you should read the entire [design](/docs/design) and [developer](/docs/developer) document sets. As well as
run the [environment setup](/docs/developer/environment-setup.md)

## Basic structure

The integration tests build and launch a dockerized network with an Ethereum node,
a 4 validator Cosmos chain running the Gravity bridge, associated Orchestrator nodes, a 'contract deployer', and a 'test runner'.

The [test runner](/orchestrator/test_runner/src/main.rs) is a single rust binary that coordinates the actual test logic.

A local Geth instance, with its version defined in the [dockerfile](/ethereum/Dockerfile).

The contract deployer contains logic for  parsing the resulting ERC20 and Gravity.sol contract addresses. This is all done before we get into starting the actual test logic.

## Adding tests

In order to add a new test define a new test_type environmental variable in the test runners `main.rs` file from there you can create a new file containing the test logic templated off of the various existing examples.

The [happy_path_test](/orchestrator/test_runner/src/happy_path.rs) for example uses several repeatable utility functions to test validator set updates.

Every test should perform some action and then meticulously verify that it actually took place. It is especially important to go off the happy path and ensure correct functionality.
