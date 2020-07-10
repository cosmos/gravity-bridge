#!/bin/bash
# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# builds the container containing various system deps
# also builds Peggy once in order to cache Go deps, this container
# is also used for the solidity tests
bash $DIR/build-container.sh

# Remove existing container instance if it exits
set +e
docker rm -f peggy_test_instance
set -e

# Solidity tests
# this only tests the solidty code using Ganahe this is sufficient
# to see if the contracts compile and test basic functionality. The
# contract is later deployed in the run-tests stage of the module tests
# and is subjected to actual operation within that container
docker run --name peggy_test_instance -it peggy-base /bin/bash /peggy/tests/solidity-tests.sh

# Module tests

# Remove existing container instance
set +e
docker rm -f peggy_test_instance
set -e

NODES=3

# Run new test container instance
docker run --name peggy_test_instance --cap-add=NET_ADMIN -it peggy-base /bin/bash /peggy/tests/all-up-test-internal.sh $NODES
