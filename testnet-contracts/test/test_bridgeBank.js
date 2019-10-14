const BridgeBank = artifacts.require("BridgeBank");
const TestToken = artifacts.require("TestToken");

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

  describe("BridgeBank deployment and basics", function() {
    beforeEach(async function() {
      this.bridgeBank = await BridgeBank.new(operator);
    });

    it("should deploy the BridgeBank and set the operator", async function() {
      this.bridgeBank.should.exist;

      const bridgeBankOperator = await this.bridgeBank.operator();
      bridgeBankOperator.should.be.equal(operator);
    });

    it("should correctly set initial values of CosmosBank and EthereumBank", async function() {
      // EthereumBank initial values
      const nonce = Number(await this.bridgeBank.nonce());
      nonce.should.be.bignumber.equal(0);

      // CosmosBank initial values
      const cosmosTokenCount = Number(await this.bridgeBank.cosmosTokenCount());
      cosmosTokenCount.should.be.bignumber.equal(0);
    });

    it("should not allow a user to send ethereum directly to the contract", async function() {
      await this.bridgeBank
        .send(Web3Utils.toWei("0.25", "ether"), { from: userOne })
        .should.be.rejectedWith(EVMRevert);
    });
  });

  describe("BridgeToken creation", function() {
    beforeEach(async function() {
      this.bridgeBank = await BridgeBank.new(operator);
      this.symbol = "ABC";
    });

    it("should not allow non-operators to create new BankTokens", async function() {
      await this.bridgeBank
        .createNewBridgeToken(this.symbol, {
          from: userOne
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow the operator to create new BankTokens", async function() {
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      }).should.be.fulfilled;
    });

    it("should emit event LogNewBankToken containing the new BankToken's address and symbol", async function() {
      //Get the BankToken's address if it were to be created
      const expectedBankTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Actually create the BankToken
      const { logs } = await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Get the event logs and compare to expected BankToken address and symbol
      const event = logs.find(e => e.event === "LogNewBankToken");
      event.args._token.should.be.equal(expectedBankTokenAddress);
      event.args._symbol.should.be.equal(this.symbol);
    });

    it("should increase the BankToken count upon creation", async function() {
      const priorTokenCount = await this.bridgeBank.cosmosTokenCount();
      Number(priorTokenCount).should.be.bignumber.equal(0);

      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      const afterTokenCount = await this.bridgeBank.cosmosTokenCount();
      Number(afterTokenCount).should.be.bignumber.equal(1);
    });

    it("should add new BankTokens to the whitelist", async function() {
      // Get the BankToken's address if it were to be created
      const bankTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the BridgeToken
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Check BankToken whitelist
      const isOnWhitelist = await this.bridgeBank.bankTokenWhitelist(
        bankTokenAddress
      );
      isOnWhitelist.should.be.equal(true);
    });

    it("should allow the creation of BankTokens with the same symbol", async function() {
      // Get the first BankToken's address if it were to be created
      const firstBankTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the first BankToken
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Get the second BankToken's address if it were to be created
      const secondBankTokenAddress = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the second BankToken
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      });

      // Check BankToken whitelist for both tokens
      const firstTokenOnWhitelist = await this.bridgeBank.bankTokenWhitelist.call(
        firstBankTokenAddress
      );
      const secondTokenOnWhitelist = await this.bridgeBank.bankTokenWhitelist.call(
        secondBankTokenAddress
      );

      // Should be different addresses
      firstBankTokenAddress.should.not.be.equal(secondBankTokenAddress);

      // Confirm whitelist status
      firstTokenOnWhitelist.should.be.equal(true);
      secondTokenOnWhitelist.should.be.equal(true);
    });
  });

  describe("BankToken minting", function() {
    beforeEach(async function() {
      this.bridgeBank = await BridgeBank.new(operator);

      // Set up our variables
      this.amount = 100;
      this.sender = web3.utils.bytesToHex([
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      ]);
      this.recipient = userOne;
      this.symbol = "ETH";
      this.bankToken = await this.bridgeBank.createNewBridgeToken.call(
        this.symbol,
        {
          from: operator
        }
      );

      // Create the BankToken, adding it to the whitelist
      await this.bridgeBank.createNewBridgeToken(this.symbol, {
        from: operator
      }).should.be.fulfilled;
    });

    // TODO: should be VALIDATORS
    it("should allow the operator to mint new BankTokens", async function() {
      await this.bridgeBank.mintBankTokens(
        this.sender,
        this.recipient,
        this.bankToken,
        this.symbol,
        this.amount,
        {
          from: operator
        }
      ).should.be.fulfilled;
    });

    it("should emit event LogBankTokenMint with correct values upon successful minting", async function() {
      const { logs } = await this.bridgeBank.mintBankTokens(
        this.sender,
        this.recipient,
        this.bankToken,
        this.symbol,
        this.amount,
        {
          from: operator
        }
      );

      const event = logs.find(e => e.event === "LogBankTokenMint");
      event.args._token.should.be.equal(this.bankToken);
      event.args._symbol.should.be.equal(this.symbol);
      Number(event.args._amount).should.be.bignumber.equal(this.amount);
      event.args._beneficiary.should.be.equal(this.recipient);
    });
  });

  describe("Ethereum/ERC20 token locking", function() {
    beforeEach(async function() {
      this.bridgeBank = await BridgeBank.new(operator);
      this.recipient = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.amount = 250;
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      // TODO: Need to set symbol correctly (use bank token?)
      this.symbol = "";
      this.ethereumToken = "0x0000000000000000000000000000000000000000";

      this.token = await TestToken.new();
      // this.gasForLock = 300000; // 300,000 Gwei

      //Load user account with ERC20 tokens for testing
      await this.token.mint(userOne, 1000, {
        from: operator
      }).should.be.fulfilled;
    });

    it("Testing ERC20 token minting", async function() {
      const bridgeBankTokenBalance = Number(
        await this.token.balanceOf(this.bridgeBank.address)
      );
      const userBalance = Number(await this.token.balanceOf(userOne));

      console.log(bridgeBankTokenBalance);
      console.log(userBalance);

      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, 100, {
        from: userOne
      }).should.be.fulfilled;

      const bridgeBankTokenAllowance = Number(
        await this.token.allowance(userOne, this.bridgeBank.address)
      );

      console.log(bridgeBankTokenAllowance);

      await this.token.transferFrom(userOne, this.bridgeBank.address, 33, {
        from: userOne
      }).should.be.fulfilled;

      // console.log(bridgeBankTokenBalance);
      // console.log(userBalance);
    });

    it("should allow users to lock ERC20 tokens", async function() {
      // Approve tokens to contract
      await this.token.approve(this.bridgeBank.address, 100, {
        from: userOne
      }).should.be.fulfilled;

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

    it("should allow users to lock ERC20 assets", async function() {
      await this.bridgeBank.lock(
        this.recipient,
        this.token.address,
        this.amount,
        {
          from: userOne
        }
      ).should.be.fulfilled;
    });

    // it("should generate unique deposit id's for a created deposit", async function() {
    //   //Simulate sha3 hash to get deposit's expected id
    //   const expectedID = Web3Utils.soliditySha3(
    //     { t: "address payable", v: userOne },
    //     { t: "bytes", v: this.recipient },
    //     { t: "address", v: this.token.address },
    //     { t: "int256", v: this.amount },
    //     { t: "int256", v: 1 }
    //   );

    //   //Get the deposit's id if it were to be created
    //   const depositID = await this.bridgeBank.lock.call(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   depositID.should.be.equal(expectedID);
    // });

    // it("should correctly mark deposits as locked", async function() {
    //   //Get the deposit's expected id, then lock funds
    //   const depositID = await this.bridgeBank.lock.call(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   await this.bridgeBank.lock(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   //Check if a deposit has been created and locked
    //   const locked = await this.bridgeBank.getEthereumDepositStatus(depositID);
    //   locked.should.be.equal(true);
    // });

    // it("should be able to access the deposit's information by its ID", async function() {
    //   const depositID = await this.bridgeBank.lock.call(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   await this.bridgeBank.lock(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   //Attempt to get an deposit's information
    //   await this.bridgeBank.viewEthereumDeposit(depositID).should.be.fulfilled;
    // });

    // it("should correctly store deposit information", async function() {
    //   //Get the deposit's expected id, then lock funds
    //   const depositID = await this.bridgeBank.lock.call(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   await this.bridgeBank.lock(
    //     userOne,
    //     this.recipient,
    //     this.token.address,
    //     this.amount
    //   );

    //   //Get the deposit's information
    //   const depositData = await this.bridgeBank.viewEthereumDeposit(depositID);

    //   //Parse each attribute
    //   const sender = depositData[0];
    //   const receiver = depositData[1];
    //   const token = depositData[2];
    //   const amount = Number(depositData[3]);
    //   const nonce = Number(depositData[4]);

    //   //Confirm that each attribute is correct
    //   sender.should.be.equal(userOne);
    //   receiver.should.be.equal(this.recipient);
    //   token.should.be.equal(this.token.address);
    //   amount.should.be.bignumber.equal(this.amount);
    //   nonce.should.be.bignumber.equal(1);
    // });
  });

  // describe("Unlocking of Ethereum deposits", function() {
  //   beforeEach(async function() {
  //     this.bridgeBank = await BridgeBank.new(operator);
  //     this.recipient = web3.utils.bytesToHex([
  //       "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
  //     ]);
  //     this.ethereumToken = "0x0000000000000000000000000000000000000000";
  //     this.symbol = "eth";
  //     this.weiAmount = web3.utils.toWei("0.25", "ether");

  //     this.depositID = await this.bridgeBank.lock.call(
  //       userOne,
  //       this.recipient,
  //       this.ethereumToken,
  //       this.weiAmount
  //     );

  //     await this.bridgeBank.lock(
  //       userOne,
  //       this.recipient,
  //       this.ethereumToken,
  //       this.weiAmount
  //     );

  //     //Load contract with ethereum so it can complete items
  //     await this.bridgeBank.send(web3.utils.toWei("1", "ether"), {
  //       from: operator
  //     }).should.be.fulfilled;
  //   });

  //   // it("should not allow for the completion of items whose value exceeds the contract's balance", async function() {
  //   //   //Create an deposit with an overlimit amount
  //   //   const overlimitAmount = web3.utils.toWei("1.25", "ether");
  //   //   const id = await this.ethereumBank.callNewEthereumDeposit.call(
  //   //     userOne,
  //   //     this.recipient,
  //   //     this.ethereumToken,
  //   //     overlimitAmount
  //   //   );
  //   //   await this.ethereumBank.callNewEthereumDeposit(
  //   //     userOne,
  //   //     this.recipient,
  //   //     this.ethereumToken,
  //   //     overlimitAmount
  //   //   );

  //   //   //Attempt to complete the deposit
  //   //   await this.ethereumBank
  //   //     .callUnlockEthereumDeposit(id)
  //   //     .should.be.rejectedWith(EVMRevert);
  //   // });

  //   it("should not allow for the unlocking of non-existant Ethereum deposits", async function() {
  //     //Generate a fake Ethereum deposit id
  //     const fakeId = Web3Utils.soliditySha3(
  //       { t: "address payable", v: userOne },
  //       { t: "bytes", v: this.recipient },
  //       { t: "address", v: this.ethereumToken },
  //       { t: "int256", v: 12 },
  //       { t: "int256", v: 1 }
  //     );

  //     await this.bridgeBank.unlock(fakeId).should.be.rejectedWith(EVMRevert);
  //   });

  //   it("should not allow an Ethereum deposit that has already been unlocked to be unlocked", async function() {
  //     //Unlock the deposit
  //     await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

  //     //Attempt to Unlock the deposit again
  //     await this.bridgeBank
  //       .unlock(this.depositID)
  //       .should.be.rejectedWith(EVMRevert);
  //   });

  //   it("should allow for an Ethereum deposit to be unlocked", async function() {
  //     await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;
  //   });

  //   it("should update lock status of Ethereum deposits upon completion", async function() {
  //     //Confirm that the deposit is locked
  //     const firstLockStatus = await this.bridgeBank.getEthereumDepositStatus(
  //       this.depositID
  //     );
  //     firstLockStatus.should.be.equal(true);

  //     //Unlock the deposit
  //     await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

  //     //Check that the deposit is unlocked
  //     const secondLockStatus = await this.bridgeBank.getEthereumDepositStatus(
  //       this.depositID
  //     );
  //     secondLockStatus.should.be.equal(false);
  //   });

  //   // TODO: Original sender VS. intended recipient
  //   it("should correctly transfer unlocked funds to the original sender", async function() {
  //     //Get prior balances of user and BridgeBank contract
  //     const beforeUserBalance = Number(await web3.eth.getBalance(userOne));
  //     const beforeContractBalance = Number(
  //       await web3.eth.getBalance(this.bridgeBank.address)
  //     );

  //     await this.bridgeBank.unlock(this.depositID).should.be.fulfilled;

  //     //Get balances after completion
  //     const afterUserBalance = Number(await web3.eth.getBalance(userOne));
  //     const afterContractBalance = Number(
  //       await web3.eth.getBalance(this.bridgeBank.address)
  //     );

  //     //Expected balances
  //     afterUserBalance.should.be.bignumber.equal(
  //       beforeUserBalance + Number(this.weiAmount)
  //     );
  //     afterContractBalance.should.be.bignumber.equal(
  //       beforeContractBalance - Number(this.weiAmount)
  //     );
  //   });
  // });
});
