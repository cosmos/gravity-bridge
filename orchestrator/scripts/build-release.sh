#!/bin/bash
set -eux
cross build --target=x86_64-unknown-linux-musl --release  --all
cross build --target=aarch64-unknown-linux-musl --release  --all

mkdir -p bins

cp target/x86_64-unknown-linux-musl/release/client bins/
cp target/x86_64-unknown-linux-musl/release/orchestrator bins/
cp target/x86_64-unknown-linux-musl/release/relayer bins/
cp target/x86_64-unknown-linux-musl/release/register-delegate-keys bins/

cp target/aarch64-unknown-linux-musl/release/client bins/client-arm
cp target/aarch64-unknown-linux-musl/release/orchestrator bins/orchestrator-arm
cp target/aarch64-unknown-linux-musl/release/relayer bins/relayer-arm
cp target/aarch64-unknown-linux-musl/release/register-delegate-keys bins/register-delegate-keys-arm
