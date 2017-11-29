pragma solidity ^0.4.11;

import "./TendermintWire.sol";

contract IAVL is TendermintWire {
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
    
    function bytes20IsZero(bytes20 arr) constant returns (bool) {
        for (uint i = 0; i < 20; i++) {
            if (arr[i] != 0) break;
        }
        return i == 20;
    }
    
    function verifyInternal(Proof memory proof, bytes key, bytes value) private pure returns (bool) {
        bytes20 hash = leafHash(ProofLeafNode(key, value));
        if (sha3(hash) != sha3(proof.leafHash)) return false;
        for (uint idx = 0; idx < proof.innerNodes.length; idx++) {
            hash = innerHash(proof.innerNodes[idx], hash);
        }
        return proof.rootHash == hash;
    }
   
    function mkProof( 
        bytes20   proofLeafHash, 
        int8[]    proofInnerHeight, 
        int[]     proofInnerSize,
        bytes20[] proofInnerHash,
        bytes     proofInnerDirection,
        bytes20   proofRootHash
    ) internal pure returns (Proof) {
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
    
    
    function leafHash(ProofLeafNode leaf) internal pure returns (bytes20) {
        return ripemd160(byte(uint8(0)), 
                         byte(uint8(1)), byte(uint8(1)), 
                         writeUvarint(uint64(20)), leaf.keyBytes,
                         writeUvarint(uint64(20)), leaf.valueBytes);
    }
    
    function innerHash(ProofInnerNode branch, bytes20 childHash) internal pure returns (bytes20) {
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
    
    
}
