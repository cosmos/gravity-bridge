'use strict';
/* Add the dependencies you're testing */
const web3 = global.web3;
const utils = require('./utils.js');
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");

contract('CosmosERC20', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2],
  };

  const _name = 'ATOMS';
  const _decimals = 18;
  const _controller = args._default;

  const _account_one = accounts[1];
  const _account_two = accounts[2];

  const numTestAccounts = 3;

  /* Functions */

  describe('CosmosERC20(address,bytes)', function() {
    let cosmosToken;

    before('Setup contract', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
    });

    it('Has a name', async function () {
      const name = await cosmosToken.name();
      assert.strictEqual(String(name), _name, "Name should match");
    });

    it('Has a symbol', async function () {
      const symbol = await cosmosToken.symbol();
      assert.strictEqual(String(symbol), _name, "Symbol should match name");
    });

    it('Has an amount of decimals', async function () {
      const decimals = await cosmosToken.decimals();
      assert.strictEqual(decimals.toNumber(), _decimals, "Decimals should match")
    });

    it('Has a controller', async function () {
      const controller = await cosmosToken.controller();
      assert.equal(controller, _controller, "Controllers should match")
    });

    it('Has an initial supply of 0', async function () {
      const totalSupply = await cosmosToken.totalSupply();
      assert.strictEqual(totalSupply.toNumber(), 0, "Initial supply should be 0")
    });
  });

  describe('mint(address,uint)', function() {
    let cosmosToken, res;

    before('Controller mints tokens', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
      res = await cosmosToken.mint(_account_one, 50, {from: _controller});
      let bal = await cosmosToken.balanceOf(args._account_one);
    });

    it('Increases balance of minting recipient', async function () {
      let bal = await cosmosToken.balanceOf(args._account_one);
    });

    it('Increases totalSupply', async function () {
      let totalSupply = await cosmosToken.totalSupply();
      assert.strictEqual(totalSupply.toNumber(), 50, "totalSupply should be increased");
    });

    it('Emits Mint event', async function () {
      assert.strictEqual(res.logs.length, 1, "Successful minting should have logged one event");
      assert.strictEqual(res.logs[0].event, "Mint", "Successful update should have logged the Mint event");
      assert.strictEqual(res.logs[0].args._to, args._account_one, "Mint event should have proper _to field");
      assert.strictEqual(res.logs[0].args._amount.toNumber(), 50, "Mint event should have proper _amount field");
    });

    it('Reverts if non-controller tries to mint', async function () {
      await utils.expectRevert(cosmosToken.mint(args._account_two, 50, {from: args._account_two}));
    });
  });

  describe('burn(address,uint)', function() {
    let cosmosToken, res;

    before('Controller burns tokens', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
      await cosmosToken.mint(args._account_one, 50, {from: _controller});
      res = await cosmosToken.burn(args._account_one, 25, {from: _controller});
    });

    it('Decreases balance of minting recipient', async function () {
      let bal = await cosmosToken.balanceOf(args._account_one);
      assert.strictEqual(bal.toNumber(), 25, "balance should be decreased");
    });

    it('Decreases totalSupply', async function () {
      let totalSupply = await cosmosToken.totalSupply();
      assert.strictEqual(totalSupply.toNumber(), 25, "totalSupply should be decreased");
    });

    it('Emits Burn event', async function () {
      assert.strictEqual(res.logs.length, 1, "Successful burning should have logged one event");
      assert.strictEqual(res.logs[0].event, "Burn", "Successful update should have logged the Burn event");
      assert.strictEqual(res.logs[0].args._from, args._account_one, "Burn event should have proper _from field");
      assert.strictEqual(res.logs[0].args._amount.toNumber(), 25, "Burn event should have proper _amount field");
    });

    it('Reverts if tries to burn more than balance', async function () {
      await utils.expectRevert(cosmosToken.burn(args._account_one, 100, {from: _controller}));
    });

    it('Reverts if non-controller tries to mint', async function () {
      await utils.expectRevert(cosmosToken.burn(args._account_one, 10, {from: _account_two}));
    });
  });


  describe('transfer(address,uint)', function () {
    let cosmosToken, res;

    beforeEach('Transfers tokens', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
      res = await cosmosToken.mint(args._account_one, 50, {from: _controller});
      res = await cosmosToken.transfer(args._account_two, 25, {from: args._account_one});
    });

    it('Correctly modifies balances', async function () {
      const senderBalance = await cosmosToken.balanceOf(args._account_one);
      assert.equal(senderBalance, 25);

      const recipientBalance = await cosmosToken.balanceOf(args._account_one);
      assert.equal(recipientBalance, 25);
    });

    it('Emits Transfer event', async function () {
      assert.equal(res.logs.length, 1);
      assert.equal(res.logs[0].event, 'Transfer');
      assert.equal(res.logs[0].args._from, args._account_one);
      assert.equal(res.logs[0].args._to, args._account_two);
      assert.equal(res.logs[0].args._value.toNumber(), 25);
    });

    it('Reverts if try to send more than balance', async function () {
      await utils.expectRevert(cosmosToken.transfer(args._account_one, 100, {from: args._account_one}));
    });

    it('Reverts if try to send to controller', async function () {
      await utils.expectRevert(cosmosToken.transfer(_controller, 10, {from: args._account_one}));
    });
  });

  describe('approve(address,uint)', function () {
    let cosmosToken, res;

    before('Gives allowance', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
      await cosmosToken.mint(_account_one, 50, {from: _controller});
      res = await cosmosToken.approve(_account_two, 25, {from: args._account_one});
    });

    it('Correctly increases allowance', async function () {
      const allowance = await cosmosToken.allowance(args._account_one, args._account_two);
      assert.equal(allowance.toNumber(), 25);
    });

    it('Emits Approval event', async function () {
      assert.equal(res.logs.length, 1);
      assert.equal(res.logs[0].event, 'Approval');
      assert.equal(res.logs[0].args._owner, args._account_one);
      assert.equal(res.logs[0].args._spender, args._account_two);
      assert.equal(res.logs[0].args._value.toNumber(), 25);
    });
  });

  describe('transferFrom(address,address,uint)', function () {
    let cosmosToken, res;

    before('Spends Allowance', async function() {
      cosmosToken = await CosmosERC20.new(_controller, _name, _decimals, {from: args._default});
      await cosmosToken.mint(_account_one, 100, {from: _controller});
      await cosmosToken.approve(_account_two, 25, {from: _account_one});
      res = await cosmosToken.transferFrom(_account_one, _account_two, 25, {from: _account_two});
    });

    it('Correctly transfers tokens', async function () {
      const senderBalance = await cosmosToken.balanceOf(_account_one);
      assert.equal(senderBalance.toNumber(), 75);

      const recipientBalance = await cosmosToken.balanceOf(_account_two);
      assert.equal(recipientBalance.toNumber(), 25);
    });

    it('Decreases allowance', async function () {
      const allowance = await cosmosToken.allowance(_account_one, _account_two);
      assert.equal(allowance.toNumber(), 0);
    });

    it('Reverts if spending more than allowance', async function () {
      await utils.expectRevert(cosmosToken.transfer(_account_one, 50, {from: _account_two}));
    });
  });
});
