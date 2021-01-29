import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash,
  examplePowers
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;


async function runTest(opts: {}) {


  // Prep and deploy contract
  // ========================
  const signers = await ethers.getSigners();
  const peggyId = ethers.utils.formatBytes32String("foo");
  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);
  const powerThreshold = 6666;
  const {
    peggy,
    testERC20,
    checkpoint: deployCheckpoint
  } = await deployContracts(peggyId, validators, powers, powerThreshold);


  // Transfer out to Cosmos, locking coins
  // =====================================
  await testERC20.functions.approve(peggy.address, 1000);
  await expect(peggy.functions.sendToCosmos(
    testERC20.address,
    ethers.utils.formatBytes32String("myCosmosAddress"),
    1000
  )).to.emit(peggy, 'SendToCosmosEvent').withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000, 
      1
    );

  expect((await testERC20.functions.balanceOf(peggy.address))[0]).to.equal(1000);
  expect((await peggy.functions.state_lastEventNonce())[0]).to.equal(1);


    
  // Do it again
  // =====================================
  await testERC20.functions.approve(peggy.address, 1000);
  await expect(peggy.functions.sendToCosmos(
    testERC20.address,
    ethers.utils.formatBytes32String("myCosmosAddress"),
    1000
  )).to.emit(peggy, 'SendToCosmosEvent').withArgs(
      testERC20.address,
      await signers[0].getAddress(),
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000, 
      2
    );

  expect((await testERC20.functions.balanceOf(peggy.address))[0]).to.equal(2000);
  expect((await peggy.functions.state_lastEventNonce())[0]).to.equal(2);
}

describe("sendToCosmos tests", function () {
  it("works right", async function () {
    await runTest({})
  });
});