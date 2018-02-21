'use strict';

/* Add the dependencies you're testing */
const utils = require('./utils.js');
const web3 = global.web3;
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");
const Peggy = artifacts.require("./../contracts/Peggy.sol");
const MockERC20Token = artifacts.require("./../contracts/MockERC20Token.sol");
const createKeccakHash = require('keccak');
const ethUtils = require('ethereumjs-util');

contract('Peggy', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2],
    _address0: "0x0000000000000000000000000000000000000000"
  };

  let validators, standardTokenMock;
  let _account_one = args._account_one;
  let _account_two = args._account_two;
  let _address0 = args._address0;


	before('Setup Validators', async function() {
    validators = utils.createValidators(20);
  });

  describe('Peggy(address[],uint64[]', function () {
    let res, peggy;

    before ('Sets up Peggy contract', async function () {
      peggy = await Peggy.new(validators.addresses, validators.powers, {from: args._default});
    });

    it ('Correctly verifies ValSet signatures', async function () {
      let hashData = String(await peggy.hashValidatorArrays.call(validators.addresses, validators.powers));
      let signatures = await utils.createSigns(validators, hashData);

      res = await peggy.verifyValidators.call(hashData, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      assert.isTrue(res, "Should have successfully verified signatures");
    });
  });

  describe('newCosmosERC20(string,uint,uint[],uint8[],bytes32[],bytes32[]', function () {
    let res, peggy, cosmosTokenAddress, cosmosToken;

    before ('Creates new Cosmos ERC20 token', async function () {
      peggy = await Peggy.new(validators.addresses, validators.powers, {from: args._default});

      let hashData = String(await peggy.hashNewCosmosERC20.call('ATOMS', 18));
      let signatures = await utils.createSigns(validators, hashData);

      cosmosTokenAddress = await peggy.newCosmosERC20.call('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      res = await peggy.newCosmosERC20('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      cosmosToken = await CosmosERC20.at(cosmosTokenAddress);
    });

    it('Adds new token to cosmosToken mapping', async function () {
      assert.equal(await peggy.getCosmosTokenAddress('ATOMS'), cosmosTokenAddress);
    });

    it('Adds address to cosmosTokensAddresses set', async function () {
      assert.isTrue(await peggy.isCosmosTokenAddress(cosmosTokenAddress));
    });

    it('Emits NewCosmosERC20 event', async function () {
      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "NewCosmosERC20", "Successful execution should have logged the NewCosmosERC20 event");
      assert.strictEqual(res.logs[0].args.name, 'ATOMS');
      assert.strictEqual(res.logs[0].args.tokenAddress, cosmosTokenAddress);
    });

    it('Is controller of new CosmosERC20', async function () {
      assert.equal(await cosmosToken.controller.call(), peggy.address);
    });

    it('Fails if same name is resubmitted', async function () {
      let hashData = String(await peggy.hashNewCosmosERC20.call('ATOMS', 10));
      let signatures = await utils.createSigns(validators, hashData);

      await utils.expectRevert(peggy.newCosmosERC20('ATOMS', 10, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray));
    });
  });


  describe('lock(bytes,address,uint64)', function () {
    let res, peggy, cosmosTokenAddress, standardTokenMock;

    beforeEach('Sets up peggy contract', async function () {
      peggy = await Peggy.new(validators.addresses, validators.powers, {from: args._default});
    });

    it('Recieves Normal ERC20 and emits Lock event', async function () {
      let standardTokenMock = await MockERC20Token.new(_account_one, 10000, {from: args._default});
      await standardTokenMock.approve(peggy.address, 1000, {from: args._account_one});
      let res = await peggy.lock("0xdeadbeef", standardTokenMock.address, 1000, {from: args._account_one});

      assert.equal((await standardTokenMock.balanceOf(peggy.address)).toNumber(), 1000);
      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Lock");
      assert.strictEqual(String(res.logs[0].args.to), '0xdeadbeef');
      assert.strictEqual(res.logs[0].args.token, standardTokenMock.address);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 1000);
    });


    it('Burns CosmosERC20 and emits Lock event', async function () {

      let hashData = String(await peggy.hashNewCosmosERC20.call('ATOMS', 18));
      let signatures = await utils.createSigns(validators, hashData);
      let cosmosTokenAddress = await peggy.newCosmosERC20.call('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      await peggy.newCosmosERC20('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      let cosmosToken = CosmosERC20.at(cosmosTokenAddress);
      hashData = await peggy.hashUnlock(_account_one, cosmosTokenAddress, 1000);
      signatures = await utils.createSigns(validators, hashData);
      await peggy.unlock(_account_one, cosmosTokenAddress, 1000, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);

      let res = await peggy.lock("0xdeadbeef", cosmosTokenAddress, 500, {from: args._account_one});

      assert.equal((await cosmosToken.balanceOf(_account_one)).toNumber(), 500);
      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Lock");
      assert.strictEqual(String(res.logs[0].args.to), '0xdeadbeef');
      assert.strictEqual(res.logs[0].args.token, cosmosTokenAddress);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 500);
    });

    it('Sends Ether when token is 0 address and emits Lock event', async function () {

      let res = await peggy.lock("0xdeadbeef", _address0, 1000, {from: args._account_one, value: 1000});

      let ethBalance = await web3.eth.getBalance(peggy.address);

      assert.equal(ethBalance.toNumber(), 1000);
      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Lock");
      assert.strictEqual(String(res.logs[0].args.to), '0xdeadbeef');
      assert.strictEqual(res.logs[0].args.token, _address0);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 1000);
    });
  });


  describe('unlock(address,address,uint64,uint[],uint8[],bytes32[],bytes32[])', function () {
    let peggy, res;

    beforeEach('Sets up peggy contract', async function () {
      peggy = await Peggy.new(validators.addresses, validators.powers, {from: args._default});
    });

    it('Sends Normal ERC20 and emits Unlock event', async function () {
      let standardTokenMock = await MockERC20Token.new(peggy.address, 10000, {from: args._default});
      let hashData = await peggy.hashUnlock(_account_one, standardTokenMock.address, 1000);
      let signatures = await utils.createSigns(validators, hashData);

      res = await peggy.unlock(args._account_one, standardTokenMock.address, 1000, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      assert.equal((await standardTokenMock.balanceOf(_account_one)).toNumber(), 1000);

      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Unlock");
      assert.strictEqual(String(res.logs[0].args.to), args._account_one);
      assert.strictEqual(res.logs[0].args.token, standardTokenMock.address);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 1000);
    });

    it('Mints Cosmos ERC20 and emits Unlock event', async function () {

      let hashData = String(await peggy.hashNewCosmosERC20.call('ATOMS', 18));
      let signatures = await utils.createSigns(validators, hashData);
      let cosmosTokenAddress = await peggy.newCosmosERC20.call('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      await peggy.newCosmosERC20('ATOMS', 18, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      let cosmosToken = CosmosERC20.at(cosmosTokenAddress);

      hashData = await peggy.hashUnlock(_account_one, cosmosTokenAddress, 1000);
      signatures = await utils.createSigns(validators, hashData);

      res = await peggy.unlock(_account_one, cosmosTokenAddress, 1000, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      assert.equal((await cosmosToken.balanceOf(_account_one)).toNumber(), 1000);

      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Unlock");
      assert.strictEqual(String(res.logs[0].args.to), args._account_one);
      assert.strictEqual(res.logs[0].args.token, cosmosTokenAddress);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 1000);
    });

    it('Sends Ether when token is address 0x0 and emits Unlock event', async function () {
      // fund the peggy contract with a little bit of ether
      await peggy.lock("0xdeadbeef", _address0, 5000, {from: args._account_two, value: 5000});

      let oldBalance = await web3.eth.getBalance(_account_one);

      let hashData = await peggy.hashUnlock(_account_one, _address0, 1000);
      let signatures = await utils.createSigns(validators, hashData);
      res = await peggy.unlock(args._account_one, args._address0, 1000, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      assert.equal(await web3.eth.getBalance(_account_one).toNumber(), oldBalance.toNumber() + 1000);

      assert.strictEqual(res.logs.length, 1);
      assert.strictEqual(res.logs[0].event, "Unlock");
      assert.strictEqual(String(res.logs[0].args.to), args._account_one);
      assert.strictEqual(res.logs[0].args.token, args._address0);
      assert.strictEqual(res.logs[0].args.value.toNumber(), 1000);
    });
  });
});
