validator_address=$(getent hosts ${VALIDATOR} | awk '{ print $1 }')
abci="http://$validator_address:26657"
grpc="http://$validator_address:9090"
ethrpc="http://$(getent hosts ethereum | awk '{ print $1 }'):8545"

COSMOS_NODE_GRPC="$grpc" COSMOS_NODE_ABCI="$abci" ETH_NODE="$ethrpc" RUST_BACKTRACE=full TEST_TYPE="VALSET_STRESS" RUST_LOG=INFO PATH=$PATH:$HOME/.cargo/bin cargo run --bin test_runner