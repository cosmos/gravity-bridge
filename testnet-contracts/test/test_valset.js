const Valset = artifacts.require("Valset");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Bank", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Valset smart contract deployment", function() {
    it("should deploy the Valset with the correct parameters", async function() {
      const initialValidators = [userOne, userTwo, userThree];
      const initialPowers = [5, 8, 12];

      this.valset = await Valset.new(initialValidators, initialPowers);
      this.valset.should.exist;

      const valsetValidators = await this.valset.validators();
      const valsetNumbValidators = await this.valset.numbValidators();

      valsetNumbValidators.should.be.equal(initialValidators.length());
      valsetValidators[0].should.be.equal(initialValidators[0]);
      valsetValidators[1].should.be.equal(initialValidators[1]);
      valsetValidators[2].should.be.equal(initialValidators[2]);

      const valsetPowers = await this.valset.powers();
      const valsetPowers = await this.valset.totalPower();

      valsetTotalPower.should.be.equal(
        initialPowers[0] + initialPowers[1] + initialPowers[2]
      );
      valsetPowers[0].should.be.equal(initialPowers[0]);
      valsetPowers[1].should.be.equal(initialPowers[1]);
      valsetPowers[2].should.be.equal(initialPowers[2]);
    });
  });

  describe("Bank token deployment", function() {
    beforeEach(async function() {
      this.bank = await Bank.new();
      this.symbol = "ABC";
    });

    it("should deploy new bank tokens", async function() {
      await this.bank.callDeployBankToken(this.symbol, {
        from: provider
      }).should.be.fulfilled;
    });

    it("should increase the token count upon new bank token deployment", async function() {
      const priorTokenCount = await this.bank.numbTokens();
      priorTokenCount.should.be.equal(0);

      await this.bank.callDeployBankToken(this.symbol, { from: provider });

      const afterTokenCount = await this.bank.numbTokens();
      afterTokenCount.should.be.equal(1);
    });

    it("should return the new bank token's address", async function() {
      const newBankTokenAddress = await this.bank.callDeployBankToken(
        this.symbol,
        {
          from: provider
        }
      );

      // TODO: Check bank token ethereum address type
    });

    it("should emit event LogTokenDeployed containing the new bank token's address", async function() {
      const expectedTokenAddress = await this.bank.callDeployBankToken(
        this.symbol,
        {
          from: provider
        }
      );

      const event = logs.find(e => e.event === "LogTokenDeploy");
      event.args._token.should.be.equal(expectedTokenAddress);
    });
  });
});
