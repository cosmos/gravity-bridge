#!/bin/bash

result=$( docker images -q peggy-base )

if [[ -n "$result" ]]; then
  echo "Container exists"
else
  # builds the container containing various system deps
  # also builds Peggy once in order to cache Go deps, this container
  # is also used for the solidity tests
  bash $DIR/build-container.sh
fi

# Remove existing container instance
set +e
docker rm -f peggy_solidity_test_instance
set -e

# Run new test container instance
docker run --name peggy_solidity_test_instance --mount type=bind,source="$(pwd)"/,target=/peggy --cap-add=NET_ADMIN -it peggy-base /bin/bash /peggy/tests/container-scripts/solidity-tests.sh