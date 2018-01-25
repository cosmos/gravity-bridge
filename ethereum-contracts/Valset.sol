pragma solidity ^0.4.11;

contract Valset {
    struct Validator {
        address Address;
        uint64 Power;
    }

    Validator[] public validators;
    uint64 public totalPower;

    uint updateSeq = 0;

    function verify(bytes32 hash, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) returns (bool) {
        if (!(idxs.length <= validators.length)) return false;
        if (!(idxs.length == v.length &&
              v.length == r.length &&
              r.length == s.length)) {
            return false;
        }

        uint64 signedPower = 0;

        for (uint i = 0; i < idxs.length; i++) {
            if (i >= 1 && idxs[i] <= idxs[i-1]) return false;
            if (ecrecover(hash, v[idxs[i]], r[idxs[i]], s[idxs[i]]) == validators[idxs[i]].Address) 
                signedPower += validators[idxs[i]].Power;
        }

        if (signedPower * 3 <= totalPower * 2) return false;

        return true;
    }

    event Update(Validator[] validators, uint indexed seq);

    function updateInternal(address[] newAddress, uint64[] newPowers) internal {
        validators = new Validator[](newAddress.length);
        totalPower = 0;
        for (uint i = 0; i < newAddress.length; i++) {
            validators[i] = Validator(newAddress[i], newPowers[i]);
            totalPower += newPowers[i];
        }

        Update(validators, updateSeq++);
    }

    function update(address[] newAddress, uint64[] newPowers, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) {
        require(newAddress.length == newPowers.length);
        
        assert(verify(keccak256(newAddress, newPowers), idxs, v, r, s)); // hashing can be changed

        updateInternal(newAddress, newPowers);
    }

    function Valset(address[] initAddress, uint64[] initPowers) {
        updateInternal(initAddress, initPowers);
    }
}
