pragma solidity ^0.4.11;

contract Valset {
    struct Validator {
        address Address;
        uint64 Power;
    }

    Validator[] public validators;
    uint64 public totalPower;

    uint updateSeq = 0;

    function verify(bytes32 hash, uint8[] v, bytes32[] r, bytes32[] s) returns (bool) {
        if (!(v.length == r.length &&
              r.length == s.length)) {
            return false;
        }

        uint64 signedPower = 0;

        for (uint i = 0; i < v.length; i++) {
            if (v[i] == 0) continue;
            if (ecrecover(hash, v[i], r[i], s[i]) == validators[i].Address) 
                signedPower += validators[i].Power;
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

    function update(address[] newAddress, uint64[] newPowers, uint8[] v, bytes32[] r, bytes32[] s) {
        require(newAddress.length == newPowers.length);
        
        assert(verify(keccak256(newAddress, newPowers), v, r, s)); // hashing can be changed

        updateInternal(newAddress, newPowers);
    }

    function Valset(address[] initAddress, uint64[] initPowers) {
        updateInternal(initAddress, initPowers);
    }
}
