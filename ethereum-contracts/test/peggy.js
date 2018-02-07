'use strict';
/* Add the dependencies you're testing */
const utils = require('./utils.js');
const web3 = global.web3;
// const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");
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
    // cosmosToken = await CosmosERC20.new(args._default, 'Cosmos', {from: args._default});
    peggy = await Peggy.new(validators.address, validators.powers, {from: args._default});
    console.log("test peggy");

  });

  /* Functions */

  describe('Locks tokens correctly', function () {
    it('Calls the Lock event on success', function (done) {
      let updateWatcher = peggy.Lock();
      peggy.lock(utils.hexToBytes(args._account_one), args._lock_amount, args._address0)
      // assert.isTrue(Boolean(peggy), "This should be true");
      // peggy.lock(, args._lock_amount)
    });

  });

});
