module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(require("../build/contracts/Peggy.json"));

  /*******************************************
   *** Constants
   ******************************************/
  // Lock transaction default params
  const DEFAULT_COSMOS_RECIPIENT =
    "0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668";
  const DEFAULT_ETH_DENOM = "0x0000000000000000000000000000000000000000";
  const DEFAULT_AMOUNT = 10;

  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const DEFAULT_PARAMS =
    process.argv[4] === "--default" ||
    (NETWORK_ROPSTEN && process.argv[6] === "--default");
  const NUM_ARGS = process.argv.length - 4;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if (NETWORK_ROPSTEN && DEFAULT_PARAMS) {
    if (NUM_ARGS !== 3) {
      return console.error(
        "Error: custom parameters are invalid on --default."
      );
    }
  } else if (NETWORK_ROPSTEN) {
    if (NUM_ARGS !== 2 || 5) {
      return console.error("Error: invalid parameters, please try again.");
    }
  } else if (DEFAULT_PARAMS) {
    if (NUM_ARGS !== 1) {
      return console.error(
        "Error: custom parameters are invalid on --default."
      );
    }
  } else {
    if (NUM_ARGS !== 3) {
      return console.error(
        "Error: must specify recipient address, token address, and amount."
      );
    }
  }

  /*******************************************
   *** Lock transaction parameters
   ******************************************/
  let cosmosRecipient = DEFAULT_COSMOS_RECIPIENT;
  let coinDenom = DEFAULT_ETH_DENOM;
  let amount = DEFAULT_AMOUNT;

  // TODO: Input validation
  if (!DEFAULT_PARAMS) {
    if (NETWORK_ROPSTEN) {
      cosmosRecipient = process.argv[6];
      coinDenom = process.argv[7];
      amount = parseInt(process.argv[8], 10);
    } else {
      cosmosRecipient = process.argv[4];
      coinDenom = process.argv[5];
      amount = parseInt(process.argv[6], 10);
    }
  }

  /*******************************************
   *** Web3 provider
   *** Set contract provider based on --network flag
   ******************************************/
  let provider;
  if (NETWORK_ROPSTEN) {
    provider = new HDWalletProvider(
      process.env.MNEMONIC,
      "https://ropsten.infura.io/".concat(process.env.INFURA_PROJECT_ID)
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

  // Send lock transaction
  const { logs } = await contract.deployed().then(function(instance) {
    return instance.lock(cosmosRecipient, coinDenom, amount, {
      from: accounts[1],
      value: coinDenom === DEFAULT_ETH_DENOM ? amount : 0,
      gas: 300000 // 300,000 Gwei
    });
  });

  // Get event logs
  const event = logs.find(e => e.event === "LogLock");

  // Parse event fields
  const lockEvent = {
    id: event.args._id,
    to: event.args._to,
    token: event.args._token,
    value: Number(event.args._value),
    nonce: Number(event.args._nonce)
  };

  console.log(lockEvent);

  return;
};
