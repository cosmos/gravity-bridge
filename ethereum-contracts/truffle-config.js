// Allows us to use ES6 in our migrations and tests.
require('babel-register');
require('babel-polyfill');

module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  networks: {
     development: {
     host: "localhost",
     port: 8545,
     network_id: "*" // Match any network id
    }
  },
  mocha: {
    useColors: true
  },
  authors: [
    "Adrian Brink <adrian@tendermint.com>",
    "Federico Kunze <federico@tendermint.com>"
  ],
  keywords: [
    "peggy",
    "peg zone",
    "Cosmos"
  ],
  dependencies: {
    "bignumber.js": "^6.0.0",
    "ethereumjs-util": "^5.1.3",
    "keccak": "^1.4.0",
    "keythereum": "^1.0.2"
  }
}
