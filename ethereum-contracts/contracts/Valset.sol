pragma solidity ^0.4.11;

contract Valset {
    address[] public addresses;
    uint64[] public powers;
    uint64 public totalPower;

    uint updateSeq = 0;

    function verify(bytes32 hash, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) public constant returns (bool) {
        if (!(idxs.length <= addresses.length)) return false;
        if (!(idxs.length == v.length &&
              v.length == r.length &&
              r.length == s.length)) {
            return false;
        }

        uint64 signedPower = 0;

        for (uint i = 0; i < idxs.length; i++) {
            if (i >= 1 && idxs[i] <= idxs[i-1]) return false;
            if (ecrecover(hash, v[idxs[i]], r[idxs[i]], s[idxs[i]]) == addresses[idxs[i]])
                signedPower += powers[idxs[i]];
        }

        if (signedPower * 3 <= totalPower * 2) return false;

        return true;
    }

    event Update(address[] newAddresses, uint64[] newPowers, uint indexed seq);

    function updateInternal(address[] newAddress, uint64[] newPowers) internal {
        require(newAddress.length == newPowers.length);
        addresses = new address[](newAddress.length);
        powers    = new uint64[](newPowers.length);
        totalPower = 0;
        for (uint i = 0; i < newAddress.length; i++) {
            addresses[i] = newAddress[i];
            powers[i]    = newPowers[i];
            totalPower  += newPowers[i];
        }

        Update(addresses, powers, updateSeq++);
    }

    /// Updates validator set. Called by the relayers.
    /*
     * @param newAddress  new validators addresses
     * @param newPower    power of each validator
     * @param idxs        indexes of each validator
     * @param v           recovery id. Used to compute ecrecover
     * @param r           output of ECDSA signature. Used to compute ecrecover
     * @param s           output of ECDSA signature.  Used to compute ecrecover
     */
/*
    function update(address[] newAddress, uint64[] newPowers, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) {
        require(newAddress.length == newPowers.length);

        assert(verify(keccak256(newAddress, newPowers), idxs, v, r, s)); // hashing can be changed

        updateInternal(newAddress, newPowers);
    }
*/
    function Valset(address[] initAddress, uint64[] initPowers) public {
        updateInternal(initAddress, initPowers);
    }
}
