#!/bin/bash
npx ts-node \
contract-deployer.ts \
--cosmos-node="http://localhost:26657" \
--eth-node="http://localhost:8545" \
--eth-privkey="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7" \
--peggy-id="defaultpeggyid" \
--contract=artifacts/Peggy.json \
--erc20-contract=artifacts/TestERC20.json \
--test-mode=true