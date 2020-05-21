// this is a command line utility for deploying the peggy solidity contract
// to a 'real' blockchain. It needs an erc20 address (or in integraiton test mode
// it will deploy bitcoinMAXXXX), a eth full node address, and a cosmos full node
// address. And finally an eth private key containing enough funds to deploy the
// contracts.
// The utility will then go and request the full validator set from the cosmos chain
// format it, sign it, then deploy to ethereum using the provided full node and provided
// eth private key funds.

// example call contract-deployer --eth-node=eth.althea.net --cosmos-node=cosmos.althea.net --integration-test --eth-privkey=0xw34234
