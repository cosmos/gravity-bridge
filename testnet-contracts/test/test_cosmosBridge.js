const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBridge", function(accounts) {
  const operator = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("CosmosBridge smart contract deployment", function() {
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
    });

    it("should deploy the CosmosBridge with the correct parameters", async function() {
      this.cosmosBridge.should.exist;

      const claimCount = await this.cosmosBridge.bridgeClaimCount();
      Number(claimCount).should.be.bignumber.equal(0);

      const cosmosBridgeValset = await this.cosmosBridge.valset();
      cosmosBridgeValset.should.be.equal(this.valset.address);
    });
  });

  describe("Creation of bridge claims", function() {
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
    });

    it("should allow for the creation of new bridge claims", async function() {
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

    it("should log an event containing the new bridge claim's information", async function() {
      const { logs } = await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewBridgeClaim");

      Number(event.args._bridgeClaimCount).should.be.bignumber.equal(1);
      Number(event.args._nonce).should.be.bignumber.equal(this.nonce);
      event.args._cosmosSender.should.be.equal(this.cosmosSender);
      event.args._ethereumReceiver.should.be.equal(this.ethereumReceiver);
      event.args._validatorAddress.should.be.equal(userOne);
      event.args._tokenAddress.should.be.equal(this.tokenAddress);
      event.args._symbol.should.be.equal(this.symbol);
      Number(event.args._amount).should.be.bignumber.equal(this.amount);
    });

    it("should increase the bridge claim count upon the creation of new a bridge claim", async function() {
      const priorBridgeClaimCount = await this.cosmosBridge.bridgeClaimCount();
      Number(priorBridgeClaimCount).should.be.bignumber.equal(0);

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
      );

      const postBridgeClaimCount = await this.cosmosBridge.bridgeClaimCount();
      Number(postBridgeClaimCount).should.be.bignumber.equal(1);
    });
  });

  describe("Bridge claim status", function() {
    beforeEach(async function() {
      // Set up BridgeClaim values
      this.nonce = 4;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "TEST";
      this.amount = 100;

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
    });
    it("should not show fake bridge claims as active", async function() {
      const bridgeClaimCount = 4;

      // Get a BridgeClaim's status
      const status = await this.cosmosBridge.isBridgeClaimActive(
        bridgeClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should not be active
      status.should.be.equal(false);
    });

    it("should allow users to check if a bridge claim is currently active", async function() {
      // Create the BridgeClaim
      const { logs } = await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewBridgeClaim");
      const bridgeClaimCount = event.args._bridgeClaimCount;

      // Get the BridgeClaim's status
      const status = await this.cosmosBridge.isBridgeClaimActive(
        bridgeClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);
    });

    it("should allow users to check if a bridge claim's original validator is currently an active validator", async function() {
      // Create the BridgeClaim
      const { logs } = await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewBridgeClaim");
      const bridgeClaimCount = event.args._bridgeClaimCount;

      // Get the BridgeClaim's status
      const status = await this.cosmosBridge.isBridgeClaimValidatorActive(
        bridgeClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);
    });
  });
});
