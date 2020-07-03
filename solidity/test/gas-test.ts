// Questions:
// - How much fixed overhead is there in a call to updateValset or submitBatch?
// - Per-validator cost in updateValset
// - Per-tx cost in submitBatch
// - How much more gas does a transaction in submitBatch take above that used by a
//   regular ERC20 transfer?
// - Is iterative hashing the most efficient way to build the digest of valsets and
//   transaction batches?

import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash
} from "../test-utils/pure";

import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";
import { Peggy } from "../typechain/Peggy";

chai.use(solidity);
const { expect } = chai;

describe("Peggy gas tests", function() {
  it("updateValset fixed overhead", async function() {
    const signers = await ethers.getSigners();
    const peggyId = ethers.utils.formatBytes32String("foo");
    const validators = [signers[1], signers[2], signers[3]];
    const powers = [60000, 20000, 20000];
    const powerThreshold = 66666;

    const { peggy, testERC20, checkpoint: deployCheckpoint } = await deployContracts(
      peggyId,
      validators,
      powers,
      powerThreshold
    );

    await updateValset(peggy, peggyId, signers, validators, powers);
  });
});

async function updateValset(
  peggy: Peggy,
  peggyId: string,
  signers: Signer[],
  validators: Signer[],
  powers: number[]
) {
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

  let sigs = await signHash(validators, checkpoint);

  let tx = await peggy.updateValset(
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

  console.log("ESTIMATE GAS", ethers.provider.estimateGas(tx));
  console.log("GAS", tx.gasLimit.toNumber());
}
