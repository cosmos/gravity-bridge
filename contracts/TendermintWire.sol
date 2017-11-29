pragma solidity ^0.4.11;

contract TendermintWire {
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


}
