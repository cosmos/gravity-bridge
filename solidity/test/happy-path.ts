import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
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

describe("Peggy happy path", function() {
  it("Happy path", async function() {
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

    expect(await peggy.functions.peggyId()).to.equal(peggyId);
    expect(await peggy.functions.powerThreshold()).to.equal(powerThreshold);
    expect(await peggy.functions.tokenContract()).to.equal(testERC20.address);
    expect(await peggy.functions.lastCheckpoint()).to.equal(deployCheckpoint);

    let newPowers = examplePowers();
    newPowers[0] -= 3;
    newPowers[1] += 3;
    let newValidators = signers.slice(0, newPowers.length);

    const currentValsetNonce = 0;
    const newValsetNonce = 1;

    const checkpoint = makeCheckpoint(
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      peggyId
    );

    let sigs = await signHash(validators, checkpoint);

    await peggy.updateValset(
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      await getSignerAddresses(validators),
      powers,
      currentValsetNonce,
      sigs.v,
      sigs.r,
      sigs.s
    );

    expect(await peggy.functions.lastCheckpoint()).to.equal(checkpoint);

    // Transferring out to Cosmos

    await testERC20.functions.approve(peggy.address, 100);

    await peggy.functions.transferOut(
      ethers.utils.formatBytes32String("myCosmosAddress"),
      100
    );

    // Transferring into ERC20 from Cosmos
    const txAmounts = [11, 22, 33];
    const txDestinations = await getSignerAddresses([
      signers[6],
      signers[7],
      signers[8]
    ]);
    const txFees = [1, 1, 1];
    const txNonces = [1, 2, 3];

    let txHash = makeTxBatchHash(
      txAmounts,
      txDestinations,
      txFees,
      txNonces,
      peggyId
    );

    sigs = await signHash(newValidators, txHash);

    await peggy.submitBatch(
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      sigs.v,
      sigs.r,
      sigs.s,
      txAmounts,
      txDestinations,
      txFees,
      txNonces
    );

    expect(
      await (
        await testERC20.functions.balanceOf(await signers[6].getAddress())
      ).toNumber()
    ).to.equal(11);
  });
});

describe("Gas tests", function() {
  it("makeCheckpoint in isolation", async function() {
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

    await peggy.testMakeCheckpoint(
      await getSignerAddresses(validators),
      powers,
      0,
      peggyId
    );
  });

  it("checkValidatorSignatures in isolation", async function() {
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

    let sigs = await signHash(
      validators,
      "0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d"
    );

    await peggy.testCheckValidatorSignatures(
      await getSignerAddresses(validators),
      powers,
      sigs.v,
      sigs.r,
      sigs.s,
      "0x7bc422a00c175cae98cf2f4c36f2f8b63ec51ab8c57fecda9bccf0987ae2d67d",
      6666
    );
  });
});
