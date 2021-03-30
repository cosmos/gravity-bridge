#!/bin/sh

CHAINID=$1
GENACCT=$2
ETHKEY=$3

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input gravity orchestrator account address..."
  exit 1
fi

if [ -z "$3" ]; then 
  echo "Need to input gravity ethereum key to delegate to..."
fi

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
gravity init --chain-id $CHAINID $CHAINID
gravity keys add validator --keyring-backend="test"
gravity add-genesis-account $(gravity keys show validator -a --keyring-backend="test") $coins
gravity add-genesis-account $GENACCT $coins
gravity gentx validator 1000000000stake $ETHKEY $GENACCT --keyring-backend="test" --chain-id $CHAINID
gravity collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.gravity/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.gravity/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.gravity/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.gravity/config/config.toml

# Start the gravity
gravity start --pruning=nothing
