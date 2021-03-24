#/bin/bash
set -eux

# Remove existing container instance
set +e
docker rm -f gravity_module_test_instance
set -e

NODES=3

# Run new test container instance
docker run --name gravity_module_test_instance --mount type=bind,source="$(pwd)"/,target=/gravity --cap-add=NET_ADMIN -it gravity-base /bin/bash /gravity/tests/container-scripts/module-only-internal.sh $NODES