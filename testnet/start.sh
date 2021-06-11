#!/bin/bash
set -eu
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

# Constants
PROJECT_DIR="$(dirname $SCRIPT_DIR)"
CHAINID="testchain"
CHAINDIR="$PROJECT_DIR/testdata"
gravity=gravity
home_dir="$CHAINDIR/$CHAINID"

# Find or install rust binary
#binaryFound=$(which register-delegate-keys 2>/dev/null || echo FALSE)
#echo "$binaryFound"
#if [ $binaryFound == "FALSE" ]
#then
#  pushd orchestrator/register_delegate_keys
#  cargo install --path .
#  popd
#else
#  echo "found binary at $(which register-delegate-keys)"
#fi

# stop processes
export DOCKER_SCAN_SUGGEST=false
docker-compose down
rm -rf $CHAINDIR

n0name="gravity0"
n1name="gravity1"
n2name="gravity2"
n3name="gravity3"

# Folders for nodes
n0dir="$home_dir/$n0name"
n1dir="$home_dir/$n1name"
n2dir="$home_dir/$n2name"
n3dir="$home_dir/$n3name"

# Home flag for folder
home0="--home $n0dir"
home1="--home $n1dir"
home2="--home $n2dir"
home3="--home $n3dir"

# Config directories for nodes
n0cfgDir="$n0dir/config"
n1cfgDir="$n1dir/config"
n2cfgDir="$n2dir/config"
n3cfgDir="$n3dir/config"

# Config files for nodes
n0cfg="$n0cfgDir/config.toml"
n1cfg="$n1cfgDir/config.toml"
n2cfg="$n2cfgDir/config.toml"
n3cfg="$n3cfgDir/config.toml"

# App config files for nodes
n0appCfg="$n0cfgDir/app.toml"
n1appCfg="$n1cfgDir/app.toml"
n2appCfg="$n2cfgDir/app.toml"
n3appCfg="$n3cfgDir/app.toml"

# Common flags
kbt="--keyring-backend test"
cid="--chain-id $CHAINID"

echo "Creating 4x $gravity validators with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Build genesis file incl account for passed address
coins="100000000000stake,100000000000footoken"

# Initialize the 3 home directories and add some keys
$gravity $home0 $cid init n0 &>/dev/null
$gravity $home0 keys add val $kbt --output json | jq . >> $n0dir/validator_key.json
$gravity $home1 $cid init n1 &>/dev/null
$gravity $home1 keys add val $kbt --output json | jq . >> $n1dir/validator_key.json
$gravity $home2 $cid init n2 &>/dev/null
$gravity $home2 keys add val $kbt --output json | jq . >> $n2dir/validator_key.json
$gravity $home3 $cid init n3 &>/dev/null
$gravity $home3 keys add val $kbt --output json | jq . >> $n3dir/validator_key.json

find $home_dir -name validator_key.json | xargs cat | jq -r '.mnemonic' > $CHAINDIR/validator-phrases

echo "Adding validator addresses to genesis files"
$gravity $home0 add-genesis-account $($gravity $home0 keys show val -a $kbt) $coins &>/dev/null
#$gravity $home0 add-genesis-account $($FED $home0 keys show feeder) $coins &>/dev/null
$gravity $home0 add-genesis-account $($gravity $home1 keys show val -a $kbt) $coins &>/dev/null
#$gravity $home0 add-genesis-account $($FED $home1 keys show feeder) $coins &>/dev/null
$gravity $home0 add-genesis-account $($gravity $home2 keys show val -a $kbt) $coins &>/dev/null
#$gravity $home0 add-genesis-account $($FED $home2 keys show feeder) $coins &>/dev/null
$gravity $home0 add-genesis-account $($gravity $home3 keys show val -a $kbt) $coins &>/dev/null
#$gravity $home0 add-genesis-account $($FED $home3 keys show feeder) $coins &>/dev/null

echo "Generating orchestrator keys"
$gravity $home0 keys add --dry-run=true --output=json orch | jq . >> $n0dir/orchestrator_key.json
$gravity $home1 keys add --dry-run=true --output=json orch | jq . >> $n1dir/orchestrator_key.json
$gravity $home2 keys add --dry-run=true --output=json orch | jq . >> $n2dir/orchestrator_key.json
$gravity $home3 keys add --dry-run=true --output=json orch | jq . >> $n3dir/orchestrator_key.json

