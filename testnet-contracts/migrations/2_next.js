var Peggy = artifacts.require("Peggy");
var TestToken = artifacts.require("TestToken");

module.exports = function(deployer, network, accounts) {
  // Declare initial validators and powers
  const initalValidators = [(accounts[1], accounts[2], accounts[3])];
  const initalPowers = [(8, 10, 15)];

  // Deploy TestToken contract
  deployer.deploy(TestToken, {
    gas: 4612388,
    from: accounts[0]
  });

  // Deploy Peggy contract
  deployer.deploy(Peggy, initalValidators, initalPowers, {
    gas: 4612388,
    from: accounts[0]
  });
};
