const Oracle = artifacts.require("TestOracle");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Oracle", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Oracle smart contract deployment", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];
      this.oracle = await Oracle.new(
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should deploy the Oracle and correctly set initial validator count", async function() {
      this.oracle.should.exist;

      const oracleNumbValidators = await this.oracle.numbValidators();

      Number(oracleNumbValidators).should.be.bignumber.equal(
        this.initialValidators.length
      );
    });

    it("should correctly set initial validators", async function() {
      const userOneValidator = await this.oracle.activeValidators(userOne);
      const userTwoValidator = await this.oracle.activeValidators(userTwo);
      const userThreeValidator = await this.oracle.activeValidators(userThree);

      userOneValidator.should.be.equal(true);
      userTwoValidator.should.be.equal(true);
      userThreeValidator.should.be.equal(true);
    });

    it("should correctly set initial validator powers ", async function() {
      const userOnePower = await this.oracle.powers(userOne);
      const userTwoPower = await this.oracle.powers(userTwo);
      const userThreePower = await this.oracle.powers(userThree);

      Number(userOnePower).should.be.bignumber.equal(this.initialPowers[0]);
      Number(userTwoPower).should.be.bignumber.equal(this.initialPowers[1]);
      Number(userThreePower).should.be.bignumber.equal(this.initialPowers[2]);
    });

    it("should correctly set the total power", async function() {
      const oracleTotalPower = await this.oracle.totalPower();

      Number(oracleTotalPower).should.be.bignumber.equal(
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2]
      );
    });
  });

  describe("Creation of OracleClaims", function() {
    beforeEach(async function() {
      // Create hash using Solidity's Sha3 hashing function
      this.cosmosBridgeNonce = 3;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.hash = Web3Utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      // Generate signature from userOne (validator)
      this.signature = await web3.eth.sign(this.hash, userOne);

      // Deploy Oracle contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.oracle = await Oracle.new(
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should allow for the creation of new OracleClaims", async function() {
      await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        provider,
        this.hash,
        this.signature,
        {
          from: provider
        }
      ).should.be.fulfilled;
    });

    it("should emit an event containing the new OracleClaim's information", async function() {
      const { logs } = await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        provider,
        this.hash,
        this.signature,
        {
          from: provider
        }
      );

      const event = logs.find(e => e.event === "LogNewOracleClaim");

      Number(event.args._cosmosBridgeClaimId).should.be.bignumber.equal(
        this.cosmosBridgeNonce
      );
      event.args._validatorAddress.should.be.equal(provider);
      event.args._contentHash.should.be.equal(this.hash);
      event.args._signature.should.be.equal(this.signature);
    });

    // TODO: Access mapped array
    // it("should index the OracleClaim by the associated CosmosBridgeClaimId", async function() {
    //   await this.oracle.oracleClaims(this.cosmosBridgeNonce[0], {
    //     from: provider
    //   }).should.be.fulfilled;
    // });

    // it("should index the OracleClaim by validator address", async function() {
    //   await this.oracle.validatorOracleClaims(provider[0], {
    //     from: provider
    //   }).should.be.fulfilled;
    // });
  });

  describe("Signature verification", function() {
    beforeEach(async function() {
      // Create hash using Solidity's Sha3 hashing function
      this.cosmosBridgeNonce = 3;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.hash = Web3Utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.oracle = await Oracle.new(
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should validate signatures containing a CosmosBridgeClaim's information", async function() {
      // Generate signature from userOne (validator)
      const signature = await web3.eth.sign(this.hash, userOne);

      // validator userOne makes a new claim
      await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        userOne,
        this.hash,
        signature,
        {
          from: provider
        }
      );

      // Parse signature
      const sig = parseSignature(signature);

      // Validate signature
      const valid = await this.oracle.isValidSignature(
        userOne,
        this.hash,
        sig.v,
        sig.r,
        sig.s
      );

      valid.should.be.equal(true);
    });
  });

  describe("Prophecy processing", function() {
    beforeEach(async function() {
      // Create hash using Solidity's Sha3 hashing function
      this.cosmosBridgeNonce = 3;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.hash = Web3Utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      // Generate signature from all active validators
      this.userOneSignature = await web3.eth.sign(this.hash, userOne);
      this.userTwoSignature = await web3.eth.sign(this.hash, userTwo);
      this.userThreeSignature = await web3.eth.sign(this.hash, userThree);

      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];
      this.totalPower =
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2];

      this.oracle = await Oracle.new(
        this.initialValidators,
        this.initialPowers
      );

      // Validator userOne makes a new claim
      await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        userOne,
        this.hash,
        this.userOneSignature,
        {
          from: provider
        }
      );

      // Validator userOne makes a new claim
      await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        userTwo,
        this.hash,
        this.userTwoSignature,
        {
          from: provider
        }
      );

      // Validator userOne makes a new claim
      await this.oracle.callNewOracleClaim(
        this.cosmosBridgeNonce,
        userThree,
        this.hash,
        this.userThreeSignature,
        {
          from: provider
        }
      );
    });

    it("should allow for the processing of prophecies", async function() {
      const sigUserOne = parseSignature(this.userOneSignature);
      const sigUserTwo = parseSignature(this.userTwoSignature);
      const sigUserThree = parseSignature(this.userThreeSignature);

      const sigV = [sigUserOne.v, sigUserTwo.v, sigUserThree.v];
      const sigR = [sigUserOne.r, sigUserTwo.r, sigUserThree.r];
      const sigS = [sigUserOne.s, sigUserTwo.s, sigUserThree.s];

      await this.oracle.callProcessProphecyClaim(
        this.cosmosBridgeNonce,
        this.hash,
        [userOne, userTwo, userThree],
        sigV,
        sigR,
        sigS
      ).should.be.fulfilled;
    });

    it("should allow nonunanimous consensus if signed power passes threshold", async function() {
      const sigUserOne = parseSignature(this.userOneSignature);
      const sigUserThree = parseSignature(this.userThreeSignature);

      const sigV = [sigUserOne.v, sigUserThree.v];
      const sigR = [sigUserOne.r, sigUserThree.r];
      const sigS = [sigUserOne.s, sigUserThree.s];

      await this.oracle.callProcessProphecyClaim(
        this.cosmosBridgeNonce,
        this.hash,
        [userOne, userThree],
        sigV,
        sigR,
        sigS
      ).should.be.fulfilled;
    });

    it("should emit an event upon successful prophecy processing", async function() {
      const sigUserOne = parseSignature(this.userOneSignature);
      const sigUserTwo = parseSignature(this.userTwoSignature);
      const sigUserThree = parseSignature(this.userThreeSignature);

      const sigV = [sigUserOne.v, sigUserTwo.v, sigUserThree.v];
      const sigR = [sigUserOne.r, sigUserTwo.r, sigUserThree.r];
      const sigS = [sigUserOne.s, sigUserTwo.s, sigUserThree.s];

      const { logs } = await this.oracle.callProcessProphecyClaim(
        this.cosmosBridgeNonce,
        this.hash,
        [userOne, userTwo, userThree],
        sigV,
        sigR,
        sigS
      );

      const event = logs.find(e => e.event === "LogProphecyProcessed");

      Number(event.args._cosmosBridgeClaimId).should.be.bignumber.equal(
        this.cosmosBridgeNonce
      );
      Number(event.args._signedPower).should.be.bignumber.equal(
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2]
      );
      Number(event.args._totalPower).should.be.bignumber.equal(this.totalPower);
      event.args._submitter.should.be.equal(provider);
    });

    it("should not include the signatures of non-active validators", async function() {
      // Generate signature from userFour (non-validator)
      const userFourSignature = await web3.eth.sign(this.hash, accounts[4]);

      const sigUserOne = parseSignature(this.userOneSignature);
      const sigUserTwo = parseSignature(this.userTwoSignature);
      const sigUserFour = parseSignature(userFourSignature);

      const sigV = [sigUserOne.v, sigUserTwo.v, sigUserFour.v];
      const sigR = [sigUserOne.r, sigUserTwo.r, sigUserFour.r];
      const sigS = [sigUserOne.s, sigUserTwo.s, sigUserFour.s];

      await this.oracle
        .callProcessProphecyClaim(
          this.cosmosBridgeNonce,
          this.hash,
          [userOne, userTwo, sigUserFour],
          sigV,
          sigR,
          sigS
        )
        .should.be.rejectedWith("invalid address");
    });

    it("should not allow validator signatures to be applied twice", async function() {
      const sigUserOne = parseSignature(this.userOneSignature);
      const sigUserTwo = parseSignature(this.userTwoSignature);

      const sigV = [sigUserOne.v, sigUserTwo.v, sigUserTwo.v];
      const sigR = [sigUserOne.r, sigUserTwo.r, sigUserTwo.r];
      const sigS = [sigUserOne.s, sigUserTwo.s, sigUserTwo.s];

      await this.oracle
        .callProcessProphecyClaim(
          this.cosmosBridgeNonce,
          this.hash,
          [userOne, userTwo, userTwo],
          sigV,
          sigR,
          sigS
        )
        .should.be.rejectedWith("does not meet the threshold");
    });
  });
});

// Helpers
const parseSignature = signature => {
  const signatureText = signature.substr(2, signature.length);

  const r = "0x" + signatureText.substr(0, 64);
  const s = "0x" + signatureText.substr(64, 64);
  const v = web3.utils.hexToNumber(signatureText.substr(128, 2)) + 27;

  return {
    v,
    r,
    s
  };
};
