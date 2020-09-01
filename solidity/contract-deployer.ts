import { Peggy } from "./typechain/Peggy";
import { TestERC20 } from "./typechain/TestERC20";
import { ethers } from "ethers";
import fs from "fs";
import commandLineArgs from "command-line-args";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";

const args = commandLineArgs([
  // the ethernum node used to deploy the contract
  { name: "eth-node", type: String },
  // the cosmos node that will be used to grab the validator set via RPC (TODO),
  { name: "cosmos-node", type: String },
  // the Ethereum private key that will contain the gas required to pay for the contact deployment
  { name: "eth-privkey", type: String },
  // the peggy contract .json file
  { name: "contract", type: String },
  // the peggy contract erc20 address for the hardcoded erc20 version, only used if test mode is not on
  { name: "erc20-address", type: String },
  // the id to be used for this version of peggy, be sure to avoid conflicts in production
  { name: "peggy-id", type: String },
  // test mode, if enabled this script deploys an erc20 script and uses that script as the contract erc20
  { name: "test-mode", type: String },
  // if test mode is enabled this contract is deployed and it's address is used as the erc20 address in the contract
  { name: "erc20-contract", type: String }
]);

// 4. Now, the deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
//     - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
// 5. The deployer submits the address of the peggy contract that it deployed to Ethereum.
//     - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.

type Valset = {
  EthAddresses: string[];
  Powers: string[];
  Nonce: number;
};
async function deploy() {
  const provider = await new ethers.providers.JsonRpcProvider(args["eth-node"]);
  let wallet = new ethers.Wallet(args["eth-privkey"], provider);
  let contract;

  if (Boolean(args["test-mode"])) {
    console.log("Test mode, deploying ERC20 contract");
    const { abi, bytecode } = getContractArtifacts(args["erc20-contract"]);
    const erc20Factory = new ethers.ContractFactory(abi, bytecode, wallet);

    const testERC20 = (await erc20Factory.deploy()) as TestERC20;

    await testERC20.deployed();
    const erc20TestAddress = testERC20.address;
    contract = erc20TestAddress;
    console.log("Successfully deployed ERC20");
  } else {
    contract = args["erc20-address"];
  }

  console.log("Starting Peggy contract deploy");
  const { abi, bytecode } = getContractArtifacts(args["contract"]);
  const factory = new ethers.ContractFactory(abi, bytecode, wallet);

  console.log("About to get latest Peggy valset");
  const latestValset = await getLatestValset(args.peggyId);

  console.log("Deploying peggy contract using valset ", latestValset);
  const peggy = (await factory.deploy(
    contract,
    // todo generate this randomly at deployment time that way we can avoid
    // anything but intentional conflicts
    "0x6c00000000000000000000000000000000000000000000000000000000000000",
    "6666",
    latestValset.EthAddresses,
    latestValset.Powers
  )) as Peggy;

  await peggy.deployed();
  console.log("Peggy deployed");
  await submitPeggyAddress(peggy.address);
}

function getContractArtifacts(path: string): { bytecode: string; abi: string } {
  var { bytecode, abi } = JSON.parse(fs.readFileSync(path, "utf8").toString());
  return { bytecode, abi };
}
async function getLatestValset(peggyId: string): Promise<Valset> {
  let request_string = args["cosmos-node"] + "/peggy/current_valset";
  let response = await axios.get(request_string);
  return response.data;
}

async function submitPeggyAddress(address: string) {}

async function main() {
  await deploy();
}

main();
