#!/bin/bash
set -eu

peggyd init --chain-id=testing local
peggyd add-genesis-account validator 1000000000stake
peggyd gentx --name validator  --amount 1000000000stake
peggyd collect-gentxs
