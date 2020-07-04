#!/bin/bash
set -eux
# WIP 
NODES=3
for i in $(seq 1 $NODES);
do
NODE_IP="7.7.7.$i"
FLAGS="--home /validator$i --keyring-backend test --from validator$i --trace --node=http://$NODE_IP:26657 --chain-id=peggy-test"
ETH_PRIVKEY=$(gen_eth_key)
echo "$ETH_PRIVKEY"
# peggycli config node http://$NODE_IP:26657
peggycli tx nameservice update-eth-addr "$ETH_PRIVKEY" $FLAGS
# peggycli tx nameservice update-eth-addr $ETH_PRIVKEY --home "/validator$i" --keyring-backend test --from "validator$i"
done

# This worked in the terminal
# peggycli tx nameservice update-eth-addr 0x6f4cf6911b895d058cea5f9d0a9d65f60c5f212da16a3d4e180f902b277dc59a --home /validator2 --keyring-backend test --from validator2 --trace --node=http://7.7.7.2:26657 --chain-id=peggy-test