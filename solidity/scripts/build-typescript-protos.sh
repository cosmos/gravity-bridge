#!/usr/bin/env bash
set -eux

# the directory of this script, useful for allowing this script
# to be run with any PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROTO_DIR="$DIR/../../module/proto/"
REPO_DIR="/tmp/ts-proto-repo"
COSMOS_VERSION="v0.40.0-rc3"
TMP_DIR="/tmp/ts-proto-compile"
OUR_DIR="$DIR/../proto"

set +e
rm -rf $TMP_DIR/
mkdir $TMP_DIR/
rm -rf $REPO_DIR/
set -e

git clone --depth 1 https://github.com/cosmos/cosmos-sdk/ $REPO_DIR
pushd $REPO_DIR
git fetch --tags origin $COSMOS_VERSION
git checkout $COSMOS_VERSION
popd

cp -r $PROTO_DIR/peggy $TMP_DIR/peggy
cp -r $REPO_DIR/third_party/proto/* $TMP_DIR/peggy/v1/
cp -r $REPO_DIR/proto/* $TMP_DIR/peggy/v1/


PROTOC_GEN_TS_PATH="$DIR/../node_modules/.bin/protoc-gen-ts"
GRPC_TOOLS_NODE_PROTOC_PLUGIN="$DIR/../node_modules/.bin/grpc_tools_node_protoc_plugin"
GRPC_TOOLS_NODE_PROTOC="$DIR/../node_modules/.bin/grpc_tools_node_protoc"

for f in $TMP_DIR/peggy/v1/; do

  # skip the non proto files
  if [ "$(basename "$f")" == "index.ts" ]; then
      continue
  fi

  # loop over all the available proto files and compile them into respective dir
  # JavaScript code generating
  ${GRPC_TOOLS_NODE_PROTOC} \
      --js_out=import_style=commonjs,binary:"${f}" \
      --grpc_out="${f}" \
      --plugin=protoc-gen-grpc="${GRPC_TOOLS_NODE_PROTOC_PLUGIN}" \
      -I "${f}" \
      "${f}"/*.proto

  ${GRPC_TOOLS_NODE_PROTOC} \
      --plugin=protoc-gen-ts="${PROTOC_GEN_TS_PATH}" \
      --ts_out="${f}" \
      -I "${f}" \
      "${f}"/*.proto

done

rm -rf $OUT_DIR
mkdir $OUT_DIR

cp $TMP_DIR/peggy/v1/*.ts $OUT_DIR/
cp $TMP_DIR/peggy/v1/*.js $OUT_DIR/

set +e
rm -rf $TMP_DIR/
rm -rf $REPO_DIR/
set -e