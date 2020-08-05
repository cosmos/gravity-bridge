#!/bin/bash
set -eu
# WIP 
NODES=3 # Permanently set to 3 for now!
QUERY_FLAGS="--home /validator1 --trace --node=http://7.7.7.1:26657 --chain-id=peggy-test -o=json"

#### valset creation test
for i in $(seq 1 $NODES);
do
    NODE_IP="7.7.7.$i"
    TX_FLAGS="--home /validator$i --keyring-backend test --from validator$i --trace --node=http://$NODE_IP:26657 --chain-id=peggy-test -y"

    ETH_PRIVKEY=$(jq .[$i] /peggy/tests/assets/eth_keys.json -r)
    peggycli tx peggy update-eth-addr $ETH_PRIVKEY $TX_FLAGS > /dev/null
done
BLOCK=$(peggycli status $QUERY_FLAGS | jq .sync_info.latest_block_height -r)

# wait for bootstrapping to finish
while [ $(peggycli status $QUERY_FLAGS | jq .sync_info.latest_block_height -r) -eq $BLOCK ]
do
sleep 0.2
done

RES=$(peggycli query peggy current-valset $QUERY_FLAGS)
GOAL='{"Nonce":"0","Powers":["100","100","100"],"EthAdresses":["0xE987c5D2CFA68CD803e720FDD40ae10cE959c47B","0xa34F8827225c7FA6565C618b01de86549e07d667","0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255"]}'

if [ $RES != $GOAL ]; then
    echo "valset test failed"
    echo $RES
    echo "is NOT equal to"
    echo $GOAL
else
    echo "valset test successful"
fi


#### valset-request test
# This is called by anyone to request that the validators save a valset for the next block
peggycli tx peggy valset-request $TX_FLAGS > /dev/null
BLOCK=$(peggycli status $QUERY_FLAGS | jq .sync_info.latest_block_height -r)

# Wait for a block to mine
while [ $(peggycli status $QUERY_FLAGS | jq .sync_info.latest_block_height -r) -eq $BLOCK ]
do
sleep 0.2
done

let "NONCE = $BLOCK + 1"
# This is called by the peggy daemons to see if a valset has been saved for a block
RES=$(peggycli query peggy valset-request $NONCE $QUERY_FLAGS)
GOAL="{\"Nonce\":\"$NONCE\",\"Powers\":[\"100\",\"100\",\"100\"],\"EthAdresses\":[\"0xE987c5D2CFA68CD803e720FDD40ae10cE959c47B\",\"0xa34F8827225c7FA6565C618b01de86549e07d667\",\"0xb462864E395d88d6bc7C5dd5F3F5eb4cc2599255\"]}"


if [ $RES != $GOAL ]; then
    echo "valset-request test failed"
    echo $RES
    echo "is NOT equal to"
    echo $GOAL
else
    echo "valset-request test successful"
fi
