'use strict';
/* Add the dependencies you're testing */
const web3 = global.web3;
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");

contract('CosmosERC20', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2],
    _zero: 0,
    _initialSupply: 0,
    _amount: 1000
  };
  let cosmosToken;
  /* Do something before every `describe` method */
	beforeEach('Setup contract', async function() {
    cosmosToken = await CosmosERC20.new(args._default, 'Cosmos', {from: args._default});
  });

  /* Functions */

  describe('Supply control', function () {
    let supply;
    beforeEach(async function() {
      supply = await cosmosToken.totalSupply.call();
    });

    it("Initial Supply is Correct", async function() {
			assert.strictEqual(supply.toNumber(), args._initialSupply, "Initial supply must be 0");
		});

    /* approve(address spender, uint tokens) */
    // allowed[tokenOwner][spender]

  describe('', function() {
    let minted;
    beforeEach('Mint', async function() {
      minted = await cosmosToken.mint(args._account_one, args._amount, {from: args._default});
    });

    /* mint(address to, uint tokens) */

    it("Can Mint Tokens", async function() {
      assert.strictEqual(minted.logs.length, 1, "Successful mint should have logged Mint event");
      assert.strictEqual(minted.logs[0].args.to, args._account_one, "'to' address parameter from Mint event should be equal to account 1");
      assert.strictEqual(minted.logs[0].args.tokens.toNumber(), args._amount, `'tokens' uint parameter from Mint event should be equal to ${args._amount}`);
			let totalSupply2 = await cosmosToken.totalSupply.call();
			assert.strictEqual(totalSupply2.toNumber(), args._amount, "Supply should increase in 1000");
      let balanceAfter = await cosmosToken.balanceOf.call(args._account_one);
      assert.strictEqual(balanceAfter.toNumber(), args._amount, "Controller's balance should increase in 1000");
		});

    /* burn(address from, uint tokens) */

    it("Can Burn tokens from a user's balance", async function() {
      // initial supply and balance = 1000
      let res = await cosmosToken.burn(args._account_one, 100, {from: args._default});
      assert.strictEqual(res.logs.length, 1, "Successful burn should have logged Burn event");
      assert.strictEqual(res.logs[0].args.from, args._account_one, "'from' address parameter from Burn event should be equal to account 1");
      assert.strictEqual(res.logs[0].args.tokens.toNumber(), 100, "'tokens' uint parameter from Burn event should be equal to 100");
      let totalSupply3 = await cosmosToken.totalSupply.call();
      assert.strictEqual(totalSupply3.toNumber(), 900, "Supply should decrease in 100");
      let balanceAfter = await cosmosToken.balanceOf.call(args._account_one);
      assert.strictEqual(balanceAfter.toNumber(), 900, "Controller's balance should decrease by 100");
    });

    describe('', function() {
      let approveEvent;
      it("Can Approve a certain amount to be spend by an user", async function() {
        let res = await cosmosToken.approve(args._account_two, args._amount, {from: args._account_one});
        assert.strictEqual(res.logs.length, 1, "Successful approve should have logged Approval event");
        assert.strictEqual(res.logs[0].args.tokenOwner, args._account_one, "'tokenOwner' address parameter from Approval event should be equal to account 1");
        assert.strictEqual(res.logs[0].args.spender, args._account_two, "'spender' address parameter from Approval event should be equal to account 1");
        assert.strictEqual(res.logs[0].args.tokens.toNumber(), args._amount, `'tokens' uint parameter from Burn event should be equal to ${args._amount}`);
        let amountAllowed = await cosmosToken.allowance(args._account_one, args._account_two);
        assert.strictEqual(Number(amountAllowed.toNumber()), args._amount, "Approved amount should be the same as the user allowed balance");
      });
    });

    /* transferFrom(address from, address to, uint tokens) */

    it("Can transfer tokens from one account to another", async function() {
      await cosmosToken.approve(args._account_two, args._amount, {from: args._account_one});
      let res = await cosmosToken.transferFrom(args._account_one, args._default, 100, {from: args._account_two});
      assert.strictEqual(res.logs.length, 1, "Successful transferFrom should have logged Transfer event");
      assert.strictEqual(res.logs[0].args.from, args._account_one, "'from' address parameter from Approval event should be equal to account 1");
      assert.strictEqual(res.logs[0].args.to, args._default, "'to' address parameter from Approval event should be equal to account 0");
      assert.strictEqual(res.logs[0].args.tokens.toNumber(), 100, "'tokens' uint parameter from Burn event should be equal to 100");
      let balanceSender = await cosmosToken.balanceOf.call(args._account_one);
      assert.strictEqual(balanceSender.toNumber(), 900, "Sender's balance should decrease by 100");
      let balanceRecipient = await cosmosToken.balanceOf.call(args._default);
      assert.strictEqual(balanceRecipient.toNumber(), 100, "Recipient's balance should increase by 100");
    });

    /* transfer(address to, uint tokens) */

    it("Can transfer tokens from caller to a recipient", async function() {
      let res = await cosmosToken.transfer(args._default, 50, {from: args._account_one});
      assert.strictEqual(res.logs.length, 1, "Successful transfer should have logged Transfer event");
      assert.strictEqual(res.logs[0].args.from, args._account_one, "'from' address parameter from Approval event should be equal to account 1");
      assert.strictEqual(res.logs[0].args.to, args._default, "'to' address parameter from Approval event should be equal to account 0");
      assert.strictEqual(res.logs[0].args.tokens.toNumber(), 50, "'tokens' uint parameter from Burn event should be equal to 50");
      let balanceSender = await cosmosToken.balanceOf.call(args._account_one);
      assert.strictEqual(balanceSender.toNumber(), 950, "Sender's balance should decrease by 50");
      let balanceRecipient = await cosmosToken.balanceOf.call(args._default);
      assert.strictEqual(balanceRecipient.toNumber(), 50, "Recipient's balance should increase by 50");
    });

  });
  });
});
