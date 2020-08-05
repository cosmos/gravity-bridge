#!/bin/bash
# the script run inside the container for all-up-test.sh
NODES=$1
bash /peggy/tests/container-scripts/setup-validators.sh $NODES

bash /peggy/tests/container-scripts/run-testnet.sh $NODES &

bash /peggy/tests/container-scripts/integration-tests.sh $NODES