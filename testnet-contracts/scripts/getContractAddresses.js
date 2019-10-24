module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  require("dotenv").config();
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");
  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(
    require("../build/contracts/BridgeRegistry.json")
  );

  /*******************************************
   *** Constants
   ******************************************/
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";

  /*******************************************
   *** Web3 provider
   *** Set contract provider based on --network flag
   ******************************************/
  let provider;
  if (NETWORK_ROPSTEN) {
    provider = new HDWalletProvider(
      process.env.MNEMONIC,
      "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
    );
  } else {
    provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  }

  const web3 = new Web3(provider);
  contract.setProvider(web3.currentProvider);

  /*******************************************
   *** Contract interaction
   ******************************************/
  const cosmosBridgeAddress = await contract
    .deployed()
    .then(function(instance) {
      return instance.cosmosBridge({
        from: accounts[0],
        value: 0,
        gas: 300000 // 300,000 Gwei;
      });
    });

  const bridgeBankAddress = await contract.deployed().then(function(instance) {
    return instance.bridgeBank({
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei;
    });
  });

  const oracleAddress = await contract.deployed().then(function(instance) {
    return instance.oracle({
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei;
    });
  });

  const valsetAddress = await contract.deployed().then(function(instance) {
    return instance.valset({
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei;
    });
  });

  console.log("CosmosBridge:", cosmosBridgeAddress);
  console.log("BridgeBank:", bridgeBankAddress);
  console.log("Oracle:", oracleAddress);
  console.log("Valset:", valsetAddress);
};
