module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(
    require("../build/contracts/CosmosBridge.json")
  );

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
  let oracleContractAddress;
  let bridgeBankContractAddress;

  // TODO: Input validation
  if (NETWORK_ROPSTEN) {
    oracleContractAddress = process.argv[6];
    bridgeBankContractAddress = process.argv[7];
  } else {
    oracleContractAddress = process.argv[4];
    bridgeBankContractAddress = process.argv[5];
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

  // Set Oracle
  const { logs: setOracleLogs } = await contract
    .deployed()
    .then(function(instance) {
      return instance.setOracle(oracleContractAddress, {
        from: accounts[0],
        value: 0,
        gas: 300000 // 300,000 Gwei
      });
    });

  // Get event logs
  const setOracleEvent = setOracleLogs.find(e => e.event === "LogOracleSet");
  console.log("CosmosBridge's Oracle set:", setOracleEvent.args._oracle);

  // Set BridgeBank
  const { logs: setBridgeBankLogs } = await contract
    .deployed()
    .then(function(instance) {
      return instance.setBridgeBank(bridgeBankContractAddress, {
        from: accounts[0],
        value: 0,
        gas: 300000 // 300,000 Gwei
      });
    });

  // Get event logs
  const setBridgeBankEvent = setBridgeBankLogs.find(
    e => e.event === "LogBridgeBankSet"
  );
  console.log(
    "CosmosBridge's BridgeBank set:",
    setBridgeBankEvent.args._bridgeBank
  );

  return;
};
