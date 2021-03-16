#!/bin/sh
# USAGE: ./two-node-net skip

# Constants
CURRENT_WORKING_DIR=$(pwd)
CHAINID="testchain"
CHAINDIR="$CURRENT_WORKING_DIR/testdata"
PEGGY=peggy
#FED=oracle-feeder
home_dir="$CHAINDIR/$CHAINID"

# stop processes
docker-compose down
rm -r $CHAINDIR

n0name="peggy_0"
n1name="peggy_1"
n2name="peggy_2"
n3name="peggy_3"

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

# Config files for feeders
#fd0cfg="$n0dir/config.yaml"
#fd1cfg="$n1dir/config.yaml"
#fd2cfg="$n2dir/config.yaml"

# Common flags
kbt="--keyring-backend test"
cid="--chain-id $CHAINID"

# Ensure user understands what will be deleted
if [[ -d $SIGNER_DATA ]] && [[ ! "$1" == "skip" ]]; then
  read -p "$0 will delete \$(pwd)/data folder. Do you wish to continue? (y/n): " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
fi

echo "Creating 4x $PEGGY validators with chain-id=$CHAINID..."
echo "Initializing genesis files"

# Build genesis file incl account for passed address
coins="100000000000stake,100000000000samoleans"

# Initialize the 3 home directories and add some keys
$PEGGY $home0 $cid init n0 &>/dev/null
echo "$n0name"_COSMOS_PHRASE=$($PEGGY $home0 keys add val $kbt --output json | jq .mnemonic) >> $n0dir/orchestrator.env
$PEGGY $home1 $cid init n1 &>/dev/null
echo "$n1name"_COSMOS_PHRASE=$($PEGGY $home1 keys add val $kbt --output json | jq .mnemonic) >> $n1dir/orchestrator.env
$PEGGY $home2 $cid init n2 &>/dev/null
echo "$n2name"_COSMOS_PHRASE=$($PEGGY $home2 keys add val $kbt --output json | jq .mnemonic) >> $n2dir/orchestrator.env
$PEGGY $home3 $cid init n3 &>/dev/null
echo "$n3name"_COSMOS_PHRASE=$($PEGGY $home3 keys add val $kbt --output json | jq .mnemonic) >> $n3dir/orchestrator.env

# Add some keys and init feeder configs
#$FED $home0 config init &>/dev/null
#$FED $home0 keys add feeder &>/dev/null
#$FED $home1 config init &>/dev/null
#$FED $home1 keys add feeder &>/dev/null
#$FED $home2 config init &>/dev/null
#$FED $home2 keys add feeder &>/dev/null

echo "Adding addresses to genesis files"
$PEGGY $home0 add-genesis-account $($PEGGY $home0 keys show val -a $kbt) $coins &>/dev/null
#$PEGGY $home0 add-genesis-account $($FED $home0 keys show feeder) $coins &>/dev/null
$PEGGY $home0 add-genesis-account $($PEGGY $home1 keys show val -a $kbt) $coins &>/dev/null
#$PEGGY $home0 add-genesis-account $($FED $home1 keys show feeder) $coins &>/dev/null
$PEGGY $home0 add-genesis-account $($PEGGY $home2 keys show val -a $kbt) $coins &>/dev/null
#$PEGGY $home0 add-genesis-account $($FED $home2 keys show feeder) $coins &>/dev/null
$PEGGY $home0 add-genesis-account $($PEGGY $home3 keys show val -a $kbt) $coins &>/dev/null
#$PEGGY $home0 add-genesis-account $($FED $home3 keys show feeder) $coins &>/dev/null

echo "Copying genesis file around to sign"
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n3cfgDir/genesis.json

