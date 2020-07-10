#!/bin/bash
# the script run inside the container for all-up-test-internal.sh
NODES=$1
bash /peggy/tests/setup-validators.sh $NODES

bash /peggy/tests/run-testnet.sh $NODES &

bash /peggy/tests/integration-tests.sh $NODES