find $home_dir -name orchestrator_key.json | xargs cat | jq -r '.mnemonic' > $CHAINDIR/orchestrator-phrases

echo "Adding orchestrator keys to genesis"
n0orchKey="$(jq .address $n0dir/orchestrator_key.json)"
n1orchKey="$(jq .address $n1dir/orchestrator_key.json)"
n2orchKey="$(jq .address $n2dir/orchestrator_key.json)"
n3orchKey="$(jq .address $n3dir/orchestrator_key.json)"
jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $n0orchKey,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $n1orchKey,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $n2orchKey,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.auth.accounts += [{\"@type\": \"/cosmos.auth.v1beta1.BaseAccount\",\"address\": $n3orchKey,\"pub_key\": null,\"account_number\": \"0\",\"sequence\": \"0\"}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.bank.balances += [{\"address\": $n0orchKey,\"coins\": [{\"denom\": \"footoken\",\"amount\": \"100000000000\"},{\"denom\": \"stake\",\"amount\": \"100000000000\"}]}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.bank.balances += [{\"address\": $n1orchKey,\"coins\": [{\"denom\": \"footoken\",\"amount\": \"100000000000\"},{\"denom\": \"stake\",\"amount\": \"100000000000\"}]}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.bank.balances += [{\"address\": $n2orchKey,\"coins\": [{\"denom\": \"footoken\",\"amount\": \"100000000000\"},{\"denom\": \"stake\",\"amount\": \"100000000000\"}]}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json
jq ".app_state.bank.balances += [{\"address\": $n3orchKey,\"coins\": [{\"denom\": \"footoken\",\"amount\": \"100000000000\"},{\"denom\": \"stake\",\"amount\": \"100000000000\"}]}]" $n0cfgDir/genesis.json | sponge $n0cfgDir/genesis.json

echo "Copying genesis file around to sign"
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n3cfgDir/genesis.json

echo "Generating ethereum keys"
$gravity $home0 eth_keys add --output=json --dry-run=true | jq . >> $n0dir/eth_key.json
$gravity $home1 eth_keys add --output=json --dry-run=true | jq . >> $n1dir/eth_key.json
$gravity $home2 eth_keys add --output=json --dry-run=true | jq . >> $n2dir/eth_key.json
$gravity $home3 eth_keys add --output=json --dry-run=true | jq . >> $n3dir/eth_key.json

find testdata -name eth_key.json | xargs cat | jq -r '.private_key' > $CHAINDIR/validator-eth-keys

echo "Copying ethereum genesis file"
cp tests/assets/ETHGenesis.json $home_dir

echo "Adding initial ethereum value"
jq ".alloc |= . + {$(jq .address $n0dir/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" $home_dir/ETHGenesis.json | sponge $home_dir/ETHGenesis.json
jq ".alloc |= . + {$(jq .address $n1dir/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" $home_dir/ETHGenesis.json | sponge $home_dir/ETHGenesis.json
jq ".alloc |= . + {$(jq .address $n2dir/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" $home_dir/ETHGenesis.json | sponge $home_dir/ETHGenesis.json
jq ".alloc |= . + {$(jq .address $n3dir/eth_key.json) : {\"balance\": \"0x1337000000000000000000\"}}" $home_dir/ETHGenesis.json | sponge $home_dir/ETHGenesis.json

echo "Creating gentxs"
$gravity $home0 gentx --ip $n0name val 100000000000stake $(jq -r .address $n0dir/eth_key.json) $(jq -r .address $n0dir/orchestrator_key.json) $kbt $cid &>/dev/null
$gravity $home1 gentx --ip $n1name val 100000000000stake $(jq -r .address $n1dir/eth_key.json) $(jq -r .address $n1dir/orchestrator_key.json) $kbt $cid &>/dev/null
$gravity $home2 gentx --ip $n2name val 100000000000stake $(jq -r .address $n2dir/eth_key.json) $(jq -r .address $n2dir/orchestrator_key.json) $kbt $cid &>/dev/null
$gravity $home3 gentx --ip $n3name val 100000000000stake $(jq -r .address $n3dir/eth_key.json) $(jq -r .address $n3dir/orchestrator_key.json) $kbt $cid &>/dev/null

echo "Collecting gentxs in $n0name"
cp $n1cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n2cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n3cfgDir/gentx/*.json $n0cfgDir/gentx/
$gravity $home0 collect-gentxs &>/dev/null

echo "Distributing genesis file into $n1name, $n2name, $n3name"
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n3cfgDir/genesis.json

# Switch sed command in the case of linux
fsed() {
  if [ `uname` = 'Linux' ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Change ports on n0 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n0cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n0cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n0cfg
fsed 's#external_address = ""#external_address = "tcp://'$n0name:26656'"#g' $n0cfg
fsed 's#enable = false#enable = true#g' $n0appCfg
fsed 's#swagger = false#swagger = true#g' $n0appCfg

# Change ports on n1 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n1cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n1cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n1cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n1cfg
fsed 's#external_address = ""#external_address = "tcp://'$n1name':26656"#g' $n1cfg
fsed 's#enable = false#enable = true#g' $n1appCfg

# Change ports on n2 val
fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n2cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n2cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n2cfg
fsed 's#external_address = ""#external_address = "tcp://'$n2name':26656"#g' $n2cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n2cfg
fsed 's#enable = false#enable = true#g' $n2appCfg

fsed "s#\"tcp://127.0.0.1:26656\"#\"tcp://0.0.0.0:26656\"#g" $n3cfg
fsed "s#\"tcp://127.0.0.1:26657\"#\"tcp://0.0.0.0:26657\"#g" $n3cfg
fsed 's#addr_book_strict = true#addr_book_strict = false#g' $n3cfg
fsed 's#external_address = ""#external_address = "tcp://'$n3name':26656"#g' $n3cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n3cfg
fsed 's#enable = false#enable = true#g' $n3appCfg

echo "Setting peers"
peer0="$($gravity $home0 tendermint show-node-id)@$n0name:26656"
peer1="$($gravity $home1 tendermint show-node-id)@$n1name:26656"
peer2="$($gravity $home2 tendermint show-node-id)@$n2name:26656"
peer3="$($gravity $home3 tendermint show-node-id)@$n3name:26656"
# First node has peers already set when collecting gentxs
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer2','$peer3'"#g' $n1cfg
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer1','$peer3'"#g' $n2cfg
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer1','$peer2'"#g' $n3cfg

echo "Writing start commands"
echo "$gravity --home home start --pruning=nothing > home.n0.log" >> $n0dir/startup.sh
echo "$gravity --home home start --pruning=nothing > home.n1.log" >> $n1dir/startup.sh
echo "$gravity --home home start --pruning=nothing > home.n2.log" >> $n2dir/startup.sh
echo "$gravity --home home start --pruning=nothing > home.n3.log" >> $n3dir/startup.sh
chmod +x $home_dir/*/startup.sh

echo "Building ethereum and validator images"
docker-compose build ethereum $n0name $n1name $n2name $n3name

echo "Starting testnet"
docker-compose up --no-start ethereum $n0name $n1name $n2name $n3name &>/dev/null
docker-compose start ethereum $n0name $n1name $n2name $n3name &>/dev/null

echo "Waiting for cosmos cluster to sync"
sleep 10

echo "Applying contracts"
docker-compose build contract_deployer
contractAddress=$(docker-compose up contract_deployer | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}')
if [[ ! $contractAddress ]]; then
  echo "contract failed to deploy."
  exit 1
fi
echo "Contract address: $contractAddress"

docker-compose logs --no-color --no-log-prefix contract_deployer > $CHAINDIR/contracts

echo "Gathering keys for orchestrators"
echo VALIDATOR=$n0name >> $n0dir/orchestrator.env
echo COSMOS_GRPC="http://$n0name:9090/" >> $n0dir/orchestrator.env
echo COSMOS_RPC="http://$n0name:1317" >> $n0dir/orchestrator.env
echo COSMOS_KEY=$(jq .priv_key.value $n0cfgDir/priv_validator_key.json) >> $n0dir/orchestrator.env
echo COSMOS_PHRASE=$(jq .mnemonic $n0dir/orchestrator_key.json) >> $n0dir/orchestrator.env
echo DENOM=stake >> $n0dir/orchestrator.env
echo ETH_RPC=http://ethereum:8545 >> $n0dir/orchestrator.env
echo ETH_PRIVATE_KEY=$(jq .private_key $n0dir/eth_key.json) >> $n0dir/orchestrator.env
echo CONTRACT_ADDR=$contractAddress >> $n0dir/orchestrator.env

echo VALIDATOR=$n1name >> $n1dir/orchestrator.env
echo COSMOS_GRPC="http://$n1name:9090/" >> $n1dir/orchestrator.env
echo COSMOS_RPC="http://$n1name:1317" >> $n1dir/orchestrator.env
echo COSMOS_KEY=$(jq .priv_key.value $n1cfgDir/priv_validator_key.json) >> $n1dir/orchestrator.env
echo COSMOS_PHRASE=$(jq .mnemonic $n1dir/orchestrator_key.json) >> $n1dir/orchestrator.env
echo DENOM=stake >> $n1dir/orchestrator.env
echo ETH_RPC=http://ethereum:8545 >> $n1dir/orchestrator.env
echo ETH_PRIVATE_KEY=$(jq .private_key $n1dir/eth_key.json) >> $n1dir/orchestrator.env
echo CONTRACT_ADDR=$contractAddress >> $n1dir/orchestrator.env

echo VALIDATOR=$n2name >> $n2dir/orchestrator.env
echo COSMOS_GRPC="http://$n2name:9090/" >> $n2dir/orchestrator.env
echo COSMOS_RPC="http://$n2name:1317" >> $n2dir/orchestrator.env
echo COSMOS_KEY=$(jq .priv_key.value $n2cfgDir/priv_validator_key.json) >> $n2dir/orchestrator.env
echo COSMOS_PHRASE=$(jq .mnemonic $n2dir/orchestrator_key.json) >> $n2dir/orchestrator.env
echo DENOM=stake >> $n2dir/orchestrator.env
echo ETH_RPC=http://ethereum:8545 >> $n2dir/orchestrator.env
echo ETH_PRIVATE_KEY=$(jq .private_key $n2dir/eth_key.json) >> $n2dir/orchestrator.env
echo CONTRACT_ADDR=$contractAddress >> $n2dir/orchestrator.env

echo VALIDATOR=$n3name >> $n3dir/orchestrator.env
echo COSMOS_GRPC="http://$n3name:9090/" >> $n3dir/orchestrator.env
echo COSMOS_RPC="http://$n3name:1317" >> $n3dir/orchestrator.env
echo COSMOS_KEY=$(jq .priv_key.value $n3cfgDir/priv_validator_key.json) >> $n3dir/orchestrator.env
echo COSMOS_PHRASE=$(jq .mnemonic $n3dir/orchestrator_key.json) >> $n3dir/orchestrator.env
echo DENOM=stake >> $n3dir/orchestrator.env
echo ETH_RPC=http://ethereum:8545 >> $n3dir/orchestrator.env
echo ETH_PRIVATE_KEY=$(jq .private_key $n3dir/eth_key.json) >> $n3dir/orchestrator.env
echo CONTRACT_ADDR=$contractAddress >> $n3dir/orchestrator.env

echo "Building orchestrators"
docker-compose --env-file $n0dir/orchestrator.env build orchestrator0
docker-compose --env-file $n1dir/orchestrator.env build orchestrator1
docker-compose --env-file $n2dir/orchestrator.env build orchestrator2
docker-compose --env-file $n3dir/orchestrator.env build orchestrator3

echo "Deploying orchestrators"
docker-compose --env-file $n0dir/orchestrator.env up --no-start orchestrator0
docker-compose --env-file $n0dir/orchestrator.env start orchestrator0
docker-compose --env-file $n1dir/orchestrator.env up --no-start orchestrator1
docker-compose --env-file $n1dir/orchestrator.env start orchestrator1
docker-compose --env-file $n2dir/orchestrator.env up --no-start orchestrator2
docker-compose --env-file $n2dir/orchestrator.env start orchestrator2
docker-compose --env-file $n3dir/orchestrator.env up --no-start orchestrator3
docker-compose --env-file $n3dir/orchestrator.env start orchestrator3

echo "Run tests"
docker-compose build test_runner
docker-compose run test_runner

echo "Done."