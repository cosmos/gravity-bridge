#!/bin/bash
set -eu

# Run test entry point script
docker exec peggy_test_instance /bin/sh -c "pushd /peggy/ && tests/integration-tests.sh"