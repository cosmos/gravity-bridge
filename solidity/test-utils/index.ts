import { Peggy } from "../typechain/Peggy";
import { TestERC20A } from "../typechain/TestERC20A";
import { ethers } from "hardhat";
import { makeCheckpoint, signHash, getSignerAddresses } from "./pure";
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
  const TestERC20 = await ethers.getContractFactory("TestERC20A");
  const testERC20 = (await TestERC20.deploy()) as TestERC20A;

  const Peggy = await ethers.getContractFactory("Peggy");

  const valAddresses = await getSignerAddresses(validators);

  const checkpoint = makeCheckpoint(valAddresses, powers, 0, peggyId);

  const peggy = (await Peggy.deploy(
    peggyId,
    powerThreshold,
    valAddresses,
    powers
  )) as Peggy;

  await peggy.deployed();

  return { peggy, testERC20, checkpoint };
}
