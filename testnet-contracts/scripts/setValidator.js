module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(require("../build/contracts/Valset.json"));

  /*******************************************
   *** Constants
   ******************************************/
  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const NUM_ARGS = process.argv.length - 4;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if (NETWORK_ROPSTEN) {
    if (NUM_ARGS !== 4) {
      return console.error(
        "Error: Oracle contract address, BridgeBank contract address required"
      );
    }
  } else {
    if (NUM_ARGS !== 2) {
      return console.error(
        "Error: Oracle contract address, BridgeBank contract address required"
      );
    }
  }

  /*******************************************
   *** Lock transaction parameters
   ******************************************/
  let newValidator;

  // TODO: Input validation
  if (NETWORK_ROPSTEN) {
    newValidator = process.argv[6];
  } else {
    newValidator = process.argv[4];
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
  contract.setProvider(web3.currentProvider);

  /*******************************************
   *** Contract interaction
   ******************************************/
  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  // TODO: Implement this once new valset is merged in
  //   // Add new validator
  //   const { logs } = await contract.deployed().then(function(instance) {
  //     return instance.addValidator(oracleContractAddress, {
  //       from: accounts[0],
  //       value: 0,
  //       gas: 300000 // 300,000 Gwei
  //     });
  //   });

  //   // Get event logs
  //   const setValidatorEvent = logs.find(e => e.event === "LogValidatorAdded");
  //   console.log("Validator added:", setValidatorEvent.args._validator);

  return;
};
