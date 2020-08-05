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
tests/container-scripts/setup-validators.sh $NODES
tests/container-scripts/run-testnet.sh $NODES
tests/container-scripts/integration-tests.sh $NODES