pragma solidity ^0.4.11;

contract TendermintUtil {
    function mkEthaddr(bytes pubkey) internal pure returns (address) {
        bytes memory sliced = new bytes(pubkey.length-1);
        for (uint i = 1; i < pubkey.length; i++) {
            sliced[i-1] = pubkey[i];
        }

        return address(uint(keccak256(sliced)) & 0x00FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
    }

    function mkMintaddr(bytes pubkey) internal pure returns (bytes20) {
        bytes memory compressed = new bytes(33);
        if (uint8(pubkey[64])%2 == 0) {
            compressed[0] = byte(2);
        } else {
            compressed[0] = byte(3);
        }

        for (uint i = 1; i < 33; i++) {
            compressed[i] = pubkey[i];
        }

        return ripemd160(sha256(compressed));
    }


} 
