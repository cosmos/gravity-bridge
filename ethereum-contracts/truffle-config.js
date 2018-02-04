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
  "Federico Kunze <federico@tendermint.com>"
  ],
  dependencies: {
    "bignumber.js": "^6.0.0"
  ,
  license: "MIT"
}
