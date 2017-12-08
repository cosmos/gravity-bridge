pragma solidity ^0.4.11;

import "./ERC20.sol";

interface ValidatorSet {
    function valset() public constant returns (address[]);
    function numValidator() public constant returns (uint);
    function getValidator(uint n) public constant returns (address);
    function isValidator(address v) public constant returns (bool);
    event Update(address[] valset);

}

contract BasicValidatorSet is ValidatorSet {
    ERC20 public token;

    address[] private vs;
    uint private minval;

    mapping (address => uint) public bonded; 
    function getBonded(address k1) public constant returns (uint) { return bonded[k1]; }

    function BasicValidatorSet(address _token, address[] _vs) {
        token = ERC20(_token);
        vs = _vs;    
    }

    // entry points

    function valset() constant public returns (address[]) {
        return vs;
    }

    function numValidator() constant public returns (uint) {
        return vs.length;
    }

    function getValidator(uint n) constant public returns (address) {
        return vs[n];
    }

    function isValidator(address v) constant public returns (bool) {
        for (uint i = 0; i < vs.length; i++) {
            if (vs[i] == v) return true;
        }
        return false;
    }

    function bond(uint amount) external {
        assert(token.transferFrom(msg.sender, this, amount));
        bonded[msg.sender] += amount;

        if (bonded[msg.sender] > bonded[vs[minval]]) {
            vs[minval] = msg.sender;
            minval = newMinval();
            Update(vs);
        } 
    }

    function unbond(uint amount) external {
        assert(bonded[msg.sender] >= amount);
        bonded[msg.sender] -= amount;
        assert(token.transfer(msg.sender, amount));
        
        if (isValidator(msg.sender)) {
            minval = newMinval();
        }
    }

    function check() external {
        if (bonded[msg.sender] > bonded[vs[minval]]) {
            vs[minval] = msg.sender;
            minval = newMinval();
            Update(vs);
        }
    }    

    // helper functions 

    function newMinval() private constant returns (uint) {
        uint mv = 0;
        for (uint i = 0; i < vs.length; i++) {
            if (vs[i] < vs[mv]) mv = i;
        }
        return mv;
    }
}

/*
contract LiquidValidatorSet is ValidatorSet {
    function LiquidValidatorSet() {

    }

    function bond(uint amount) {

    }

    function unbond(uint amount) {

    }

    function delegate(address to, uint amount) {

    }

    function undelegate(address to, uint amount) {

    }
}

contract StaticValidatorSet is ValidatorSet {
    function StaticValidatorSet(address[] _vs) {

    }

    function valset() public returns 

}

contract DynamicValidatorSet is ValidatorSet {
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
    struct Commit {
        BlockID blockID;
        Vote[] precommits;
    }

    struct PartSetHeader {
        uint total;
        bytes20 hash;
    }
     
    struct BlockID {
        bytes20 hash;
        PartSetHeader partsHeader;
    }
     
    struct Vote {
        address validatorAddress;
        int validatorIndex;
        int height;
        int round;
        bytes20 blockID;
    }



    struct Certifier {
        string chainID;
        address[] vSet;
        int lastHeight;
        bytes32 vHash;
    }

    Certifier c;
   
    function updateCertifierInternal(Checkpoint check, Validator[] vset) private returns (bool) {
        assert(check.header.height > c.lastHeight);
        assert(validateCheckpoint(check, c.chainID));
        
    }
    struct Checkpoint {
        Header header;
        Commit commit;
    }    
 
    function validateCheckpoint(Checkpoint check, string chainID) private constant returns (bool) {
        if (keccak256(check.header.chainID) != keccak256(chainID)) return false;
        if (check.header.height != commitHeight(check.commit)) return false;
        if (headerHash(check.header) != check.commit.blockID.hash) return false; 
        return validateCommit(check.commit);
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

*/
