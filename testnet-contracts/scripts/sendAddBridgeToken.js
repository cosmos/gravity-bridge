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
  let symbol;
  let NETWORK_ROPSTEN
  try {
    NETWORK_ROPSTEN =
      process.argv[4] === "--network" && process.argv[5] === "ropsten";

    /*******************************************
     *** checkBridgeProphecy transaction parameters
    ******************************************/
    if (NETWORK_ROPSTEN) {
      symbol = process.argv[6].toString();
    } else {
      symbol = process.argv[4].toString();
    }
  }catch (error) {
    console.log({error})
    return
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
  try {
    /*******************************************
     *** Contract interaction
    ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();

    console.log("Attempting to send createNewBridgeToken() tx with symbol: '" + symbol + "'...");

    // Get the bridge token's address if it were to be created
    const bridgeTokenAddress = await contract.deployed().then(function(instance) {
      return instance.createNewBridgeToken.call(symbol, {
        from: accounts[0],
        value: 0,
        gas: 3000000 // 300,000 Gwei
      });
    });
    console.log(`from ${accounts[0]}`)
    console.log('Should deploy to ' + bridgeTokenAddress)

    //  Create the bridge token
    await contract.deployed().then(function(instance) {
      return instance.createNewBridgeToken(symbol, {
        from: accounts[0],
        value: 0,
        gas: 3000000 // 300,000 Gwei
      });
    });

    console.log("")
    // Check bridge token whitelist
    const isOnWhiteList = await contract.deployed().then(function(instance) {
      return instance.bridgeTokenWhitelist(bridgeTokenAddress, {
        from: accounts[0],
        value: 0,
        gas: 3000000 // 300,000 Gwei
      });
    });

    if (isOnWhiteList) {
      console.log(
        'Bridge Token "' + symbol + '" created at address ' + bridgeTokenAddress + ' and added to whitelist'
      );
    } else {
      console.log(
        "Error: Bridge Token creation and whitelisting was not successful"
      );
    }
  } catch (error) {
    console.error({error})
  }

  return;
};
