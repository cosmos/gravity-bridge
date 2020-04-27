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
    const tokenContract = truffleContract(
      require("../build/contracts/BridgeToken.json")
    );
  
    const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";
  
    /*******************************************
     *** Constants
     ******************************************/
    // Burn transaction default params
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
     *** Burn transaction parameters
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
    tokenContract.setProvider(web3.currentProvider);
    try {
      /*******************************************
       *** Contract interaction
       ******************************************/
      // Get current accounts
      const accounts = await web3.eth.getAccounts();

        // Send approve transaction
      if(coinDenom != "eth") {
        const bridgeContractAddress = await contract
        .deployed()
        .then(function(instance) {
          return instance.address;
        });

        instance = await tokenContract.at(coinDenom)
        const { logs } = await instance.approve(bridgeContractAddress, amount, {
          from: accounts[9],
          value: 0,
          gas: 300000 // 300,000 Gwei
        });
  
        console.log("Sent approval...");
          
        // Get event logs
        const eventA = logs.find(e => e.event === "Approval");
  
        // Parse event fields
        const approvalEvent = {
          owner: eventA.args.owner,
          spender: eventA.args.spender,
          value: Number(eventA.args.value)
        };
  
        console.log(approvalEvent);
      }
     
      // Send Burn transaction
      console.log("Connecting to contract....");
      const { logs: logs2 } = await contract.deployed().then(function (instance) {
        console.log("Connected to contract, sending burn...");
        return instance.burn(cosmosRecipient, coinDenom, amount, {
          from: accounts[9],
          value: coinDenom === NULL_ADDRESS ? amount : 0,
          gas: 300000 // 300,000 Gwei
        });
      });
  
      console.log("Sent burn...");
  
      // Get event logs
      const eventB = logs2.find(e => e.event === "LogBurn");
  
      // Parse event fields
      const burnEvent = {
        to: eventB.args._to,
        from: eventB.args._from,
        symbol: eventB.args._symbol,
        token: eventB.args._token,
        value: Number(eventB.args._value),
        nonce: Number(eventB.args._nonce)
      };
  
      console.log(burnEvent);
    } catch (error) {
      console.error({ error });
    }
    return;
  };
  