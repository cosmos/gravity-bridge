const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeToken = artifacts.require("BridgeToken");
const BridgeBank = artifacts.require("BridgeBank");

const { toEthSignedMessageHash, fixSignature } = require("./helpers/helpers");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("BridgeBank", function(accounts) {
  const operator = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("BridgeBank deployment and basics", function() {
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

    it("should deploy the BridgeBank, correctly setting the operator and valset", async function() {
      this.bridgeBank.should.exist;

      const bridgeBankOperator = await this.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator);

      const bridgeBankOracle = await this.bridgeBank.oracle();
      bridgeBankOracle.should.be.equal(this.oracle.address);
    });

    it("should correctly set initial values of CosmosBank and EthereumBank", async function() {
      // EthereumBank initial values
      const nonce = Number(await this.bridgeBank.nonce());
      nonce.should.be.bignumber.equal(0);

      // CosmosBank initial values
      const bridgeTokenCount = Number(await this.bridgeBank.bridgeTokenCount());
      bridgeTokenCount.should.be.bignumber.equal(0);
    });

    it("should not allow a user to send ethereum directly to the contract", async function() {
      await this.bridgeBank
        .send(Web3Utils.toWei("0.25", "ether"), { from: userOne })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  describe("BridgeToken creation (Cosmos assets)", function() {
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
      this.symbol = "ABC";
    });

    it("should not allow non-operators to create new bridge tokens", async function() {
      await this.bridgeBank
        .createNewBridgeToken(this.symbol, {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow the operator to create new bridge token", async function() {
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      }).should.be.fulfilled;
    });

    it("should emit event LogNewBridgeToken containing the new bridge token's address and symbol", async function() {
      //Get the bridge token's address if it were to be created
      const expectedBridgeTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Actually create the bridge token
      const { logs } = await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Get the event logs and compare to expected bridge token address and symbol
      const event = logs.find(e => e.event === "LogNewBridgeToken");
      event.args._token.should.be.equal(expectedBridgeTokenAddress);
      event.args._symbol.should.be.equal(this.symbol);
    });

    it("should increase the bridge token count upon creation", async function() {
      const priorTokenCount = await this.bridgeBank.bridgeTokenCount();
      Number(priorTokenCount).should.be.bignumber.equal(0);

      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      const afterTokenCount = await this.bridgeBank.bridgeTokenCount();
      Number(afterTokenCount).should.be.bignumber.equal(1);
    });

    it("should add the new bridge token to the whitelist", async function() {
      // Get the bridge token's address if it were to be created
      const bridgeTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the bridge token
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Check bridge token whitelist
      const isOnWhitelist = await this.bridgeBank.bridgeTokenWhitelist(
        bridgeTokenAddress
      );
      isOnWhitelist.should.be.equal(true);
    });

    it("should allow the creation of bridge tokens with the same symbol", async function() {
      // Get the first BridgeToken's address if it were to be created
      const firstBridgeTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the first bridge token
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Get the second BridgeToken's address if it were to be created
      const secondBridgeTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the second bridge token
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Check bridge token whitelist for both tokens
      const firstTokenOnWhitelist = await this.bridgeBank.bridgeTokenWhitelist.call(
        firstBridgeTokenAddress
      );
      const secondTokenOnWhitelist = await this.bridgeBank.bridgeTokenWhitelist.call(
        secondBridgeTokenAddress
      );

      // Should be different addresses
      firstBridgeTokenAddress.should.not.be.equal(secondBridgeTokenAddress);

      // Confirm whitelist status
      firstTokenOnWhitelist.should.be.equal(true);
      secondTokenOnWhitelist.should.be.equal(true);
    });
  });

  describe("Bridge token minting (Cosmos assets)", function() {
    beforeEach(async function() {
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
        this.cosmosBridge.address
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );

      // Set up our variables
      this.amount = 100;
      this.sender = web3.utils.bytesToHex([
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      ]);
      this.recipient = userThree;
      this.symbol = "TEST";
      this.nonce = 17;

      // Create the bridge token, adding it to the whitelist
      const { logs: firstLogs } = await this.bridgeBank.createNewBridgeToken(
        this.symbol,
        {
          from: operator
        }
      );

      // Get the event logs and compare to expected bridge token address and symbol
      const eventLogNewBridgeToken = firstLogs.find(
        e => e.event === "LogNewBridgeToken"
      );
      this.bridgeToken = eventLogNewBridgeToken.args._token;

      // Operator sets Oracle
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      });

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });

      // Submit a new bridge claim to the CosmosBridge to make oracle claims upon
      const { logs: secondLogs } = await this.cosmosBridge.newBridgeClaim(
        this.nonce,
        this.sender,
        this.recipient,
        this.bridgeToken,
        this.symbol,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      // Get the new BridgeClaim's id
      const eventLogNewBridgeClaim = secondLogs.find(
        e => e.event === "LogNewBridgeClaim"
      );
      this.bridgeClaimID = eventLogNewBridgeClaim.args._bridgeClaimCount;

      // Create hash using Solidity's Sha3 hashing function
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.bridgeClaimID },
        { t: "bytes", v: this.sender },
        { t: "uint256", v: this.nonce }
      );

      // Generate signatures from active validators userOne, userTwo, userThree
      this.userOneSignature = fixSignature(
        await web3.eth.sign(this.message, userOne)
      );

      // Validator userOne makes a valid OracleClaim
      await this.oracle.newOracleClaim(
        this.bridgeClaimID,
        toEthSignedMessageHash(this.message),
        this.userOneSignature,
        {
          from: userOne
        }
      );
    });

    it("should mint bridge tokens upon the successful processing of a prophecy claim", async function() {
      const bridgeTokenInstance = await BridgeToken.at(this.bridgeToken);

      // Confirm that the user does not hold any bridge tokens of this type
      const priorUserBalance = Number(
        await bridgeTokenInstance.balanceOf(this.recipient)
      );
      priorUserBalance.should.be.bignumber.equal(0);

      // Process the prophecy claim
      await this.oracle.processProphecyClaim(
        this.bridgeClaimID
      ).should.be.fulfilled;

      // Confirm that the user has been minted the correct token
      const afterUserBalance = Number(
        await bridgeTokenInstance.balanceOf(this.recipient)
      );
      afterUserBalance.should.be.bignumber.equal(this.amount);
    });
  });

  describe("Bridge token deposit locking (Ethereum/ERC20 assets)", function() {
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

      this.recipient = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
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

    it("should allow users to lock ERC20 tokens", async function() {
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

    it("should allow users to lock Ethereum", async function() {
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

    it("should generate unique deposit ID for a new deposit", async function() {
      //Simulate sha3 hash to get deposit's expected id
      const expectedID = Web3Utils.soliditySha3(
        { t: "address payable", v: userOne },
        { t: "bytes", v: this.recipient },
        { t: "address", v: this.token.address },
        { t: "int256", v: this.amount },
        { t: "int256", v: 1 }
      );

      //Get the deposit's id if it were to be created
      const depositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      depositID.should.be.equal(expectedID);
    });

    it("should correctly mark new deposits as locked", async function() {
      //Get the deposit's expected id, then lock funds
      const depositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      //Check if a deposit has been created and locked
      const locked = await this.bridgeBank.getEthereumDepositStatus(depositID);
      locked.should.be.equal(true);
    });

    it("should be able to access the deposit's information by its ID", async function() {
      //Get the deposit's expected id, then lock funds
      const depositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      //Attempt to get an deposit's information
      await this.bridgeBank.viewEthereumDeposit(depositID).should.be.fulfilled;
    });

    it("should correctly store deposit information", async function() {
      //Get the deposit's expected id, then lock funds
      const depositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

      //Get the deposit's information
      const depositData = await this.bridgeBank.viewEthereumDeposit(depositID);

      //Parse each attribute
      const sender = depositData[0];
      const receiver = depositData[1];
      const token = depositData[2];
      const amount = Number(depositData[3]);
      const nonce = Number(depositData[4]);

      //Confirm that each attribute is correct
      sender.should.be.equal(userOne);
      receiver.should.be.equal(this.recipient);
      token.should.be.equal(this.token.address);
      amount.should.be.bignumber.equal(this.amount);
      nonce.should.be.bignumber.equal(1);
    });
  });

  describe("Bridge token deposit unlocking (Ethereum/ERC20 assets)", function() {
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

      this.recipient = web3.utils.bytesToHex([
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      ]);
      // This is for Ethereum deposits
      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      // This is for ERC20 deposits
      this.symbol = "TEST";
      this.token = await BridgeToken.new(this.symbol);
      this.amount = 100;

      //Load contract with ethereum so it can complete items
      await this.bridgeBank.send(web3.utils.toWei("1", "ether"), {
        from: operator
      }).should.be.fulfilled;

      //Get the Ethereum deposit's expected id, then lock funds
      this.depositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        {
          from: userOne,
          value: this.weiAmount
        }
      );

      await this.bridgeBank.lock(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        {
          from: userOne,
          value: this.weiAmount
        }
      );

      //Load user account with ERC20 tokens for testing
      await this.token.mint(userOne, 1000, {
        from: operator
      }).should.be.fulfilled;

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, this.amount, {
        from: userOne
      }).should.be.fulfilled;

      //Get the deposit's expected id, then lock funds
      this.erc20DepositID = await this.bridgeBank.lock.call(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne,
          value: 0
        }
      );

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

    it("should allow for an Ethereum deposit to be unlocked", async function() {
      await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;
    });

    it("should allow for an ERC20 deposit to be unlocked", async function() {
      await this.bridgeBank.unlock(this.erc20DepositID).should.be.fulfilled;
    });

    it("should not allow for the unlocking of non-existant deposits", async function() {
      //Generate a fake Ethereum deposit id
      const fakeId = Web3Utils.soliditySha3(
        { t: "address payable", v: userOne },
        { t: "bytes", v: this.recipient },
        { t: "address", v: this.ethereumToken },
        { t: "int256", v: 12 },
        { t: "int256", v: 1 }
      );

      await this.bridgeBank.unlock(fakeId).should.be.rejectedWith(EVMRevert);
    });

    it("should not allow an unlocked deposit to be unlocked again", async function() {
      //Unlock the deposit
      await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

      //Attempt to Unlock the deposit again
      await this.bridgeBank
        .unlock(this.depositID)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should update lock status of deposits upon completion", async function() {
      //Confirm that the deposit is locked
      const firstLockStatus = await this.bridgeBank.getEthereumDepositStatus(
        this.depositID
      );
      firstLockStatus.should.be.equal(true);

      //Unlock the deposit
      await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

      //Check that the deposit is unlocked
      const secondLockStatus = await this.bridgeBank.getEthereumDepositStatus(
        this.depositID
      );
      secondLockStatus.should.be.equal(false);
    });

    it("should emit an event upon unlock with the correct deposit information", async function() {
      //Get the event logs of an unlock
      const { logs } = await this.bridgeBank.unlock(this.erc20DepositID);
      const event = logs.find(e => e.event === "LogUnlock");

      event.args._to.should.be.equal(userOne);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(this.amount);
      Number(event.args._nonce).should.be.bignumber.equal(2);
    });

    // TODO: Original sender VS. intended recipient
    it("should correctly transfer unlocked Ethereum", async function() {
      //Get prior balances of user and BridgeBank contract
      const beforeUserBalance = Number(await web3.eth.getBalance(userOne));
      const beforeContractBalance = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );

      await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

      //Get balances after completion
      const afterUserBalance = Number(await web3.eth.getBalance(userOne));
      const afterContractBalance = Number(
        await web3.eth.getBalance(this.bridgeBank.address)
      );

      //Expected balances
      afterUserBalance.should.be.bignumber.equal(
        beforeUserBalance + Number(this.weiAmount)
      );
      afterContractBalance.should.be.bignumber.equal(
        beforeContractBalance - Number(this.weiAmount)
      );
    });

    it("should correctly transfer unlocked ERC20 tokens", async function() {
      //Confirm that the tokens are locked on the contract
      const beforeBridgeBankBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const beforeUserBalance = Number(await this.token.balanceOf(userOne));

      beforeBridgeBankBalance.should.be.bignumber.equal(this.amount);
      beforeUserBalance.should.be.bignumber.equal(900);

      await this.bridgeBank.unlock(this.erc20DepositID);

      //Confirm that the tokens have been unlocked and transfered
      const afterBridgeBankBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const afterUserBalance = Number(await this.token.balanceOf(userOne));

      afterBridgeBankBalance.should.be.bignumber.equal(0);
      afterUserBalance.should.be.bignumber.equal(1000);
    });
  });
});
