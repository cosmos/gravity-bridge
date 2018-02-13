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
    _amount: 100000,
    _address0: "0x0000000000000000000000000000000000000000"
  };
  let peggy, cosmosToken;
  let valSet, totalGas, gasPrice;
  let addresses, powers, first_element, second_element, totalPower, validators, totalValidators, res;

	before('Setup contract', async function() {
    totalValidators = 10;
    validators = utils.assignPowersToAccounts(accounts);
    cosmosToken = await CosmosERC20.new(args._default, web3.fromAscii("Cosmos"), {from: args._default});
    peggy = await Peggy.new(validators.addresses, validators.powers, {from: args._default});
  });

  describe('Locks tokens correctly', function () {
    it('Locks tokens on Ethereum chain', async function() {
      let bytesToParam = web3.fromAscii(args._account_one);
      res = await peggy.lock(bytesToParam, args._lock_amount, args._address0, {from: args._default, value: args._lock_amount});
      assert.isAtLeast(res.logs.length, 1, "Successful lock initialization should have logged Lock event");
      assert.strictEqual(res.logs[0].event, "Lock", "On success it should have thrown Lock event");
      assert.strictEqual(res.logs[0].args.to, args._account_one, "'to' bytes parameter from Lock event should be equal to the bytes representation of the destination address");
      assert.strictEqual(res.logs[0].args.value.toNumber(), args._lock_amount, `'value' uint64 parameter from Lock event should be equal to the lock amount ${args._lock_amount}`);
      assert.equal(res.logs[0].args.token, args._address0, `'token' address param from Lock event should be equal to ${args._address0}`);
    });

    it('Locks tokens on Cosmos chain', async function() {
      let bytesToParam = web3.fromAscii(args._account_one);
      let minted = await cosmosToken.mint(args._account_one, args._amount, {from: args._default});
      let approve = await cosmosToken.approve(peggy.address, args._lock_amount, {from: args._account_one});
      res = await peggy.lock(bytesToParam, args._lock_amount, cosmosToken.address, {from: args._account_one, value: args._lock_amount});
      assert.isAtLeast(res.logs.length, 1, "Successful lock initialization should have logged Lock event");
      assert.strictEqual(res.logs[0].event, "Lock", "On success it should have thrown Lock event");
      assert.strictEqual(res.logs[0].args.to, args._account_one, "'to' bytes parameter from Lock event should be equal to the bytes representation of the destination address");
      assert.strictEqual(res.logs[0].args.value.toNumber(), args._lock_amount, `'value' uint64 parameter from Lock event should be equal to the lock amount ${args._lock_amount}`);
      assert.equal(res.logs[0].args.token, cosmosToken.address, `'token' address param from Lock event should be equal to ${cosmosToken.address}`);
    });

  });


  describe('Unlocks tokens from locked account in sidechain', function () {
    let prevAddresses, prevPowers, newValidators, res, signs, signature, signature2, signedPower, totalPower, msg, prefix, prefixedMsg, hashData;
    let vArray = [], rArray = [], sArray = [], signers = [];

    beforeEach('Create new validator set and get previous validator data', async function() {
      vArray = [], rArray = [], sArray = [], signers = [];
      totalPower = 0, signedPower = 0;
      validators = utils.assignPowersToAccounts(accounts);
      msg = new Buffer(accounts.concat(validators.powers));
      hashData = web3.sha3(accounts.concat(validators.powers));
      prefix = new Buffer("\x19Ethereum Signed Message:\n");
      prefixedMsg = ethUtils.sha3(
        Buffer.concat([prefix, new Buffer(String(msg.length)), msg])
      );
      for (var i = 0; i < 10; i++) {
        signs = (Math.random() <= 0.95764); // two std
        totalPower += validators.powers[i];
        if (signs) {
          signature = await web3.eth.sign(validators.addresses[i], '0x' + msg.toString('hex'));
          let ethSignature = await web3.eth.sign(validators.addresses[i], hashData).slice(2);
          const rpcSignature = ethUtils.fromRpcSig(signature);
          const pubKey  = ethUtils.ecrecover(prefixedMsg, rpcSignature.v, rpcSignature.r, rpcSignature.s);
          const addrBuf = ethUtils.pubToAddress(pubKey);
          const addr    = ethUtils.bufferToHex(addrBuf);
          vArray.push(web3.toDecimal(ethSignature.slice(128, 130)) + 27);
          rArray.push('0x' + ethSignature.slice(0, 64));
          sArray.push('0x' + ethSignature.slice(64, 128));
          signers.push(i);
          signedPower += validators.powers[i];
        }
      }
    });

    it('Calls the Unlock event on success', async function() {
      let res = await peggy.unlock(args._address0, args._account_one, args._lock_amount, signers, vArray, rArray, sArray, {from: args._default});
      console.log(res);
      assert.isAtLeast(res.logs.length, 1, "Successful lock initialization should have logged Unlock event");
      assert.equal(res.logs[0].args.to, bytesToParam, "'to' address parameter from Unlock event should be equal to the generated validators addreses");
      assert.strictEqual(res.logs[0].args.value.toNumber(), args._lock_amount, `'value' uint64 parameter from Unlock event should be equal to the unlock amount ${args._lock_amount}`);
      assert.equal(res.logs[0].args.token, args._address0, `'token' address param from Unlock event should be equal to ${args._address0}`);
    });
  });

});
