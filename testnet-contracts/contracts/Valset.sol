pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/cryptography/ECDSA.sol";

contract Valset {

    using SafeMath for uint256;
    using ECDSA for address;

    uint256 public numbValidators;
    uint256 public totalPower;
    uint256 public seqCounter = 0;

    mapping(address => bool) public activeValidators;
    address[] public validators;
    uint256[] public powers;

    event Update(
        address[] newValidators,
        uint256[] newPowers,
        uint256 seqCounter
    );

    modifier isActiveValidator(
        address _potentialValidator
    ) {
        require(
            activeValidators[_potentialValidator],
            "Must be a validator to make a claim"
        );
        _;
    }

    // Constructor which takes initial validator addresses and their powers
    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
    {
        numbValidators = 0;

        setValidatorsPower(
            initValidatorAddresses,
            initValidatorPowers
        );
    }

    function setValidatorsPower(
        address[] memory newValidators,
        uint256[] memory newPowers
    )
        internal
        returns (bool)
    {
        require(
            newValidators.length == newPowers.length,
            "Each validator must have a corresponding power"
        );

        // Reset active validators mapping
         for (uint256 i = 0; i < numbValidators; i++) {
             address priorValidator = validators[i];
             delete(validators[i]);
             activeValidators[priorValidator] = false;
         }

        // Reset validator count, validators array, powers array, and total power
        numbValidators = newValidators.length;
        validators = new address[](numbValidators);
        powers = new uint256[](newPowers.length);
        totalPower = 0;

        for (uint256 i = 0; i < numbValidators; i++) {
            // Set each new validator and their power
            validators[i] = newValidators[i];
            powers[i] = newPowers[i];
            activeValidators[newValidators[i]] = true;

            // Increment validator count and total power
            numbValidators = numbValidators.add(1);
            totalPower = totalPower.add(newPowers[i]);
        }

        // Increment the sequence counter
        seqCounter = seqCounter.add(1);

        emit Update(
            validators,
            powers,
            seqCounter
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

            // Only add active validators' powers
            if(activeValidators[signerAddr] && signerAddr == validators[signers[i]]) {
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