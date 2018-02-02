module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  authors: [
    "Adrian Brink <adrian@tendermint.com>",
    "Federico Kunze <federico@tendermint.com>"
  ],
  keywords: [
    "peggy",
    "peg zone",
    "Cosmos"
  ],
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
  dependencies: {
    "bignumber.js": "^6.0.0"
  }
};
