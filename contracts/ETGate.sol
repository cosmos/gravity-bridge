pragma solidity ^0.4.11;

contract ERC20 {
    function totalSupply() constant returns (uint totalSupply);
    function balanceOf(address _owner) constant returns (uint balance);
    function transfer(address _to, uint _value) returns (bool success);
    function transferFrom(address _from, address _to, uint _value) returns (bool success);
    function approve(address _spender, uint _value) returns (bool success);
    function allowance(address _owner, address _spender) constant returns (uint remaining);
    event Transfer(address indexed _from, address indexed _to, uint _value);
    event Approval(address indexed _owner, address indexed _spender, uint _value);
}

contract IAVL {
    struct IAVLProofInnerNode {
        int8 height;
        int size;
        bytes20 left;
        bytes20 right;
    }
    
    struct IAVLProofLeafNode {
        bytes keyBytes;
        bytes valueBytes;
    }
    
    struct IAVLProof {
        bytes20 leafHash;
        IAVLProofInnerNode[] innerNodes;
        bytes20 rootHash;
    }
    
    function verify(IAVLProof memory proof, bytes key, bytes value, bytes20 root) internal returns (bool) {
        if (sha3(proof.rootHash) != sha3(root)) return false;
        bytes20 hash = leafHash(IAVLProofLeafNode(key, value));
        if (sha3(hash) != sha3(proof.leafHash)) return false;
        for (uint idx = 0; idx < proof.innerNodes.length; idx++) {
            hash = innerHash(proof.innerNodes[idx], hash);
        }
        return proof.rootHash == hash;
    }
    
    function getProof( 
        bytes20 iavlProofLeafHash, 
        int8[] iavlProofInnerHeight, 
        int[] iavlProofInnerSize,
        bytes20[] iavlProofInnerLeft,
        bytes20[] iavlProofInnerRight,
        bytes20 iavlProofRootHash
    ) internal returns (IAVL.IAVLProof) {
        uint l = iavlProofInnerHeight.length;
        require(l == iavlProofInnerSize.length &&
                l == iavlProofInnerLeft.length &&
                l == iavlProofInnerRight.length);
        IAVL.IAVLProofInnerNode[] memory innerNodes = new IAVL.IAVLProofInnerNode[](l);
        for (uint i = 0; i < l; i++) {
            innerNodes[i] = IAVL.IAVLProofInnerNode(iavlProofInnerHeight[i],
                                                    iavlProofInnerSize[i],
                                                    iavlProofInnerLeft[i],
                                                    iavlProofInnerRight[i]);
        }
        
        return IAVL.IAVLProof(iavlProofLeafHash, 
                              innerNodes, 
                              iavlProofRootHash);
    }
    
    function verifyRaw(
        bytes20 iavlProofLeafHash, 
        int8[] iavlProofInnerHeight, 
        int[] iavlProofInnerSize,
        bytes20[] iavlProofInnerLeft,
        bytes20[] iavlProofInnerRight,
        bytes20 iavlProofRootHash,
        bytes key, 
        bytes value, 
        bytes20 root
    ) {
        IAVLProof memory proof = getProof(iavlProofLeafHash,
                                          iavlProofInnerHeight,
                                          iavlProofInnerSize,
                                          iavlProofInnerLeft,
                                          iavlProofInnerRight,
                                          iavlProofRootHash);
                                          
        assert(verify(proof, key, value, root));
    }
    
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

    
    function writeUvarint(uint64 i) returns (bytes)  {
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
                         writeUvarint(uint64(20)), leaf.keyBytes,
                         writeUvarint(uint64(20)), leaf.valueBytes);
    }
    
    function innerHash(IAVLProofInnerNode branch, bytes20 childHash) internal returns (bytes20) {
        if (branch.left.length == 0) {
            return ripemd160(byte(branch.height),
                             writeUvarint(uint64(branch.size)),
                             writeUvarint(uint64(20)), childHash,
                             writeUvarint(uint64(branch.right.length)), branch.right);
        } else {
            return ripemd160(byte(branch.height),
                             writeUvarint(uint64(branch.size)),
                             writeUvarint(uint64(branch.left.length)), branch.left,
                             writeUvarint(uint64(20)), childHash);
        }
    }
    
    
}

