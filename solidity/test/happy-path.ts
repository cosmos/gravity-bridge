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

    expect(await peggy.functions.state_peggyId()).to.equal(peggyId);
    expect(await peggy.functions.state_powerThreshold()).to.equal(
      powerThreshold
    );
    expect(await peggy.functions.state_tokenContract()).to.equal(
      testERC20.address
    );
    expect(await peggy.functions.state_lastCheckpoint()).to.equal(
      deployCheckpoint
    );

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

    expect(await peggy.functions.state_lastCheckpoint()).to.equal(checkpoint);

    // Transferring out to Cosmos

    await testERC20.functions.approve(peggy.address, 1000);

    await peggy.functions.transferOut(
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000
    );

    const numTxs = 100;
    const txDestinationsInt = new Array(numTxs);
    const txFees = new Array(numTxs);
    const txNonces = new Array(numTxs);
    const txAmounts = new Array(numTxs);
    for (let i = 0; i < numTxs; i++) {
      txNonces[i] = i + 1;
      txFees[i] = 1;
      txAmounts[i] = 1;
      txDestinationsInt[i] = signers[i + 5];
    }
    // Transferring into ERC20 from Cosmos
    const txDestinations = await getSignerAddresses(txDestinationsInt);

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
    ).to.equal(1);
  });
});

describe("Peggy happy path with combination method", function() {
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

    // Make new valset
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

    // Transfer out to Cosmos, locking coins
    await testERC20.functions.approve(peggy.address, 1000);
    await peggy.functions.transferOut(
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000
    );

    // Transferring into ERC20 from Cosmos
    const numTxs = 100;
    const txDestinationsInt = new Array(numTxs);
    const txFees = new Array(numTxs);
    const txNonces = new Array(numTxs);
    const txAmounts = new Array(numTxs);
    for (let i = 0; i < numTxs; i++) {
      txNonces[i] = i + 1;
      txFees[i] = 1;
      txAmounts[i] = 1;
      txDestinationsInt[i] = signers[i + 5];
    }

    const txDestinations = await getSignerAddresses(txDestinationsInt);

    const methodName = ethers.utils.formatBytes32String(
      "valsetAndTransactionBatch"
    );

    let abiEncoded = ethers.utils.defaultAbiCoder.encode(
      [
        "bytes32",
        "bytes32",
        "uint256[]",
        "address[]",
        "uint256[]",
        "uint256[]",
        "bytes32"
      ],
      [
        peggyId,
        methodName,
        txAmounts,
        txDestinations,
        txFees,
        txNonces,
        checkpoint
      ]
    );

    let digest = ethers.utils.keccak256(abiEncoded);

    let sigs = await signHash(validators, digest);

    await peggy.updateValsetAndSubmitBatch(
      await getSignerAddresses(validators),
      powers,
      currentValsetNonce,
      sigs.v,
      sigs.r,
      sigs.s,
      await getSignerAddresses(newValidators),
      newPowers,
      newValsetNonce,
      txAmounts,
      txDestinations,
      txFees,
      txNonces
    );

    expect(await peggy.functions.state_lastCheckpoint()).to.equal(checkpoint);

    expect(
      await (
        await testERC20.functions.balanceOf(await signers[6].getAddress())
      ).toNumber()
    ).to.equal(1);
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
