module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require('web3');
  const HDWalletProvider = require('@truffle/hdwallet-provider');
  const BigNumber = require('bignumber.js');

  const truffleContract = require('truffle-contract');
  const contract = truffleContract(require('../../build/contracts/BridgeBank.json'));
  const daiJSON = require('./dai.json');

  const DAI_ADDRESS = '0xad6d458402f60fd3bd25163575031acdce07538d';

  /*******************************************
   *** Constants
   ******************************************/
  // Lock transaction default params
  const DEFAULT_COSMOS_RECIPIENT = Web3.utils.utf8ToHex(
    'cosmos1pgkwvwezfy3qkh99hjnf35ek3znzs79mwqf48y'
  );
  const DEFAULT_DAI_DENOM = 'dai';
  const DEFAULT_AMOUNT = '1000000000000000000';

  // Config values
  const NETWORK_ROPSTEN = process.argv[4] === '--network' && process.argv[5] === 'ropsten';
  const DEFAULT_PARAMS =
    process.argv[4] === '--default' || (NETWORK_ROPSTEN && process.argv[6] === '--default');
  const NUM_ARGS = process.argv.length - 4;

  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  if (NETWORK_ROPSTEN && DEFAULT_PARAMS) {
    if (NUM_ARGS !== 3) {
      return console.error('Error: custom parameters are invalid on --default.');
    }
  } else if (NETWORK_ROPSTEN) {
    if (NUM_ARGS !== 2 && NUM_ARGS !== 5) {
      return console.error('Error: invalid number of parameters, please try again.');
    }
  } else if (DEFAULT_PARAMS) {
    if (NUM_ARGS !== 1) {
      return console.error('Error: custom parameters are invalid on --default.');
    }
  } else {
    if (NUM_ARGS !== 3) {
      return console.error('Error: must specify recipient address, token address, and amount.');
    }
  }

  /*******************************************
   *** Lock transaction parameters
   ******************************************/
  let cosmosRecipient = DEFAULT_COSMOS_RECIPIENT;
  let coinDenom = DEFAULT_DAI_DENOM;
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

  // Convert default 'dai' coin denom into dai address
  if (coinDenom == 'dai') {
    coinDenom = DAI_ADDRESS;
  }

  /*******************************************
   *** Web3 provider
   *** Set contract provider based on --network flag
   ******************************************/
  let provider;
  if (NETWORK_ROPSTEN) {
    provider = new HDWalletProvider(
      process.env.MNEMONIC,
      'https://ropsten.infura.io/v3/'.concat(process.env.INFURA_PROJECT_ID)
    );
  } else {
    provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  }
  console.log(process.env.INFURA_PROJECT_ID);

  const web3 = new Web3(provider);
  contract.setProvider(web3.currentProvider);

  const daiContract = new web3.eth.Contract(daiJSON, DAI_ADDRESS);

  try {
    /*******************************************
     *** Contract interaction
     ******************************************/
    // Get current accounts

    // sender approve Bridge Bank use DAI
    console.log('Aprove for Bridge Bank');
    await daiContract.methods.approve(contract.networks['3'].address, amount).send({
      from: '0x8f287eA4DAD62A3A626942d149509D6457c2516C',
      value: 0,
      gas: 300000 // 300,000 Gwei
    });

    console.log(approveLogs);

    // Send lock transaction
    console.log('Connecting to contract....');
    const { logs } = await contract.deployed().then(function (instance) {
      console.log('Connected to contract, sending lock...');
      return instance.lock(cosmosRecipient, coinDenom, amount, {
        from: '0x8f287eA4DAD62A3A626942d149509D6457c2516C',
        value: 0,
        gas: 300000 // 300,000 Gwei
      });
    });

    console.log('Sent lock...');

    // Get event logs
    const event = logs.find(e => e.event === 'LogLock');
    console.log(event);

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
    process.exit(0);
  } catch (error) {
    console.error({ error });
    process.exit(1);
  }
};
