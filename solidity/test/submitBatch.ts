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
  malformedCurrentValset?: boolean;
  malformedTxBatch?: boolean;
  nonMatchingCurrentValset?: boolean;
  nonceNotHigher?: boolean;
  nonceNotIncreasing?: boolean;
  badValidatorSig?: boolean;
  zeroedValidatorSig?: boolean;
  notEnoughPower?: boolean;
}) {
  const signers = await ethers.getSigners();
  const peggyId = ethers.utils.formatBytes32String("foo");

  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);

  if (opts.malformedCurrentValset) {
    validators = signers.slice(0, powers.length + 1);
  }

  const powerThreshold = 6666;

  const {
    peggy,
    testERC20,
    checkpoint: deployCheckpoint
  } = await deployContracts(peggyId, validators, powers, powerThreshold);

  // Transferring out to Cosmos

  await testERC20.functions.approve(peggy.address, 1000);

  await peggy.functions.transferOut(
    ethers.utils.formatBytes32String("myCosmosAddress"),
    1000
  );

  const txDestinationsInt = new Array(100);
  const txFees = new Array(100);
  const txNonces = new Array(100);
  const txAmounts = new Array(100);
  for (let i = 0; i < 100; i++) {
    txNonces[i] = i + 1;
    if (opts.nonceNotHigher) {
      txNonces[i] = i;
    }
    if (opts.nonceNotIncreasing) {
      txNonces[i] = 1;
    }
    txFees[i] = 1;
    txAmounts[i] = 1;
    txDestinationsInt[i] = signers[i + 5];
  }
  // Transferring into ERC20 from Cosmos
  const txDestinations = await getSignerAddresses(txDestinationsInt);

  if (opts.malformedTxBatch) {
    txFees.pop();
  }

  let currentValsetNonce = 0;
  if (opts.nonMatchingCurrentValset) {
    currentValsetNonce = 420;
  }

  let txHash = makeTxBatchHash(
    txAmounts,
    txDestinations,
    txFees,
    txNonces,
    peggyId
  );

  let sigs = await signHash(validators, txHash);
  if (opts.badValidatorSig) {
    // Switch the first sig for the second sig to screw things up
    sigs.v[1] = sigs.v[0];
    sigs.r[1] = sigs.r[0];
    sigs.s[1] = sigs.s[0];
  }

  if (opts.zeroedValidatorSig) {
    // Switch the first sig for the second sig to screw things up
    sigs.v[1] = sigs.v[0];
    sigs.r[1] = sigs.r[0];
    sigs.s[1] = sigs.s[0];
    // Then zero it out to skip evaluation
    sigs.v[1] = 0;
  }

  if (opts.notEnoughPower) {
    // zero out enough signatures that we dip below the threshold
    sigs.v[1] = 0;
    sigs.v[2] = 0;
    sigs.v[3] = 0;
    sigs.v[5] = 0;
    sigs.v[6] = 0;
    sigs.v[7] = 0;
    sigs.v[9] = 0;
    sigs.v[11] = 0;
    sigs.v[13] = 0;
  }

  await peggy.submitBatch(
    await getSignerAddresses(validators),
    powers,
    currentValsetNonce,
    sigs.v,
    sigs.r,
    sigs.s,
    txAmounts,
    txDestinations,
    txFees,
    txNonces
  );
}

describe("submitBatch tests", function() {
  it("throws on malformed current valset", async function() {
    await expect(runTest({ malformedCurrentValset: true })).to.be.revertedWith(
      "Malformed current validator set"
    );
  });

  it("throws on malformed txbatch", async function() {
    await expect(runTest({ malformedTxBatch: true })).to.be.revertedWith(
      "Malformed batch of transactions"
    );
  });

  it("throws on non matching checkpoint for current valset", async function() {
    await expect(
      runTest({ nonMatchingCurrentValset: true })
    ).to.be.revertedWith(
      "Supplied current validators and powers do not match checkpoint"
    );
  });

  it("throws on tx nonces not high enough", async function() {
    await expect(runTest({ nonceNotHigher: true })).to.be.revertedWith(
      "Transaction nonces in batch must be higher than last transaction nonce and strictly increasing"
    );
  });

  it("throws on tx nonces not strictly increasing", async function() {
    await expect(runTest({ nonceNotIncreasing: true })).to.be.revertedWith(
      "Transaction nonces in batch must be higher than last transaction nonce and strictly increasing"
    );
  });

  it("throws on bad validator sig", async function() {
    await expect(runTest({ badValidatorSig: true })).to.be.revertedWith(
      "Validator signature does not match"
    );
  });

  it("allows zeroed sig", async function() {
    await runTest({ zeroedValidatorSig: true });
  });

  it("throws on not enough signatures", async function() {
    await expect(runTest({ notEnoughPower: true })).to.be.revertedWith(
      "Submitted validator set signatures do not have enough power"
    );
  });
});
