#!/bin/bash
npx ts-node \
contract-deployer.ts \
--cosmos-node="http://peggy_0:26657" \
--eth-node="http://ethereum:8545" \
--eth-privkey="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7" \
--contract=Peggy.json \
--test-mode=true