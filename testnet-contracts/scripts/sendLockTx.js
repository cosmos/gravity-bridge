module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  let Web3 = require("web3");
  var HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");
  let contract = truffleContract(require("../build/contracts/Peggy.json"));

  /*******************************************
   *** Constants
   ******************************************/
  const COSMOS_RECIPIENT =
    "0x636f736d6f7331706a74677530766175326d35326e72796b64707a74727438383761796b756530687137646668";
  const COIN_DENOM = "0x0000000000000000000000000000000000000000";
  const AMOUNT = 10;
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
      "https://ropsten.infura.io/".concat(process.env.INFURA_PROJECT_ID)
    );
  } else {
    provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  }

  var web3 = new Web3(provider);
  contract.setProvider(web3.currentProvider);

  /*******************************************
   *** Contract interaction
   ******************************************/
  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  // Send lock transaction
  const { logs } = await contract.deployed().then(function(instance) {
    return instance.lock(COSMOS_RECIPIENT, COIN_DENOM, AMOUNT, {
      from: accounts[1],
      value: AMOUNT,
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
