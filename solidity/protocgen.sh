#!/usr/bin/env bash

set -eox pipefail

proto_dirs=$(find ../module/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
mkdir -p ./gen/
for dir in $proto_dirs; do
  buf protoc \
  -I "../module/proto" \
  -I "../module/third_party/proto" \
  --ts_proto_out=./gen/ \
  --plugin=./node_modules/.bin/protoc-gen-ts_proto  \
  --ts_proto_opt="esModuleInterop=true,forceLong=long,useOptionals=true" \
  $(find "${dir}" -maxdepth 2 -name '*.proto')
done