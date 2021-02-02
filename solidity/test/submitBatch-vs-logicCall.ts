import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";
import { TestTokenBatchMiddleware } from "../typechain/TestTokenBatchMiddleware";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  examplePowers
} from "../test-utils/pure";
import { Signer } from "ethers";
import { Peggy } from "../typechain/Peggy";
import { TestERC20A } from "../typechain/TestERC20A";

chai.use(solidity);
const { expect } = chai;

async function prepareTxBatch(batchSize: number, signers: Signer[]) {
    const numTxs = batchSize;
    const destinations = new Array(numTxs);
    const fees = new Array(numTxs);
    const amounts = new Array(numTxs);
    for (let i = 0; i < numTxs; i++) {
      fees[i] = 1;
      amounts[i] = 1;
      destinations[i] = await signers[i + 5].getAddress();
    }

    return {
      numTxs,
      destinations,
      fees,
      amounts
    }
}

async function sendToCosmos(peggy: Peggy, testERC20: TestERC20A, numCoins: number) {
    // Transfer out to Cosmos, locking coins
    // =====================================
    await testERC20.functions.approve(peggy.address, numCoins);
    await peggy.functions.sendToCosmos(
      testERC20.address,
      ethers.utils.formatBytes32String("myCosmosAddress"),
      numCoins
    );
}

async function prep() {
    // Deploy contracts
    // ================

    const signers = await ethers.getSigners();
    const peggyId = ethers.utils.formatBytes32String("foo");

    let powers = examplePowers();
    let validators = signers.slice(0, powers.length);

    const powerThreshold = 6666;

    const {
      peggy,
      testERC20,
    } = await deployContracts(peggyId, validators, powers, powerThreshold);

    return {
      signers,
      peggyId,
      powers,
      validators,
      peggy,
      testERC20,
    }
}

async function runSubmitBatchTest(opts: {
  batchSize: number
}) {
    const {
      signers,
      peggyId,
      powers,
      validators,
      peggy,
      testERC20,
    } = await prep()



    // Lock tokens in peggy
    // ====================
    await sendToCosmos(peggy, testERC20, 1000)

    expect(
      (await testERC20.functions.balanceOf(peggy.address))[0].toNumber(),
      "peggy does not have correct balance after sendToCosmos"
    ).to.equal(1000);

    expect(
      (await testERC20.functions.balanceOf(await signers[0].getAddress()))[0].toNumber(),
      "msg.sender does not have correct balance after sendToCosmos"
    ).to.equal(9000);




    // Prepare tx batch
    // ================
    const txBatch = await prepareTxBatch(opts.batchSize, signers)
    const batchNonce = 1
    const batchTimeout = 10000




    // Using submitBatch method
    // ========================
    let methodName = ethers.utils.formatBytes32String(
      "transactionBatch"
    );

    let digest = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
        [
          "bytes32",
          "bytes32",
          "uint256[]",
          "address[]",
          "uint256[]",
          "uint256",
          "address",
          "uint256"
        ],
        [
          peggyId,
          methodName,
          txBatch.amounts,
          txBatch.destinations,
          txBatch.fees,
          batchNonce,
          testERC20.address,
          batchTimeout
        ]
      ));

    let sigs = await signHash(validators, digest);

    await peggy.submitBatch(
      await getSignerAddresses(validators),
      powers,
      0,

      sigs.v,
      sigs.r,
      sigs.s,

      txBatch.amounts,
      txBatch.destinations,
      txBatch.fees,
      1,
      testERC20.address,
      batchTimeout
    );

    expect(
      (await testERC20.functions.balanceOf(await signers[5].getAddress()))[0].toNumber(),
      "first address in tx batch does not have correct balance after submitBatch"
    ).to.equal(1);

    expect(
      (await testERC20.functions.balanceOf(await signers[5 + txBatch.numTxs - 1].getAddress()))[0].toNumber(),
      "last address in tx batch does not have correct balance after submitBatch"
    ).to.equal(1);

    expect(
      (await testERC20.functions.balanceOf(peggy.address))[0].toNumber(),
      "peggy does not have correct balance after submitBatch"
      // Each tx in batch is worth 1 coin sent + 1 coin fee
    ).to.equal(1000 - (txBatch.numTxs * 2));

    expect(
      (await testERC20.functions.balanceOf(await signers[0].getAddress()))[0].toNumber(),
      "msg.sender does not have correct balance after submitBatch"
      // msg.sender has received 1 coin in fees for each tx
    ).to.equal(9000 + txBatch.numTxs);
}

