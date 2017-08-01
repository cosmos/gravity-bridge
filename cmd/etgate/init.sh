CLIENT=~/.etgate/client
SERVER=~/.etgate/server
etcli="basecli --home $CLIENT"
etgate="./etgate --home $SERVER"

CHAINID="etgate-chain"

PORT_PREFIX=1234
RPC_PORT=${PORT_PREFIX}7

echo "Make new account"

$etcli keys new money
MONEY=$($etcli keys get money | awk '{print $2}')

sleep 2

echo "etgate init"

$etgate init --chain-id $CHAINID $MONEY

sed -ie "s/4665/$PORT_PREFIX/" $SERVER/config.toml

sleep 2

echo "etgate start"

$etgate start &> etgate.log &

sleep 2

echo "etcli init"

$etcli init --node=tcp://localhost:${RPC_PORT} --genesis=${SERVER}/genesis.json

RELAY_KEY=$SERVER/key.json
RELAY_ADDR=$(cat $RELAY_KEY | jq .address | tr -d \")

sleep 2

echo "etcli tx send"

$etcli tx send --amount=100000mycoin --sequence=1 --to=$RELAY_ADDR --name=money

cp ../../static/abimap.json $SERVER

sleep 2

echo "etgate gate init"

$etgate gate init --chain-id=$CHAINID --nodeaddr=tcp://localhost:${RPC_PORT} ../../static/example.json
