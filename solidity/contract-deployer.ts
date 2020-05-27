// this is a command line utility for deploying the peggy solidity contract
// to a 'real' blockchain. It needs an erc20 address (or in integraiton test mode
// it will deploy bitcoinMAXXXX), a eth full node address, and a cosmos full node
// address. And finally an eth private key containing enough funds to deploy the
// contracts.
// The utility will then go and request the full validator set from the cosmos chain
// format it, sign it, then deploy to ethereum using the provided full node and provided
// eth private key funds.

// example call contract-deployer --eth-node=eth.althea.net --cosmos-node=cosmos.althea.net --integration-test --eth-privkey=0xw34234
import commandLineArgs from "command-line-args";

const options = commandLineArgs([
  { name: "eth-node", type: String },
  { name: "cosmos-node", type: String },
  { name: "eth-privkey", type: String },
  { name: "erc20", type: Boolean }
]) as {
  "eth-node": string;
  "cosmos-node": string;
  "eth-privkey": string;
  erc20: string;
};

// - Get the validator and powers list from the predetermined height, by hitting
//   the cosmos full node api.
// - We now need a signature from every validator on the list, over the list
// - Those signatures will have been committed to the cosmos chain in the consensus
//   state by the validators.
// - We hit the cosmos full node and access the peggy api to get the signatures

console.log(options);
