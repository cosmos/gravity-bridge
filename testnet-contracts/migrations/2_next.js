const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");

module.exports = function(deployer, network, accounts) {
  const operator = accounts[0];
  const initialValidators = [accounts[1], accounts[2], accounts[3]];
  const initialPowers = [5, 8, 12];

  // Deploy Valset contract
  deployer
    .deploy(Valset, operator, initialValidators, initialPowers, {
      gas: 6721975, // Cost: 1529823
      from: operator
    })
    .then(function() {
      // Deploy CosmosBridge contract
      return deployer
        .deploy(CosmosBridge, Valset.address, {
          gas: 6721975, // Cost: 1201274
          from: operator
        })
        .then(function() {
          // Deploy Oracle contract
          return deployer
            .deploy(Oracle, operator, Valset.address, CosmosBridge.address, {
              gas: 6721975, // Cost: 1455275
              from: operator
            })
            .then(function() {
              // Deploy BridgeBank contract
              return deployer.deploy(
                BridgeBank,
                operator,
                Oracle.address,
                CosmosBridge.address,
                {
                  gas: 6721975, // Cost: 5257988
                  from: operator
                }
              );
            });
        });
    });
};
