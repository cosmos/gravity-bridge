#!/bin/bash
# Builds and runs Solidity tests within a container

pushd /gravity/solidity/
rm -rf node_modules
HUSKY_SKIP_INSTALL=1 npm install
npm run typechain
npm run evm &
npm run test