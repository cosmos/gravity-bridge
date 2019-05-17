const Peggy = artifacts.require('Peggy');
const TestToken = artifacts.require('TestToken');

const Web3Utils = require('web3-utils');
const EVMRevert = 'revert';
const BigNumber = web3.BigNumber;

require('chai')
  .use(require('chai-as-promised'))
  .use(require('chai-bignumber')(BigNumber))
  .should();

contract('Peggy', function (accounts) {

  const relayer = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe('Peggy smart contract deployment', function(){

    beforeEach(async function() {
      this.peggy = await Peggy.new();
    });

    it('should deploy the peggy contract with the correct parameters', async function () {
      this.peggy.should.exist;

      const peggyRelayer = (await this.peggy.relayer());
      peggyRelayer.should.be.equal(relayer);
    });

  });

  describe('Locking funds', function(){

    beforeEach(async function() {
      this.peggy = await Peggy.new();
      this.ethereumToken = '0x0000000000000000000000000000000000000000';
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex('985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy');
      this.gasForLock = 300000; // 300,000 Gwei

      //Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 1000, { from: relayer }).should.be.fulfilled;
    });

    it('should not allow a user to send ethereum directly to the contract', async function () {
      await this.peggy.send(this.weiAmount, { from: userOne}).should.be.rejectedWith(EVMRevert);
    });

    it('should not allow users to lock funds if the contract is paused', async function () {
      //Confirm that the processor contract is paused
      await this.peggy.pauseLocking({ from:relayer }).should.be.fulfilled;
      
      const depositStatus = await this.peggy.active();
      depositStatus.should.be.equal(false);

      //User attempt to lock ethereum/erc20
      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.rejectedWith(EVMRevert);
    });

    it('should allow users to lock erc20 tokens if it meets validation requirements', async function () {
      //Confirm that the contract is active
      const depositStatus = await this.peggy.active();
      depositStatus.should.be.equal(true);

      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;

      //Get the contract and user token balance after the rescue
      const peggyBalance = Number(await this.token.balanceOf(this.peggy.address));
      const userBalance = Number(await this.token.balanceOf(userOne));

      //Confirm that the tokens have been locked
      peggyBalance.should.be.bignumber.equal(100);
      userBalance.should.be.bignumber.equal(900);
    });

    it('should allow users to lock ethereum if it meets validation requirements', async function () {
      //Confirm that the contract is active
      const depositStatus = await this.peggy.active();
      depositStatus.should.be.equal(true);

      await this.peggy.lock(this.recipient, this.ethereumToken, this.weiAmount, { from: userOne, value: this.weiAmount, gas: this.gasForLock }).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(this.peggy.address);
      const contractBalance = web3.utils.fromWei(contractBalanceWei, "ether");

      contractBalance.should.be.bignumber.equal(web3.utils.fromWei(this.weiAmount, "ether"));

    });

    it('should emit an event upon lock containing the new ecrow\'s information', async function () {
      const userBalance = Number(await this.token.balanceOf(userOne));
      userBalance.should.be.bignumber.equal(1000);

      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      
      //Get the event logs of a token deposit
      const expectedId = await this.peggy.lock.call(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      const {logs} = await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      const event = logs.find(e => e.event === 'LogLock');

      event.args._id.should.be.equal(expectedId);
      event.args._to.should.be.equal(this.recipient);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(100);
      Number(event.args._nonce).should.be.bignumber.equal(1);
    });

  });

  describe('Access to information', function(){

    const cosmosAddr = '77m5cfkop78sruko3ud4wjp83kuc9rmw15rqtzlp';

    beforeEach(async function() {
      this.peggy = await Peggy.new();
      this.ethereumToken = '0x0000000000000000000000000000000000000000';
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex(cosmosAddr);
      this.gasForLock = 300000; // 300,000 Gwei

      // Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 100, { from: relayer }).should.be.fulfilled;

      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      this.itemId = await this.peggy.lock.call(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
    });

    it('should allow for public viewing of a locked item\'s information', async function () {
      //Get the item struct's information
      const itemInfo = await this.peggy.viewItem(this.itemId, { from: relayer }).should.be.fulfilled;

      //Parse each attribute
      const sender = itemInfo[0];
      const receiver = itemInfo[1];
      const token = itemInfo[2];
      const amount = Number(itemInfo[3]);
      const nonce = Number(itemInfo[4]);

      //Confirm that each attribute is correct
      sender.should.be.equal(userOne);
      receiver.should.be.equal(this.recipient);
      token.should.be.equal(this.token.address);
      amount.should.be.bignumber.equal(100);
      nonce.should.be.bignumber.equal(1);
    });

    it('should correctly encode and decode the intended recipient\'s address', async function () {
      //Get the item struct's information
      const itemInfo = await this.peggy.viewItem(this.itemId, { from: relayer }).should.be.fulfilled;

      //Decode the stored recipient's address and compare it the original
      const receiver = web3.utils.hexToUtf8(itemInfo[1]);
      receiver.should.be.equal(cosmosAddr);
    });

  });

  describe('Unlocking of itemized ethereum', function(){

    beforeEach(async function() {
      this.peggy = await Peggy.new();
      this.ethereumToken = '0x0000000000000000000000000000000000000000';
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex('cosmosaccaddr985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy');
      this.gasForLock = 300000; // 300,000 Gwei

      // Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 100, { from: relayer }).should.be.fulfilled;

      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      this.itemId = await this.peggy.lock.call(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;

  });

    it('should allow the relayer to unlock itemized ethereum', async function () {
      const id = await this.peggy.lock.call(this.recipient, this.ethereumToken, this.weiAmount, { from: userOne, value: this.weiAmount, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.ethereumToken, this.weiAmount, { from: userOne, value: this.weiAmount, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.unlock(id, { from: relayer, gas: this.gasForLock }).should.be.fulfilled;
    });

    it('should allow the relayer to unlock itemized erc20 tokens', async function () {
      await this.peggy.unlock(this.itemId, { from: relayer, gas: this.gasForLock }).should.be.fulfilled;
    });

    it('should correctly transfer funds to intended recipient upon unlock', async function () {
      //Confirm that the tokens are locked on the contract
      const beforePeggyBalance = Number(await this.token.balanceOf(this.peggy.address));
      const beforeUserBalance = Number(await this.token.balanceOf(userOne));

      beforePeggyBalance.should.be.bignumber.equal(100);
      beforeUserBalance.should.be.bignumber.equal(0);

      await this.peggy.unlock(this.itemId, { from: relayer, gas: this.gasForLock });
      
      //Confirm that the tokens have been unlocked and transfered
      const afterPeggyBalance = Number(await this.token.balanceOf(this.peggy.address));
      const afterUserBalance = Number(await this.token.balanceOf(userOne));

      afterPeggyBalance.should.be.bignumber.equal(0);
      afterUserBalance.should.be.bignumber.equal(100);
    });

    it('should emit an event upon unlock containing the ecrow\'s recipient, token, amount, and nonce', async function () {
      //Get the event logs of an unlock
      const {logs} = await this.peggy.unlock(this.itemId, { from: relayer, gas: this.gasForLock });
      const event = logs.find(e => e.event === 'LogUnlock');

      event.args._to.should.be.equal(userOne);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(100);
      Number(event.args._nonce).should.be.bignumber.equal(1);
    });

    it('should update item lock statusit has been unlocked', async function () {
      const startingLockStatus = await this.peggy.getStatus(this.itemId);
      startingLockStatus.should.be.equal(true);
      
      await this.peggy.unlock(this.itemId, { from: relayer, gas: this.gasForLock });

      const endingLockStatus = await this.peggy.getStatus(this.itemId);
      endingLockStatus.should.be.equal(false);

    });

  });

  describe('Withdrawal of items by sender', function(){

    beforeEach(async function() {
      this.peggy = await Peggy.new();
      this.zeroxToken = '0xE41d2489571d322189246DaFA5ebDe1F4699F498';
      this.ethereumToken = '0x0000000000000000000000000000000000000000';
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.utf8ToHex('cosmosaccaddr985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy');
      this.gasForLock = 300000; // 300,000 Gwei

      // Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 100, { from: relayer }).should.be.fulfilled;

      await this.token.approve(this.peggy.address, 100, { from: userOne }).should.be.fulfilled;
      this.itemId = await this.peggy.lock.call(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.token.address, 100, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;

      //Lower and upper bound to allow results in range of 0.01%
      this.lowerBound = 0.9999;
      this.upperBound = 1.0001;

      //Set up gas constants
      this.gasForLock = 500000; //500,000 Gwei
      this.gasForWithdraw = 200000; //200,000 Gwei
      this.gasPrice = 200000000000; //From truffle config
    });

    it('should not allow non-senders to withdraw other\'s items', async function () {
      await this.peggy.withdraw(this.itemId, { from: userThree, gas: this.gasForLock }).should.be.rejectedWith(EVMRevert);
    });

    it('should allow senders to withdraw their own items', async function () {
      await this.peggy.withdraw(this.itemId, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
    });

    it('should return erc20 to user upon withdrawal of itemized funds', async function () {
      //Confirm that the tokens are locked on the contract
      const beforePeggyBalance = Number(await this.token.balanceOf(this.peggy.address));
      const beforeUserBalance = Number(await this.token.balanceOf(userOne));

      beforePeggyBalance.should.be.bignumber.equal(100);
      beforeUserBalance.should.be.bignumber.equal(0);

      await this.peggy.withdraw(this.itemId, { from: userOne, gas: this.gasForWithdraw }).should.be.fulfilled;
      
      //Confirm that the tokens have been unlocked and transfered
      const afterPeggyBalance = Number(await this.token.balanceOf(this.peggy.address));
      const afterUserBalance = Number(await this.token.balanceOf(userOne));

      afterPeggyBalance.should.be.bignumber.equal(0);
      afterUserBalance.should.be.bignumber.equal(100);
    });

    it('should return ethereum to user upon withdrawal of itemized funds', async function () {
      const id = await this.peggy.lock.call(this.recipient, this.ethereumToken, this.weiAmount, { from: userTwo, value: this.weiAmount, gas: this.gasForLock }).should.be.fulfilled;
      await this.peggy.lock(this.recipient, this.ethereumToken, this.weiAmount, { from: userTwo, value: this.weiAmount, gas: this.gasForLock }).should.be.fulfilled;

      //Set up accepted withdrawal balance bounds
      const lowestAcceptedBalance = this.weiAmount * this.lowerBound;
      const highestAcceptedBalance = this.weiAmount * this.upperBound;

      //Get prior balances of user and peggy contract
      const beforeUserBalance = Number(await web3.eth.getBalance(userTwo));
      const beforeContractBalance = Number(await web3.eth.getBalance(this.peggy.address));

      //Send withdrawal transaction and save gas expenditure
      const txHash = await this.peggy.withdraw(id, { from: userTwo, gas: this.gasForWithdraw}).should.be.fulfilled;
      const gasCost = this.gasPrice * txHash.receipt.gasUsed;

      //Get balances after withdrawal
      const afterUserBalance = Number(await web3.eth.getBalance(userTwo));
      const afterContractBalance = Number(await web3.eth.getBalance(this.peggy.address));

      //Check user's balances to confirm withdrawal
      const userDifference = (afterUserBalance - beforeUserBalance) + gasCost;
      userDifference.should.be.bignumber.within(lowestAcceptedBalance, highestAcceptedBalance);

      //Check contracts's balances to confirm withdrawal
      const contractDifference = (beforeContractBalance - afterContractBalance);
      contractDifference.should.be.bignumber.within(lowestAcceptedBalance, highestAcceptedBalance);
    });

    it('should emit an event upon user withdrawal containing the item\'s id', async function () {      
      //Get the event logs of a token withdrawl
      const {logs} = await this.peggy.withdraw(this.itemId, { from: userOne, gas: this.gasForLock }).should.be.fulfilled;
      const event = logs.find(e => e.event === 'LogWithdraw');

      //Check the event's parameters
      event.args._id.should.be.bignumber.equal(this.itemId);
      event.args._to.should.be.equal(userOne);
      event.args._token.should.be.equal(this.token.address);
      Number(event.args._value).should.be.bignumber.equal(100);
      Number(event.args._nonce).should.be.bignumber.equal(1);
    });

  });
  
});