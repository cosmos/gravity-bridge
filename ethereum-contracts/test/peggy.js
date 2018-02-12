'use strict';
/* Add the dependencies you're testing */
const utils = require('./utils.js');
const web3 = global.web3;
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");
const Peggy = artifacts.require("./../contracts/Peggy.sol");
const createKeccakHash = require('keccak');
const ethUtils = require('ethereumjs-util');

contract('Peggy', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2],
    _lock_amount: 1000,
    _address0: "0x0"
  };
  let peggy, cosmosToken;
  let valSet, totalGas, gasPrice;
  let addresses, powers, first_element, second_element, totalPower, validators, totalValidators;

	before('Setup contract', async function() {
    totalValidators = utils.randomIntFromInterval(1, 100); // 1-100 validators
    validators = utils.createValidators(totalValidators);
    cosmosToken = await CosmosERC20.new(args._default, web3.fromAscii("Cosmos"), {from: args._default});
    peggy = await Peggy.new(validators.address, validators.powers, {from: args._default});
  });

  describe('Locks tokens correctly', function () {
    it('Calls the Lock event on success', function (done) {
      let bytesToParam = utils.hexToBytes(args._account_one);
      let res = peggy.lock(bytesToParam, args._lock_amount, args._address0, {from: args._default});
      console.log('event log: ', res.logs);
      assert.strictEqual(res.logs.length, 1, "Successful lock initialization should have logged Lock event");
      assert.equal(res.logs[0].args.to, bytesToParam, "'to' bytes parameter from Lock event should be equal to the bytes representation of the destination address");
      assert.strictEqual(res.logs[0].args.value.toNumber(), args._lock_amount, `'value' uint64 parameter from Lock event should be equal to the lock amount ${args._lock_amount}`);
      assert.equal(res.logs[0].args.token, args._address0, `'token' address param from Lock event should be equal to ${args._address0}`);
    });
  });

  describe('Unlocks tokens from locked account in sidechain', function () {
    let vArray = [], rArray = [], sArray = [], signers = [];
    let totalPower, hashData, signs;
    before("Get validators' signature", function() {
      vArray = [], rArray = [], sArray = [], signers = [];
      totalPower = 0, signedPower = 0;
      totalValidators = utils.randomIntFromInterval(1, 100); // 1-100 validators
      validators = utils.createValidators(totalValidators);
      for (var i = 0; i < totalValidators; i++) {
        signs = (Math.random() <= 0.95764); // two std
        totalPower += validators.powers[i];
        if (signs) {
          signature = ethUtils.ecsign(ethUtils.hashPersonalMessage(ethUtils.toBuffer(validators.addresses, validators.powers)), ethUtils.toBuffer(validators.privateKeys[i]));
          vArray.push(ethUtils.bufferToInt(signature.v));
          rArray.push(web3.fromAscii(ethUtils.addHexPrefix(ethUtils.bufferToHex(signature.r)).substring(2)));
          sArray.push(web3.fromAscii(ethUtils.addHexPrefix(ethUtils.bufferToHex(signature.s)).substring(2)));
          signers.push(i);
          signedPower += validators.powers[i];
        }
      }
    })

    it('Calls the Unlock event on success', function (done) {
      let res = peggy.unlock(args._address0, args._account_one, args._lock_amount, {from: args._default});
      assert.strictEqual(res.logs.length, 1, "Successful lock initialization should have logged Unlock event");
      assert.equal(res.logs[0].args.to, bytesToParam, "'to' address parameter from Unlock event should be equal to the generated validators addreses");
      assert.strictEqual(res.logs[0].args.value.toNumber(), args._lock_amount, `'value' uint64 parameter from Unlock event should be equal to the unlock amount ${args._lock_amount}`);
      assert.equal(res.logs[0].args.token, args._address0, `'token' address param from Unlock event should be equal to ${args._address0}`);
    });
  });

});
