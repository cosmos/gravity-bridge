pragma solidity ^0.4.11;

import "TendermintWire.sol";

contract SimpleTree is TendermintWire {
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

    function simpleHashFromTwoHashes(bytes20 left, bytes20 right) private pure returns (bytes20) {
        return ripemd160(byte(uint8(1)), byte(uint8(20)), left,
                         byte(uint8(1)), byte(uint8(20)), right);
    }
    
    function computeHashFromAunts(int index, int total, bytes20 leafHash, bytes20[] innerHashes, int innerIndex) private constant returns (bytes20) {
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

   function verifySimple(int index, int total, bytes20 leafHash, bytes20 rootHash, bytes20[] innerHashes) internal constant returns (bool) {
       // assert(uint(int(innerHashes.length)) == innerHashes.length);
       return computeHashFromAunts(index, total, leafHash, innerHashes, int(innerHashes.length-1)) == rootHash;
   }
}
