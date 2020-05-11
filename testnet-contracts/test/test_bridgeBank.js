const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeToken = artifacts.require("BridgeToken");
const BridgeBank = artifacts.require("BridgeBank");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("BridgeBank", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold of 70%
  const consensusThreshold = 70;

  describe("BridgeBank deployment and basics", function () {
    beforeEach(async function () {
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

    it("should deploy the BridgeBank, correctly setting the operator and valset", async function () {
      this.bridgeBank.should.exist;

      const bridgeBankOperator = await this.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator);

      const bridgeBankOracle = await this.bridgeBank.oracle();
      bridgeBankOracle.should.be.equal(this.oracle.address);
    });

    it("should correctly set initial values", async function () {
      // EthereumBank initial values
      const bridgeLockBurnNonce = Number(await this.bridgeBank.lockBurnNonce());
      bridgeLockBurnNonce.should.be.bignumber.equal(0);

      // CosmosBank initial values
      const bridgeTokenCount = Number(await this.bridgeBank.bridgeTokenCount());
      bridgeTokenCount.should.be.bignumber.equal(0);
    });

    it("should not allow a user to send ethereum directly to the contract", async function () {
      await this.bridgeBank
        .send(Web3Utils.toWei("0.25", "ether"), { from: userOne })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  describe("Bridge token minting (for burned Cosmos assets)", function () {
    beforeEach(async function () {
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

      // This is for ERC20 deposits
      this.sender = web3.utils.bytesToHex([
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      ]);
      this.recipient = userThree;
      this.symbol = "TEST";
      this.token = await BridgeToken.new(this.symbol);
      this.amount = 100;

      //Load user account with ERC20 tokens for testing
      await this.token.mint(userOne, 1000, {
        from: operator
      }).should.be.fulfilled;

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;

      // Lock tokens on contract
      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      ).should.be.fulfilled;

    });

    it("should mint bridge tokens upon the successful processing of a burn prophecy claim", async function () {
      // Submit a new prophecy claim to the CosmosBridge to make oracle claims upon
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's id
      const event = logs.find(
        e => e.event === "LogNewProphecyClaim"
      );
      const prophecyID = event.args._prophecyID;
      const cosmosSender = event.args._cosmosSender;
      const ethereumReceiver = event.args._ethereumReceiver;
      const amount = event.args._amount;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID },
        { t: "bytes", v: cosmosSender },
        {
          t: "address payable",
          v: ethereumReceiver
        },
        { t: "uint256", v: amount }
      );

      // Generate signatures from active validator userOne
      this.userOneSignature = await web3.eth.sign(this.message, userOne);

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        prophecyID,
        this.message,
        this.userOneSignature,
        {
          from: userOne
        }
      );

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });
  });

  describe("Bridge token deposit locking (Ethereum/ERC20 assets)", function () {
    beforeEach(async function () {
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

      this.recipient = web3.utils.utf8ToHex(
        "cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh"
      );
      // This is for Ethereum deposits
      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      // This is for ERC20 deposits
      this.symbol = "TEST";
      this.token = await BridgeToken.new(this.symbol);
      this.amount = 100;

      //Load user account with ERC20 tokens for testing
      await this.token.mint(userOne, 1000, {
        from: operator
      }).should.be.fulfilled;

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;
    });

    it("should allow users to lock ERC20 tokens", async function () {
      // Attempt to lock tokens
      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      ).should.be.fulfilled;

      //Get the user and BridgeBank token balance after the transfer
      const bridgeBankTokenBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const userBalance = Number(await this.token.balanceOf(userOne));

      //Confirm that the tokens have been locked
      bridgeBankTokenBalance.should.be.bignumber.equal(100);
      userBalance.should.be.bignumber.equal(900);
    });

    it("should allow users to lock Ethereum", async function () {
      await this.bridgeBank.lock(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        { from: userOne, value: this.weiAmount }
      ).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(
        this.bridgeBank.address
      );
      const contractBalance = Web3Utils.fromWei(contractBalanceWei, "ether");

      contractBalance.should.be.bignumber.equal(
        Web3Utils.fromWei(this.weiAmount, "ether")
      );
    });

    it("should increment the token amount in the contract's locked funds mapping", async function () {
      // Confirm locked balances prior to lock
      const priorLockedTokenBalance = await this.bridgeBank.lockedFunds(
        this.token.address
      );
      Number(priorLockedTokenBalance).should.be.bignumber.equal(0);

      // Lock the tokens
      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      // Confirm deposit balances after lock
      const postLockedTokenBalance = await this.bridgeBank.lockedFunds(
        this.token.address
      );
      Number(postLockedTokenBalance).should.be.bignumber.equal(this.amount);
    });
  });

  describe("Ethereum/ERC20 token unlocking (for burned Cosmos assets)", function () {
    beforeEach(async function () {
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

      // Lock an Ethereum deposit
      this.sender = web3.utils.utf8ToHex(
        "cosmos1pjtgu0vau2m52nrykdpztrt887aykue0hq7dfh"
      );
      this.recipient = accounts[4];
      this.ethereumSymbol = "PEGGYETH";
      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.halfWeiAmount = web3.utils.toWei("0.125", "ether");

      //Load contract with ethereum so it can complete items
      await this.bridgeBank.send(web3.utils.toWei("1", "ether"), {
        from: operator
      }).should.be.fulfilled;

      // Lock Ethereum (this is to increase contract's balances and locked funds mapping)
      await this.bridgeBank.lock(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        {
          from: userOne,
          value: this.weiAmount
        }
      );

      // Lock an ERC20 deposit
      this.symbol = "TEST";
      this.token = await BridgeToken.new(this.symbol);
      this.amount = 100;

      //Load user account with ERC20 tokens for testing
      await this.token.mint(userOne, 1000, {
        from: operator
      }).should.be.fulfilled;

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;

      // Lock ERC20 tokens (this is to increase contract's balances and locked funds mapping)
      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );
    });

    it("should unlock Ethereum upon the processing of a burn prophecy", async function () {
      // Submit a new prophecy claim to the CosmosBridge for the Ethereum deposit
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.ethereumSymbol,
        this.weiAmount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's id
      const eventLogNewProphecyClaim = logs.find(
        e => e.event === "LogNewProphecyClaim"
      );
      const prophecyID = eventLogNewProphecyClaim.args._prophecyID;

      // Create hash using Solidity's Sha3 hashing function
      const message = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID },
        { t: "address payable", v: this.recipient },
        { t: "uint256", v: this.weiAmount }
      );

      // Generate signatures from active validator userOne
      const userOneSignature = await web3.eth.sign(message, userOne);

      // Get prior balances of user and BridgeBank contract
      const beforeUserBalance = Number(await web3.eth.getBalance(accounts[4]));
      const beforeContractBalance = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(prophecyID, message, userOneSignature, {
        from: userOne
      });

      // Get balances after prophecy processing
      const afterUserBalance = Number(await web3.eth.getBalance(accounts[4]));
      const afterContractBalance = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );

      // Calculate and check expected balances
      afterUserBalance.should.be.bignumber.equal(
        beforeUserBalance + Number(this.weiAmount)
      );
      afterContractBalance.should.be.bignumber.equal(
        beforeContractBalance - Number(this.weiAmount)
      );
    });

    it("should unlock and transfer ERC20 tokens upon the processing of a burn prophecy", async function () {
      // Submit a new prophecy claim to the CosmosBridge for the Ethereum deposit
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's information
      const eventLogNewProphecyClaim = logs.find(
        e => e.event === "LogNewProphecyClaim"
      );
      const prophecyID = eventLogNewProphecyClaim.args._prophecyID;
      const ethereumReceiver = eventLogNewProphecyClaim.args._ethereumReceiver;
      const amount = Number(eventLogNewProphecyClaim.args._amount);

      // Create hash using Solidity's Sha3 hashing function
      const message = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID },
        { t: "address payable", v: ethereumReceiver },
        { t: "uint256", v: amount }
      );

      // Generate signatures from active validator userOne
      const userOneSignature = await web3.eth.sign(message, userOne);

      // Get Bridge and user's token balance prior to unlocking
      const beforeBridgeBankBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const beforeUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      beforeBridgeBankBalance.should.be.bignumber.equal(this.amount);
      beforeUserBalance.should.be.bignumber.equal(0);

      // Validator userOne makes a valid oracle claim, processing the prophecy claim
      await this.oracle.newOracleClaim(prophecyID, message, userOneSignature, {
        from: userOne
      });

      //Confirm that the tokens have been unlocked and transfered
      const afterBridgeBankBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const afterUserBalance = Number(
        await this.token.balanceOf(this.recipient)
      );
      afterBridgeBankBalance.should.be.bignumber.equal(0);
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });

    it("should allow locked funds to be unlocked incrementally by successive burn prophecies", async function () {
      // -------------------------------------------------------
      // First burn prophecy
      // -------------------------------------------------------
      // Submit a new prophecy claim to the CosmosBridge for the Ethereum deposit
      const { logs: claimLogs1 } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.ethereumSymbol,
        this.halfWeiAmount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's id
      const eventLogNewProphecyClaim1 = claimLogs1.find(
        e => e.event === "LogNewProphecyClaim"
      );

      const prophecyID1 = eventLogNewProphecyClaim1.args._prophecyID;

      // Create hash using Solidity's Sha3 hashing function
      const message1 = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID1 },
        { t: "address payable", v: this.recipient },
        { t: "uint256", v: this.halfWeiAmount }
      );

      // Generate signatures from active validator userOne
      const userOneSignature1 = await web3.eth.sign(message1, userOne);

      // Get pre-claim processed balances of user and BridgeBank contract
      const beforeContractBalance1 = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );
      const beforeUserBalance1 = Number(
        await web3.eth.getBalance(this.recipient)
      );

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        prophecyID1,
        message1,
        userOneSignature1,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get post-claim processed balances of user and BridgeBank contract
      const afterBridgeBankBalance1 = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );
      const afterUserBalance1 = Number(
        await web3.eth.getBalance(this.recipient)
      );

      //Confirm that HALF the amount has been unlocked and transfered
      afterBridgeBankBalance1.should.be.bignumber.equal(
        Number(beforeContractBalance1) - Number(this.halfWeiAmount)
      );
      afterUserBalance1.should.be.bignumber.equal(
        Number(beforeUserBalance1) + Number(this.halfWeiAmount)
      );

      // -------------------------------------------------------
      // Second burn prophecy
      // -------------------------------------------------------
      // Submit a new prophecy claim to the CosmosBridge for the Ethereum deposit
      const { logs: claimLogs2 } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.ethereumSymbol,
        this.halfWeiAmount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's id
      const eventLogNewProphecyClaim2 = claimLogs2.find(
        e => e.event === "LogNewProphecyClaim"
      );

      const prophecyID2 = eventLogNewProphecyClaim2.args._prophecyID;

      // Create hash using Solidity's Sha3 hashing function
      const message2 = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID2 },
        { t: "address payable", v: this.recipient },
        { t: "uint256", v: this.halfWeiAmount }
      );

      // Generate signatures from active validator userOne
      const userOneSignature2 = await web3.eth.sign(message2, userOne);

      // Get pre-claim processed balances of user and BridgeBank contract
      const beforeContractBalance2 = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );
      const beforeUserBalance2 = Number(
        await web3.eth.getBalance(this.recipient)
      );

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        prophecyID2,
        message2,
        userOneSignature2,
        {
          from: userOne
        }
      );

      // Get post-claim processed balances of user and BridgeBank contract
      const afterBridgeBankBalance2 = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );
      const afterUserBalance2 = Number(
        await web3.eth.getBalance(this.recipient)
      );

      //Confirm that HALF the amount has been unlocked and transfered
      afterBridgeBankBalance2.should.be.bignumber.equal(
        Number(beforeContractBalance2) - Number(this.halfWeiAmount)
      );
      afterUserBalance2.should.be.bignumber.equal(
        Number(beforeUserBalance2) + Number(this.halfWeiAmount)
      );

      // Now confirm that the total wei amount has been unlocked and transfered
      afterBridgeBankBalance2.should.be.bignumber.equal(
        Number(beforeContractBalance1) - Number(this.weiAmount)
      );
      afterUserBalance2.should.be.bignumber.equal(
        Number(beforeUserBalance1) + Number(this.weiAmount)
      );
    });

    it("should not allow burn prophecies to be processed twice", async function () {
      // Submit a new prophecy claim to the CosmosBridge for the Ethereum deposit
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new ProphecyClaim's id
      const eventLogNewProphecyClaim = logs.find(
        e => e.event === "LogNewProphecyClaim"
      );
      const prophecyID = eventLogNewProphecyClaim.args._prophecyID;

      // Create hash using Solidity's Sha3 hashing function
      const message = web3.utils.soliditySha3(
        { t: "uint256", v: prophecyID },
        { t: "address payable", v: this.recipient },
        { t: "uint256", v: this.amount }
      );

      // Generate signatures from active validator
      const userOneSignature = await web3.eth.sign(message, userOne);

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(prophecyID, message, userOneSignature, {
        from: userOne
      }).should.be.fulfilled;

      // Attempt to process the same prophecy should be rejected
      await this.oracle
        .processBridgeProphecy(prophecyID)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not accept burn claims for token amounts that exceed the contract's available locked funds", async function () {
      // There are 1,000 TEST tokens approved to the contract, but only 100 have been locked
      const OVERLIMIT_TOKEN_AMOUNT = 500;

      // Attempt to submit a new prophecy claim with overlimit amount is rejected
      await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.sender,
        this.recipient,
        this.symbol,
        OVERLIMIT_TOKEN_AMOUNT,
        {
          from: userOne
        }
      ).should.be.rejectedWith(EVMRevert);
    });
  });
});