const EthereumBank = artifacts.require("TestEthereumBank");
const TestToken = artifacts.require("TestToken");

const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("EthereumBank", function(accounts) {
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("EthereumBank contract deployment", function() {
    beforeEach(async function() {
      this.ethereumBank = await EthereumBank.new();
    });

    it("should deploy the EthereumBank with the correct parameters", async function() {
      this.ethereumBank.should.exist;

      const nonce = Number(await this.ethereumBank.nonce());
      nonce.should.be.bignumber.equal(0);
    });
  });

  describe("Locking of Ethereum deposits", function() {
    beforeEach(async function() {
      this.ethereumBank = await EthereumBank.new();
      this.recipient = web3.utils.bytesToHex(["20bytestring"]);
      this.amount = 250;

      //Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 1000, {
        from: accounts[0]
      }).should.be.fulfilled;
    });

    it("should allow for the creation of Ethereum deposits", async function() {
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      ).should.be.fulfilled;
    });

    it("should generate unique deposit id's for a created deposit", async function() {
      //Simulate sha3 hash to get deposit's expected id
      const expectedId = Web3Utils.soliditySha3(
        { t: "address payable", v: userOne },
        { t: "bytes", v: this.recipient },
        { t: "address", v: this.token.address },
        { t: "int256", v: this.amount },
        { t: "int256", v: 1 }
      );

      //Get the deposit's id if it were to be created
      const id = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );
      id.should.be.equal(expectedId);
    });

    it("should allow access to an Ethereum depoit's information given it's unique id", async function() {
      const id = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );

      //Attempt to get an deposit's information
      await this.ethereumBank.callGetEthereumDeposit(id).should.be.fulfilled;
    });

    it("should correctly identify the existence of items in memory", async function() {
      //Get the deposit's expected id then lock funds
      const id = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      ).should.be.fulfilled;

      //Check if a deposit has been created and locked
      const locked = await this.ethereumBank.callIsLockedEthereumDeposit(id);
      locked.should.be.equal(true);
    });

    it("should store items with the correct parameters", async function() {
      //Create the deposit and store its id
      const id = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.token.address,
        this.amount
      );

      //Get the deposit's information
      const depositData = await this.ethereumBank.callGetEthereumDeposit(id);

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

  describe("Unlocking of Ethereum deposits", function() {
    beforeEach(async function() {
      this.ethereumBank = await EthereumBank.new();
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.bytesToHex(["20bytestring"]);
      this.ethereumToken = "0x0000000000000000000000000000000000000000";

      //Load contract with ethereum so it can complete items
      await this.ethereumBank.send(web3.utils.toWei("1", "ether"), {
        from: accounts[0]
      }).should.be.fulfilled;

      this.ethereumDepositID = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.ethereumToken,
        this.weiAmount
      );
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.ethereumToken,
        this.weiAmount
      );
    });

    it("should not allow for the completion of items whose value exceeds the contract's balance", async function() {
      //Create an deposit with an overlimit amount
      const overlimitAmount = web3.utils.toWei("1.25", "ether");
      const id = await this.ethereumBank.callNewEthereumDeposit.call(
        userOne,
        this.recipient,
        this.ethereumToken,
        overlimitAmount
      );
      await this.ethereumBank.callNewEthereumDeposit(
        userOne,
        this.recipient,
        this.ethereumToken,
        overlimitAmount
      );

      //Attempt to complete the deposit
      await this.ethereumBank
        .callUnlockEthereumDeposit(id)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow for the unlocking of non-existant Ethereum deposits", async function() {
      //Generate a fake Ethereum deposit id
      const fakeId = Web3Utils.soliditySha3(
        { t: "address payable", v: userOne },
        { t: "bytes", v: this.recipient },
        { t: "address", v: this.ethereumToken },
        { t: "int256", v: 12 },
        { t: "int256", v: 1 }
      );

      await this.ethereumBank
        .callComplete(fakeId)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should not allow an Ethereum deposit that has already been unlocked to be unlocked", async function() {
      //Complete the deposit
      await this.ethereumBank.callUnlockEthereumDeposit(
        this.ethereumDepositID
      ).should.be.fulfilled;

      //Attempt to complete the deposit again
      await this.ethereumBank
        .callComplete(this.ethereumDepositID)
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow for an Ethereum deposit to be unlocked", async function() {
      await this.ethereumBank.callUnlockEthereumDeposit(
        this.ethereumDepositID
      ).should.be.fulfilled;
    });

    it("should update lock status of Ethereum deposits upon completion", async function() {
      //Confirm that the deposit is active
      const startingLockStatus = await this.ethereumBank.callIsLockedEthereumDeposit(
        this.ethereumDepositID
      );
      startingLockStatus.should.be.equal(true);

      //Complete the deposit
      await this.ethereumBank.callUnlockEthereumDeposit(
        this.ethereumDepositID
      ).should.be.fulfilled;

      //Check if the deposit still exists
      const completedDeposit = await this.ethereumBank.callIsLockedEthereumDeposit(
        this.ethereumDepositID
      );
      completedDeposit.should.be.equal(false);
    });

    it("should correctly transfer locked funds to the original sender", async function() {
      //Get prior balances of user and peggy contract
      const beforeUserBalance = Number(await web3.eth.getBalance(userOne));
      const beforeContractBalance = Number(
        await web3.eth.getBalance(this.ethereumBank.address)
      );

      await this.ethereumBank.callUnlockEthereumDeposit(
        this.ethereumDepositID
      ).should.be.fulfilled;

      //Get balances after completion
      const afterUserBalance = Number(await web3.eth.getBalance(userOne));
      const afterContractBalance = Number(
        await web3.eth.getBalance(this.ethereumBank.address)
      );

      //Expected balances
      afterUserBalance.should.be.bignumber.equal(
        beforeUserBalance + Number(this.weiAmount)
      );
      afterContractBalance.should.be.bignumber.equal(
        beforeContractBalance - Number(this.weiAmount)
      );
    });
  });
});
