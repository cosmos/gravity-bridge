var _ = require("lodash");
var Promise = require("bluebird");
const keythereum = require("keythereum");
const ethUtils = require('ethereumjs-util');
const Hash = require("eth-lib/lib/hash");

module.exports = {
  randomIntFromInterval: function(min,max) {
      return Math.floor(Math.random()*(max-min+1)+min);
  },
  sumArrayValues: function(total, uint64) {
    return total + uint64;
  },
  createValidators: function(size) {
    var newValidators = {
      addresses: [],
      pubKeys: [],
      privateKeys: [],
      powers: []
    };
    let privateKey,hexPrivate, pubKey, address;
    for (var i = 0; i < size; i++) {
      privateKey = keythereum.create().privateKey;
      hexPrivate = ethUtils.bufferToHex(privateKey);
      address = ethUtils.addHexPrefix(ethUtils.bufferToHex(ethUtils.privateToAddress(privateKey)));
      pubKey = ethUtils.bufferToHex(ethUtils.privateToPublic(privateKey));
      newValidators.addresses.push(address);
      newValidators.privateKeys.push(hexPrivate);
      newValidators.pubKeys.push(pubKey);
      newValidators.powers.push(this.randomIntFromInterval(1, 50)); // 1-50 power
      }
    return newValidators;
  },
  assignPowersToAccounts: function(accounts) {
    var newValidators = {
      addresses: accounts,
      powers: []
    };
    for (var i = 0; i < accounts.length; i++) {
      newValidators.powers.push(this.randomIntFromInterval(1, 50)); // 1-50 power
    }
    return newValidators;
  },
  assertEvent: function(contract, filter) {
    return new Promise((resolve, reject) => {
      var event = contract[filter.event]();
      event.watch();
      event.get((error, logs) => {
        var log = _.filter(logs, filter);
        if (log) {
          resolve(log);
        } else {
          throw Error("Failed to find filtered event for " + filter.event);
        }
      });
      event.stopWatching();
    });
  },
  expectThrow: async function (promise) {
    const errMsg = 'Expected throw not received';
    try {
      await promise;
    } catch (err) {
      assert(err.toString().includes('invalid opcode'), errMsg);
      return;
    }
    assert.fail(errMsg);
  },
}
