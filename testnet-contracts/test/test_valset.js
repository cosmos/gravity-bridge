const Valset = artifacts.require("Valset");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Valset", function(accounts) {
  const operator = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Valset contract deployment", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should deploy the Valset and correctly set initial validator count", async function() {
      this.valset.should.exist;

      const valsetValidatorCount = await this.valset.validatorCount();

      Number(valsetValidatorCount).should.be.bignumber.equal(
        this.initialValidators.length
      );
    });

    it("should correctly set initial validators", async function() {
      const userOneValidator = await this.valset.activeValidators(userOne);
      const userTwoValidator = await this.valset.activeValidators(userTwo);
      const userThreeValidator = await this.valset.activeValidators(userThree);

      userOneValidator.should.be.equal(true);
      userTwoValidator.should.be.equal(true);
      userThreeValidator.should.be.equal(true);
    });

    it("should correctly set initial validator powers ", async function() {
      const userOnePower = await this.valset.powers(userOne);
      const userTwoPower = await this.valset.powers(userTwo);
      const userThreePower = await this.valset.powers(userThree);

      Number(userOnePower).should.be.bignumber.equal(this.initialPowers[0]);
      Number(userTwoPower).should.be.bignumber.equal(this.initialPowers[1]);
      Number(userThreePower).should.be.bignumber.equal(this.initialPowers[2]);
    });

    it("should correctly set the total power", async function() {
      const valsetTotalPower = await this.valset.totalPower();

      Number(valsetTotalPower).should.be.bignumber.equal(
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2]
      );
    });
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

      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should validate signatures", async function() {
      const prefix = "\x19Ethereum Signed Message:\n" + this.hash.length;
      const prefixedMessageHash = web3.utils.sha3(prefix + this.hash);

      // Generate signature from active validator userOne
      const signature = await web3.eth.sign(prefixedMessageHash, userOne);

      let originalSigner = await web3.eth.accounts.recover(
        prefixedMessageHash,
        signature
      );

      //   const r1 = "0x" + sig.slice(2, 64 + 2);
      //   const s1 = "0x" + sig.slice(64 + 2, 128 + 2);
      //   const v1 = "0x" + sig.slice(128 + 2, 130 + 2);

      // Validate userOne's signature
      // const originalSigner = await this.valset.testRecovery(
      //   prefixedMessageHash,
      //   signature
      // );

      originalSigner.should.be.equal(userOne);
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
