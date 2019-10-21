var Peggy = artifacts.require("Peggy");
var TestToken = artifacts.require("TestToken");

module.exports = function(deployer, network, accounts) {
  deployer.deploy(TestToken, { gas: 4612388, from: accounts[0] });
  deployer.deploy(Peggy, { gas: 4612388, from: accounts[0] });
};
