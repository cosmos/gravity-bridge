'use strict';
/* Add the dependencies you're testing */
const web3 = global.web3;
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");

contract('CosmosERC20', function(accounts) {
  const args = {
    _default: accounts[0],
    _other: accounts[1],
    _zero: 0,
    _initialSupply: 0,
    _amount: 1000
  };
  let cosmosToken;
  /* Do something before every `describe` method */
	beforeEach('Setup contract', async function() {
    cosmosToken = await CosmosERC20.new(args._default, 'Cosmos', {from: accounts[0]});
	});

  /* Functions */

  describe('Supply control', function () {
    let supply;
    beforeEach(async function() {
      supply = await cosmosToken.totalSupply.call();
    });

    it("Initial Supply is Correct", async function() {
			supply = await cosmosToken.totalSupply.call()
			assert.strictEqual(supply.toNumber(), args._initialSupply, "Initial supply must be 0");
		});

    /* approve(address spender, uint tokens) */
    // allowed[tokenOwner][spender]

  describe('', function() {
    let minted;
    beforeEach('Mint', async function() {
      minted = await cosmosToken.mint(args._other, args._amount, {from: accounts[0]});
    });

    /* mint(address to, uint tokens) */

    it("Can Mint Tokens", async function() {
      assert.isTrue(Boolean(minted.receipt.status), "Successful mint should return true");
			let totalSupply2 = await cosmosToken.totalSupply.call();
			assert.strictEqual(totalSupply2.toNumber(), args._amount, "Supply should increase in 1000");
      let balanceAfter = await cosmosToken.balanceOf.call(args._other);
      assert.strictEqual(balanceAfter.toNumber(), args._amount, "Controller's balance should increase in 1000");
		});

    /* burn(address from, uint tokens) */

    it("Can Burn tokens from a user's balance", async function() {
      // initial supply and balance = 1000
      let res = await cosmosToken.burn(args._other, 100, {from: accounts[0]});
      assert.isTrue(Boolean(res.receipt.status), "Successful burning should return true");
      let totalSupply3 = await cosmosToken.totalSupply.call();
      assert.strictEqual(totalSupply3.toNumber(), 900, "Supply should decrease in 100");
      let balanceAfter = await cosmosToken.balanceOf.call(args._other);
      assert.strictEqual(balanceAfter.toNumber(), 900, "Controller's balance should decrease by 100");
    });


    /* transfer(address to, uint tokens) */

    it("Can transfer tokens from caller to a recipient", async function() {
      let res = await cosmosToken.transfer(args._default, 50, {from: accounts[1]});
      assert.isTrue(Boolean(res.receipt.status), "Successful transfer should return true");
      let balanceSender = await cosmosToken.balanceOf.call(args._other);
      assert.strictEqual(balanceSender.toNumber(), 950, "Sender's balance should decrease by 50");
      let balanceRecipient = await cosmosToken.balanceOf.call(args._default);
      assert.strictEqual(balanceRecipient.toNumber(), 50, "Recipient's balance should increase by 50");
    });

    describe('', function() {
      let approved;
      beforeEach('Mint', async function() {
        approved = await cosmosToken.approve(accounts[2], args._amount, {from: accounts[1]});
      });

      it("Can Approve a certain amount to be spend by an user", async function() {
        assert.isTrue(Boolean(approved.receipt.status), "Successful approval should always return true");
        let amountAllowed = await cosmosToken.allowance.call(accounts[1], accounts[2]);
        assert.strictEqual(Number(amountAllowed.toNumber()), args._amount, "Approved amount should be the same as the user allowed balance");
      });

      /* transferFrom(address from, address to, uint tokens) */

      it("Can transfer tokens from one account to another", async function() {
        let res = await cosmosToken.transferFrom(accounts[1], args._default, 100, {from: accounts[2]});
        assert.isTrue(Boolean(res.receipt.status), "Successful transfer should return true");
        let balanceSender = await cosmosToken.balanceOf.call(args._other);
        assert.strictEqual(balanceSender.toNumber(), 900, "Sender's balance should decrease by 100");
        let balanceRecipient = await cosmosToken.balanceOf.call(args._default);
        assert.strictEqual(balanceRecipient.toNumber(), 100, "Recipient's balance should increase by 100");
      });
    });

  });


  });
});
