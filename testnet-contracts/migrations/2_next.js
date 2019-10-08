var Peggy = artifacts.require("Peggy");
var TestToken = artifacts.require("TestToken");

module.exports = function(deployer, network, accounts) {
  // Declare initial validators and powers
  const initialValidators = [accounts[1], accounts[2], accounts[3]];
  const initialPowers = [8, 10, 15];

  // Deploy TestToken contract
  deployer.deploy(TestToken, {
    gas: 4612388,
    from: accounts[0]
  });

  // Deploy Peggy contract
  // Gas deployment cost: 5,183,260 (without BankToken contract deployment in Bank's deployBankToken() method)
  deployer.deploy(Peggy, initialValidators, initialPowers, {
    gas: 6721975,
    from: accounts[0]
  });
};
