pragma solidity ^0.4.11;

library IAVL {
    struct IAVLProofInnerNode {
        int8 height;
        int size;
        bytes left;
        bytes right;
    }
    
    struct IAVLProofLeafNode {
        bytes keyBytes;
        bytes valueBytes;
    }
    
    struct IAVLProof {
        bytes leafHash;
        IAVLProofInnerNode[] innerNodes;
        bytes rootHash;
    }
    
    function verify(IAVLProof proof, bytes key, bytes value, bytes root) internal returns (bool) {
        if (sha3(proof.rootHash) != sha3(root)) return false;
        bytes20 hash = leafHash(IAVLProofLeafNode(key, value));
        if (sha3(hash) != sha3(proof.leafHash)) return false;
        for (uint idx = 0; idx < proof.innerNodes.length; idx++) {
            hash = innerHash(proof.innerNodes[idx], hash);
        }
        return sha3(proof.rootHash) == sha3(hash);
    }
    
    event Test(bytes a, bytes b);
    
    function uvarintSize(uint64 i) internal returns (uint8) {
        if (i == 0) return 0;
        if (i < 1<<8) return 1;
        if (i < 1<<16) return 2;
        if (i < 1<<24) return 3;
        if (i < 1<<32) return 4;
        if (i < 1<<40) return 5;
        if (i < 1<<48) return 6;
        if (i < 1<<56) return 7;
        return 8;
    }

    
    function writeVarint(uint64 i) returns (bytes)  {
        uint8 size = uvarintSize(i);
        bytes memory buf = new bytes(size+1);
        buf[0] = byte(size);
        for (uint idx = 0; idx < size; idx++) {
            buf[idx+1] = byte(uint8(i/(2**(8*((size-1)-idx)))));
        }
        return buf;
    }
    
    function leafHash(IAVLProofLeafNode leaf) internal returns (bytes20) {
        return ripemd160(byte(uint8(0)), 
                         byte(uint8(1)), byte(uint8(1)), 
                         writeVarint(uint64(leaf.keyBytes.length)), leaf.keyBytes,
                         writeVarint(uint64(leaf.valueBytes.length)), leaf.valueBytes);
    }
    
    function innerHash(IAVLProofInnerNode branch, bytes20 childHash) internal returns (bytes20) {
        if (branch.left.length == 0) {
            return ripemd160(byte(branch.height),
                             writeVarint(uint64(branch.size)),
                             writeVarint(uint64(20)), childHash,
                             writeVarint(uint64(branch.right.length)), branch.right);
        } else {
            return ripemd160(byte(branch.height),
                             writeVarint(uint64(branch.size)),
                             writeVarint(uint64(branch.left.length)), branch.left,
                             writeVarint(uint64(20)), childHash);
        }
    }
}
