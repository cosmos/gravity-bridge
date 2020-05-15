import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";

import { Greeter } from "../typechain/Greeter";
import { Peggy } from "../typechain/Peggy";
import { BitcoinMAX } from "../typechain/BitcoinMAX";
import { SigningTest } from "../typechain/SigningTest";
import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;

describe("Peggy happy path", function() {
  it.only("Happy path", async function() {
    const signers = await ethers.getSigners();
    const peggyId = ethers.utils.formatBytes32String("foo");
    const validators = [signers[1], signers[2], signers[3]];
    const powers = [60000, 20000, 20000];
    const powerThreshold = 66666;

    const { peggy, max, checkpoint: deployCheckpoint } = await deployContracts(
      peggyId,
      validators,
      powers,
      powerThreshold
    );

    expect(await peggy.functions.peggyId()).to.equal(peggyId);
    expect(await peggy.functions.powerThreshold()).to.equal(powerThreshold);
    expect(await peggy.functions.tokenContract()).to.equal(max.address);
    expect(await peggy.functions.lastCheckpoint()).to.equal(deployCheckpoint);

    const newValidators = [signers[1], signers[2], signers[3], signers[4]];
    const newPowers = [50000, 20000, 20000, 10000];
    const currentValsetNonce = 0;
    const newValsetNonce = 1;

    const checkpoint = makeCheckpoint(
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      peggyId
    );

    console.log("checkpoint calling updateValset", checkpoint);

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

    const txAmounts = [10, 20, 30];
    const txDestinations = await getSignerAddresses([
      signers[6],
      signers[7],
      signers[8]
    ]);
    const txFees = [1, 1, 1];
    const txNonces = [0, 1, 2];

    let txHash = makeTxBatchHash(
      txAmounts,
      txDestinations,
      txFees,
      txNonces,
      peggyId
    );

    sigs = await signHash(newValidators, txHash);

    // await peggy.submitBatch(
    //   await getSignerAddresses(newValidators),
    //   newPowers,
    //   currentValsetNonce,
    //   sigs.v,
    //   sigs.r,
    //   sigs.s,
    //   txAmounts,
    //   txDestinations,
    //   txFees,
    //   txNonces
    // );
  });
});
