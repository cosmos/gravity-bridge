var _ = require("lodash");
var Promise = require("bluebird");
const keythereum = require("keythereum");
const ethUtils = require('ethereumjs-util');

module.exports = {
  // Methods from web3 1.0:  https://github.com/ethereum/web3.js/blob/1.0/packages/web3-utils/src/utils.js
  /**
   * Check if string is HEX, requires a 0x in front
   *
   * @method isHexStrict
   * @param {String} hex to be checked
   * @returns {Boolean}
   */
  isHexStrict: function (hex) {
      return ((_.isString(hex) || _.isNumber(hex)) && /^(-)?0x[0-9a-f]*$/i.test(hex));
  },
  /**
   * Convert a hex string to a byte array
   *
   * Note: Implementation from crypto-js
   *
   * @method hexToBytes
   * @param {string} hex
   * @return {Array} the byte array
   */
  hexToBytes: function(hex) {
      hex = hex.toString(16);

      if (!isHexStrict(hex)) {
          throw new Error('Given value "'+ hex +'" is not a valid hex string.');
      }

      hex = hex.replace(/^0x/i,'');

      for (var bytes = [], c = 0; c < hex.length; c += 2)
          bytes.push(parseInt(hex.substr(c, 2), 16));
      return bytes;
  },
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

      // console.log("Keys: \n\tPrivate: " + hexPrivate + "\n\tPublic:" + pubKey + "\n\Address:" + address);
      newValidators.addresses.push(address);
      newValidators.privateKeys.push(hexPrivate);
      newValidators.pubKeys.push(pubKey);
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
