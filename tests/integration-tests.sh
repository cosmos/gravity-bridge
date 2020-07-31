#!/bin/bash
set -eu
# WIP 
NODES=3 # Permanently set to 3 for now!
for i in $(seq 1 $NODES);
do
NODE_IP="7.7.7.$i"
TX_FLAGS="--home /validator$i --keyring-backend test --from validator$i --trace --node=http://$NODE_IP:26657 --chain-id=peggy-test -y"

ETH_PRIVKEY=$(jq .[$i] ./tests/eth_keys.json -r)
peggycli tx nameservice update-eth-addr $ETH_PRIVKEY $TX_FLAGS > /dev/null
done

QUERY_FLAGS="--home /validator1 --trace --node=http://7.7.7.1:26657 --chain-id=peggy-test"

sleep 5 # Wait for a block to mine (there must be a better way)

RES=$(peggycli query nameservice current-valset $QUERY_FLAGS -o=json)
GOAL='{"Nonce":"0","Powers":["100","100","100"],"EthAdresses":["0xE987c5D2CFA68CD803e720FDD40ae10cE959c47B","0xa34F8827225c7FA6565C618b01de86549e07d667","0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"]}'

if [ $RES != $GOAL ]; then
    echo "valset test failed"
    echo $RES
    echo "is NOT equal to"
    echo $GOAL
    exit
else
    echo "valset test successful"
fi
