#!/bin/bash
# Builds and runs Solidity tests within a container

pushd /peggy/solidity/
rm -rf node_modules
npm install
npm run typechain
npm run evm &
npm run test