echo "Creating gentxs and collect them in $n0name"
$PEGGY $home0 gentx --ip $n0name val 100000000000stake $kbt $cid &>/dev/null
$PEGGY $home1 gentx --ip $n1name val 100000000000stake $kbt $cid &>/dev/null
$PEGGY $home2 gentx --ip $n2name val 100000000000stake $kbt $cid &>/dev/null
$PEGGY $home3 gentx --ip $n3name val 100000000000stake $kbt $cid &>/dev/null
cp $n1cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n2cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n3cfgDir/gentx/*.json $n0cfgDir/gentx/
$PEGGY $home0 collect-gentxs &>/dev/null
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
#$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $n0cfg
#$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n0cfg
#$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n0cfg
fsed -i '' 's#external_address = ""#external_address = "tcp://'$n0name:26657'"#g' $n0cfg

# Change ports on n1 val
#$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26667"#g' $n1cfg
#$SED 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26666"#g' $n1cfg
#$SED 's#"localhost:6060"#"localhost:6061"#g' $n1cfg
#$SED 's#"0.0.0.0:9090"#"0.0.0.0:9091"#g' $n1app
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n1cfg
#$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n1cfg
fsed 's#external_address = ""#external_address = "tcp://'$n1name':26657"#g' $n1cfg
#$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n1cfg

# Change ports on n1 feeder
#$SED 's#http://localhost:9090#http://localhost:9091#g' $fd1cfg
#$SED 's#http://http://localhost:26657#http://http://localhost:26667#g' $fd1cfg

# Change ports on n2 val
#$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n2cfg
fsed 's#external_address = ""#external_address = "tcp://'$n2name':26657"#g' $n2cfg
#$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26677"#g' $n2cfg
#$SED 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26676"#g' $n2cfg
#$SED 's#"localhost:6060"#"localhost:6062"#g' $n2cfg
#$SED 's#"0.0.0.0:9090"#"0.0.0.0:9092"#g' $n2app
#$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n2cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n2cfg

fsed 's#external_address = ""#external_address = "tcp://'$n3name':26657"#g' $n3cfg
fsed 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n3cfg

# Change ports on n2 feeder
#$SED 's#http://localhost:9090#http://localhost:9092#g' $fd1cfg
#$SED 's#http://http://localhost:26657#http://http://localhost:26677#g' $fd1cfg

echo "Setting peers"
peer0="$($PEGGY $home0 tendermint show-node-id)@$n0name:26656"
peer1="$($PEGGY $home1 tendermint show-node-id)@$n1name:26656"
peer2="$($PEGGY $home2 tendermint show-node-id)@$n2name:26656"
peer3="$($PEGGY $home3 tendermint show-node-id)@$n3name:26656"
# First node has peers already set when collecting gentxs
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer2','$peer3'"#g' $n1cfg
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer1','$peer3'"#g' $n2cfg
fsed 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer1','$peer2'"#g' $n3cfg

echo "Writing start commands"
echo "$PEGGY --home home start --pruning=nothing --grpc.address="$n0name:9090" > home.n0.log" >> $n0dir/startup.sh
echo "$PEGGY --home home start --pruning=nothing --grpc.address="$n1name:9090" > home.n1.log" >> $n1dir/startup.sh
echo "$PEGGY --home home start --pruning=nothing --grpc.address="$n2name:9090" > home.n2.log" >> $n2dir/startup.sh
echo "$PEGGY --home home start --pruning=nothing --grpc.address="$n3name:9090" > home.n3.log" >> $n3dir/startup.sh
chmod +x $home_dir/*/startup.sh

echo "Gathering keys for orchestrator"
#--cosmos-key=<ckey>          The Cosmos private key of the validator
#            --ethereum-key=<ekey>        The Ethereum private key of the validator
#            --cosmos-legacy-rpc=<curl>   The Cosmos RPC url, usually the validator
#            --cosmos-grpc=<gurl>         The Cosmos gRPC url, usually the validator
#            --ethereum-rpc=<eurl>        The Ethereum RPC url, should be a self hosted node
#            --fees=<denom>               The Cosmos Denom in which to pay Cosmos chain fees
#            --contract-address=<addr>    The Ethereum contract address for Peggy, this is temporary

echo "$n0name"_COSMOS_GRPC="http://$n0name:9090" >> $n0dir/orchestrator.env
echo "$n0name"_COSMOS_RPC="http://$n0name:26657" >> $n0dir/orchestrator.env
echo "$n0name"_COSMOS_KEY=$(jq .priv_key.value $n0cfgDir/priv_validator_key.json) >> $n0dir/orchestrator.env
echo "$n0name"_DENOM=stake >> $n0dir/orchestrator.env
echo "$n0name"_ETH_RPC=http://ethereum:8545 >> $n0dir/orchestrator.env

echo "$n1name"_COSMOS_GRPC="http://$n1name:9090" >> $n1dir/orchestrator.env
echo "$n1name"_COSMOS_RPC="http://$n1name:26657" >> $n1dir/orchestrator.env
echo "$n1name"_COSMOS_KEY=$(jq .priv_key.value $n1cfgDir/priv_validator_key.json) >> $n1dir/orchestrator.env
echo "$n1name"_DENOM=stake >> $n1dir/orchestrator.env
echo "$n1name"_ETH_RPC=http://ethereum:8545 >> $n1dir/orchestrator.env

echo "$n2name"_COSMOS_GRPC="http://$n2name:9090" >> $n2dir/orchestrator.env
echo "$n2name"_COSMOS_RPC="http://$n2name:26657" >> $n2dir/orchestrator.env
echo "$n2name"_COSMOS_KEY=$(jq .priv_key.value $n2cfgDir/priv_validator_key.json) >> $n2dir/orchestrator.env
echo "$n2name"_DENOM=stake >> $n2dir/orchestrator.env
echo "$n2name"_ETH_RPC=http://ethereum:8545 >> $n2dir/orchestrator.env

echo "$n3name"_COSMOS_GRPC="http://$n3name:9090" >> $n3dir/orchestrator.env
echo "$n3name"_COSMOS_RPC="http://$n3name:26657" >> $n3dir/orchestrator.env
echo "$n3name"_COSMOS_KEY=$(jq .priv_key.value $n3cfgDir/priv_validator_key.json) >> $n3dir/orchestrator.env
echo "$n3name"_DENOM=stake >> $n3dir/orchestrator.env
echo "$n3name"_ETH_RPC=http://ethereum:8545 >> $n3dir/orchestrator.env

exit 0
echo "Building images"
docker-compose build

echo "Starting testnet"
docker-compose up --no-start
docker-compose start