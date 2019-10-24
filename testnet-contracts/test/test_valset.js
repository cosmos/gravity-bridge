const Valset = artifacts.require("Valset");

const { toEthSignedMessageHash, fixSignature } = require("./helpers/helpers");

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

  const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

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
      this.message = web3.utils.soliditySha3(
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

    it("should correctly validate signatures", async function() {
      // Create the signature
      const signature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );

      // Recover the signer address from the generated message and signature.
      const signer = await this.valset.recover(
        toEthSignedMessageHash(this.message),
        signature
      );

      signer.should.be.equal(userOne);
    });

    it("should not validate signatures on a different hashed message", async function() {
      // Create the signature
      const signature = await web3.eth.sign(this.message, userOne);

      // Set up a different message (has increased nonce)
      const differentMessage = web3.utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce + 1 }
      );

      // Recover the signer address from a different message
      const signer = await this.valset.recover(differentMessage, signature);

      signer.should.be.equal(ZERO_ADDRESS);
    });

    it("should not validate signatures from a different address", async function() {
      // Create the signature
      const signature = fixSignature(
        await web3.eth.sign(this.message, userTwo)
      );

      // Recover the signer address from the generated message and signature.
      const signer = await this.valset.recover(
        toEthSignedMessageHash(this.message),
        signature
      );

      signer.should.not.be.equal(userOne);
    });
  });
});
