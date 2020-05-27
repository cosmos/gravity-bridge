import { Peggy } from "./typechain/Peggy";
import { ethers } from "ethers";
import fs from "fs";
import commandLineArgs from "command-line-args";

const args = commandLineArgs([
  { name: "eth-node", type: String },
  { name: "cosmos-node", type: String },
  { name: "eth-privkey", type: String },
  { name: "contract", type: String },
  { name: "peggy-id", type: String }
]) as {
  ethNode: string;
  cosmosNode: string;
  ethPrivkey: string;
  peggyId: string;
  contract: string;
};

// 4. Now, the deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
//     - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
// 5. The deployer submits the address of the peggy contract that it deployed to Ethereum.
//     - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.

type Valset = {
  addresses: string[];
  powers: string[];
  currentValsetNonce: number;
  r: string[];
  s: string[];
  v: number[];
  erc20: string;
  powerThreshold: number;
};
async function deploy() {
  const provider = await new ethers.providers.JsonRpcProvider(args.ethNode);
  let wallet = new ethers.Wallet(args.ethPrivkey, provider);

  const { abi, bytecode } = getContractArtifacts();
  const factory = new ethers.ContractFactory(abi, bytecode, wallet);

  const latestValset = await getLatestValset(args.peggyId);

  const peggy = (await factory.deploy(
    latestValset.erc20,
    args.peggyId,
    latestValset.powerThreshold,
    latestValset.addresses,
    latestValset.powers,
    latestValset.v,
    latestValset.r,
    latestValset.s
  )) as Peggy;

  await peggy.deployed();
  await submitPeggyAddress(peggy.address);
}

function getContractArtifacts(): { bytecode: string; abi: string } {
  var { bytecode, abi } = JSON.parse(
    fs.readFileSync(args.contract, "utf8").toString()
  );
  return { bytecode, abi };
}
async function getLatestValset(peggyId: string): Promise<Valset> {
  return {
    addresses: [],
    powers: [],
    currentValsetNonce: 0,
    r: [],
    s: [],
    v: [],
    erc20: "0xERC20...",
    powerThreshold: 66666
  };
}

async function submitPeggyAddress(address: string) {}
