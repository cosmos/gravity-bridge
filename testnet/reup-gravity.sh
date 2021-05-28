#!/bin/bash
set -eu
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

# Constants
PROJECT_DIR="$(dirname $SCRIPT_DIR)"
CHAINID="testchain"
CHAINDIR="$PROJECT_DIR/testdata"
gravity=gravity
home_dir="$CHAINDIR/$CHAINID"

n0name="gravity0"
n1name="gravity1"
n2name="gravity2"
n3name="gravity3"

# Folders for nodes
n0dir="$home_dir/$n0name"
n1dir="$home_dir/$n1name"
n2dir="$home_dir/$n2name"
n3dir="$home_dir/$n3name"

echo "Removing orchestrators"
docker-compose rm --force --stop orchestrator{0..3}

# echo "Building orchestrators"
docker-compose --env-file $n0dir/orchestrator.env build orchestrator0
docker-compose --env-file $n1dir/orchestrator.env build orchestrator1
docker-compose --env-file $n2dir/orchestrator.env build orchestrator2
docker-compose --env-file $n3dir/orchestrator.env build orchestrator3

# echo "Deploying orchestrators"
docker-compose --env-file $n0dir/orchestrator.env up --no-start orchestrator0
docker-compose --env-file $n1dir/orchestrator.env up --no-start orchestrator1
docker-compose --env-file $n2dir/orchestrator.env up --no-start orchestrator2
docker-compose --env-file $n3dir/orchestrator.env up --no-start orchestrator3

docker-compose --env-file $n0dir/orchestrator.env start orchestrator0
docker-compose --env-file $n1dir/orchestrator.env start orchestrator1
docker-compose --env-file $n2dir/orchestrator.env start orchestrator2
docker-compose --env-file $n3dir/orchestrator.env start orchestrator3