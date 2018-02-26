module.exports = {
  // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  keywords: [
    "peggy",
    "peg zone",
    "Cosmos"
  ],
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*", // Match any network id
      gas: 4700000,
      solc: { optimizer: { enabled: true, runs: 200 } }
    },
    ganache: {
      host: "127.0.0.1",
      port: 7545,
      network_id: "*"
    },
    ropsten:  {
      network_id: 3,
      host: "localhost",
      port: 8545,
      gas: 4700000,
      solc: { optimizer: { enabled: true, runs: 200 } }
   }
  },
  rpc: {
    host: 'localhost',
    post:8080
  },
  mocha: {
    useColors: true
  },
};