async function runLogicCallTest(opts: {
  batchSize: number
}) {
    const {
      signers,
      peggyId,
      powers,
      validators,
      peggy,
      testERC20,
    } = await prep()

    const TestTokenBatchMiddleware = await ethers.getContractFactory("TestTokenBatchMiddleware");
    const tokenBatchMiddleware = (await TestTokenBatchMiddleware.deploy()) as TestTokenBatchMiddleware;
    await tokenBatchMiddleware.transferOwnership(peggy.address);


    // Lock tokens in peggy
    // ====================
    await sendToCosmos(peggy, testERC20, 1000)

    expect(
      (await testERC20.functions.balanceOf(peggy.address))[0].toNumber(),
      "peggy does not have correct balance after sendToCosmos"
    ).to.equal(1000);

    expect(
      (await testERC20.functions.balanceOf(await signers[0].getAddress()))[0].toNumber(),
      "msg.sender does not have correct balance after sendToCosmos"
    ).to.equal(9000);




    // Preparing tx batch
    // ===================================
    const txBatch = await prepareTxBatch(opts.batchSize, signers)
    const batchNonce = 1




    // Using logicCall method
    // ========================
    const methodName = ethers.utils.formatBytes32String(
        "logicCall"
      );

    let logicCallArgs = {
      transferAmounts: [txBatch.numTxs], // transferAmounts
      transferTokenContracts: [testERC20.address], // transferTokenContracts
      feeAmounts: [txBatch.numTxs], // feeAmounts
      feeTokenContracts: [testERC20.address], // feeTokenContracts
      logicContractAddress: tokenBatchMiddleware.address, // logicContractAddress
      payload: tokenBatchMiddleware.interface.encodeFunctionData("submitBatch",[txBatch.amounts, txBatch.destinations, testERC20.address]), // payload
      timeOut: 4766922941000, // timeOut, Far in the future
      invalidationId: ethers.utils.hexZeroPad(testERC20.address, 32), // invalidationId
      invalidationNonce: 1 // invalidationNonce
    }
  
    const digest = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
          [
            "bytes32", // peggyId
            "bytes32", // methodName
            "uint256[]", // transferAmounts
            "address[]", // transferTokenContracts
            "uint256[]", // feeAmounts
            "address[]", // feeTokenContracts
            "address", // logicContractAddress
            "bytes", // payload
            "uint256", // timeOut
            "bytes32", // invalidationId
            "uint256" // invalidationNonce
          ],
          [
            peggyId,
            methodName,
            logicCallArgs.transferAmounts,
            logicCallArgs.transferTokenContracts,
            logicCallArgs.feeAmounts,
            logicCallArgs.feeTokenContracts,
            logicCallArgs.logicContractAddress,
            logicCallArgs.payload,
            logicCallArgs.timeOut,
            logicCallArgs.invalidationId,
            logicCallArgs.invalidationNonce
          ]
        ));
  
      const sigs = await signHash(validators, digest);
  
      await peggy.submitLogicCall(
        await getSignerAddresses(validators),
        powers,
        0,
  
        sigs.v,
        sigs.r,
        sigs.s,
        logicCallArgs
      );

      expect(
        (await testERC20.functions.balanceOf(await signers[5].getAddress()))[0].toNumber(),
        "first address in tx batch does not have correct balance after submitLogicCall"
      ).to.equal(1);
  
      expect(
        (await testERC20.functions.balanceOf(await signers[5 + txBatch.numTxs - 1].getAddress()))[0].toNumber(),
        "last address in tx batch does not have correct balance after submitLogicCall"
      ).to.equal(1);
  
      expect(
        (await testERC20.functions.balanceOf(peggy.address))[0].toNumber(),
        "peggy does not have correct balance after submitLogicCall"
        // Each tx in batch is worth 1 coin sent + 1 coin fee
      ).to.equal(1000 - (txBatch.numTxs * 2));
  
      expect(
        (await testERC20.functions.balanceOf(await signers[0].getAddress()))[0].toNumber(),
        "msg.sender does not have correct balance after submitLogicCall"
        // msg.sender has received 1 coin in fees for each tx
      ).to.equal(9000 + txBatch.numTxs);
}

describe("Compare gas usage of old submitBatch method vs new logicCall method submitting one batch", function () {
  it("Large batch", async function () {
    await runSubmitBatchTest({batchSize: 10})
    await runLogicCallTest({batchSize: 10})
  });

  it("Small batch", async function () {
    await runSubmitBatchTest({batchSize: 1})
    await runLogicCallTest({batchSize: 1})
  });
});
