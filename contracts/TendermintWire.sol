pragma solidity ^0.4.11;

contract TendermintWire {
    function openBrace(bytes o, uint n) internal pure returns (bytes, uint) {
        o[n++] = '{';
        return (o, n);
    }

    function openBraceLen() internal pure returns (uint) {
        return 1;
    }

    function closBrace(bytes o, uint n) internal pure returns (bytes, uint) {
        o[n++] = '}';
        return (o, n);
    }

    function closBraceLen() internal pure returns (uint) {
        return 1;
    }

    function objectCma(bytes o, uint n) internal pure returns (bytes, uint) {
        o[n++] = ',';
        return (o, n);
    }
    
    function objectKey(bytes o, uint n, bytes k) internal pure returns (bytes, uint) {
        o[n++] = '"';
        for (uint i = 0; i < k.length; i++) {
            o[n++] = k[i];
        }
        o[n++] = '"';
        o[n++] = ':';
        return (o, n);
    }

    function objectKeyLen(bytes k) internal pure returns (uint) {
        return k.length + 3;
    }

    function objectStr(bytes o, uint n, bytes s) internal pure returns (bytes, uint) {
        o[n++] = '"';
        for (uint i = 0; i < s.length; i++) {
            o[n++] = s[i];
        }
        o[n++] = '"';
        return (o, n);
    }

    function objectStrLen(bytes s) internal pure returns (uint) {
        return s.length + 2;
    }

    function objectInt(bytes o, uint n, uint v) internal pure returns (bytes, uint) {
        bytes memory buf = new bytes(20); // uint64 = [0, 1.8844e+19]
        uint i = 0;
        if (v == 0) {
            buf[i++] = byte(48);
        }
        while (v != 0) {
            buf[i++] = byte(v%10+48); // ascii '0' is 48
            v /= 10;
        }
        while(i != 0) {
            o[n++] = buf[--i];
        }

        return (o, n);
    }

    function objectB20(bytes o, uint n, bytes20 b) internal pure returns (bytes, uint) {
        for (uint i = 0; i < 20; i++) {
            o[n++] = b[i];
        }

        return (o, n);
    }

    function objectIntLen(uint v) internal pure returns (uint) {
        if (v == 0) return 1;
        for (uint i = 0; v != 0; i++) {
            v /= 10;
        }
        return i;
    }
   
    function uvarintSize(uint64 i) internal pure returns (uint8) {
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

    
    function writeUvarint(uint64 i) internal pure returns (bytes)  {
        uint8 size = uvarintSize(i);
        bytes memory buf = new bytes(size+1);
        buf[0] = byte(size);
        for (uint idx = 0; idx < size; idx++) {
            buf[idx+1] = byte(uint8(i/(2**(8*((size-1)-idx)))));
        }
        return buf;
    }

    // wont work on cosmos-sdk
    function toKey(bytes chain, bytes seq) constant returns (bytes) {
        bytes memory prefix = "etgate,withdraw,";
        bytes memory res = new bytes(prefix.length + chain.length + 1 + seq.length);
        for (uint i = 0; i < prefix.length; i++) {
            res[i] = prefix[i];
        }
        for (i = 0; i < chain.length; i++) {
            res[prefix.length + i] = chain[i];
        }
        res[prefix.length + chain.length] = ',';
        for (i = 0; i < seq.length; i++) {
            res[prefix.length + chain.length + 1 + i] = seq[i];
        }
        return res;
    }
    
    // https://ethereum.stackexchange.com/questions/884/how-to-convert-an-address-to-bytes-in-solidity
    function toBytes(address a) constant returns (bytes b){
        assembly {
            let m := mload(0x40)
            mstore(add(m, 20), xor(0x140000000000000000000000000000000000000000, a))
            mstore(0x40, add(m, 52))
            b := m
        }
    }
    
    function toVal(address to, uint64 value, address token) returns (bytes) {
        bytes memory bto = toBytes(to);
        bytes memory btolen = writeUvarint(uint64(bto.length));
        bytes memory bvalue = writeUvarint(value);
        bytes memory btoken = toBytes(token);
        bytes memory btokenlen = writeUvarint(uint64(btoken.length));
        bytes memory res = new bytes(btolen.length + 
                                     bto.length +
                                     bvalue.length +
                                     btokenlen.length + 
                                     btoken.length);
        
        uint idx = 0;
        uint i;
        for (i = 0; i < btolen.length; i++) {
            res[idx++] = btolen[i];
        }
        for (i = 0; i < bto.length; i++) {
            res[idx++] = bto[i];
        }
        for (i = 0; i < bvalue.length; i++) {
            res[idx++] = bvalue[i];
        }
        for (i = 0; i < btokenlen.length; i++) {
            res[idx++] = btokenlen[i];
        }
        for (i = 0; i < btoken.length; i++) {
            res[idx++] = btoken[i];
        }
        return res;
    }

    function writeByteSlice(bytes o, uint n, bytes k) internal pure returns (bytes, uint) {

    }
}
