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

export function examplePowers(): number[] {
  return [
    707,
    621,
    608,
    439,
    412,
    407,
    319,
    312,
    311,
    303,
    246,
    241,
    224,
    213,
    194,
    175,
    173,
    170,
    154,
    149,
    139,
    123,
    119,
    113,
    110,
    107,
    105,
    104,
    92,
    90,
    88,
    88,
    88,
    85,
    85,
    84,
    82,
    70,
    67,
    64,
    59,
    58,
    56,
    55,
    52,
    52,
    52,
    50,
    49,
    44,
    42,
    40,
    39,
    38,
    37,
    37,
    36,
    35,
    34,
    33,
    33,
    33,
    32,
    31,
    30,
    30,
    29,
    28,
    27,
    26,
    25,
    24,
    23,
    23,
    22,
    22,
    22,
    21,
    21,
    20,
    19,
    18,
    17,
    16,
    14,
    14,
    13,
    13,
    11,
    10,
    10,
    10,
    10,
    10,
    9,
    8,
    8,
    7,
    7,
    7,
    6,
    6,
    5,
    5,
    5,
    5,
    5,
    5,
    4,
    4,
    3,
    2,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1,
    1
  ];
}
