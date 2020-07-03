#!/bin/bash
set -eu
# WIP 

NODE_1_IP="7.7.7.1"
NODE_2_IP="7.7.7.2"
NODE_3_IP="7.7.7.3"

ETH_PRIVKEY_1="0x41b6fe18ea396208ab7dc526ca1cc59942b24df2d6436b6970fbe3a5d0c947a8"

peggycli config node http://$NODE_1_IP:26657
# Do node 1 crap
peggycli tx nameservice update-eth-addr $ETH_PRIVKEY_1