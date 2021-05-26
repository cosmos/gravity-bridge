#bin/sh

validator_address=$(getent hosts ${VALIDATOR} | awk '{ print $1 }')
rpc="http://$validator_address:1317"
grpc="http://$validator_address:9090"
ethrpc="http://$(getent hosts ethereum | awk '{ print $1 }'):8545"

echo orchestrator \
    --address-prefix=cosmos \
    --contract-address="${CONTRACT_ADDR}" \
    --cosmos-grpc="$grpc" \
    --cosmos-phrase="${COSMOS_PHRASE}" \
    --ethereum-key="${ETH_PRIVATE_KEY}" \
    --ethereum-rpc="$ethrpc" \
    --fees="${DENOM}"

orchestrator \
    --address-prefix=cosmos \
    --contract-address="${CONTRACT_ADDR}" \
    --cosmos-grpc="$grpc" \
    --cosmos-phrase="${COSMOS_PHRASE}" \
    --ethereum-key="${ETH_PRIVATE_KEY}" \
    --ethereum-rpc="$ethrpc" \
    --fees="${DENOM}"