pragma solidity ^0.4.11;

contract Valset {

    /* Variables */

    address[] public addresses;
    uint64[] public powers;
    uint64 public totalPower;
    uint internal updateSeq = 0;

    /* modifiers */

    modifier indexDoesntOverflow(uint index, uint length) {
      require(index < length);
      _;
    }

    modifier equalSizeArrays(uint validatorsLength, uint powersLenght) {
      require(validatorsLength == powersLenght);
      _;
    }

    modifier validatorSizeAtMost100(uint length) {
      require(length <= 100);
      _;
    }

    modifier valSetLargerThanSigners(uint signersLen, uint valLen) {
      require((signersLen <= valLen) && (valLen <= 100));
      _;
    }

    modifier equalSignatureLen(uint16 signersLen, uint16 vLen, uint16 rLen, uint16 sLen) {
      require((signersLen == vLen) && (vLen == rLen) && (rLen == sLen));
      _;
    }

    /* Events */

    event Update(address[] newAddresses, uint64[] newPowers, uint indexed seq);
    event Verify(uint16[] signers);
    event NoSupermajority();
    event NoLen(uint16[] signers);
    event InvalidSignature(uint16 validatorIdx, uint8 v, bytes32 r, bytes32 s);

    /* Functions */

    function getTotalPower() public constant returns (uint64) {
      return totalPower;
    }

    function getValidator(uint index)
      indexDoesntOverflow(index, addresses.length)
      public
      constant
      returns (address)
    {
      return addresses[index];
    }

    function getPower(uint index)
      indexDoesntOverflow(index, powers.length)
      public
      constant
      returns (uint64)
    {
      return powers[index];
    }


    function verify(bytes32 hash, uint8 v, bytes32 r, bytes32 s, address valAddress)
      internal
      pure
      returns(bool)
    {
      bytes memory prefix = "\x19Ethereum Signed Message:\n32";
      bytes32 prefixedHash = keccak256(prefix, hash);
      return ecrecover(prefixedHash, v, r, s) == (valAddress);
    }

    function verifyValidators(bytes32 hash, uint16[] signers, uint8[] v, bytes32[] r, bytes32[] s)
      valSetLargerThanSigners(signers.length, addresses.length)
      /* equalSignatureLen(signers.length, v.length, r.length, s.length) */
      public
      returns (bool)
    {
      uint64 signedPower = 0;
      uint16 currentIdx;
      if (signers.length == 0) {
        NoLen(signers);
      }
      for (uint i = 0; i < signers.length; i++) {
          currentIdx = signers[i];
          if (i >= 1 && currentIdx <= signers[i-1]) return false; // validators can't sign more than once
          if (verify(hash, v[currentIdx], r[currentIdx], s[currentIdx], getValidator(currentIdx))) {
            signedPower += powers[signers[i]];
          } else {
            InvalidSignature(currentIdx, v[currentIdx], r[currentIdx], s[currentIdx]);
            return false;
          }
      }
      if (signedPower * 3 < totalPower * 2) {
        NoSupermajority();
        return false;
      }
      Verify(signers);
      return true;
    }

    function updateInternal(address[] newAddress, uint64[] newPowers)
      equalSizeArrays(newAddress.length, newPowers.length)
      validatorSizeAtMost100(newAddress.length)
      internal
      returns (bool)
      {
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

    function update(address[] newAddress, uint64[] newPowers, uint16[] signers, uint8[] v, bytes32[] r, bytes32[] s)
      /* equalSizeArrays(newAddress.length, newPowers.length) */
      valSetLargerThanSigners(signers.length, newAddress.length)
      /* equalSignatureLen(signers.length, v.length, r.length, s.length) */
      public
      returns (bool)
    {
        bytes32 hashData = keccak256(newAddress, newPowers);
        assert(verifyValidators(hashData, signers, v, r, s)); // hashing can be changed
        if (updateInternal(newAddress, newPowers)) {
          return true;
        } else {
          return false;
        }
    }

    function Valset(address[] initAddress, uint64[] initPowers) public {
        updateInternal(initAddress, initPowers);
    }
}
