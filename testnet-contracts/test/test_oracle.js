const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Oracle", function(accounts) {
  const operator = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Oracle smart contract deployment", function() {
    beforeEach(async function() {
      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
      );
    });

    it("should deploy the Oracle, correctly setting the operator and valset", async function() {
      this.oracle.should.exist;

      const oracleOperator = await this.oracle.operator();
      oracleOperator.should.be.equal(operator);

      const oracleValset = await this.oracle.valset();
      oracleValset.should.be.equal(this.valset.address);
    });
  });

  describe("Creation of oracle claims", function() {
    beforeEach(async function() {
      this.bridgeClaimID = 1;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.amount = 100;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.bridgeClaimID },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
      );

      // Submit a new bridge claim to the CosmosBridge to make oracle claims upon
      await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should not allow oracle claims upon inactive bridge claims", async function() {
      const inactiveBridgeClaimID = this.bridgeClaimID + 5;

      // Create hash using Solidity's Sha3 hashing function
      const inactiveBridgeClaimMessage = web3.utils.soliditySha3(
        { t: "uint256", v: inactiveBridgeClaimID },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );
      // Generate signature from userOne (validator)
      const signature = fixSignature(
        await web3.eth.sign(inactiveBridgeClaimMessage, userOne)
      );

      await this.oracle
        .newOracleClaim(
          inactiveBridgeClaimID,
          toEthSignedMessageHash(inactiveBridgeClaimMessage),
          signature,
          {
            from: userOne
          }
        )
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow non-validators to make oracle claims", async function() {
      // Generate signature from userOne (validator)
      const signature = fixSignature(
        await web3.eth.sign(this.message, accounts[6])
      );

      await this.oracle
        .newOracleClaim(
          this.bridgeClaimID,
          toEthSignedMessageHash(this.message),
          signature,
          {
            from: accounts[6]
          }
        )
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow validators to make OracleClaims with invalid signatures", async function() {
      const badMessage = web3.utils.soliditySha3(
        { t: "uint256", v: 20 },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      // Generate signature from userTwo (validator) on bad message
      const signature = fixSignature(await web3.eth.sign(badMessage, userTwo));

      // userTwo submits the expected message with an invalid signature
      await this.oracle
        .newOracleClaim(
          this.bridgeClaimID,
          toEthSignedMessageHash(this.message),
          signature,
          {
            from: userTwo
          }
        )
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow validators to make OracleClaims with another validator's signature", async function() {
      // Generate signature from userOne (validator)
      const signature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );

      // userTwo submits the expected message with userOne's valid signature
      await this.oracle
        .newOracleClaim(
          this.bridgeClaimID,
          toEthSignedMessageHash(this.message),
          signature,
          {
            from: userTwo
          }
        )
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow validators to make OracleClaims with their own signatures", async function() {
      // Generate signature from userOne (validator)
      const signature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );

      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        signature,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should emit an event containing the new OracleClaim's information", async function() {
      // Generate signature from userOne (validator)
      const signature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );

      // Get the logs from a new OracleClaim
      const { logs } = await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        signature,
        {
          from: userOne
        }
      );
      const event = logs.find(e => e.event === "LogNewOracleClaim");

      // Confirm that the event data is correct
      Number(event.args._bridgeClaimID).should.be.bignumber.equal(
        this.bridgeClaimID
      );
      event.args._validatorAddress.should.be.equal(userOne);
      event.args._message.should.be.equal(toEthSignedMessageHash(this.message));
      event.args._signature.should.be.equal(signature);
    });
  });

  describe("Prophecy processing", function() {
    beforeEach(async function() {
      this.bridgeClaimID = 1;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.amount = 100;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.bridgeClaimID },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Set up total power
      this.totalPower =
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2];

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
      );

      // Submit a new bridge claim to the CosmosBridge to make oracle claims upon
      await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Generate signatures from active validators userOne, userTwo, userThree
      this.userOneSignature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );
      this.userTwoSignature = fixSignature(
        await web3.eth.sign(this.message, userTwo)
      );
      this.userThreeSignature = fixSignature(
        await web3.eth.sign(this.message, userThree)
      );
    });

    it("should allow for the processing of prophecies", async function() {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userOneSignature,
        {
          from: userOne
        }
      );
      // Validator userTwo makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userTwoSignature,
        {
          from: userTwo
        }
      );
      // Validator userThree makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      await this.oracle.processProphecyClaim(
        this.bridgeClaimID
      ).should.be.fulfilled;
    });

    it("should allow non-unanimous consensus if signed power passes threshold", async function() {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userOneSignature,
        {
          from: userOne
        }
      );
      // Validator userThree makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      // Confirm that our validators' powers are sufficient to pass the threshold
      const signedPowerWeighted =
        (this.initialPowers[0] + this.initialPowers[2]) * 3;
      const totalPowerWeighted = this.totalPower * 2;

      signedPowerWeighted.should.be.bignumber.greaterThan(totalPowerWeighted);

      // Process prophecy should be fulfilled
      await this.oracle.processProphecyClaim(
        this.bridgeClaimID
      ).should.be.fulfilled;
    });

    it("should emit an event upon successful prophecy processing", async function() {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userOneSignature,
        {
          from: userOne
        }
      );
      // Validator userTwo makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userTwoSignature,
        {
          from: userTwo
        }
      );
      // Validator userThree makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      const submitter = accounts[7];

      const { logs } = await this.oracle.processProphecyClaim(
        this.bridgeClaimID,
        {
          from: submitter
        }
      );

      const event = logs.find(e => e.event === "LogProphecyProcessed");
      Number(event.args._cosmosBridgeClaimId).should.be.bignumber.equal(
        this.bridgeClaimID
      );
      Number(event.args._signedPower).should.be.bignumber.equal(
        this.totalPower * 3
      );
      Number(event.args._totalPower).should.be.bignumber.equal(
        this.totalPower * 2
      );
      event.args._submitter.should.be.equal(submitter);
    });

    // TODO: should not include the signatures of non-active validators
    // TODO: should not allow for the processing of bridge claims whose original validator is no longer active
    // TODO: should not allow bridge claims to be processed twice (e.g. update BridgeClaim status to inactive)
  });
});

// Helpers
function fixSignature(signature) {
  // in geth its always 27/28, in ganache its 0/1. Change to 27/28 to prevent
  // signature malleability if version is 0/1
  // see https://github.com/ethereum/go-ethereum/blob/v1.8.23/internal/ethapi/api.go#L465
  let v = parseInt(signature.slice(130, 132), 16);
  if (v < 27) {
    v += 27;
  }
  const vHex = v.toString(16);
  return signature.slice(0, 130) + vHex;
}

function toEthSignedMessageHash(messageHex) {
  const messageBuffer = Buffer.from(messageHex.substring(2), "hex");
  const prefix = Buffer.from(
    `\u0019Ethereum Signed Message:\n${messageBuffer.length}`
  );
  return web3.utils.sha3(Buffer.concat([prefix, messageBuffer]));
}
