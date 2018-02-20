var Valset = artifacts.require("./Valset.sol");
var Peggy = artifacts.require("./Peggy.sol");
var CosmosERC20 = artifacts.require("./CosmosERC20.sol");

module.exports = function(deployer) {
  deployer.deploy(Valset);
  deployer.deploy(Peggy);
  deployer.deploy(CosmosERC20);
};
