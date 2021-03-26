cargo run --release orchestrator \
    --cosmos-phrase="\"${COSMOS_PHRASE}\"" \
    --ethereum-key="\"${ETH_PRIVATE_KEY}\"" \
    --cosmos-legacy-rpc="\"${COSMOS_RPC}\"" \
    --cosmos-grpc="\"${COSMOS_GRPC}\"" \
    --ethereum-rpc="\"${ETH_RPC}\"" \
    --fees="${DENOM}" \
    --contract-address="\"${CONTRACT_ADDR}\""