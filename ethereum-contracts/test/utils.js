'use strict';

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
  seededRandomInt: function(min, max, seed) {
    seed = (seed * 9301 + 49297) % 233280;
    var rnd = seed / 233280;

    return Math.floor(min + rnd * (max - min));
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
      powers: [],
      totalPower: 0
    };

    let privateKey,hexPrivate, pubKey, address, power;

    for (var i = 0; i < size; i++) {
      privateKey = keythereum.create().privateKey;
      hexPrivate = ethUtils.bufferToHex(privateKey);
      address = ethUtils.bufferToHex(ethUtils.privateToAddress(privateKey));
      pubKey = ethUtils.bufferToHex(ethUtils.privateToPublic(privateKey));
      power = this.seededRandomInt(1, 50, i);

      newValidators.addresses.push(address);
      newValidators.privateKeys.push(hexPrivate);
      newValidators.pubKeys.push(pubKey);
      newValidators.powers.push(power);
      newValidators.totalPower += power;
    }

    return newValidators;
  },
  assignPowersToAccounts: function(accounts) {
    var newValidators = {
      addresses: accounts,
      powers: []
    };
    for (var i = 0; i < accounts.length; i++) {
      newValidators.powers.push(this.seededRandomInt(1, 50, i));
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
  createSigns: async function (validators, data, percentSign) {
    var vArray = [], rArray = [], sArray = [], signers = [];
    var signedPower = 0;
    if (!percentSign) {
      percentSign = 0.95
    }
    for (var i = 0; i < validators.addresses.length; i++) {
      if (this.seededRandomInt(1, 100, i) <= percentSign * 100) {

        let signature = await ethUtils.ecsign(ethUtils.toBuffer(data), ethUtils.toBuffer(validators.privateKeys[i]));

        vArray.push(signature.v);
        rArray.push(ethUtils.bufferToHex(signature.r));
        sArray.push(ethUtils.bufferToHex(signature.s));
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
  },
  expectThrow: async function (promise) {
    try {
      await promise;
    } catch (error) {
      const revert = error.message.search('revert') >= 1;
      const invalidOpcode = error.message.search('invalid opcode') >= 0;
      const outOfGas = error.message.search('out of gas') >= 0;
      assert(
        invalidOpcode || outOfGas || revert,
        'Expected throw, got \'' + error + '\' instead',
      );
      return;
    }
    assert.fail('Expected throw not received');
  },
  hexToBool: function (hexBool) {
    if (hexBool == '0x01') {
      return true;
    } else if (hexBool == '0x00') {
      return false;
    } else throw `StatusException: ${hexBool} is not a valid transaction receipt status`;
  },
  expectRevert: async function (promise) {
    try {
      let response = await promise;
      if (typeof(response) == 'object') {
        assert.isFalse(this.hexToBool(response.receipt.status), "Should not execute properly.");
      } else {
        assert.isFalse(response, "Should not execute properly.");
      }
    } catch (error) {
      assert.isTrue(Boolean(error.message.search('revert')) || error.message.startsWith('Invalid JSON RPC response:'), 'Expected revert, got \'' + error + '\' instead');
      return;
    }
  }
}
