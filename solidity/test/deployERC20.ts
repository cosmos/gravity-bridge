import chai from "chai";
import { ethers } from "hardhat";
import { solidity } from "ethereum-waffle";

import { deployContracts } from "../test-utils";
import {
  getSignerAddresses,
  makeCheckpoint,
  signHash,
  makeTxBatchHash,
  examplePowers
} from "../test-utils/pure";
import {ContractTransaction, utils} from 'ethers';
import { BigNumber } from "ethers";
chai.use(solidity);
const { expect } = chai;

async function parseEvent(contract: any, txPromise: Promise<ContractTransaction>, eventOrder: number) {
  const tx = await txPromise
  const receipt = await contract.provider.getTransactionReceipt(tx.hash!)
  let args = (contract.interface as utils.Interface).parseLog(receipt.logs![eventOrder]).args

  // Get rid of weird quasi-array keys
  const acc: any = {}
  args = Object.keys(args).reduce((acc, key) => {
    if (Number.isNaN(parseInt(key, 10)) && key !== 'length') {
      acc[key] = args[key]
    }
    return acc
  }, acc)

  return args
}

async function runTest(opts: {}) {



  // Prep and deploy Gravity contract
  // ========================
  const signers = await ethers.getSigners();
  const gravityId = ethers.utils.formatBytes32String("foo");
  // This is the power distribution on the Cosmos hub as of 7/14/2020
  let powers = examplePowers();
  let validators = signers.slice(0, powers.length);
  const powerThreshold = 6666;
  const {
    gravity,
    testERC20,
    checkpoint: deployCheckpoint
  } = await deployContracts(gravityId, validators, powers, powerThreshold);




  // Deploy ERC20 contract representing Cosmos asset
  // ===============================================
  const eventArgs = await parseEvent(gravity, gravity.deployERC20('uatom', 'Atom', 'ATOM', 6), 1)

  expect(eventArgs).to.deep.equal({
    _cosmosDenom: 'uatom',
    _tokenContract: eventArgs._tokenContract, // We don't know this ahead of time
    _name: 'Atom',
    _symbol: 'ATOM',
    _decimals: 6,
    _eventNonce: BigNumber.from(2)
  })




  // Connect to deployed contract for testing
  // ========================================
  let ERC20contract = new ethers.Contract(eventArgs._tokenContract, [
    "function balanceOf(address account) view returns (uint256 balance)"
  ], gravity.provider);


  const maxUint256 = BigNumber.from(2).pow(256).sub(1)

  // Check that gravity balance is correct
  expect((await ERC20contract.functions.balanceOf(gravity.address)).toString()).to.equal(maxUint256.toString())


  // Prepare batch
  // ===============================
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
  let batchNonce = 1
  let batchTimeout = 10000




  // Call method
  // ===========
  const methodName = ethers.utils.formatBytes32String(
    "transactionBatch"
  );
  let abiEncoded = ethers.utils.defaultAbiCoder.encode(
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
      gravityId,
      methodName,
      txAmounts,
      txDestinations,
      txFees,
      batchNonce,
      eventArgs._tokenContract,
      batchTimeout
    ]
  );
  let digest = ethers.utils.keccak256(abiEncoded);
  let sigs = await signHash(validators, digest);
  let currentValsetNonce = 0;

  await gravity.submitBatch(
    await getSignerAddresses(validators),
    powers,
    currentValsetNonce,

    sigs.v,
    sigs.r,
    sigs.s,

    txAmounts,
    txDestinations,
    txFees,
    batchNonce,
    eventArgs._tokenContract,
    batchTimeout
  );

  // Check that Gravity's balance is correct
  expect((await ERC20contract.functions.balanceOf(gravity.address)).toString()).to.equal(maxUint256.sub(200).toString())

  // Check that one of the recipient's balance is correct
  expect((await ERC20contract.functions.balanceOf(await signers[6].getAddress())).toString()).to.equal('1')
}

describe("deployERC20 tests", function () {
  // There is no way for this function to throw so there are
  // no throwing tests
  it("runs", async function () {
    await runTest({})
  });
});
