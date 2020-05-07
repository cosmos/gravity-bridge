import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { deployContract, getWallets, solidity } from "ethereum-waffle";

// import GreeterArtifact from "../artifacts/Greeter.json";
import { Greeter } from "../typechain/Greeter";

chai.use(solidity);
const { expect } = chai;

describe("Greeter", function() {
  let greeter: Greeter;
  it("Should return the new greeting once it's changed", async function() {
    const Greeter = await ethers.getContractFactory("Greeter");
    greeter = (await Greeter.deploy("Hello, world!")) as Greeter;

    await greeter.deployed();
    expect(await greeter.functions.greet()).to.equal("Hello, world!");

    await greeter.functions.setGreeting("Hola, mundo!");
    expect(await greeter.functions.greet()).to.equal("Hola, mundo!");
  });
});
