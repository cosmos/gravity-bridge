#!/bin/bash
set -eux

# Remove existing container instance
set +e
docker rm -f peggy_test_instance
set -e

NODES=3

# Run new test container instance
docker run --name peggy_test_instance --mount type=bind,source="$(pwd)"/,target=/peggy --cap-add=NET_ADMIN -it peggy-base /bin/bash /peggy/tests/reload-code.sh $NODES