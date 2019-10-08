const Oracle = artifacts.require("Oracle");

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

    it("should deploy the Oracle and correctly set initial validators", async function() {
      this.oracle.should.exist;

      const oracleValidators = await this.oracle.validators();
      const oracleNumbValidators = await this.oracle.numbValidators();

      oracleNumbValidators.should.be.equal(this.initialValidators.length());
      oracleValidators[0].should.be.equal(this.initialValidators[0]);
      oracleValidators[1].should.be.equal(this.initialValidators[1]);
      oracleValidators[2].should.be.equal(this.initialValidators[2]);
    });

    it("should deploy the Oracle and correctly set initial validator powers", async function() {
      this.oracle.should.exist;

      const oraclePowers = await this.oracle.powers();
      const oracleTotalPower = await this.oracle.totalPower();

      oracleTotalPower.should.be.equal(
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2]
      );
      oraclePowers[0].should.be.equal(this.initialPowers[0]);
      oraclePowers[1].should.be.equal(this.initialPowers[1]);
      oraclePowers[2].should.be.equal(this.initialPowers[2]);
    });
  });

  // TODO: Testing of Signature verification
  describe("Signature verification", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      //   this.signedHash = Web3Utils.soliditySha3(
      //     { t: "address payable", v: userOne },
      //     { t: "bytes", v: this.recipient },
      //     { t: "address", v: this.token.address },
      //     { t: "int256", v: this.amount },
      //     { t: "int256", v: 1 }
      //   );

      this.oracle = await Oracle.new(
        this.initialValidators,
        this.initialPowers
      );
    });
  });
});
