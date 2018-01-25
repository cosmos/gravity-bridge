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
}
