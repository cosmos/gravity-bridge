pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/cryptography/ECDSA.sol";

contract Valset {

    using SafeMath for uint256;
    using ECDSA for address;

    address[] public validators;
    uint256[] public powers;
    uint256 public totalPower;
    uint256 public seqCounter = 0;

    event Update(
        address[] newValidators,
        uint256[] newPowers,
        uint256 seqCounter
    );

    // Constructor which takes initial validator addresses and their powers
    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
    {
        setValidatorsPower(
            initValidatorAddresses,
            initValidatorPowers
        );
    }

    function setValidatorsPower(
        address[] memory newAddress,
        uint256[] memory newPowers
    )
        internal
        returns (bool)
    {
        validators = new address[](newAddress.length);
        powers = new uint256[](newPowers.length);
        totalPower = 0;

        for (uint256 i = 0; i < newAddress.length; i++) {
            validators[i] = newAddress[i];
            powers[i] = newPowers[i];
            totalPower = totalPower.add(newPowers[i]);
        }

        // Increment and set the sequence counter
        seqCounter = seqCounter.add(1);
        uint256 updateCount = seqCounter;

        emit Update(
            validators,
            powers,
            updateCount
        );

        return true;
    }

    // TODO: signed hash must include nonce to prevent replay attack
    function verifyValidators(
        bytes32 signedHash,
        uint[] memory signers,
        bytes[] memory signatures
    )
        public
        view
        returns (bool)
    {
        uint256 signedPower = 0;

        // Iterate over the signers array
        for (uint i = 0; i < signers.length; i = i.add(1)) {
            // Recover the original signature's signing address
            address signerAddr = ECDSA.recover(
                signedHash,
                signatures[i]
            );

            // Only add valid validators' powers
            if(signerAddr == validators[signers[i]]) {
                signedPower = signedPower.add(powers[signers[i]]);
            }
        }

        require(
            signedPower.mul(3) > totalPower.mul(2),
            "The cumulative power of signatory validators does not meet the threshold"
        );

        return true;
    }

    // TODO: These getter methods should be available automatically and are likely redundant
    function getValidators()
        public
        view
        returns (address[] memory)
    {
        return validators;
    }

    function getPowers()
        public
        view
        returns (uint256[] memory)
    {
        return powers;
    }

    function getTotalPower()
        public
        view
        returns (uint256)
    {
        return totalPower;
    }
}