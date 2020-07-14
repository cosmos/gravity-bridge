import { ethers } from "@nomiclabs/buidler";
import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";

export async function getSignerAddresses(signers: Signer[]) {
  return await Promise.all(signers.map(signer => signer.getAddress()));
}

export function makeCheckpoint(
  validators: string[],
  powers: BigNumberish[],
  valsetNonce: BigNumberish,
  peggyId: string
) {
  const methodName = ethers.utils.formatBytes32String("checkpoint");

  let checkpoint = ethers.utils.solidityKeccak256(
    ["bytes32", "bytes32", "uint256", "address[]", "uint256[]"],
    [peggyId, methodName, valsetNonce, validators, powers]
  );

  return checkpoint;
}

export async function signHash(signers: Signer[], hash: string) {
  let v: number[] = [];
  let r: string[] = [];
  let s: string[] = [];

  for (let i = 0; i < signers.length; i = i + 1) {
    const sig = await signers[i].signMessage(ethers.utils.arrayify(hash));
    const address = await signers[i].getAddress();

    const splitSig = ethers.utils.splitSignature(sig);
    v.push(splitSig.v!);
    r.push(splitSig.r);
    s.push(splitSig.s);
  }

  return { v, r, s };
}

// bytes32 methodName = 0x7472616e73616374696f6e426174636800000000000000000000000000000000;
// bytes32 transactionsHash = keccak256(
//   abi.encodePacked(peggyId, methodName, _amounts, _destinations, _fees, _nonces)
// );
export function makeTxBatchHash(
  amounts: number[],
  destinations: string[],
  fees: number[],
  nonces: number[],
  peggyId: string
) {
  const methodName = ethers.utils.formatBytes32String("transactionBatch");

  let txHash = ethers.utils.solidityKeccak256(
    ["bytes32", "bytes32", "uint256[]", "address[]", "uint256[]", "uint256[]"],
    [peggyId, methodName, amounts, destinations, fees, nonces]
  );

  return txHash;
}
