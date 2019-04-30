const EscrowProcessor = artifacts.require('TestEscrowProcessor');
const TestToken = artifacts.require('TestToken');

const Web3Utils = require('web3-utils');
const EVMRevert = 'revert';
const BigNumber = web3.BigNumber;

require('chai')
  .use(require('chai-as-promised'))
  .use(require('chai-bignumber')(BigNumber))
  .should();

contract('TestEscrowProcessor', function (accounts) {

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe('Escrow Processor contract deployment', function() {

    beforeEach(async function() {
      this.processor = await EscrowProcessor.new();
    });

    it('should deploy the escrow processor with the correct parameters', async function () {
      this.processor.should.exist;
      
      const nonce = Number(await this.processor.nonce());
      nonce.should.be.bignumber.equal(0);
    });

  });

  describe('Escrow creation', function() {

    beforeEach(async function() {
      this.processor = await EscrowProcessor.new();
      this.recipient = web3.utils.bytesToHex(['20bytestring'])
      this.amount = 250;

      //Load user account with tokens for testing
      this.token = await TestToken.new();
      await this.token.mint(userOne, 1000, { from: accounts[0] }).should.be.fulfilled;
    });

 
    it('should allow for the creation of escrows', async function () {
      await this.processor.callCreateEscrow(userOne, this.recipient, this.token.address, this.amount).should.be.fulfilled;
    });

    it('should generate unique escrow id\'s for a created escrow', async function () {
      //Simulate sha3 hash to get escrow's expected id
      const expectedId = Web3Utils.soliditySha3(
        {t: 'address payable', v: userOne},
        {t: 'bytes', v: this.recipient},
        {t: 'address', v: this.token.address},
        {t: 'int256', v:this.amount},
        {t: 'int256', v:1});

      //Get the escrow's id if it were to be created
      const id = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.token.address, this.amount);
      id.should.be.equal(expectedId);
    });

    it('should allow access to an escrow\'s information given it\'s unique id', async function () {
      const id = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreateEscrow(userOne, this.recipient, this.token.address, this.amount);

      //Attempt to get an escrow's information
      await this.processor.callGetEscrow(id).should.be.fulfilled;
    });

    it('should correctly identify the existence of escrows in memory', async function () {
      //Get the escrow's expected id then lock funds
      const id = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreateEscrow(userOne, this.recipient, this.token.address, this.amount).should.be.fulfilled;
      
      //Check if escrow has been created
      const isEscrow = await this.processor.callIsEscrow(id);
      isEscrow.should.be.equal(true);
    });

    it('should store escrows with the correct parameters', async function () {
      //Create the escrow and store its id
      const id = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.token.address, this.amount);
      await this.processor.callCreateEscrow(userOne, this.recipient, this.token.address, this.amount);

      //Get the escrow's information
      const escrowInfo = await this.processor.callGetEscrow(id);

      //Parse each attribute
      const sender = escrowInfo[0];
      const receiver = escrowInfo[1];
      const token = escrowInfo[2];
      const amount = Number(escrowInfo[3]);
      const nonce = Number(escrowInfo[4]);

      //Confirm that each attribute is correct
      sender.should.be.equal(userOne);
      receiver.should.be.equal(this.recipient);
      token.should.be.equal(this.token.address);
      amount.should.be.bignumber.equal(this.amount);
      nonce.should.be.bignumber.equal(1);
    });

  });

  describe('Escrow completion', function() {

    beforeEach(async function() {
      this.processor = await EscrowProcessor.new();
      this.weiAmount = web3.utils.toWei("0.25", "ether");
      this.recipient = web3.utils.bytesToHex(['20bytestring'])
      this.ethereumToken = '0x0000000000000000000000000000000000000000';

      //Load contract with ethereum so it can complete escrows
      await this.processor.send(web3.utils.toWei("1", "ether"), { from: accounts[0]}).should.be.fulfilled;

      this.escrowId = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.ethereumToken, this.weiAmount);
      await this.processor.callCreateEscrow(userOne, this.recipient, this.ethereumToken, this.weiAmount);
    });

    it('should not allow for the completion of escrows whose value exceeds the contract\'s balance', async function () {
      //Create an escrow with an overlimit amount
      const overlimitAmount = web3.utils.toWei("1.25", "ether");
      const escrowId = await this.processor.callCreateEscrow.call(userOne, this.recipient, this.ethereumToken, overlimitAmount);
      await this.processor.callCreateEscrow(userOne, this.recipient, this.ethereumToken, overlimitAmount);

      //Attempt to complete the escrow
      await this.processor.callCompleteEscrow(escrowId).should.be.rejectedWith(EVMRevert);
    });

    it('should not allow for the completion of non-escrows', async function () {
      //Generate a false escrow id
      const fakeId = Web3Utils.soliditySha3(
        {t: 'address payable', v: userOne},
        {t: 'bytes', v: this.recipient},
        {t: 'address', v: this.ethereumToken},
        {t: 'int256', v:12},
        {t: 'int256', v:1});

      await this.processor.callCompleteEscrow(fakeId).should.be.rejectedWith(EVMRevert);
    
    });

    it('should not allow for the completion of an escrow that has already been completed', async function () {
      //Complete the escrow
      await this.processor.callCompleteEscrow(this.escrowId).should.be.fulfilled;

      //Attempt to complete the escrow again
      await this.processor.callCompleteEscrow(this.escrowId).should.be.rejectedWith(EVMRevert);
    });

    it('should allow for an escrow to be completed', async function () {
      await this.processor.callCompleteEscrow(this.escrowId).should.be.fulfilled;
    });

    it('should delete escrows upon completion', async function () {
      //Confirm that the escrow is active
      const escrowExists = await this.processor.callIsEscrow(this.escrowId);
      escrowExists.should.be.equal(true);

      //Complete the escrow
      await this.processor.callCompleteEscrow(this.escrowId).should.be.fulfilled;

      //Check if the escrow still exists
      const completedEscrow = await this.processor.callIsEscrow(this.escrowId);
      completedEscrow.should.be.equal(false);
    });

    it('should correctly transfer escrowed funds to the original sender', async function () {
      //Get prior balances of user and peggy contract
      const beforeUserBalance = Number(await web3.eth.getBalance(userOne));
      const beforeContractBalance = Number(await web3.eth.getBalance(this.processor.address));

      await this.processor.callCompleteEscrow(this.escrowId).should.be.fulfilled;

      //Get balances after completion
      const afterUserBalance = Number(await web3.eth.getBalance(userOne));
      const afterContractBalance = Number(await web3.eth.getBalance(this.processor.address));

      //Expected balances 
      afterUserBalance.should.be.bignumber.equal(beforeUserBalance + Number(this.weiAmount)); 
      afterContractBalance.should.be.bignumber.equal(beforeContractBalance - Number(this.weiAmount));
    });

  });

});