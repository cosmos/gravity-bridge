import { Gravity } from "./typechain/Gravity";
import { ethers } from "ethers";
import fs from "fs";

function getContractArtifacts(path: string): { bytecode: string; abi: string } {
    var { bytecode, abi } = JSON.parse(fs.readFileSync(path, "utf8").toString());
    return { bytecode, abi };
}


(async function () {
    const ethNode = "http://localhost:8545"
    const contract = "0xD7600ae27C99988A6CD360234062b540F88ECA43"

    const provider = new ethers.providers.JsonRpcProvider(ethNode);

    const { abi } = getContractArtifacts("artifacts/contracts/Gravity.sol/Gravity.json");

    const gravity = (new ethers.Contract(contract, abi, provider.getSigner()) as any) as Gravity;

    const events = await gravity.queryFilter({})
    console.log(events)

    // gravity.on({}, function () {
    //     console.log(arguments)
    // });
})()