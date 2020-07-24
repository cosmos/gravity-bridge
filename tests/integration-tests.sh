#!/bin/bash
set -eux
# WIP 
NODES=$1
for i in $(seq 1 $NODES);
do
NODE_IP="7.7.7.$i"
TX_FLAGS="--home /validator$i --keyring-backend test --from validator$i --trace --node=http://$NODE_IP:26657 --chain-id=peggy-test -y"

ETH_PRIVKEY=$(gen_eth_key)
echo "$ETH_PRIVKEY"
peggycli tx nameservice update-eth-addr "$ETH_PRIVKEY" $TX_FLAGS
done

QUERY_FLAGS="--home /validator1 --trace --node=http://7.7.7.1:26657 --chain-id=peggy-test"

sleep 10 # Wait for a block to mine (there must be a better way)

peggycli query nameservice valset $QUERY_FLAGS

# This worked in the terminal
# peggycli tx nameservice update-eth-addr 0x6f4cf6911b895d058cea5f9d0a9d65f60c5f212da16a3d4e180f902b277dc59a --home /validator2 --keyring-backend test --from validator2 --trace --node=http://7.7.7.2:26657 --chain-id=peggy-test