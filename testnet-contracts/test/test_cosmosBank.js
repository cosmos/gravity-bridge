const CosmosBank = artifacts.require("TestCosmosBank");
const CosmosToken = artifacts.require("CosmosToken");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBank", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("CosmosBank smart contract deployment", function() {
    it("should deploy the CosmosBank", async function() {
      this.cosmosBank = await CosmosBank.new();
      this.cosmosBank.should.exist;
    });
  });

  describe("Cosmos token deployment", function() {
    beforeEach(async function() {
      this.cosmosBank = await CosmosBank.new();
      this.symbol = "ABC";
    });

    it("should deploy new bank tokens", async function() {
      await this.cosmosBank.callDeployNewCosmosToken(this.symbol, {
        from: provider
      }).should.be.fulfilled;
    });

    // TODO: Once Peggy gas deployment issues are resolved, uncomment this test
    // it("should increase the token count upon new bank token deployment", async function() {
    //   const priorTokenCount = await this.cosmosBank.numbTokens();
    //   priorTokenCount.should.be.equal(0);

    //   await this.cosmosBank.callDeployBankToken(this.symbol, { from: provider });

    //   const afterTokenCount = await this.cosmosBank.numbTokens();
    //   afterTokenCount.should.be.equal(1);
    // });

    it("should return the new Cosmos token's address", async function() {
      const newCosmosTokenAddress = await this.cosmosBank.callDeployNewCosmosToken(
        this.symbol,
        {
          from: provider
        }
      );

      // TODO: Check bank token ethereum address type
    });

    // TODO: Once Peggy gas deployment issues are resolved, uncomment this test
    // it("should emit event LogTokenDeployed containing the new bank token's address", async function() {
    //   const expectedTokenAddress = await this.cosmosBank.callDeployBankToken(
    //     this.symbol,
    //     {
    //       from: provider
    //     }
    //   );

    //   const event = logs.find(e => e.event === "LogBankTokenDeploy");
    //   event.args._token.should.be.equal(expectedTokenAddress);
    // });
  });

  describe("Cosmos token minting", function() {
    beforeEach(async function() {
      this.cosmosBank = await CosmosBank.new();

      this.token = "0x0000000000000000000000000000000000000000";
      this.symbol = "ABC";
      this.amount = 100;
    });

    // TODO: Once Peggy gas deployment issues are resolved, uncomment this test
    // it("should mint new bank tokens", async function() {
    //   await this.cosmosBank.callDeliver(this.token, this.symbol, 100, userOne, {
    //     from: provider
    //   }).should.be.fulfilled;
    // });

    // TODO: Once Peggy gas deployment issues are resolved, uncomment this test
    // it("should emit event LogBankTokenMint upon successful minting of bank tokens", async function() {
    //   await this.cosmosBank.callDeliver(
    //     this.token,
    //     this.symbol,
    //     this.amount,
    //     userOne,
    //     {
    //       from: provider
    //     }
    //   );

    //   const event = logs.find(e => e.event === "LogBankTokenMint");
    //   event.args._token.should.be.equal(this.token);
    //   event.args._symbol.should.be.equal(this.symbol);
    //   Number(event.args._amount).should.be.bignumber.equal(this.amount);
    //   event.args._beneficiary.should.be.equal(userOne.address);
    // });
  });
});
