pragma solidity ^0.4.11;

import "./IAVL.sol";

contract TendermintLC is IAVL {
    // entry points

    // updates certifier (validator set)
    function update() internal {

    }

    // updates header with last known valset
    function certify()
         internal {
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

    uint private nextSeq = 0;

    mapping (uint => bool) private submitted;

    mapping (uint => bytes20) private apphash;

    struct Certifier {
         string chainID;
         Validator[] vSet;
         int lastHeight;
         bytes32 vHash;
     }
 
     Certifier c;
 
     struct Validator {
         address ethaddr;
         bytes20 mintaddr;
         bytes pubkey; // uncompressed
         uint votingPower;
         uint accum;
     }
     
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
 
     struct Checkpoint {
         Header header;
         Commit commit;
     }    
 
     function validateCheckpoint(Checkpoint check, string chainID) private view returns (bool) {
         if (keccak256(check.header.chainID) != keccak256(chainID)) return false;
         if (check.header.height != commitHeight(check.commit)) return false;
         if (headerHash(check.header) != check.commit.blockID.hash) return false; 
         return validateCommit(check.commit);
     }
     
     function updateCertifierInternal(Checkpoint check, Validator[] vset) private returns (bool) {
         assert(check.header.height > c.lastHeight);
         assert(validateCheckpoint(check, c.chainID));
         
     }
 
 
} 
