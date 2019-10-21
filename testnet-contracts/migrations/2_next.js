require("dotenv").config();

const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");

module.exports = function(deployer, network, accounts) {
  // Required initial paramater variables
  let operator;
  let initialValidators = [];
  let initialPowers = [];

  // Development network deployment param parsing/setting
  if (network === "develop") {
    // Operator
    operator = accounts[0];
    // Initial validators
    const localValidatorCount = Number(process.env.LOCAL_VALIDATOR_COUNT);
    if (localValidatorCount <= 0 || localValidatorCount > 9) {
      return console.error(
        "Must provide an initial validator count between 1-8 for local deployment."
      );
    } else {
      initialValidators = accounts.slice(1, localValidatorCount + 1);
    }
    // Initial validator power
    if (process.env.LOCAL_INITIAL_VALIDATOR_POWERS.length === 0) {
      return console.error(
        "Must provide initial local validator powers as environment variable."
      );
    } else {
      initialPowers = process.env.LOCAL_INITIAL_VALIDATOR_POWERS.split(",");
    }
    // Testnet/mainnet network deployment param parsing/setting
  } else {
    // Operator
    if (process.env.OPERATOR.length === 0) {
      return console.error(
        "Must provide operator address as environment variable."
      );
    } else {
      operator = process.env.OPERATOR;
    }
    // Initial validators
    if (process.env.INITIAL_VALIDATOR_ADDRESSES.length === 0) {
      return console.error(
        "Must provide initial validator addresses as environment variable."
      );
    } else {
      initialValidators = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
    }
    // Initial validator powers
    if (process.env.INITIAL_VALIDATOR_POWERS.length === 0) {
      return console.error(
        "Must provide initial validator powers as environment variable."
      );
    } else {
      initialPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");
    }
  }

  // Check that each initial validator has a power
  if (initialValidators.length !== initialPowers.length) {
    return console.error(
      "Each initial validator must have a corresponding power specified."
    );
  }

  // 1. Deploy Valset contract
  deployer
    .deploy(Valset, operator, initialValidators, initialPowers, {
      gas: 6721975, // Cost: 1,529,823
      from: operator
    })
    .then(function() {
      // 2. Deploy CosmosBridge contract
      return deployer
        .deploy(CosmosBridge, operator, Valset.address, {
          gas: 6721975, // Cost: 1,201,274
          from: operator
        })
        .then(function() {
          // 3. Deploy Oracle contract
          return deployer
            .deploy(Oracle, operator, Valset.address, CosmosBridge.address, {
              gas: 6721975, // Cost: 1,455,275
              from: operator
            })
            .then(function() {
              // 4. Deploy BridgeBank contract
              return deployer.deploy(
                BridgeBank,
                operator,
                Oracle.address,
                CosmosBridge.address,
                {
                  gas: 6721975, // Cost: 5,257,988
                  from: operator
                }
              );
            });
        });
    });
};
