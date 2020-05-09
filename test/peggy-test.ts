import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { solidity } from "ethereum-waffle";

import { Greeter } from "../typechain/Greeter";
import { BitcoinMAX } from "../typechain/BitcoinMAX";

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
  });
});
