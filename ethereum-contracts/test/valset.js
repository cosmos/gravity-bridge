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
    console.log("\t - Initial validators size: ", validators.addresses.length);
    valSet = await Valset.new(validators.addresses, validators.powers, {from: args._default});
  });

  describe('Constructor function', function() {

    // Proved by induction
    it("Saves initial validators' address in array", async function() {
      first_element = await valSet.getValidator.call(0);
      second_element = await valSet.getValidator.call(1);
      // returns the string Address of the elements and check if they exist
      console.log("First Validator: ",String(first_element));
      console.log("Second Validator: ",String(second_element));
      assert.isTrue(((String(first_element) == validators.addresses[0]) && (String(second_element) == validators.addresses[1])), "Initial validators' addresses array should be equal as the saved one");
    });

    // Proved by induction
    it("Saves initial validators' powers in array", async function() {
      first_element = await valSet.getPower.call(0);
      second_element = await valSet.getPower.call(1);
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

  describe('Updates the Validator set', function() {
    let prevAddresses, prevPowers, newValidators, response, signs, signature;
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
        signs = (Math.random() <= 0.682); // one std range from 0.5
        if (signs) {
          signature = ethUtils.ecsign(ethUtils.hashPersonalMessage(ethUtils.toBuffer(validators.addresses, validators.powers)), ethUtils.toBuffer(validators.privateKeys[i]));
          vArray.push(signature.v);
          rArray.push(signature.r);
          sArray.push(signature.s);
          signers.push(i);
        }
      }
      prevAddresses = await valSet.addresses;
      prevPowers = await valSet.powers;
    });

    it("Get validator signature set", async function() {
      assert.isAtMost(signers.length, validators.addresses.length, "Signers set can't be higher than validator set");
    });

    it("Returns a successful response on update", async function() {
      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
      assert.isTrue(Boolean(response.receipt.status), "Succesful update should return true");
    });

    // Proved by induction
    it("Saves updated validators' address in array", async function() {

      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
      first_element = await valSet.getValidator.call(0);
      second_element = await valSet.getValidator.call(1);
      console.log(first_element, validators.addresses[0]);
      console.log(second_element, validators.addresses[1]);
      assert.isTrue(((String(first_element) == validators.addresses[0]) && (String(second_element) == validators.addresses[1])), "Initial validators' addresses array should be equal as the saved one");
    });

    // Proved by induction
    it("Changes the validators' addresses", async function() {
      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
      first_element = await valSet.getValidator.call(0);
      second_element = await valSet.getValidator.call(1);
      assert.isFalse(((String(first_element) == prevAddresses[0]) || (String(second_element) == prevAddresses[1])), "New validators' addresses should be disctinct as the previous validator set addresses");
    });

    // Proved by induction
    it("Saves updated validators' powers in array", async function() {
      first_element = await valSet.getPower.call(0);
      second_element = await valSet.getPower.call(1);
      console.log(first_element.toNumber(), validators.powers[0]);
      console.log(second_element.toNumber(), validators.powers[1]);
      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
      first_element = await valSet.getPower.call(0);
      second_element = await valSet.getPower.call(1);
      console.log("new power0: ",first_element.toNumber(), validators.powers[0]);
      console.log("new power1: ", second_element.toNumber(), validators.powers[1]);
      assert.isTrue(((first_element.toNumber() == validators.powers[0]) && (second_element.toNumber() == validators.powers[1])), "Initial validators' powers array should be equal as the saved one");
    });

    // Proved by induction
    it("Changes the validators' powers", async function() {
      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
      first_element = await valSet.getPower(0);
      second_element = await valSet.getPower(1);
      assert.isFalse(((String(first_element) == prevPowers[0]) || (String(second_element) == prevPowers[1])), "New validators' powers should be disctinct as the previous validator set powers");
    });
  });



});
