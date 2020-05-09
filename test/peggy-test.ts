import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";

import { Greeter } from "../typechain/Greeter";
import { Peggy } from "../typechain/Peggy";
import { BitcoinMAX } from "../typechain/BitcoinMAX";
import { BigNumberish } from "ethers/utils";

chai.use(solidity);
const { expect } = chai;

describe("Test", function() {
  it("Smoke test", async function() {
    const Greeter = await ethers.getContractFactory("Greeter");
    const greeter = (await Greeter.deploy("Hello, world!")) as Greeter;

    const signers = await ethers.getSigners();

    await greeter.deployed();

    await greeter.setGreeting("one");

    await greeter.setGreeting("two");

    await greeter.connect(signers[0]).setGreeting("three");

    await greeter.connect(signers[1]).setGreeting("four");

    expect(await greeter.greet()).to.equal("four");
  });

  it("Coin test", async function() {
    const BitcoinMAX = await ethers.getContractFactory("BitcoinMAX");
    const max = (await BitcoinMAX.deploy()) as BitcoinMAX;

    const signers = await ethers.getSigners();

    await max.deployed();

    await max.connect(signers[1]).transfer(await signers[2].getAddress(), 99);

    console.log(
      (await max.balanceOf(await signers[2].getAddress())).toString()
    );

    console.log(
      (await max.balanceOf(await signers[1].getAddress())).toString()
    );

    await max
      .connect(signers[1])
      .transfer(await signers[2].getAddress(), 9999999999);
  });

  it("Peggy test", async function() {
    const signers = await ethers.getSigners();

    const BitcoinMAX = await ethers.getContractFactory("BitcoinMAX");
    const max = (await BitcoinMAX.deploy()) as BitcoinMAX;

    const Peggy = await ethers.getContractFactory("Peggy");

    const peggy = (await Peggy.deploy(
      max.address,
      ethers.utils.formatBytes32String("foo"),
      66666,
      [signers[1], signers[2], signers[3]],
      [60, 20, 20]
    )) as Peggy;
  });
});

// function makeCheckpoint(
//   address[] memory _newValidators,
//   uint256[] memory _newPowers,
//   uint256 _newValsetNonce
// ) public view returns (bytes32) {

// bytes32 encoding of "checkpoint"
// bytes32 methodName = 0x636865636b706f696e7400000000000000000000000000000000000000000000;
// bytes32 newCheckpoint = keccak256(abi.encodePacked(peggyId, methodName, _newValsetNonce));

// {
// 	for (uint256 i = 0; i < _newValidators.length; i = i.add(1)) {
// 		// - Check that validator powers are decreasing or equal (this allows the next
// 		//   caller to break out of signature evaluation ASAP to save more gas)
// 		if (i != 0) {
// 			require(
// 				!(_newPowers[i] > _newPowers[i - 1]),
// 				"Validator power must not be higher than previous validator in batch"
// 			);
// 		}
// 		newCheckpoint = keccak256(
// 			abi.encodePacked(newCheckpoint, _newValidators[i], _newPowers[i])
// 		);
// 	}
// }

// return newCheckpoint;

function makeCheckpoint(
  newValidators: string[],
  newPowers: BigNumberish[],
  newValsetNonce: BigNumberish,
  peggyId: string
) {
  const methodName = ethers.utils.formatBytes32String("checkpoint");

  let checkpoint = ethers.utils.solidityKeccak256(
    ["bytes32", "bytes32", "uint256"],
    [peggyId, methodName, newValsetNonce]
  );

  for (let i = 0; i < newValidators.length; i = i + 1) {
    checkpoint = ethers.utils.solidityKeccak256(
      ["bytes32", "address", "uint256"],
      [checkpoint, newValidators[i], newPowers[i]]
    );
  }

  return checkpoint;
}
