'use strict';
// var createValidators = require('./valset/createValidators');
const web3 = global.web3;
const Valset = artifacts.require("./../contracts/Valset.sol");
const createKeccakHash = require('keccak');
const keythereum = require("keythereum");
const ethUtils = require('ethereumjs-util');

function sumArrayValues(total, uint64) {
  return total + uint64;
}

function createValidators(size) {
  var newValidators = {
    addresses: [],
    pubKeys: [],
    privateKeys: [],
    powers: []
  };
  let privateKey,hexPrivate, pubKey, address;
  for (var i=0; i< size; i++) {
    privateKey = keythereum.create().privateKey;
    hexPrivate = ethUtils.bufferToHex(privateKey);
    address = ethUtils.addHexPrefix(ethUtils.bufferToHex(ethUtils.privateToAddress(privateKey)));
    pubKey = ethUtils.bufferToHex(ethUtils.privateToPublic(privateKey));

    // console.log("Keys: \n\tPrivate: " + hexPrivate + "\n\tPublic:" + pubKey + "\n\Address:" + address);
    newValidators.addresses.push(address);
    newValidators.privateKeys.push(hexPrivate);
    newValidators.pubKeys.push(pubKey);
    newValidators.powers.push(Math.floor((Math.random() * 50) + 1));
    }
  return newValidators;
}

contract('Valset', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2]
  };
  let valSet, totalGas, gasPrice;
  let addresses, powers, first_element, second_element, totalPower;
  let totalValidators = (Math.random() * 100) + 1; // 1-100 validators
  let validators = createValidators(totalValidators);
  beforeEach('Setup contract', async function() {
    valSet = await Valset.new(validators.addresses, validators.powers, {from: args._default});
  });

  describe('Constructor function', function() {

    // Proved by induction
    it("Saves initial validators' address in array", async function() {
      first_element = await valSet.getValidator(0);
      second_element = await valSet.getValidator(1);
      // returns the string Address of the elements and check if they exist
      // console.log("First Validator: ",String(first_element));
      // console.log("Second Validator: ",String(second_element));
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
      let accumulatedPower = validators.powers.reduce(sumArrayValues);
      assert.strictEqual(totalPower.toNumber(), accumulatedPower, "totalSum should the sum of each individual validator's power")
    });
  });

  describe('Updates the Validator set', function() {
    let prevAddresses, prevPowers, newValidators, response, signs, signature;
    let vArray = [];
    let rArray = [];
    let sArray = [];
    let signers = [];

    before('Get previous validator data', async function() {
      totalValidators = (Math.random() * 100) + 1; // 1-100 validators
      validators = createValidators(totalValidators);
      // let msgHash = {
      //   addresses: validators.addresses,
      //   powers: validators.powers
      // };
      for (var i = 0; i < totalValidators; i++) {
        signs = Math.random() <= 0.682; // one std range from 0.5
        if (signs) {
          signature = ethUtils.ecsign(ethUtils.hashPersonalMessage(ethUtils.toBuffer(validators.addresses, validators.powers)), ethUtils.toBuffer(validators.privateKeys[i]));
          vArray.push(signature.v);
          rArray.push(signature.r);
          sArray.push(signature.s);
          signers.push(i);
        }
      }
      // console.log('Test');
      // console.log(validators);
      console.log("Val size: ", validators.addresses.length)
      prevAddresses = await valSet.addresses;
      prevPowers = await valSet.powers;
      response = await valSet.update(validators.addresses, validators.powers, signers, vArray, rArray, sArray);
    });


    /* TODO check signatures */

    it("Returns a successful response on update", async function() {
      assert.isTrue(Boolean(response.receipt.status), "Succesful update should return true");
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
