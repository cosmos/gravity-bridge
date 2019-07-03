module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*",
      gas: 6000000,
      gasPrice: 200000000000,
      solc: {
        version: "0.5.0",
        optimizer: {
          enabled: true,
          runs: 200
        }
      }
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
      gas: 6000000,
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
