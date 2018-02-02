'use strict';

const Valset = artifacts.require("./../contracts/Valset.sol");

function sumArrayValues(total, uint64) {
  return total + uint64;
}

contract('Valset', function(accounts) {
  const args = {
    _default: accounts[0],
    _account_one: accounts[1],
    _account_two: accounts[2]
  };
  let valSet;
  let initialValidators = [
    "0xe81e0f466dc44478a4db02d21e10680bd794b549",
    "0x36e6068382b6c51e3861cef20fb9c1199c42fd5d",
    "0xf2de4bfb3919b9bfbce3122c992d2b6b6dd55f68",
    "0x6a6a13a4f861e3d728d6afea34b83fbc938d7135",
    "0x7b5749433eea79dff3b16317942261edf2bec622",
    "0xac650b4029cdbff2f272d068b340da8a47849250",
    "0x05df294534a201a5fbbfb75a4e2c337f0d3000c8",
    "0x66190eb0a5f1161729bcf1ba4b3631c75264e043",
    "0x8125648effea25d483412886741d0630f7693499",
    "0xb2c1bafa9419f03e08cffa9b86c3bfe8e3c068dc"
  ];
  let initialPowers = [ 9, 15, 9, 13, 19, 11, 16, 13, 11, 12 ];
  let newValidators = [
    "0x669ef6f50d14490584e98d6c8a2b29e55842da97",
    "0x81f2a757e418b484950a3f2d3b15b56a218fbca1",
    "0x99ebe5e9934e92eb8320d110415344ee6f3955d0",
    "0xd23d8393a7258bcea86d7c7a1abd299af8bc130b",
    "0xfd30d0acf42667cd5f893301dc23ee7c558a3e14",
    "0xdb231190506e958de59b9e24af9bdf31ae73eadb",
    "0x0f3918e784a94e1ed6fb13f9d8db1cf1793c80c4",
    "0x03542245d1b555bcec48f1692bf49c1e8640f186",
    "0xfca7d6dba53c9c29d8b0acb951fe87245229d56d",
    "0x7e82bb8859b766f716128291299f28828430f6c7"
  ];
  let initialPowers = [ 23, 11, 23, 13, 3, 6, 13, 12, 18, 17 ];

  beforeEach('Setup contract', async function() {
    valSet = await Valset.new(initialValidators, initialPowers, {from: args._default});
  });

  describe('Constructor function', function() {
    let addresses, powers, first_element, second_element, totalPower;

    // Proved by induction
    it("Saves initial validators' address in array", async function() {
      first_element = await valSet.getValidator.call(0);
      second_element = await valSet.getValidator.call(1);
      assert.isTrue(Boolean(String(first_element) && String(second_element)), "Initial validators' addresses array should be equal as the saved one");
    });

    // Proved by induction
    it("Saves initial validators' powers in array", async function() {
      first_element = await valSet.getPower.call(0);
      second_element = await valSet.getPower.call(1);
      assert.isTrue(Boolean(first_element.toNumber() && second_element.toNumber()), "Initial validators' powers array should be equal as the saved one");
    });

    it("Checks that addresses and powers arrays have the same length", async function() {
      addresses = await valSet.addresses;
      powers = await valSet.powers;
      let powersLength = powers.length;
      assert.lengthOf(addresses, powersLength, "Both initial arrays must have the same length");
    });

    it("Number of validator is below 100", async function() {
      addresses = await valSet.addresses;
      let valLength = addresses.length;
      assert.isAtMost(valLength, 100, "Validator set should not be larger than 100")
    });

    it("Sums totalPower correctly", async function() {
      totalPower = await valSet.getTotalPower.call();
      let accumulatedPower = initialPowers.reduce(sumArrayValues);
      assert.strictEqual(totalPower.toNumber(), accumulatedPower, "totalSum should the sum of each individual validator's power")
    });


  })



});
