'use strict';

const Valset = artifacts.require("./../contracts/Valset.sol");

contract('Valset', function(accounts) {
  const args = {_default: accounts[0], _owner: accounts[1]};

  let valSet;

  beforeEach('Setup contract', async function() {
    valSet = await Valset.new();
	});

});
