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

echo "Removing validators"
docker-compose rm --force --stop gravity{0..3}

echo "Building validators"
docker-compose build gravity{0..3}

echo "Deploying validators"
docker-compose up --no-start gravity{0..3}
docker-compose start gravity{0..3}