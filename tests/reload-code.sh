#!/bin/bash
set -eux

# Number of validators to start
NODES=$1

# Stop any currently running peggy and eth processes
pkill peggyd || true # allowed to fail
pkill geth || true # allowed to fail

# Wipe filesystem changes
for i in $(seq 1 $NODES);
do
    rm -rf "/validator$i"
done


cd /peggy/module/
make
make install
cd /peggy/
tests/setup-validators.sh $NODES
tests/run-testnet.sh $NODES

# This keeps the script open to prevent Docker from stopping the container
# immediately if the nodes are killed by a different process
read -p "Press Return to Close..."