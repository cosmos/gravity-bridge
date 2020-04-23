const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Oracle", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];
  const userFour = accounts[4];

  // User account
  const userSeven = accounts[7];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold of 70%
  const consensusThreshold = 70;

  describe("Oracle smart contract deployment", function () {
    beforeEach(async function () {
      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      );
    });

    it("should deploy the Oracle, correctly setting the operator and valset", async function () {
      this.oracle.should.exist;

      const oracleOperator = await this.oracle.operator();
      oracleOperator.should.be.equal(operator);

      const oracleValset = await this.oracle.valset();
      oracleValset.should.be.equal(this.valset.address);
    });
  });

  describe("Creation of oracle claims", function () {
    beforeEach(async function () {
      this.prophecyID = 1;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.symbol = "TEST";
      this.amount = 100;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.prophecyID },
        { t: "bytes", v: this.cosmosSender },
        {
          t: "address payable",
          v: this.ethereumReceiver
        },
        { t: "uint256", v: this.amount }
      );

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );

      // Operator sets Oracle
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      });

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should not allow oracle claims upon inactive prophecy claims", async function () {
      const inactiveBridgeClaimID = this.prophecyID + 50;

      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // Validator userOne submits an oracle claim
      await this.oracle
        .newOracleClaim(inactiveBridgeClaimID, this.message, signature, {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow non-validators to make oracle claims", async function () {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // Validator userOne submits an oracle claim
      await this.oracle
        .newOracleClaim(this.prophecyID, this.message, signature, {
          from: userSeven
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow validators to make OracleClaims with invalid signatures", async function () {
      const badMessage = web3.utils.soliditySha3(
        {
          t: "uint256",
          v: 20
        },
        {
          t: "bytes",
          v: this.cosmosSender
        },
        {
          t: "address payable",
          v: this.ethereumReceiver
        },
        {
          t: "uint256",
          v: this.amount
        }
      );

      // Generate signature from userTwo (validator) on bad message
      const signature = await web3.eth.sign(badMessage, userTwo);

      // Validator userOne submits an oracle claim
      await this.oracle
        .newOracleClaim(this.prophecyID, this.message, signature, {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow validators to make OracleClaims with another validator's signature", async function () {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // userTwo submits the expected message with userOne's valid signature
      await this.oracle
        .newOracleClaim(this.prophecyID, this.message, signature, {
          from: userTwo
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow valid OracleClaims", async function () {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // Validator makes an oracle claim with their signature
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        signature,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should not allow validators to make duplicate OracleClaims", async function () {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // Validator makes the first oracle claim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        signature,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Validator attempts to make a second oracle claim on the same bridge claim
      await this.oracle
        .newOracleClaim(this.prophecyID, this.message, signature, {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should emit an event containing the new OracleClaim's information", async function () {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.message, userOne);

      // Get the logs from a new OracleClaim
      const { logs } = await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        signature,
        {
          from: userOne
        }
      ).should.be.fulfilled;
      const event = logs.find(e => e.event === "LogNewOracleClaim");

      // Confirm that the event data is correct
      Number(event.args._prophecyID).should.be.bignumber.equal(this.prophecyID);
      event.args._validatorAddress.should.be.equal(userOne);
      event.args._message.should.be.equal(this.message);
      event.args._signature.should.be.equal(signature);
    });
  });

  describe("Prophecy processing", function () {
    beforeEach(async function () {
      this.prophecyID = 1;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.symbol = "TEST";
      this.amount = 100;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.prophecyID },
        { t: "bytes", v: this.cosmosSender },
        {
          t: "address payable",
          v: this.ethereumReceiver
        },
        { t: "uint256", v: this.amount }
      );

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Set up total power
      this.totalPower = this.initialPowers.reduce(function (a, b) {
        return a + b;
      }, 0);

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );

      // Operator sets Oracle
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      });

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Generate signatures from active validators userOne, userTwo, userThree
      this.userOneSignature = await web3.eth.sign(this.message, userOne);
      this.userTwoSignature = await web3.eth.sign(this.message, userTwo);
      this.userThreeSignature = await web3.eth.sign(this.message, userThree);
    });

    it("should not process the prophecy if signed power does not pass the required threshold power", async function () {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );

      const isActive = await this.cosmosBridge.isProphecyClaimActive(
        this.prophecyID
      );
      isActive.should.be.equal(true);
    });

    it("should allow for the processing of prophecies", async function () {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );

      // Validator userTwo makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userTwoSignature,
        {
          from: userTwo
        }
      );
      // Validator userThree makes a valid OracleClaim
      const { logs } = await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      const event = logs.find(e => e.event === "LogProphecyProcessed");
      const prophecyID = Number(event.args._prophecyID);

      const isActive = await this.cosmosBridge.isProphecyClaimActive(
        prophecyID
      );
      isActive.should.be.equal(false);
    });

    it("should process prophecies if signed power passes threshold", async function () {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );

      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userTwoSignature,
        {
          from: userTwo
        }
      );

      // Validator userThree makes a valid OracleClaim
      const { logs } = await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userThreeSignature,
        {
          from: userThree
        }
      ).should.be.fulfilled;

      const event = logs.find(e => e.event === "LogProphecyProcessed");
      const processedProphecyID = Number(event.args._prophecyID);
      processedProphecyID.should.be.bignumber.equal(this.prophecyID);

      // Confirm that our validators' powers are sufficient to pass the threshold
      const processedPowerCurrent = Number(event.args._prophecyPowerCurrent);
      const processedPowerThreshold = Number(
        event.args._prophecyPowerThreshold
      );

      processedPowerThreshold.should.be.bignumber.equal(
        Number(this.totalPower) * consensusThreshold
      );

      const expectedCurrentPower =
        Number(this.initialPowers[0]) +
        Number(this.initialPowers[1]) +
        Number(this.initialPowers[2]);

      processedPowerCurrent.should.be.bignumber.equal(
        expectedCurrentPower * Number(100)
      );
    });

    it("should not allow a prophecy to be processed twice", async function () {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );

      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userTwoSignature,
        {
          from: userTwo
        }
      );

      // Validator userThree makes a valid OracleClaim, processing the prophecy claim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      // Attempt to process the same prophecy should be rejected
      await this.oracle
        .processBridgeProphecy(this.prophecyID)
        .should.be.rejectedWith(EVMRevert);
    });

    // TODO: Add these tests once Valset has been with dynamic validator set:
    // 1. Should not include the signatures of non-active validators
    // 2. Should not allow for the processing of bridge claims whose original validator is no longer active

    it("should emit an event upon successful prophecy processing", async function () {
      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );
      // Validator userTwo makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userTwoSignature,
        {
          from: userTwo
        }
      );
      // Validator userThree makes a valid OracleClaim
      const { logs } = await this.oracle.newOracleClaim(
        this.prophecyID,
        this.message,
        this.userThreeSignature,
        {
          from: userThree
        }
      );

      const event = logs.find(e => e.event === "LogProphecyProcessed");
      Number(event.args._prophecyID).should.be.bignumber.equal(this.prophecyID);
      Number(event.args._prophecyPowerThreshold).should.be.bignumber.equal(
        this.totalPower * consensusThreshold
      );

      const expectedCurrentPower =
        Number(this.initialPowers[0]) +
        Number(this.initialPowers[1]) +
        Number(this.initialPowers[2]);

      Number(event.args._prophecyPowerCurrent).should.be.bignumber.equal(
        expectedCurrentPower * Number(100)
      );

      event.args._submitter.should.be.equal(userThree);
    });
  });
});