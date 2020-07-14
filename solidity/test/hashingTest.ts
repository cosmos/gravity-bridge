import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";
import { HashingTest } from "../typechain/HashingTest";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;

describe.only("Hashing test", function() {
  it("Hashing test", async function() {
    const signers = await ethers.getSigners();
    const peggyId = ethers.utils.formatBytes32String("foo");

    let validators = [];
    let powers = [];

    for (let i = 0; i < 100; i++) {
      validators.push(signers[i]);
      powers.push(5000);
    }

    const HashingTest = await ethers.getContractFactory("HashingTest");

    const hashingContract = (await HashingTest.deploy()) as HashingTest;

    await hashingContract.deployed();

    await hashingContract.IterativeHash(
      await getSignerAddresses(validators),
      powers,
      1,
      peggyId
    );

    await hashingContract.ConcatHash(
      await getSignerAddresses(validators),
      powers,
      1,
      peggyId
    );

    await hashingContract.ConcatHash2(
      await getSignerAddresses(validators),
      powers,
      1,
      peggyId
    );
  });
});
