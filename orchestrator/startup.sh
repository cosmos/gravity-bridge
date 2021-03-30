#bin/sh

validator_address=$(getent hosts ${VALIDATOR} | awk '{ print $1 }')
rpc="http://$validator_address:1317"
grpc="http://$validator_address:9090"
ethrpc="http://$(getent hosts ethereum | awk '{ print $1 }'):8545"

echo cargo run --bin orchestrator -- \
    --cosmos-phrase="${COSMOS_PHRASE}" \
    --ethereum-key="${ETH_PRIVATE_KEY}" \
    --cosmos-legacy-rpc="$rpc" \
    --cosmos-grpc="$grpc" \
    --ethereum-rpc="$ethrpc" \
    --fees="${DENOM}" \
    --contract-address="${CONTRACT_ADDR}"

cargo run --bin orchestrator -- \
    --cosmos-phrase="${COSMOS_PHRASE}" \
    --ethereum-key="${ETH_PRIVATE_KEY}" \
    --cosmos-legacy-rpc="$rpc" \
    --cosmos-grpc="$grpc" \
    --ethereum-rpc="$ethrpc" \
    --fees="${DENOM}" \
    --contract-address="${CONTRACT_ADDR}"