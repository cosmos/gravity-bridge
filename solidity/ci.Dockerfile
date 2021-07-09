FROM alpine:edge

COPY contract-deployer /usr/bin/contract-deployer
COPY contracts contracts

CMD contract-deployer \
    --cosmos-node="http://gravity0:26657" \
    --eth-node="http://ethereum:8545" \
    --eth-privkey="0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7" \
    --contract=artifacts/contracts/Gravity.sol/Gravity.json \
    --test-mode=true