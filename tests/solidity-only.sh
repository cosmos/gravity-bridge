#!/bin/bash

# Remove existing container instance
set +e
docker rm -f peggy_solidity_test_instance
set -e

# Run new test container instance
docker run --name peggy_solidity_test_instance --mount type=bind,source="$(pwd)"/,target=/peggy --cap-add=NET_ADMIN -it peggy-base /bin/bash /peggy/tests/container-scripts/solidity-tests.sh