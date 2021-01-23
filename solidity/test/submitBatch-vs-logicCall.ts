import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";
import { TestTokenBatchMiddleware } from "../typechain/TestTokenBatchMiddleware";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  examplePowers
} from "../test-utils/pure";

chai.use(solidity);
const { expect } = chai;


describe.only("Compare gas usage of old submitBatch method vs new logicCall method submitting one batch", function () {
  it("Full batch", async function () {

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
      checkpoint: deployCheckpoint
    } = await deployContracts(peggyId, validators, powers, powerThreshold);

    const TestTokenBatchMiddleware = await ethers.getContractFactory("TestTokenBatchMiddleware");
    const tokenBatchMiddleware = (await TestTokenBatchMiddleware.deploy()) as TestTokenBatchMiddleware;
    await tokenBatchMiddleware.transferOwnership(peggy.address);

    // Transfer out to Cosmos, locking coins
    // =====================================
    await testERC20.functions.approve(peggy.address, 1000);
    await peggy.functions.sendToCosmos(
      testERC20.address,
      ethers.utils.formatBytes32String("myCosmosAddress"),
      1000
    );




    // Preparing tx batch
    // ===================================
    const numTxs = 100;
    const txDestinationsInt = new Array(numTxs);
    const txFees = new Array(numTxs);
    const txAmounts = new Array(numTxs);
    for (let i = 0; i < numTxs; i++) {
      txFees[i] = 1;
      txAmounts[i] = 1;
      txDestinationsInt[i] = signers[i + 5];
    }

    const txDestinations = await getSignerAddresses(txDestinationsInt);

    const batchNonce = 1




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
          "address"
        ],
        [
          peggyId,
          methodName,
          txAmounts,
          txDestinations,
          txFees,
          batchNonce,
          testERC20.address
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

      txAmounts,
      txDestinations,
      txFees,
      1,
      testERC20.address
    );

    expect(
      (await testERC20.functions.balanceOf(await signers[6].getAddress())).toNumber()
    ).to.equal(1);


    // Using logicCall method
    // ========================
    methodName = ethers.utils.formatBytes32String(
        "logicCall"
      );

    let logicCallArgs = {
      transferAmounts: [100], // transferAmounts
      transferTokenContracts: [testERC20.address], // transferTokenContracts
      feeAmounts: [100], // feeAmounts
      feeTokenContracts: [testERC20.address], // feeTokenContracts
      logicContractAddress: tokenBatchMiddleware.address, // logicContractAddress
      payload: tokenBatchMiddleware.interface.functions.submitBatch.encode([txAmounts, txDestinations, testERC20.address]), // payload
      timeOut: 4766922941000, // timeOut, Far in the future
      invalidationId: ethers.utils.hexZeroPad(testERC20.address, 32), // invalidationId
      invalidationNonce: 1 // invalidationNonce
    }
  
    digest = ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(
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
  
      sigs = await signHash(validators, digest);
  
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
        (await testERC20.functions.balanceOf(await signers[6].getAddress())).toNumber()
      ).to.equal(2);
  });
});
