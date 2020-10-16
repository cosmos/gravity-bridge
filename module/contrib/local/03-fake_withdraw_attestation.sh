#!/bin/bash
set -eux

echo "--> Needs to run after deposit\n"
echo "## Add ETH key"
peggycli tx peggy update-eth-addr 0xb8662f35f9de8720424e82b232e8c98d15399490adae9ca993f5ef1dc4883690 --from validator  --chain-id=testing -b block -y

echo "## Add ETH withdraw to pool"
peggycli tx peggy withdraw validator 0xc783df8a850f42e7f7e57013759c285caa701eb6 1peggyf01b315c8e 0peggyf01b315c8e --from validator --chain-id=testing -b block -y

echo "## Request a batch for outgoing TX"
peggycli tx peggy build-batch peggyf01b315c8e --from validator --chain-id=testing -b block -y

echo "## Query pending request nonce"
nonce=$(peggycli q peggy pending-batch-request $(peggycli keys show validator -a) -o json | jq -r ".value.nonce")

echo "## Approve pending request"
peggycli tx peggy approved batch-confirm  "$nonce" 0xb8662f35f9de8720424e82b232e8c98d15399490adae9ca993f5ef1dc4883690 --from validator --chain-id=testing -b block -y

echo "## Submit observation"
# chain id: 1
# bridge contract address: 0x8858eeb3dfffa017d4bce9801d340d36cf895ccf
peggycli tx peggy observed withdrawal 1 0x8858eeb3dfffa017d4bce9801d340d36cf895ccf  "$nonce" --from validator --chain-id=testing -b block -y

echo "## Query balance"
peggycli q account $(peggycli keys show validator -a)
echo "## Query last observed state"
peggycli q peggy observed nonces -o json