module.exports = async () => {
  /*******************************************
   *** Set up
   ******************************************/
  const Web3 = require("web3");
  const HDWalletProvider = require("@truffle/hdwallet-provider");

  // Contract abstraction
  const truffleContract = require("truffle-contract");

  const tokenContract = truffleContract(
    require("../build/contracts/BridgeToken.json")
  );

  /*******************************************
   *** Constants
   ******************************************/
  // Config values
  const NETWORK_ROPSTEN =
    process.argv[4] === "--network" && process.argv[5] === "ropsten";
  const NUM_ARGS = process.argv.length - 4;

  // Mint transaction parameters
  const TOKEN_AMOUNT = (1).toString().padEnd(20, "0")
    console.log({TOKEN_AMOUNT})
  /*******************************************
   *** Command line argument error checking
   ***
   *** truffle exec lacks support for dynamic command line arguments:
   *** https://github.com/trufflesuite/truffle/issues/889#issuecomment-522581580
   ******************************************/
  // if (NETWORK_ROPSTEN) {
  //   return console.error(
  //     "Error: token minting on Ropsten network is not supported at this time"
  //   );
  // } else {
  //   if (NUM_ARGS !== 0) {
  //     return console.error(
  //       "Error: custom parameters for token minting are not supported at this time."
  //     );
  //   }
  // }

  /*******************************************
   *** Web3 provider
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
  // const provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
  const web3 = new Web3(provider);
  tokenContract.setProvider(web3.currentProvider);
  try {
  /*******************************************
   *** Contract interaction
   ******************************************/
  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  // Send mint transaction
  const { logs } = await tokenContract.deployed().then(function(instance) {
    return instance.mint(accounts[0], TOKEN_AMOUNT, {
      from: accounts[0],
      value: 0,
      gas: 300000 // 300,000 Gwei
    });
  });

  // Get event logs
  const event = logs.find(e => e.event === "Transfer");

  // Parse event fields
  const transferEvent = {
    from: event.args.from,
    to: event.args.to,
    value: Number(event.args.value)
  };

  console.log(transferEvent);
} catch (error) {
  console.error({error})
}
  return;
};
