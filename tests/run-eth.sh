#!/bin/bash
# Starts the Ethereum testnet chain in the background

# init the genesis block
geth --identity "PeggyTestnet" \
--nodiscover \
--networkid 15 init /peggy/tests/ETHGenesis.json 

# etherbase is where rewards get sent
geth --identity "PeggyTestnet" --nodiscover \
--networkid 15 \
--mine \
--minerthreads=1 \
--verbosity "0" \
--etherbase=0xb2958b1537f37f5ab92f719d9a33ab0c79f8f8db