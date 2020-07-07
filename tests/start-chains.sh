#!/bin/bash
set -eux

# Build test container
docker build --no-cache -f ./tests/tests.Dockerfile -t peggy-test .

# Remove existing container instance
set +e
docker rm -f peggy_test_instance
set -e

# Run new test container instance
docker run --name peggy_test_instance --mount type=bind,source="$(pwd)"/,target=/peggy --cap-add=NET_ADMIN -it peggy-test