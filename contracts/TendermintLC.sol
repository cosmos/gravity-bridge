pragma solidity ^0.4.11;

import "./IAVL.sol";
import "./SimpleTree.sol";
import "./ValidatorSet.sol";

contract TendermintLC is IAVL, SimpleTree {
    function TendermintLC(address _vs) {
        vs = ValidatorSet(_vs);
    }

    // updates header with last known valset
    // https://github.com/tendermint/tendermint/blob/master/types/validator_set.go verifycommit
    function certify(
        // signbytes (extracted from commit)
        uint signlen,
        bytes20 voteHash,
        bytes20 partsHash,
        uint partsTotal,
        uint height,
        uint round,
        // signs
        uint8[] v,
        bytes32[] r,
        bytes32[] s,
        // apphash
        bytes20 appHash,
        // apphash simpletree merkle proof
        bytes20[] appHashInner
    ) external {
        require(v.length == r.length && r.length == s.length);

        bytes memory o = new bytes(signlen);
        uint n = 0;
        // {"chain_id":"test-chain","vote":{"block_id":{"hash":"68617368","parts":{"hash":"70617274735F68617368","total":1000000}},"height":12345,"round":23456,"type":2}}
        // a lot of function call and array passing is here, just hope that the optimizer will inline these...
        (o, n) = openBrace(o, n);
        (o, n) = objectKey(o, n, "chain_id");
        (o, n) = objectStr(o, n, chainid);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "vote");
        (o, n) = openBrace(o, n);
        (o, n) = objectKey(o, n, "block_id");
        (o, n) = openBrace(o, n);
        (o, n) = objectKey(o, n, "hash");
        (o, n) = objectB20(o, n, voteHash);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "parts");
        (o, n) = openBrace(o, n);
        (o, n) = objectKey(o, n, "hash");
        (o, n) = objectB20(o, n, partsHash);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "total");
        (o, n) = objectInt(o, n, partsTotal);
        (o, n) = closBrace(o, n);
        (o, n) = closBrace(o, n);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "height");
        (o, n) = objectInt(o, n, height);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "round");
        (o, n) = objectInt(o, n, round);
        (o, n) = objectCma(o, n);
        (o, n) = objectKey(o, n, "type");
        (o, n) = objectInt(o, n, 3); // https://github.com/tendermint/tendermint/blob/master/types/priv_validator.go
        (o, n) = closBrace(o, n);
        (o, n) = closBrace(o, n);

        assert(n == signlen);

        bytes32 hash = sha256(o); // or keccak256? tendermint uses golang's crypto/sha256, not sure it is identital with sol's sha256

        uint sum = 0;

        for (uint i = 0; i < v.length; i++) {
            address signer = ecrecover(hash, v[i], r[i], s[i]);
            if (vs.isValidator(signer)) sum++;
        }

        assert(sum * 3 >= vs.numValidator() * 2);

        // and verify the merkle proof of apphash, push it to apphash[height]

        bytes memory appHashBytes = new bytes(20);
        for (i = 0; i < 20; i++) {
            appHashBytes[i] = appHash[i];
        }

        assert(verifySimple(0, 9, kvPairHash("App", appHashBytes), voteHash, appHashInner));
       
        apphash[height] = appHash;
    }

    // verify a key -value pair with a known header
    function verify(
        bytes     key, 
        bytes     value, 
        uint      height,
        int8[]    proofInnerHeight,
        int[]     proofInnerSize,
        bytes20[] proofInnerHash,
        bytes     proofInnerDirection
    ) internal {
        bytes20 proofRootHash = apphash[height]; 

        IAVL.verify(key, 
                    value, 
                    proofInnerHeight, 
                    proofInnerSize, 
                    proofInnerHash,
                    proofInnerDirection,
                    proofRootHash);
    }

    // check the header is submitted
    function available(uint height) constant internal returns (bool) {
        return submitted[height];
    }

    // check the packet sequence is continuous
    function continuous(uint seq) constant internal returns (bool) {
        return seq == nextSeq;
    }

    // structs and state variables

    bytes public chainid;

    uint private nextSeq = 0;

    mapping (uint => bool) private submitted;

    mapping (uint => bytes20) private apphash;

    ValidatorSet vs;
     
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
        int height;
        bytes20 timeHash;
        uint numTxs;
        BlockID lastBlockID;
        bytes20 lastCommitHash;
        bytes20 dataHash;
        bytes20 validatorsHash;
        bytes20 appHash;
    }
    function headerHash(Header header) internal returns (bytes20) {
        return 0x00;
    }
    struct Vote {
        address validatorAddress;
        int validatorIndex;
        int height;
        int round;
        bytes20 blockID;
    }
    struct Commit {
        BlockID blockID;
        Vote[] precommits;
    }
    function validateCommit(Commit commit) internal returns (bool) {
        
    }
 
    function commitHeight(Commit commit) internal pure returns (int){
        if (commit.precommits.length == 0) 
            return 0;
        else                               
            return commit.precommits[0].height;
    }
 
} 
