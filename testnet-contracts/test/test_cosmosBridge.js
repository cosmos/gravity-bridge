const CosmosBridge = artifacts.require("TestCosmosBridge");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBridge", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("CosmosBridge smart contract deployment", function() {
    beforeEach(async function() {
      this.cosmosBridge = await CosmosBridge.new();
    });

    it("should deploy the CosmosBridge with the correct parameters", async function() {
      this.cosmosBridge.should.exist;

      const bridgeNonce = await this.cosmosBridge.cosmosBridgeNonce();
      Number(bridgeNonce).should.be.bignumber.equal(0);
    });
  });

  describe("Creation of CosmosBridgeClaims", function() {
    beforeEach(async function() {
      // Set up CosmosBridgeClaim values
      this.nonce = 4;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.amount = 100;

      // Deploy new CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new();
    });

    it("should allow for the creation of new CosmosBridgeClaims", async function() {
      await this.cosmosBridge.callNewCosmosBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: provider
        }
      ).should.be.fulfilled;
    });

    it("should log an event containing the new CosmosBridgeClaim's information", async function() {
      const { logs } = await this.cosmosBridge.callNewCosmosBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: provider
        }
      );

      const event = logs.find(e => e.event === "LogNewCosmosBridgeClaim");

      Number(event.args._cosmosBridgeNonce).should.be.bignumber.equal(1);
      Number(event.args._nonce).should.be.bignumber.equal(this.nonce);
      event.args._cosmosSender.should.be.equal(this.cosmosSender);
      event.args._ethereumReceiver.should.be.equal(this.ethereumReceiver);
      event.args._validatorAddress.should.be.equal(provider);
      event.args._tokenAddress.should.be.equal(this.tokenAddress);
      event.args._symbol.should.be.equal(this.symbol);
      Number(event.args._amount).should.be.bignumber.equal(this.amount);
    });

    it("should increase the CosmosBridge nonce upon the creation of new a CosmosBridgeClaim", async function() {
      const priorBridgeNonce = await this.cosmosBridge.cosmosBridgeNonce();
      Number(priorBridgeNonce).should.be.bignumber.equal(0);

      await this.cosmosBridge.callNewCosmosBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: provider
        }
      );

      const postBridgeNonce = await this.cosmosBridge.cosmosBridgeNonce();
      Number(postBridgeNonce).should.be.bignumber.equal(1);
    });
  });

  describe("CosmosBridgeClaim status", function() {
    beforeEach(async function() {
      // Set up CosmosBridgeClaim values
      this.nonce = 4;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.amount = 100;

      // Deploy new CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new();
    });

    it("should allow anyone to check the status of a CosmosBridgeClaim", async function() {
      // Create the CosmosBridgeClaim
      const { logs } = await this.cosmosBridge.callNewCosmosBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: provider
        }
      );

      const event = logs.find(e => e.event === "LogNewCosmosBridgeClaim");
      const cosmosBridgeClaimNonce = event.args._cosmosBridgeNonce;

      // Get the CosmosBridgeClaim's status
      const status = await this.cosmosBridge.getCosmosBridgeClaimStatus(
        cosmosBridgeClaimNonce,
        {
          from: provider
        }
      );

      // Solidity enums are represented as integers
      Number(status).should.be.bignumber.equal(0);
    });
  });
});
