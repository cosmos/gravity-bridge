require("dotenv").config();

var HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
  networks: {
    develop: {
      host: "localhost",
      port: 7545,
      network_id: "*",
      gas: 6721975,
      gasPrice: 200000000000,
      solc: {
        version: "0.5.0",
        optimizer: {
          enabled: true,
          runs: 200
        }
      }
    },
    ropsten: {
      provider: function() {
        return new HDWalletProvider(
          process.env.MNEMONIC,
          "https://ropsten.infura.io/".concat(process.env.INFURA_PROJECT_ID)
        );
      },
      network_id: 3,
      gas: 6000000
    }
  },
  rpc: {
    host: "localhost",
    post: 8080
  },
  mocha: {
    useColors: true
  }
};
