var Valset = artifacts.require("./Valset.sol");
var Peggy = artifacts.require("./Peggy.sol");
// var ERC20 = artifacts.require("./ERC20.sol");
var CosmosERC20 = artifacts.require("./CosmosERC20.sol");

/* Abstract Contracts from https://solidity.readthedocs.io/en/develop/contracts.html#abstract-contracts

(C)ontracts cannot be compiled (even if they contain implemented functions
alongside non-implemented functions), but they can be used as base contracts

*/

module.exports = function(deployer) {
  deployer.deploy(Valset);
  // deployer.deploy(ERC20); // Abstract contract
  deployer.deploy(Peggy);
  deployer.deploy(CosmosERC20);
};
