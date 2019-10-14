const Peggy = artifacts.require("Peggy");
const TestToken = artifacts.require("TestToken");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Peggy", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Peggy smart contract deployment", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.peggy = await Peggy.new(this.initialValidators, this.initialPowers);
    });

    it("should deploy the peggy contract with the correct parameters", async function() {
      this.peggy.should.exist;

      const peggyProvider = await this.peggy.provider();
      peggyProvider.should.be.equal(provider);
    });
  });

  describe("Locking of Ethereum assets", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.peggy = await Peggy.new(this.initialValidators, this.initialPowers);

      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.gasForLock = 300000; // 300,000 Gwei

      //Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 1000, {
        from: provider
      }).should.be.fulfilled;
    });

    // it("should not allow a user to send ethereum directly to the contract", async function() {
    //   await this.peggy
    //     .send(this.weiAmount, { from: userOne })
    //     .should.be.rejectedWith(EVMRevert);
    // });

    it("should not allow users to lock funds if the contract is paused", async function() {
      //Confirm that the processor contract is paused
      await this.peggy.pauseLocking({ from: provider }).should.be.fulfilled;

      const depositStatus = await this.peggy.active();
      depositStatus.should.be.equal(false);

      //User attempt to lock ethereum/erc20
      await this.token.approve(this.peggy.address, 100, {
        from: userOne
      }).should.be.fulfilled;
      await this.peggy
        .lock(this.recipient, this.token.address, 100, {
          from: userOne,
          gas: this.gasForLock
        })
        .should.be.rejectedWith(EVMRevert);
    });

    // it("should allow users to lock erc20 tokens if it meets validation requirements", async function() {
    //   //Confirm that the contract is active
    //   const depositStatus = await this.peggy.active();
    //   depositStatus.should.be.equal(true);

    //   await this.token.approve(this.peggy.address, 100, {
    //     from: userOne
    //   }).should.be.fulfilled;
    //   await this.peggy.lock(this.recipient, this.token.address, 100, {
    //     from: userOne,
    //     gas: this.gasForLock
    //   }).should.be.fulfilled;

    //   //Get the contract and user token balance after the rescue
    //   const peggyBalance = Number(
    //     await this.token.balanceOf(this.peggy.address)
    //   );
    //   const userBalance = Number(await this.token.balanceOf(userOne));

    //   //Confirm that the tokens have been locked
    //   peggyBalance.should.be.bignumber.equal(100);
    //   userBalance.should.be.bignumber.equal(900);
    // });

    it("should allow users to lock ethereum if it meets validation requirements", async function() {
      //Confirm that the contract is active
      const depositStatus = await this.peggy.active();
      depositStatus.should.be.equal(true);

      await this.peggy.lock(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        { from: userOne, value: this.weiAmount, gas: this.gasForLock }
      ).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(this.peggy.address);
      const contractBalance = web3.utils.fromWei(contractBalanceWei, "ether");

      contractBalance.should.be.bignumber.equal(
        web3.utils.fromWei(this.weiAmount, "ether")
      );
    });

    it("should emit an event upon lock containing the new ecrow's information", async function() {
      const userBalance = Number(await this.token.balanceOf(userOne));
      userBalance.should.be.bignumber.equal(1000);

      await this.token.approve(this.peggy.address, 100, {
        from: userOne
      }).should.be.fulfilled;

      //Get the event logs of a token deposit
      const expectedId = await this.peggy.lock.call(
        this.recipient,
        this.token.address,
        100,
        { from: userOne, gas: this.gasForLock }
      ).should.be.fulfilled;
      const { logs } = await this.peggy.lock(
        this.recipient,
        this.token.address,
        100,
        { from: userOne, gas: this.gasForLock }
      ).should.be.fulfilled;
      const event = logs.find(e => e.event === "LogLock");

      event.args._id.should.be.equal(expectedId);
      event.args._to.should.be.equal(this.recipient);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(100);
      Number(event.args._nonce).should.be.bignumber.equal(1);
    });
  });

  describe("Access to information", function() {
    const cosmosAddr = "77m5cfkop78sruko3ud4wjp83kuc9rmw15rqtzlp";

    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.peggy = await Peggy.new(this.initialValidators, this.initialPowers);

      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex(cosmosAddr);
      this.gasForLock = 300000; // 300,000 Gwei

      // Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 100, {
        from: provider
      }).should.be.fulfilled;

      await this.token.approve(this.peggy.address, 100, {
        from: userOne
      }).should.be.fulfilled;
      this.depositId = await this.peggy.lock.call(
        this.recipient,
        this.token.address,
        100,
        { from: userOne, gas: this.gasForLock }
      ).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, {
        from: userOne,
        gas: this.gasForLock
      }).should.be.fulfilled;
    });

    it("should allow for public viewing of a locked ethereum deposit's information", async function() {
      //Get the ethereum deposit struct's information
      const depositInfo = await this.peggy.viewEthereumDeposit(this.depositId, {
        from: provider
      }).should.be.fulfilled;

      //Parse each attribute
      const sender = depositInfo[0];
      const receiver = depositInfo[1];
      const token = depositInfo[2];
      const amount = Number(depositInfo[3]);
      const nonce = Number(depositInfo[4]);

      //Confirm that each attribute is correct
      sender.should.be.equal(userOne);
      receiver.should.be.equal(this.recipient);
      token.should.be.equal(this.token.address);
      amount.should.be.bignumber.equal(100);
      nonce.should.be.bignumber.equal(1);
    });

    it("should correctly encode and decode the intended recipient's address", async function() {
      //Get the ethereum deposit struct's information
      const depositInfo = await this.peggy.viewEthereumDeposit(this.depositId, {
        from: provider
      }).should.be.fulfilled;

      //Decode the stored recipient's address and compare it the original
      const receiver = web3.utils.hexToUtf8(depositInfo[1]);
      receiver.should.be.equal(cosmosAddr);
    });
  });

  describe("Unlocking of itemized ethereum", function() {
    beforeEach(async function() {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.peggy = await Peggy.new(this.initialValidators, this.initialPowers);

      this.ethereumToken = "0x0000000000000000000000000000000000000000";
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex(
        "cosmosaccaddr985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.gasForLock = 300000; // 300,000 Gwei

      // Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 100, {
        from: provider
      }).should.be.fulfilled;

      await this.token.approve(this.peggy.address, 100, {
        from: userOne
      }).should.be.fulfilled;
      this.depositId = await this.peggy.lock.call(
        this.recipient,
        this.token.address,
        100,
        { from: userOne, gas: this.gasForLock }
      ).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, {
        from: userOne,
        gas: this.gasForLock
      }).should.be.fulfilled;
    });

    it("should allow the provider to unlock itemized ethereum", async function() {
      const id = await this.peggy.lock.call(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        { from: userOne, value: this.weiAmount, gas: this.gasForLock }
      ).should.be.fulfilled;
      await this.peggy.lock(
        this.recipient,
        this.ethereumToken,
        this.weiAmount,
        { from: userOne, value: this.weiAmount, gas: this.gasForLock }
      ).should.be.fulfilled;
      await this.peggy.unlock(id, {
        from: provider,
        gas: this.gasForLock
      }).should.be.fulfilled;
    });

    it("should allow the provider to unlock itemized erc20 tokens", async function() {
      await this.peggy.unlock(this.depositId, {
        from: provider,
        gas: this.gasForLock
      }).should.be.fulfilled;
    });

    it("should correctly transfer funds to intended recipient upon unlock", async function() {
      //Confirm that the tokens are locked on the contract
      const beforePeggyBalance = Number(
        await this.token.balanceOf(this.peggy.address)
      );
      const beforeUserBalance = Number(await this.token.balanceOf(userOne));

      beforePeggyBalance.should.be.bignumber.equal(100);
      beforeUserBalance.should.be.bignumber.equal(0);

      await this.peggy.unlock(this.depositId, {
        from: provider,
        gas: this.gasForLock
      });

      //Confirm that the tokens have been unlocked and transfered
      const afterPeggyBalance = Number(
        await this.token.balanceOf(this.peggy.address)
      );
      const afterUserBalance = Number(await this.token.balanceOf(userOne));

      afterPeggyBalance.should.be.bignumber.equal(0);
      afterUserBalance.should.be.bignumber.equal(100);
    });

    it("should emit an event upon unlock containing the ecrow's recipient, token, amount, and nonce", async function() {
      //Get the event logs of an unlock
      const { logs } = await this.peggy.unlock(this.depositId, {
        from: provider,
        gas: this.gasForLock
      });
      const event = logs.find(e => e.event === "LogUnlock");

      event.args._to.should.be.equal(userOne);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(100);
      Number(event.args._nonce).should.be.bignumber.equal(1);
    });

    it("should update deposit lock status upon unlock", async function() {
      const startingLockStatus = await this.peggy.getEthereumDepositStatus(
        this.depositId
      );
      startingLockStatus.should.be.equal(true);

      await this.peggy.unlock(this.depositId, {
        from: provider,
        gas: this.gasForLock
      });

      const endingLockStatus = await this.peggy.getEthereumDepositStatus(
        this.depositId
      );
      endingLockStatus.should.be.equal(false);
    });
  });
});
