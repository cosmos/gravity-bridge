'use strict';
// var createValidators = require('./valset/createValidators');
const web3 = global.web3;
const Valset = artifacts.require("./../contracts/Valset.sol");
const utils = require('./utils.js');
const ethUtils = require('ethereumjs-util');

contract('Valset', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2]
  };
  let valSet, totalGas, gasPrice;
  let addresses, powers, first_element, second_element, totalPower, validators, totalValidators;
  before('Setup contract', async function() {
    totalValidators = utils.randomIntFromInterval(1, 100); // 1-100 validators
    validators = utils.createValidators(totalValidators);
    valSet = await Valset.new(validators.addresses, validators.powers, {from: args._default});
  });

  describe('Constructor function', function() {

    // Proved by induction
    it("Saves initial validators' address in array", async function() {
      first_element = await valSet.getValidator(0);
      second_element = await valSet.getValidator(1);
      assert.isTrue(((String(first_element) == validators.addresses[0]) && (String(second_element) == validators.addresses[1])), "Initial validators' addresses array should be equal as the saved one");
    });

    // Proved by induction
    it("Saves initial validators' powers in array", async function() {
      first_element = await valSet.getPower(0);
      second_element = await valSet.getPower(1);
      assert.isTrue(((first_element.toNumber() == validators.powers[0]) && (second_element.toNumber() == validators.powers[1])), "Initial validators' powers array should be equal as the saved one");
    });

    it("Checks that addresses and powers arrays have the same length", async function() {
      addresses = await valSet.addresses;
      powers = await valSet.powers;
      let powersLength = powers.length;
      assert.lengthOf(addresses, powersLength, "Both initial arrays must have the same length");
    });

    it("Number of validators is below 100", async function() {
      addresses = await valSet.addresses;
      let valLength = addresses.length;
      assert.isAtMost(valLength, 100, "Validator set should not be larger than 100")
    });

    it("Sums totalPower correctly", async function() {
      totalPower = await valSet.getTotalPower();
      let accumulatedPower = validators.powers.reduce(utils.sumArrayValues);
      assert.strictEqual(totalPower.toNumber(), accumulatedPower, "totalSum should the sum of each individual validator's power")
    });
  });

  describe("Verifies validators' signatures", function() {

    let prevAddresses, prevPowers, newValidators, res, signs, signature, signedPower, totalPower, hashData;
    let vArray = [], rArray = [], sArray = [], signers = [];
    beforeEach('Create new validator set and get previous validator data', async function() {
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
      hashData = web3.fromAscii(ethUtils.addHexPrefix(ethUtils.bufferToHex(ethUtils.hashPersonalMessage(ethUtils.toBuffer(validators.addresses, validators.powers)))).substring(2));
    });

    it('Signature data arrays and signers array have the same size', function () {
      assert.isTrue((vArray.length === signers.length) && (vArray.length === rArray.length) && (sArray.length === rArray.length), "Arrays should have the same size");
    });

    it('Expects to throw if super majority is not reached', async function() {
      res = await valSet.verifyValidators(hashData, signers, vArray, rArray, sArray);
      assert.isAtLeast(res.logs.length, 1, "Successful call should have logged at least one event");
      if(signedPower * 3 < totalPower * 2) {
        assert.strictEqual(res.logs[0].event, "NoSupermajority", "Should have thrown the NoSupermajority event");
      } else {
        assert.notEqual(res.logs[0].event, "NoSupermajority", "Shouldn't have thrown the NoSupermajority event");
      }
    })

    it('Signatures are correct', async function () {
      res = await valSet.verifyValidators(hashData, signers, vArray, rArray, sArray);
      assert.isAtLeast(res.logs.length, 1, "Successful verification should have logged at least one event (1 on success and more than 1 if it fails)");
      assert.strictEqual(res.logs[0].event, "Verify", "On success it should have thrown Verify event");
      assert.deepEqual(res.logs[0].args.signers, signers, "'signers' uint16[] parameter from Verify event should be equal to the signers from the validator set");
    });

  });

  describe('Updates the Validator set', function() {
    let prevAddresses, prevPowers, newValidators, res, signs, signature;
    let vArray = [];
    let rArray = [];
    let sArray = [];
    let signers = [];
    beforeEach('Create new validator set and get previous validator data', async function() {
      vArray = [];
      rArray = [];
      sArray = [];
      signers = [];
      totalValidators = utils.randomIntFromInterval(1, 100); // 1-100 validators
      validators = utils.createValidators(totalValidators);
      for (var i = 0; i < totalValidators; i++) {
        signs = (Math.random() <= 0.95764); // two std
        if (signs) {
          signature = ethUtils.ecsign(ethUtils.hashPersonalMessage(ethUtils.toBuffer(validators.addresses, validators.powers)), ethUtils.toBuffer(validators.privateKeys[i]));
          vArray.push(ethUtils.bufferToInt(signature.v));
          rArray.push(web3.fromAscii(ethUtils.addHexPrefix(ethUtils.bufferToHex(signature.r)).substring(2)));
          sArray.push(web3.fromAscii(ethUtils.addHexPrefix(ethUtils.bufferToHex(signature.s)).substring(2)));
          signers.push(i);
        }
      }
      prevAddresses = await valSet.addresses;
      prevPowers = await valSet.powers;
      res = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
    });

    it('Should updated the validator set', async function () {
      assert.strictEqual(res.logs.length, 1, "Successful update should have logged Update event");
      assert.deepEqual(res.logs[0].args.newAddresses, validators.addresses, "'newAddresses' address[] parameter from Update event should be equal to the generated validators addreses");
      assert.deepEqual(res.logs[0].args.newPowers, validators.powers, "'newPowers' uint64[] parameter from Update event should be equal to the generated validators addreses");
      assert.isNumber(res.logs[0].args.seq.toNumber(), "Update event should return 'seq' param in the log");
    });

    it("Get validator signature set", async function() {
      assert.isAtMost(signers.length, validators.addresses.length, "Signers set can't be higher than validator set");
    });

    // Proved by induction
    it("Saves updated validators' address in array", async function() {
      first_element = await valSet.getValidator(0);
      second_element = await valSet.getValidator(1);
      assert.isTrue(((String(first_element) == validators.addresses[0]) && (String(second_element) == validators.addresses[1])), "Initial validators' addresses array should be equal as the saved one");
    });

    // Proved by induction
    it("Changes the validators' addresses", async function() {
      first_element = await valSet.getValidator(0);
      second_element = await valSet.getValidator(1);
      assert.isFalse(((String(first_element) == prevAddresses[0]) || (String(second_element) == prevAddresses[1])), "New validators' addresses should be disctinct as the previous validator set addresses");
    });

    // Proved by induction
    it("Saves updated validators' powers in array", async function() {
      first_element = await valSet.getPower(0);
      second_element = await valSet.getPower(1);
      assert.isTrue(((first_element.toNumber() == validators.powers[0]) && (second_element.toNumber() == validators.powers[1])), "Initial validators' powers array should be equal as the saved one");
    });

    // Proved by induction
    it("Changes the validators' powers", async function() {
      first_element = await valSet.getPower(0);
      second_element = await valSet.getPower(1);
      assert.isFalse(((String(first_element) == prevPowers[0]) || (String(second_element) == prevPowers[1])), "New validators' powers should be disctinct as the previous validator set powers");
    });
  });



});
