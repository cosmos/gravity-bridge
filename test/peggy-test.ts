import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import { deployContract, getWallets, solidity } from "ethereum-waffle";

// import GreeterArtifact from "../artifacts/Greeter.json";
import { Peggy } from "../typechain/Peggy";

chai.use(solidity);
const { expect } = chai;

describe("Peggy", function() {
  it("Smoke test", async function() {
    const peggy = (await (await ethers.getContractFactory("Peggy")).deploy(
      "Hello, world!"
    )) as Peggy;

    await peggy.deployed();

    expect(
      await peggy.functions.checkCheckpoint(
        [
          "0xc783df8a850f42e7f7e57013759c285caa701eb6",
          "0xead9c93b79ae7c1591b1fb5323bd777e86e150d4",
          "0xe5904695748fe4a84b40b3fc79de2277660bd1d3",
        ],
        [20, 40, 40]
      )
    ).to.equal("Hello, world!");

    // expect(await peggy.functions.greet()).to.equal("Hello, world!");

    // await peggy.functions.setGreeting("Hola, mundo!");
    // expect(await peggy.functions.greet()).to.equal("Hola, mundo!");
  });
});
