'use strict';

const web3 = global.web3;
const Valset = artifacts.require("./../contracts/Valset.sol");
const utils = require('./utils.js');
const ethUtils = require('ethereumjs-util');
const soliditySha3 = require('./soliditySha3.js');

contract('Valset', function(accounts) {
  const args = { _default: accounts[0] };

  let valSet, totalGas, gasPrice;

  let validators, totalValidators;

  before('Setup contract', async function() {
    validators = utils.createValidators(20);
    valSet = await Valset.new(validators.addresses, validators.powers, {from: args._default});
  });

  describe('ValSet(address[],uint64[])', function() {
    it("Saves initial validators' addresses and powers in array", async function() {

      let contractAddresses = await valSet.getAddresses.call();
      let contractPowers = await valSet.getPowers.call();
      let contractTotalPower = await valSet.getTotalPower.call();

      assert.strictEqual(contractAddresses.length, contractPowers.length, "Both initial arrays must have the same length");
      assert.isAtMost(contractAddresses.length, 100, "Validator set should not be larger than 100");

      for (var i = 0; i < validators.addresses.length; i++) {
        assert.strictEqual(String(contractAddresses[i]), validators.addresses[i], "Initial validators' addresses array should be equal as the saved one")
        assert.strictEqual(contractPowers[i].toNumber(), validators.powers[i], "Initial validators' powers array should be equal as the saved one");
      }

      assert.strictEqual(contractTotalPower.toNumber(), validators.totalPower, "totalSum should the sum of each individual validator's power")
    });
  });

  describe("verifyValidators(bytes32,uint[],uint8[],bytes32[],bytes32[])", function() {
    let hashData;

    before('Hashes the data', async function() {
      hashData = await valSet.hashValidatorArrays.call(validators.addresses, validators.powers);
    });

    // it('Hashes validators array correctly', function () {
    //   let hashObj = soliditySha3({t: 'address', v: validators.addresses}, {t: 'uint64', v: validators.powers});
    //   assert.strictEqual(hashObj, hashData, "keccak256 hashes should be equal");
    // });

    it('Correctly verifies signatures', async function () {
      let signatures = await utils.createSigns(validators, hashData);
      assert.isAtLeast(signatures.signedPower * 3, validators.totalPower * 2, "Did not have supermajority. Try increasing signProbability threshhold.");

      let res = await valSet.verifyValidators.call(hashData, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);
      assert.isTrue(res, "Should have successfully verified signatures");
    });

    it('Throws if super majority is not reached', async function() {
      let signatures = await utils.createSigns(validators, hashData, 0.25);
      assert.isBelow(signatures.signedPower * 3, validators.totalPower * 2, "Still had supermajority. Try lowering signProbability threshhold.");

      await utils.expectRevert(valSet.verifyValidators.call(hashData, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray));
    })

    it('Throws if invalid signature is included', async function() {
      let signatures = await utils.createSigns(validators, hashData);
      assert.isAtLeast(signatures.signedPower * 3, validators.totalPower * 2, "Did not have supermajority. Try increasing signProbability threshhold.");
      signatures.rArray[0] = signatures.rArray[1];

      await utils.expectRevert(valSet.verifyValidators.call(hashData, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray));
    })
  });


  describe('update(address[],uint64[],uint[],uint8[],bytes32[],bytes32[])', function() {
    // let prevAddresses, prevPowers, newValidators, res, signs, signature, signature2, signedPower, totalPower, msg, prefix, prefixedMsg, hashData;
    let newValidators, hashData, signatures, res;

    before('Generates new validator set and signs it', async function() {
      newValidators = utils.createValidators(30);
      hashData = String(await valSet.hashValidatorArrays.call(newValidators.addresses, newValidators.powers));
      signatures = await utils.createSigns(validators, hashData);
    });

    it('Successfully updates the validator set', async function () {
      res = await valSet.update(newValidators.addresses, newValidators.powers, signatures.signers, signatures.vArray, signatures.rArray, signatures.sArray);

      let contractAddresses = await valSet.getAddresses.call();
      let contractPowers = await valSet.getPowers.call();
      let contractTotalPower = await valSet.getTotalPower.call();

      assert.strictEqual(contractAddresses.length, contractPowers.length, "Both contract arrays must have the same length");
      assert.strictEqual(contractAddresses.length, newValidators.addresses.length, "contract addresses array length should be same as passed in addresses array");
      assert.strictEqual(contractPowers.length, newValidators.powers.length, "contract powers array length should be same as passed in powers array");

      for (var i = 0; i < newValidators.addresses.length; i++) {
        assert.strictEqual(String(contractAddresses[i]), newValidators.addresses[i], "New validators' addresses array should be equal as the saved one")
        assert.strictEqual(contractPowers[i].toNumber(), newValidators.powers[i], "New validators' powers array should be equal as the saved one");
      }

      assert.strictEqual(contractTotalPower.toNumber(), newValidators.totalPower, "totalSum should the sum of each individual validator's power")
    });

    it('Successfully logs the Update event', async function () {
      assert.strictEqual(res.logs.length, 1, "Successful update should have logged one event");
      assert.strictEqual(res.logs[0].event, "Update", "Successful update should have logged the Update event");
      assert.strictEqual(res.logs[0].args.newAddresses.length, res.logs[0].args.newPowers.length, "Both contract arrays must have the same length");
      assert.strictEqual(res.logs[0].args.newAddresses.length, newValidators.addresses.length, "Event addresses array length should be same as passed in addresses array");
      assert.strictEqual(res.logs[0].args.newPowers.length, newValidators.powers.length, "Event powers array length should be same as passed in power array");

      for (var i = 0; i < newValidators.addresses.length; i++) {
        assert.strictEqual(String(res.logs[0].args.newAddresses[i]), newValidators.addresses[i], "newAddresses' address[] parameter from Update event should be equal to the generated validators addreses");
        assert.strictEqual(res.logs[0].args.newPowers[i].toNumber(), newValidators.powers[i], "'newPowers' uint64[] parameter from Update event should be equal to the generated validators addreses");
      }

      assert.isNumber(res.logs[0].args.seq.toNumber(), "Update event should return 'seq' param in the log");
    });
  });
});
