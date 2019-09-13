var Peggy = artifacts.require("Peggy");

module.exports = function(deployer, network, accounts) {
  deployer.deploy(Peggy, { gas: 4612388, from: accounts[0] });
};
