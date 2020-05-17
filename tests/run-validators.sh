#!/bin/bash
set -eux
# your gaiad binary name
BIN=peggyd
CLI=peggycli

NODES=$1

for i in $(seq 1 $NODES);
do
# add this ip for loopback dialing
ip addr add 7.7.7.$i/32 dev eth0

GAIA_HOME="--home /validator$i"
# this implicitly caps us at ~6000 nodes for this sim
# note that we start on 26656 the idea here is that the first
# node (node 1) is at the expected contact address from the gentx
# faciliating automated peer exchange
# not sure what this one does but we need to set it or we'll
# see port conflicts
LISTEN_ADDRESS="--address tcp://7.7.7.$i:26655"
RPC_ADDRESS="--rpc.laddr tcp://7.7.7.$i:26657"
P2P_ADDRESS="--p2p.laddr tcp://7.7.7.$i:26656"
ARGS="$GAIA_HOME $LISTEN_ADDRESS $RPC_ADDRESS $P2P_ADDRESS"
if [ $i -le $NODES ]; then
$BIN $ARGS start &
fi
if [ $i -gt $(($NODES - 1)) ]; then
$BIN $ARGS start
fi
done