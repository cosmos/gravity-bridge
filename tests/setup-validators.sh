#!/bin/bash
set -eux
# your gaiad binary name
BIN=peggyd
CLI=peggycli
# Sets up an arbitrary number of validators on a single machine by manipulating
# the --home parameter on gaiad
# todo add git hash to chain name
$BIN init --chain-id=peggy-test validator1
$CLI keys add --keyring-backend test validator1
KEY=$($CLI keys show validator1 -a --keyring-backend test)
$BIN add-genesis-account --keyring-backend test $KEY 1000000000stake,1000000000footoken
$BIN gentx --keyring-backend test --name validator1
$BIN collect-gentxs test