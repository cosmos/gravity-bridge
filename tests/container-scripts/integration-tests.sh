#!/bin/bash
set -eu

set +e
killall -9 test-runner
set -e

pushd /peggy/orchestrator
RUST_BACKTRACE=full RUST_LOG=INFO PATH=$PATH:$HOME/.cargo/bin cargo run --release --bin test-runner
