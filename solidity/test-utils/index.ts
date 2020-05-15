import { Peggy } from "../typechain/Peggy";
import { BitcoinMAX } from "../typechain/BitcoinMAX";
import { ethers } from "@nomiclabs/buidler";
import { makeCheckpoint, signHash, getSignerAddresses } from "./pure";
import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";

export async function deployContracts(
  peggyId: string = "foo",
  validators: Signer[],
  powers: number[],
  powerThreshold: number
) {
  const BitcoinMAX = await ethers.getContractFactory("BitcoinMAX");
  const max = (await BitcoinMAX.deploy()) as BitcoinMAX;

  const Peggy = await ethers.getContractFactory("Peggy");

  const valAddresses = await getSignerAddresses(validators);

  const checkpoint = makeCheckpoint(valAddresses, powers, 0, peggyId);

  const theHash = ethers.utils.solidityKeccak256(
    ["bytes32", "address", "bytes32", "uint256"],
    [checkpoint, max.address, peggyId, powerThreshold]
  );

  const { v, r, s } = await signHash(validators, theHash);

  const peggy = (await Peggy.deploy(
    max.address,
    peggyId,
    powerThreshold,
    valAddresses,
    powers,
    v,
    r,
    s
  )) as Peggy;

  await peggy.deployed();

  return { peggy, max };
}
