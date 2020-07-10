#!/bin/bash
# Builds and runs Solidity tests within a container

pushd /peggy/solidity/
npm install
npm run typechain
npm run evm &
npm run test