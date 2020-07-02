#!/bin/bash
set -eux

# Build test container
docker build --no-cache -f ./tests/tests.Dockerfile -t peggy-test .

# Remove existing container instance
docker rm -f peggy_test_instance

# Run new test container instance
docker run --name peggy_test_instance --cap-add=NET_ADMIN -it peggy-test