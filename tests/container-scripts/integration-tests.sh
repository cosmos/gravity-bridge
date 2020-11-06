#!/bin/bash
set -eu
# WIP 
NODES=3 # Permanently set to 3 for now!
QUERY_FLAGS="--home /validator1 --trace --node=http://7.7.7.1:26657 --chain-id=peggy-test -o=json"

set +e
killall -9 test-runner
set -e

pushd /peggy/orchestrator
RUST_BACKTRACE=full RUST_LOG=INFO PATH=$PATH:$HOME/.cargo/bin cargo run --release --bin test-runner
