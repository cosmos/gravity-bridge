#!/bin/bash

result=$( docker images -q gravity-base )

if [[ -n "$result" ]]; then
  echo "Container exists"
else
  # builds the container containing various system deps
  # also builds Gravity once in order to cache Go deps, this container
  # is also used for the solidity tests
  bash $DIR/build-container.sh
fi

# Remove existing container instance
set +e
docker rm -f gravity_solidity_test_instance
set -e

# Run new test container instance
docker run --name gravity_solidity_test_instance --mount type=bind,source="$(pwd)"/,target=/gravity --cap-add=NET_ADMIN -it gravity-base /bin/bash /gravity/tests/container-scripts/solidity-tests.sh