module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");
  const BigNumber = require("bignumber.js");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  const contract = truffleContract(
    require("../build/contracts/BridgeBank.json")
  );

  const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

  /*******************************************
   *** Constants
   ******************************************/
  // Lock transaction default params
  const DEFAULT_COSMOS_RECIPIENT = Web3.utils.utf8ToHex(
    "cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh"
  );
  const DEFAULT_ETH_DENOM = "eth";
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
    if (NUM_ARGS !== 2 && NUM_ARGS !== 5) {
      return console.error(
        "Error: invalid number of parameters, please try again."
      );
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

  if (!DEFAULT_PARAMS) {
    if (NETWORK_ROPSTEN) {
      cosmosRecipient = Web3.utils.utf8ToHex(process.argv[6]);
      coinDenom = process.argv[7];
      amount = new BigNumber(process.argv[8]);
    } else {
      cosmosRecipient = Web3.utils.utf8ToHex(process.argv[4]);
      coinDenom = process.argv[5];
      amount = new BigNumber(process.argv[6]);
    }
  }

  // Convert default 'eth' coin denom into null address
  if (coinDenom == "eth") {
    coinDenom = NULL_ADDRESS;
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
  try {
    /*******************************************
     *** Contract interaction
     ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();

    // Send lock transaction
    console.log("Connecting to contract....");
    const { logs } = await contract.deployed().then(function (instance) {
      console.log("Connected to contract, sending lock...");
      return instance.lock(cosmosRecipient, coinDenom, amount, {
        from: accounts[0],
        value: coinDenom === NULL_ADDRESS ? amount : 0,
        gas: 300000 // 300,000 Gwei
      });
    });

    console.log("Sent lock...");

    // Get event logs
    const event = logs.find(e => e.event === "LogLock");

    // Parse event fields
    const lockEvent = {
      to: event.args._to,
      from: event.args._from,
      symbol: event.args._symbol,
      token: event.args._token,
      value: Number(event.args._value),
      nonce: Number(event.args._nonce)
    };

    console.log(lockEvent);
  } catch (error) {
    console.error({ error });
  }
  return;
};
