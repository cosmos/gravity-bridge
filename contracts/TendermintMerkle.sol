pragma solidity ^0.4.11;

import "./TendermintWire.sol";

contract TendermintMerkle is TendermintWire {
    struct ProofInnerNode {
        int8 height;
        byte direction;
        int size;
        bytes20 hash;
    }
    
    struct ProofLeafNode {
        bytes keyBytes;
        bytes valueBytes;
    }
    
    struct Proof {
        bytes20 leafHash;
        ProofInnerNode[] innerNodes;
        bytes20 rootHash;
    }
    
    // entry points

    function verify(
        bytes     key, 
        bytes     value, 
        int8[]    proofInnerHeight, 
        int[]     proofInnerSize,
        bytes20[] proofInnerHash,
        bytes     proofInnerDirection,
        bytes20   proofRootHash
    ) internal pure returns (bool) {
        bytes20 proofLeafHash = leafHash(ProofLeafNode(key, value));

        Proof memory proof = mkProof(proofLeafHash,
                                     proofInnerHeight,
                                     proofInnerSize,
                                     proofInnerHash,
                                     proofInnerDirection,
                                     proofRootHash);
                                          
        return verifyInternal(proof, key, value);
    }
    
    function verifySimple(int index, int total, bytes20 leafHash, bytes20 rootHash, bytes20[] innerHashes) internal pure returns (bool) {
       // assert(uint(int(innerHashes.length)) == innerHashes.length);
       return computeHashFromAunts(index, total, leafHash, innerHashes, int(innerHashes.length-1)) == rootHash;
    }

    // https://github.com/tendermint/tmlibs/blob/master/merkle/simple_tree.go

    function kvPairHash(bytes key, bytes value) internal pure returns (bytes20) {
        // hash(len(len(key)) ++ len(key) ++ key ++ len(len(value)) ++ len(value) ++ value)

        assert(uint64(key.length)   == key.length);
        assert(uint64(value.length) == value.length);

        bytes memory keylen   = writeUvarint(uint64(key.length));
        bytes memory valuelen = writeUvarint(uint64(value.length));

        return ripemd160(keylen.length,   keylen,   key,
                         valuelen.length, valuelen, value);
    } 

    // helper functions

    function leafHash(ProofLeafNode leaf) private pure returns (bytes20) {
        return ripemd160(byte(uint8(0)), 
                         byte(uint8(1)), byte(uint8(1)), 
                         writeUvarint(uint64(20)), leaf.keyBytes,
                         writeUvarint(uint64(20)), leaf.valueBytes);
    }
    
    function innerHash(ProofInnerNode branch, bytes20 childHash) private pure returns (bytes20) {
        if (branch.direction == byte(0x00)) {
            return ripemd160(byte(branch.height),
                             writeUvarint(uint64(branch.size)),
                             writeUvarint(uint64(20)), childHash,
                             writeUvarint(uint64(20)), branch.hash);
        }
        if (branch.direction == byte(0x01)) {
            return ripemd160(byte(branch.height),
                             writeUvarint(uint64(branch.size)),
                             writeUvarint(uint64(20)), branch.hash,
                             writeUvarint(uint64(20)), childHash);
        } 
        revert();
    }
    
    function mkProof( 
        bytes20   proofLeafHash, 
        int8[]    proofInnerHeight, 
        int[]     proofInnerSize,
        bytes20[] proofInnerHash,
        bytes     proofInnerDirection,
        bytes20   proofRootHash
    ) private pure returns (Proof) {
        uint l = proofInnerHeight.length;
        require(l == proofInnerSize.length &&
                l == proofInnerHash.length &&
                l == proofInnerDirection.length);
        ProofInnerNode[] memory innerNodes = new ProofInnerNode[](l);
        for (uint i = 0; i < l; i++) {
            innerNodes[i] = ProofInnerNode(proofInnerHeight[i],
                                           proofInnerDirection[i],
                                           proofInnerSize[i],
                                           proofInnerHash[i]);
        }
        
        return Proof(proofLeafHash, 
                     innerNodes, 
                     proofRootHash);
    }
  
    function computeHashFromAunts(int index, int total, bytes20 leafHash, bytes20[] innerHashes, int innerIndex) private pure returns (bytes20) {
        assert(total < 1024); // prevent deep recursion
        assert(index < total);
        assert(total >= 1);

        if (total == 1) {
             return leafHash;
        }

        assert(innerHashes.length != 0);
        assert(innerIndex < int(innerHashes.length) && innerIndex >= 0);

        int numLeft = (total + 1) / 2;

        if (index < numLeft) {
            bytes20 leftHash = computeHashFromAunts(index, numLeft, leafHash, innerHashes, innerIndex-1);
            return simpleHashFromTwoHashes(leftHash, innerHashes[uint(innerIndex)]);
        } else {
            bytes20 rightHash = computeHashFromAunts(index-numLeft, total-numLeft, leafHash, innerHashes, innerIndex-1);
            return simpleHashFromTwoHashes(innerHashes[uint(innerIndex)], rightHash);
        }
    }

      
    function verifyInternal(Proof memory proof, bytes key, bytes value) private pure returns (bool) {
        bytes20 hash = leafHash(ProofLeafNode(key, value));
        if (sha3(hash) != sha3(proof.leafHash)) return false;
        for (uint idx = 0; idx < proof.innerNodes.length; idx++) {
            hash = innerHash(proof.innerNodes[idx], hash);
        }
        return proof.rootHash == hash;
    }

    function simpleHashFromTwoHashes(bytes20 left, bytes20 right) private pure returns (bytes20) {
        return ripemd160(byte(uint8(1)), byte(uint8(20)), left,
                         byte(uint8(1)), byte(uint8(20)), right);
    }
     
}
