const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;
const Web3Utils = require("web3-utils");

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
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold
  const consensusThreshold = 70;

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
        this.cosmosBridge.address,
        consensusThreshold
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
  });

  describe("Creation of prophecy claims", function() {
    beforeEach(async function() {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userOne;
      this.tokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "eth";
      this.amountWei = 100;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [50, 1, 1];
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
    });

    it("should allow for the creation of new lock prophecy claims", async function() {
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amountWei,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    it("should allow for the creation of new burn prophecy claims", async function() {
      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
        this.ethereumReceiver,
        this.tokenAddress,
        this.amountWei,
        {
          from: userOne,
          value: this.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(
        this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      contractBalanceWei.should.be.bignumber.equal(this.amountWei);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.ethereumReceiver,
        this.tokenAddress,
        this.symbol,
        this.amountWei,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });
  });
});