contract ETGate is IAVL {
    mapping (uint => Header) headers;
    mapping (uint => uint) updated;
    mapping (bytes => bool) used;
    mapping (bytes => mapping (address => uint)) deposited;
    
    uint delay = 50;
    
    struct AppHash {
        bytes20 hash;
        uint block;
    }
    
    struct Validator {
        address ethaddr;
        bytes20 mintaddr;
        bytes pubkey; // uncompressed
        uint votingPower;
        uint accum;
    }
    
    struct BlockchainState {
        bytes32 chainID;
        Validator[] validators;
        uint lastBlockHeight;
        uint totalVotingPower;
    }
    
    BlockchainState public chainState;
    
    struct PartSetHeader {
        uint total;
        bytes20 hash;
    }
    
    struct BlockID {
        bytes20 hash;
        PartSetHeader partsHeader;
    }
    
    struct Header {
        string chainID;
        uint height;
        uint time;
        uint numTxs;
        BlockID lastBlockID;
        bytes20 lastCommitHash;
        bytes20 dataHash;
        bytes20 validatorHash;
        bytes20 appHash;
    }
    
    Header lastHeader;
    
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
        bytes memory btolen = IAVL.writeUvarint(uint64(bto.length));
        bytes memory bvalue = IAVL.writeUvarint(value);
        bytes memory btoken = toBytes(token);
        bytes memory btokenlen = IAVL.writeUvarint(uint64(btoken.length));
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
    
    uint64 accSeq;
    
    event Deposit(bytes to, uint64 value, address token, bytes chain, uint64 seq);
    
    function deposit(bytes to, uint64 value, address token, bytes chain) payable {
        if (token == 0) assert(value == msg.value);
        else assert(ERC20(token).transferFrom(msg.sender, this, value));
        deposited[chain][token] += value;
        Deposit(to, value, token, chain, accSeq++);
    }
    
    function depositEther(bytes to, uint64 value, bytes chain) payable { deposit(to, value, 0, chain); } 
    
    
    
    event Withdraw(address to, uint64 value, address token, bytes chain);
    
    function verifyWithdraw(IAVLProof proof, uint height, address to, uint64 value, address token, bytes chain, bytes seq) internal {
        bytes memory key = toKey(chain, seq);
        assert(!used[key]);
        used[key] = true;
        assert(verify(proof, key, toVal(to, value, token), headers[height].appHash));
    }
    
    function withdraw( // use when withdrawable(height)
        uint height, 
        bytes20 iavlProofLeafHash, 
        int8[] iavlProofInnerHeight, 
        int[] iavlProofInnerSize,
        bytes20[] iavlProofInnerLeft,
        bytes20[] iavlProofInnerRight,
        bytes20 iavlProofRootHash,
        address to,
        uint64 value,
        address token,
        bytes chain,
        bytes seq
    ) {
        require(withdrawable(height, chain, token, value));
        IAVLProof memory proof = getProof(iavlProofLeafHash,
                                          iavlProofInnerHeight,
                                          iavlProofInnerSize,
                                          iavlProofInnerLeft,
                                          iavlProofInnerRight,
                                          iavlProofRootHash);
                                          
        verifyWithdraw(proof, height, to, value, token, chain, seq);
               
        if (token == 0) {
            assert(to.send(value));
        } else {
            assert(ERC20(token).transfer(to, value));
        }
        
        Withdraw(to, value, token, chain);
    }
    
    function withdrawable(uint height, bytes chain, address token, uint64 value) constant returns (bool) {
        return updated[height] != 0 && updated[height] < block.number-delay &&
               deposited[chain][token] >= value;
    }
    
    modifier onlyValidator() {
        for (uint i = 0; i < chainState.validators.length; i++) {
            if (msg.sender == chainState.validators[i].ethaddr) break;
        }
        assert(i != chainState.validators.length);
        _;
    }
    
    function updateHeader(Header header) internal {
        headers[header.height] = header;
        updated[header.height] = block.number;
    }
    
    // update() accepts header submission that is from any of the known validators.
    // during the challange period, conflicting headers will be validated with 
    // VerifyCommit(https://github.com/tendermint/tendermint/blob/master/types/validator_set.go)
    function update(
        string _chainID,
        uint _height,
        uint _time,
        uint _numTxs,
        bytes20 _blockIDHash,
        uint _blockIDPartSetHeaderTotal,
        bytes20 _blockIDPartSetHEaderHash,
        bytes20 _lastCommitHash,
        bytes20 _dataHash,
        bytes20 _validatorHash,
        bytes20 _appHash
    ) onlyValidator {
        if (updated[_height] != 0) {
            // Verify
            revert();
        }
        updateHeader(Header(
            _chainID,
            _height,
            _time,
            _numTxs,
            BlockID (
                _blockIDHash,
                PartSetHeader(
                    _blockIDPartSetHeaderTotal,
                    _blockIDPartSetHEaderHash
                )
            ),
            _lastCommitHash,
            _dataHash,
            _validatorHash,
            _appHash
        ));
    }

    function newValidator(bytes pubkey, uint votingPower) internal returns (Validator) {
        address ethaddr = address(uint(keccak256(pubkey)) & 0x00FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF);
        bytes memory compressed = new bytes(33);
        for (uint i = 0; i < 32; i++) {
            compressed[i] = pubkey[i];
        }
        compressed[32] = pubkey[64];
        bytes20 mintaddr = ripemd160(sha256(compressed));
        return Validator(ethaddr, mintaddr, pubkey, votingPower, 0);
    }
    
    function ETGate(
        bytes32 _chainID,
        bytes _pubkey,
        uint[] _votingPower
    ) {
        require(_pubkey.length == _votingPower.length * 63);
        Validator[] storage validators;
        uint totalVotingPower = 0;
        for (uint i = 0; i < _votingPower.length; i++) {
            bytes memory pubkey = new bytes(63);
            for (uint j = 0; j < 63; j++) {
                pubkey[j] = _pubkey[i*63+j];
            }   
            validators.push(newValidator(pubkey, _votingPower[i]));
            totalVotingPower += _votingPower[i];
        }
        chainState.chainID = _chainID;
        chainState.validators = validators;
        chainState.lastBlockHeight = 0;
        chainState.totalVotingPower = totalVotingPower;
    }
}
