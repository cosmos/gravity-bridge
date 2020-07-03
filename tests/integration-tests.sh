#!/bin/bash
set -eux
# WIP 
NODES=3
for i in $(seq 1 $NODES);
do
NODE_IP="7.7.7.$i"
ETH_PRIVKEY=$(gen_eth_key)
echo "$ETH_PRIVKEY"
peggycli config node http://$NODE_IP:26657
peggycli tx nameservice update-eth-addr $ETH_PRIVKEY
done