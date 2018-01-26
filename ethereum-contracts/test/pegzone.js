'use strict';

/* Add the dependencies you're testing */
const CosmosERC20 = artifacts.require("./../contracts/CosmosERC20.sol");
const Peggy = artifacts.require("./../contracts/Peggy.sol");
// const Valset = artifacts.require("./../contracts/Valset.sol");

// function stringToBytes(string) {
//   var bytes = [];
//   for (var i = 0; i < string.length; ++i)
//   {
//     bytes.push(string.charCodeAt(i));
//   }
//   return bytes;
// }


contract('CosmosERC20', function(accounts) {
  const args = {
    _default: accounts[0],
    _owner: accounts[1],
    _zero: 0,
    _amount: 1000
  };
  let cosmosToken;
  /* Do something before every `describe` method */
	beforeEach(async function() {
    cosmosToken = CosmosERC20.new(args.default, 'Cosmos');
	});

  describe('Deployment', function() {
    it("The contract can be deployed", function() {
  		return CosmosERC20.new()
  		.then(function(instance) {
  			assert.ok(instance.address);
  		});
  	});
  });

  // describe('Functions', function() {
  //   describe('totalsupply()', function {
  //     it("returns a number", function() {
  //       cosmosToken.totalSupply.call().then((supply) => {
  //         assert.isNumber(supply, 'return value should be of type number');
  //       });
  //     });
  //
  //     it("Supply is not zero", function() {
  //       cosmosToken.totalSupply.call().then((supply) => {
  //         return assert.isAtLeast(supply, args._zero, 'Supply must be ≥ 0');
  //       })
  //     });
  //
  //   });
  //   describe('transfer(address to, uint tokens)', function {
  //     it("returns a number", function() {
  //       cosmosToken.transfer.call(args._default {to: args._owner, args._amount }).then((supply) => {
  //         assert.isNumber(supply, 'return value should be of type number');
  //       });
  //     });
  //
  //
  //
  //   });
  // });
  //
  // it("Initial Supply is Correct", async function() {
  //   return cosmosToken.totalSupply.call().then(supply => {
  //     assert.equal(supply, 100);
  //   });
  // });
  // it("Can Mint Tokens", async function() {
  //   return cosmosToken.mint.call(100, {from: accounts[0]}).then(_ => {
  //     return cosmosToken.totalSupply.call().then(supply => {
  //       assert.equal(supply, 200);
  //     });
  //   });
  // });
  // it("Can Burn Tokens", async function() {
  //   return cosmosToken.burn.call(90, {from: accounts[0]}).then(() => {
  //     return cosmosToken.totalSupply.call().then(supply => {
  //       assert.equal(supply, 10);
  //     });
  //   });
  // });
  // it("Can retrieve the balance from an account", async function() {
	// 	return token.balanceOf.call(args._default).then(balance => {
  //     console.log(balance);
	// 		return assert.isAtLeast(balance, args._zero, 'balance ≥ 0');
	// 	});
  // });
  // it("Can transfer tokens from one account to another", async function() {
	// 	return token.balanceOf.call(args._default).then(balance => {
  //     console.log(balance);
	// 		return assert.isAtLeast(balance, args._zero, 'balance ≥ 0');
	// 	});
  // });

});

contract('Peggy', function(accounts) {
  const args = {_default: accounts[0], _owner: accounts[1],
  _null_address: '0x0000000000000000000000000000000000000000'};

  describe('Deployment', function() {
  	it("The contract can be deployed", function() {
  		return Peggy.new()
  		.then(function(instance) {
  			assert.ok(instance.address);
  		});
  	});
  });

});
