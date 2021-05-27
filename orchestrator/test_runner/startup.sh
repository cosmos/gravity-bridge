validator_address=$(getent hosts gravity0 | awk '{ print $1 }')
abci="http://$validator_address:26657"
grpc="http://$validator_address:9090"
ethrpc="http://$(getent hosts ethereum | awk '{ print $1 }'):8545"

COSMOS_NODE_GRPC="$grpc" COSMOS_NODE_ABCI="$abci" ETH_NODE="$ethrpc" PATH=$PATH:$HOME/.cargo/bin test_runner