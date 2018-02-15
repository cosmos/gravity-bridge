var _ = require("lodash");
var Promise = require("bluebird");
const keythereum = require("keythereum");
const BN = require('bn.js');
const utf8 = require('utf8');
const ethUtils = require('ethereumjs-util');
const Hash = require("eth-lib/lib/hash");

module.exports = {
  isHexStrict: function (hex) {
    return ((_.isString(hex) || _.isNumber(hex)) && /^(-)?0x[0-9a-f]*$/i.test(hex));
  },
  randomIntFromInterval: function(min,max) {
      return Math.floor(Math.random()*(max-min+1)+min);
  },
  sumArrayValues: function(total, uint64) {
    return total + uint64;
  },
  utf8ToHex: function(str) {
    str = utf8.encode(str);
    var hex = "";

    // remove \u0000 padding from either side
    str = str.replace(/^(?:\u0000)*/,'');
    str = str.split("").reverse().join("");
    str = str.replace(/^(?:\u0000)*/,'');
    str = str.split("").reverse().join("");

    for(var i = 0; i < str.length; i++) {
        var code = str.charCodeAt(i);
        // if (code !== 0) {
        var n = code.toString(16);
        hex += n.length < 2 ? '0' + n : n;
        // }
    }

    return "0x" + hex;
  },
  isBN: function (object) {
    return object instanceof BN ||
        (object && object.constructor && object.constructor.name === 'BN');
  },
  isBigNumber: function (object) {
    return object && object.constructor && object.constructor.name === 'BigNumber';
  },
  leftPad: function (string, chars, sign) {
    var hasPrefix = /^0x/i.test(string) || typeof string === 'number';
    string = string.toString(16).replace(/^0x/i,'');

    var padding = (chars - string.length + 1 >= 0) ? chars - string.length + 1 : 0;

    return (hasPrefix ? '0x' : '') + new Array(padding).join(sign ? sign : "0") + string;
  },
  rightPad: function (string, chars, sign) {
    var hasPrefix = /^0x/i.test(string) || typeof string === 'number';
    string = string.toString(16).replace(/^0x/i,'');

    var padding = (chars - string.length + 1 >= 0) ? chars - string.length + 1 : 0;

    return (hasPrefix ? '0x' : '') + string + (new Array(padding).join(sign ? sign : "0"));
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
  createSigns: async function (validators, data) {
    var vArray = [], rArray = [], sArray = [], signers = [];
    var signedPower = 0;
    for (var i = 0; i < validators.addresses.length; i++) {
      let signs = (Math.random() <= 0.95764);
      if (signs) {
        let signature = await web3.eth.sign(validators.addresses[i], data).slice(2);
        vArray.push(web3.toDecimal(signature.slice(128, 130)) + 27);
        rArray.push('0x'+signature.slice(0, 64));
        sArray.push('0x'+signature.slice(64, 128));
        signers.push(i);
        signedPower += validators.powers[i];
      }
    } 
    return {
        signers: signers,
        vArray: vArray,
        rArray: rArray,
        sArray: sArray,
        signedPower: signedPower
    } 
  } 
}
