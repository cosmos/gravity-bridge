pragma solidity ^0.4.17;

contract Valset {

    /* Variables */

    address[] public addresses;
    uint64[] public powers;
    uint64 public totalPower;
    uint internal updateSeq = 0;


    /* Events */

    event Update(address[] newAddresses, uint64[] newPowers, uint indexed seq);


    /* Getters (These are supposed to be auto implemented by solidity but aren't ¯\_(ツ)_/¯) */

    function getAddresses() public view returns (address[]) {
      return addresses;
    }

    function getPowers() public view returns (uint64[]) {
      return powers;
    }

    function getTotalPower() public view returns (uint64) {
      return totalPower;
    }


    /* Functions */

    function hashValidatorArrays(address[] addressesArr, uint64[] powersArr) public pure returns (bytes32 hash) {
      return keccak256(addressesArr, powersArr);
    }

    function verifyValidators(bytes32 hash, uint[] signers, uint8[] v, bytes32[] r, bytes32[] s) public constant returns (bool) {
      uint64 signedPower = 0;
      for (uint i = 0; i < signers.length; i++) {
        if (i > 0) {
          require(signers[i] > signers[i-1]);
        }
        address recAddr = ecrecover(hash, v[i], r[i], s[i]);
        require(recAddr == addresses[signers[i]]);

        signedPower += powers[signers[i]];
      }
      require(signedPower * 3 > totalPower * 2);
      return true;
    }


    function updateInternal(address[] newAddress, uint64[] newPowers) internal returns (bool) {
        addresses = new address[](newAddress.length);
        powers    = new uint64[](newPowers.length);
        totalPower = 0;
        for (uint i = 0; i < newAddress.length; i++) {
            addresses[i] = newAddress[i];
            powers[i]    = newPowers[i];
            totalPower  += newPowers[i];
        }
        uint updateCount = updateSeq;
        Update(addresses, powers, updateCount);
        updateSeq++;
        return true;
    }


    /// Updates validator set. Called by the relayers.
    /*
     * @param newAddress  new validators addresses
     * @param newPower    power of each validator
     * @param signers     indexes of each signer validator
     * @param v           recovery id. Used to compute ecrecover
     * @param r           output of ECDSA signature. Used to compute ecrecover
     * @param s           output of ECDSA signature.  Used to compute ecrecover
     */
    function update(address[] newAddress, uint64[] newPowers, uint[] signers, uint8[] v, bytes32[] r, bytes32[] s) public {
        bytes32 hashData = keccak256(newAddress, newPowers);
        require(verifyValidators(hashData, signers, v, r, s)); // hashing can be changed
        require(updateInternal(newAddress, newPowers));
    }

    function Valset(address[] initAddress, uint64[] initPowers) public {
        updateInternal(initAddress, initPowers);
    }
}
