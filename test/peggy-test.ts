import chai from "chai";
import { ethers } from "@nomiclabs/buidler";
import {
  deployContract,
  getWallets,
  MockProvider,
  solidity
} from "ethereum-waffle";

// import GreeterArtifact from "../artifacts/Greeter.json";
import { Peggy } from "../typechain/Peggy";
import { BitcoinZirconium } from "../typechain/BitcoinZirconium";
import { Greeter } from "../typechain/Greeter";

chai.use(solidity);
const { expect } = chai;
const provider = new MockProvider({
  accounts: [
    {
      balance: 999999999999999999999999999999,
      secretKey:
        "0xc5e8f61d1ab959b397eecc0a37a6517b8e67a0e7cf1f4bce5591f3ed80199122"
    },
    {
      balance: 999999999999999999999999999999,
      secretKey:
        "0xd49743deccbccc5dc7baa8e69e5be03298da8688a15dd202e20f15d5e0e9a9fb"
    },
    {
      balance: 999999999999999999999999999999,
      secretKey:
        "0x23c601ae397441f3ef6f1075dcb0031ff17fb079837beadaf3c84d96c6f3e569"
    }
  ]
});
const [wallet, otherWallet] = provider.getWallets();

describe("Peggy", function() {
  it("Smoke test", async function() {
    console.log(wallet.address);

    const greeter = (await (await ethers.getContractFactory("Greeter")).deploy(
      "Fuuuck"
    )) as Greeter;

    await greeter.deployed();

    await greeter.setGreeting("jimbo");

    await greeter.setGreeting("jumbo");

    console.log(
      "afeaf",
      await greeter.connect(otherWallet).signer.getAddress()
    );

    await greeter.connect(otherWallet).setGreeting("gumbo");

    // const zrc = (await (
    //   await ethers.getContractFactory("BitcoinZirconium")
    // ).deploy()) as BitcoinZirconium;

    // await zrc.deployed();

    // const tx = await zrc.connect(wallet).transfer(otherWallet.address, 9);

    // console.log(tx);

    // console.log((await zrc.balanceOf(otherWallet.address)).toString());
    // console.log((await zrc.totalSupply()).toString());

    // const peggy = (await (await ethers.getContractFactory("Peggy")).deploy(
    //   "Hello, world!"
    // )) as Peggy;

    // await peggy.deployed();

    // expect(
    //   await peggy.functions.checkCheckpoint(
    //     [
    //       "0xc783df8a850f42e7f7e57013759c285caa701eb6",
    //       "0xead9c93b79ae7c1591b1fb5323bd777e86e150d4",
    //       "0xe5904695748fe4a84b40b3fc79de2277660bd1d3"
    //     ],
    //     [20, 40, 40]
    //   )
    // ).to.equal("Hello, world!");

    // expect(await peggy.functions.greet()).to.equal("Hello, world!");

    // await peggy.functions.setGreeting("Hola, mundo!");
    // expect(await peggy.functions.greet()).to.equal("Hola, mundo!");
  });
});
