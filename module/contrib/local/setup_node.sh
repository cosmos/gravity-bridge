#!/bin/bash
set -eu

peggy init --chain-id=testing local
peggy add-genesis-account validator 1000000000stake
peggy gentx --name validator --amount 1000000000stake
peggy collect-gentxs
