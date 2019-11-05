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

contract("CosmosBridge", function(accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 0;
  const CLAIM_TYPE_LOCK = 1;

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
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );
    });

    it("should deploy the CosmosBridge with the correct parameters", async function() {
      this.cosmosBridge.should.exist;

      const claimCount = await this.cosmosBridge.prophecyClaimCount();
      Number(claimCount).should.be.bignumber.equal(0);

      const cosmosBridgeValset = await this.cosmosBridge.valset();
      cosmosBridgeValset.should.be.equal(this.valset.address);
    });

    it("should allow the operator to set the Oracle", async function() {
      this.oracle.should.exist;

      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      }).should.be.fulfilled;

      const bridgeOracle = await this.cosmosBridge.oracle();
      bridgeOracle.should.be.equal(this.oracle.address);
    });

    it("should not allow the operator to update the Oracle once it has been set", async function() {
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      }).should.be.fulfilled;

      await this.cosmosBridge
        .setOracle(this.oracle.address, {
          from: operator
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow the operator to set the Bridge Bank", async function() {
      this.bridgeBank.should.exist;

      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      }).should.be.fulfilled;

      const bridgeBank = await this.cosmosBridge.bridgeBank();
      bridgeBank.should.be.equal(this.bridgeBank.address);
    });

    it("should not allow the operator to update the Bridge Bank once it has been set", async function() {
      await this.cosmosBridge.setBridgeBank(this.oracle.address, {
        from: operator
      }).should.be.fulfilled;

      await this.cosmosBridge
        .setBridgeBank(this.oracle.address, {
          from: operator
        })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  describe("Creation of prophecy claims", function() {
    beforeEach(async function() {
      // Set up ProphecyClaim values
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
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
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
    });

    it("should allow for the creation of new burn prophecy claims", async function() {
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
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

    it("should allow for the creation of new lock prophecy claims", async function() {
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
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

    it("should log an event containing the new prophecy claim's information", async function() {
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewProphecyClaim");

      Number(event.args._prophecyID).should.be.bignumber.equal(1);
      Number(event.args._claimType).should.be.bignumber.equal(CLAIM_TYPE_LOCK);
      event.args._cosmosSender.should.be.equal(this.cosmosSender);
      event.args._ethereumReceiver.should.be.equal(this.ethereumReceiver);
      event.args._validatorAddress.should.be.equal(userOne);
      event.args._tokenAddress.should.be.equal(this.tokenAddress);
      event.args._symbol.should.be.equal(this.symbol);
      Number(event.args._amount).should.be.bignumber.equal(this.amount);
    });

    it("should increase the prophecy claim count upon the creation of new a prophecy claim", async function() {
      const priorProphecyClaimCount = await this.cosmosBridge.prophecyClaimCount();
      Number(priorProphecyClaimCount).should.be.bignumber.equal(0);

      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const postProphecyClaimCount = await this.cosmosBridge.prophecyClaimCount();
      Number(postProphecyClaimCount).should.be.bignumber.equal(1);
    });
  });

  describe("Bridge claim status", function() {
    beforeEach(async function() {
      // Set up ProphecyClaim values
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
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address
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
    });

    it("should not show fake prophecy claims as active", async function() {
      const prophecyClaimCount = 4;

      // Get a prophecy claim's status
      const status = await this.cosmosBridge.isProphecyClaimActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should not be active
      status.should.be.equal(false);
    });

    it("should allow users to check if a prophecy claim is currently active", async function() {
      // Create the prophecy claim
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const prophecyClaimCount = event.args._prophecyID;

      // Get the ProphecyClaim's status
      const status = await this.cosmosBridge.isProphecyClaimActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);
    });

    it("should allow users to check if a prophecy claim's original validator is currently an active validator", async function() {
      // Create the ProphecyClaim
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      );

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const prophecyClaimCount = event.args._prophecyID;

      // Get the ProphecyClaim's status
      const status = await this.cosmosBridge.isProphecyClaimValidatorActive(
        prophecyClaimCount,
        {
          from: accounts[7]
        }
      );

      // Bridge claim should be active
      status.should.be.equal(true);
    });
  });
});
