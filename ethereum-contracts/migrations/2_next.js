var CosmosERC20.sol = artifacts.require("CosmosERC20")
var ERC20.sol = artifacts.require("ERC20")
var Peggy.sol = artifacts.require("Peggy")
var Valset.sol = artifacts.require("Valset")


module.exports = function(deployer) {
  // Use deployer to state migration tasks.
  deployer.deploy(ERC20);
  deployer.deploy(Valset);
  deployer.deploy(Peggy);
  deployer.deploy(CosmosERC20);
};
