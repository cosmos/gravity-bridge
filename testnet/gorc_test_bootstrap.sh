#!/busybox/sh

set -ex

# import orchestrator Sommelier key
gorc --config=/root/gorc/config.toml keys cosmos recover orch-key "$ORCH_MNEMONIC"

# import orchestrator Ethereum key
gorc --config=/root/gorc/config.toml keys eth import orch-eth-key $ETH_PRIV_KEY

# start gorc orchestrator
gorc --config=/root/gorc/config.toml orchestrator start --cosmos-key=orch-key --ethereum-key=orch-eth-key
