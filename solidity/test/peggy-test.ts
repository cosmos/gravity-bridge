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
  signHash
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

    const { peggy, max } = await deployContracts(
      peggyId,
      validators,
      powers,
      powerThreshold
    );

    expect(await peggy.functions.peggyId()).to.equal(peggyId);
    expect(await peggy.functions.powerThreshold()).to.equal(powerThreshold);
    expect(await peggy.functions.tokenContract()).to.equal(max.address);

    const newValidators = [signers[1], signers[2], signers[3], signers[4]];
    const newPowers = [50000, 20000, 20000, 10000];
    const currentValsetNonce = 0;
    const newValsetNonce = 1;

    const checkpoint = makeCheckpoint(
      await getSignerAddresses(validators),
      powers,
      0,
      peggyId
    );

    const { v, r, s } = await signHash(validators, checkpoint);

    await peggy.updateValset(
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      await getSignerAddresses(validators),
      powers,
      currentValsetNonce,
      v,
      r,
      s
    );
  });
});
