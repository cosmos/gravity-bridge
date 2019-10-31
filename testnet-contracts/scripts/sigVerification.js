module.exports = async () => {
  // Import test helpers
  const {
    toEthSignedMessageHash,
    fixSignature
  } = require("../test/helpers/helpers");

  // Import Web3
  const Web3 = require("web3");

  // Set up accounts
  let provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);

  const web3 = new Web3(provider);

  console.log("web3");
  console.log(web3);

  // Get current accounts
  const accounts = await web3.eth.getAccounts();

  console.log("ACCOUNTS");
  console.log(accounts[0]);

  // Set up data
  const prophecyID = 8;
  const sender = "cosmos1qwnw2r9ak79536c4dqtrtk2pl2nlzpqh763rls";
  const recipient = "0x7B95B6EC7EbD73572298cEf32Bb54FA408207359";
  const token = "0xbEDdB076fa4dF04859098A9873591dcE3E9C404d";
  const amount = 1;
  const validator = "0x22448F19a3Cb4EE570C2ec6d3323761d3399bbbD";

  // Create hash using Solidity's Sha3 hashing function
  const message = web3.utils.soliditySha3(
    { t: "uint256", v: prophecyID },
    { t: "bytes", v: sender },
    { t: "address payable", v: recipient },
    { t: "address", v: token },
    { t: "uint256", v: amount },
    { t: "address", v: validator }
  );

  console.log("Signing hash:", message);

  const userTwoSig = fixSignature(await web3.eth.sign(message, accounts[2]));

  console.log("Signature", userTwoSig);

  return;
};
