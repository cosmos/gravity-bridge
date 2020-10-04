#!/bin/bash
set -eu

echo "## Add ETH key"
peggycli tx peggy update-eth-addr 0xb8662f35f9de8720424e82b232e8c98d15399490adae9ca993f5ef1dc4883690 --from validator  --chain-id=testing -b block -y
echo "## Request valset update"
peggycli tx peggy valset-request --from validator --chain-id=testing -b block -y
echo "## Query pending request nonce"
nonce=$(peggycli q peggy pending-valset-request $(peggycli keys show validator -a) -o json | jq -r ".value.nonce")

echo "## Approve pending request"
peggycli tx peggy approved valset-confirm  "$nonce" 0xb8662f35f9de8720424e82b232e8c98d15399490adae9ca993f5ef1dc4883690 --from validator --chain-id=testing -b block -y

echo "## View attestations"
peggycli q peggy attestation orchestrator_signed_multisig_update $nonce -o json | jq

echo "## Submit observation"
# chain id: 1
# bridge contract address: 0x8858eeb3dfffa017d4bce9801d340d36cf895ccf
#peggycli tx peggy observed  multisig-update 1 0x8858eeb3dfffa017d4bce9801d340d36cf895ccf  "$nonce" --from validator --chain-id=testing -b block -y
echo "## Query last observed state"
peggycli q peggy observed nonces -o json