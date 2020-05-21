import { ethers } from "@nomiclabs/buidler";
import { BigNumberish } from "ethers/utils";
import { Signer } from "ethers";

export async function getSignerAddresses(signers: Signer[]) {
  return await Promise.all(signers.map(signer => signer.getAddress()));
}

export function makeCheckpoint(
  newValidators: string[],
  newPowers: BigNumberish[],
  newValsetNonce: BigNumberish,
  peggyId: string
) {
  const methodName = ethers.utils.formatBytes32String("checkpoint");

  let checkpoint = ethers.utils.solidityKeccak256(
    ["bytes32", "bytes32", "uint256"],
    [peggyId, methodName, newValsetNonce]
  );

  for (let i = 0; i < newValidators.length; i = i + 1) {
    checkpoint = ethers.utils.solidityKeccak256(
      ["bytes32", "address", "uint256"],
      [checkpoint, newValidators[i], newPowers[i]]
    );
  }

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

// // bytes32 encoding of "transactionBatch"
// bytes32 methodName = 0x7472616e73616374696f6e426174636800000000000000000000000000000000;
// bytes32 transactionsHash = keccak256(abi.encodePacked(peggyId, methodName));

// uint256 lastTxNonceTemp = lastTxNonce;
// {
// 	for (uint256 i = 0; i < _amounts.length; i = i.add(1)) {
// 		require(
// 			_nonces[i] > lastTxNonceTemp,
// 			"Transaction nonces in batch must be strictly increasing"
// 		);
// 		lastTxNonceTemp = _nonces[i];

// 		transactionsHash = keccak256(
// 			abi.encodePacked(
// 				transactionsHash,
// 				_amounts[i],
// 				_destinations[i],
// 				_fees[i],
// 				_nonces[i]
// 			)
// 		);
// 	}
// }
export function makeTxBatchHash(
  amounts: number[],
  destinations: string[],
  fees: number[],
  nonces: number[],
  peggyId: string
) {
  const methodName = ethers.utils.formatBytes32String("transactionBatch");

  let txHash = ethers.utils.solidityKeccak256(
    ["bytes32", "bytes32"],
    [peggyId, methodName]
  );

  for (let i = 0; i < amounts.length; i = i + 1) {
    txHash = ethers.utils.solidityKeccak256(
      ["bytes32", "uint256", "address", "uint256", "uint256"],
      [txHash, amounts[i], destinations[i], fees[i], nonces[i]]
    );
  }

  return txHash;
}
