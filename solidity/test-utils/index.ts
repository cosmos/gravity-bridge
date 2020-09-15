import { Peggy } from "../typechain/Peggy";
import { TestERC20 } from "../typechain/TestERC20";
import { ethers } from "@nomiclabs/buidler";
import { makeCheckpoint, signHash, getSignerAddresses } from "./pure";
import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";

type DeployContractsOptions = {
  corruptSig?: boolean;
};

export async function deployContracts(
  peggyId: string = "foo",
  validators: Signer[],
  powers: number[],
  powerThreshold: number,
  opts?: DeployContractsOptions
) {
  const TestERC20 = await ethers.getContractFactory("TestERC20");
  const testERC20 = (await TestERC20.deploy()) as TestERC20;

  const Peggy = await ethers.getContractFactory("Peggy");

  const valAddresses = await getSignerAddresses(validators);

  const checkpoint = makeCheckpoint(valAddresses, powers, 0, peggyId);

  const peggy = (await Peggy.deploy(
    testERC20.address,
    peggyId,
    powerThreshold,
    valAddresses,
    powers
  )) as Peggy;

  await peggy.deployed();

  return { peggy, testERC20, checkpoint };
}
