cargo run --release orchestrator \
    --cosmos-key=${COSMOS_KEY} \
    --ethereum-key=${ETH_PRIVATE_KEY} \
    --cosmos-legacy-rpc=${COSMOS_RPC} \
    --cosmos-grpc=${COSMOS_GRPC} \
    --ethereum-rpc=${ETH_RPC} \
    --fees=${DENOM} \
    --contract-address=${CONTRACT_ADDR}