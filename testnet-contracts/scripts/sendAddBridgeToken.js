module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");

  const contract = truffleContract(
    require("../build/contracts/BridgeBank.json")
  );

  /*******************************************
   *** Constants
   ******************************************/
  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";

  /*******************************************
   *** checkBridgeProphecy transaction parameters
   ******************************************/
  let symbol;

  if (NETWORK_ROPSTEN) {
    symbol = process.argv[6].toString();
  } else {
    symbol = process.argv[4].toString();
  }

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

  console.log("Fetching BridgeBank contract...");
  contract.setProvider(web3.currentProvider);

  /*******************************************
   *** Contract interaction
   ******************************************/
  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  console.log("Attempting to send createNewBridgeToken() tx...");

  // Get the bridge token's address if it were to be created
  const bridgeTokenAddress = await contract.deployed().then(function(instance) {
    return instance.createNewBridgeToken.call(symbol, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });
  });

  //  Create the bridge token
  await contract.deployed().then(function(instance) {
    return instance.createNewBridgeToken(symbol, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });
  });

  // Check bridge token whitelist
  const isOnWhiteList = await contract.deployed().then(function(instance) {
    return instance.bridgeTokenWhitelist(bridgeTokenAddress, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });
  });

  if (isOnWhiteList) {
    console.log(
      'Bridge Token"' + symbol + '" created at address:',
      bridgetokenAddress
    );
  } else {
    console.log(
      "Error: Bridge Token creation and whitelisting was not successful"
    );
  }

  return;
};
