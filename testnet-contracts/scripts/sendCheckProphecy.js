module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");

  const oracleContract = truffleContract(
    require("../build/contracts/Oracle.json")
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
  let prophecyID;

  if (NETWORK_ROPSTEN) {
    prophecyID = Number(process.argv[6]);
  } else {
    prophecyID = Number(process.argv[4]);
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

  console.log("Fetching Oracle contract...");
  oracleContract.setProvider(web3.currentProvider);

  /*******************************************
   *** Contract interaction
   ******************************************/
  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  console.log("Attempting to send checkBridgeProphecy() tx...");

  const instance = await oracleContract.deployed()
  let result
  try {
    result = await instance.checkBridgeProphecy(prophecyID, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });
  } catch (error) {
    console.log(error.message)
    return
  }

  const valid = result[0];
  const prophecyPowerCurrent = Number(result[1]);
  const prophecyPowerThreshold = Number(result[2]);
  if (result) {
    console.log(`\n\tProphecy ${prophecyID}`);
    console.log("----------------------------------------");
    console.log(`Current prophecy power:\t ${prophecyPowerCurrent}`);
    console.log(`Prophecy power threshold:\t ${prophecyPowerThreshold}`);
    console.log(`Reached threshold:\t ${valid}`);
    console.log("----------------------------------------");
  } else {
    console.error("Error: no result from transaction!");
  }

  return;
};
