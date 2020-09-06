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

async function runTest(opts: {
  malformedNewValset?: boolean;
  malformedCurrentValset?: boolean;
  nonMatchingCurrentValset?: boolean;
  nonceNotIncremented?: boolean;
  badValidatorSig?: boolean;
}) {
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

  let newPowers = examplePowers();
  newPowers[0] -= 3;
  newPowers[1] += 3;

  let newValidators = signers.slice(0, newPowers.length);
  if (opts.malformedNewValset) {
    // Validators and powers array don't match
    newValidators = signers.slice(0, newPowers.length - 1);
  }

  let currentValsetNonce = 0;
  if (opts.nonMatchingCurrentValset) {
    currentValsetNonce = 420;
  }
  let newValsetNonce = 1;
  if (opts.nonceNotIncremented) {
    newValsetNonce = 420;
  }

  const checkpoint = makeCheckpoint(
    await getSignerAddresses(newValidators),
    newPowers,
    newValsetNonce,
    peggyId
  );

  let sigs = await signHash(validators, checkpoint);
  if (opts.badValidatorSig) {
    // Switch the first sig for the second sig to screw things up
    sigs.v[1] = sigs.v[0];
    sigs.r[1] = sigs.r[0];
    sigs.s[1] = sigs.s[0];
  }

  if (opts.malformedCurrentValset) {
    // Remove one of the powers to make the length not match
    powers.pop();
  }

  return peggy.updateValset(
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
}

describe("updateValset tests", function() {
  it("throws on malformed new valset", async function() {
    await expect(runTest({ malformedNewValset: true })).to.be.revertedWith(
      "Malformed new validator set"
    );
  });

  it("throws on malformed current valset", async function() {
    await expect(runTest({ malformedCurrentValset: true })).to.be.revertedWith(
      "Malformed current validator set"
    );
  });

  it("throws on non matching checkpoint for current valset", async function() {
    await expect(
      runTest({ nonMatchingCurrentValset: true })
    ).to.be.revertedWith(
      "Supplied current validators and powers do not match checkpoint"
    );
  });

  it("throws on new valset nonce not incremented", async function() {
    await expect(runTest({ nonceNotIncremented: true })).to.be.revertedWith(
      "Valset nonce must be incremented by one"
    );
  });

  it("throws on bad validator sig", async function() {
    await expect(runTest({ badValidatorSig: true })).to.be.revertedWith(
      "Validator signature does not match"
    );
  });
});
