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

    address[] public validators;
    mapping(address => bool) public activeValidators;
    mapping(address => uint256) public powers;
    // uint256[] public powers;

    event LogUpdateValidatorSet(
        address[] _newValidators,
        uint256 _totalPower,
        uint256 _seqCounter
    );

    modifier isValidator(
        address _potentialValidator
    )
    {
        require(
            activeValidators[_potentialValidator],
            "Must be an active validator"
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

        updateValidatorsPower(
            initValidatorAddresses,
            initValidatorPowers
        );
    }

    function updateValidatorsPower(
        address[] memory newValidators,
        uint256[] memory newPowers
    )
        public
        returns (bool)
    {
        require(
            newValidators.length == newPowers.length,
            "Each validator must have a corresponding power"
        );

        // Reset active validators mapping and powers mapping
         for (uint256 i = 0; i < numbValidators; i++) {
             address priorValidator = validators[i];
             delete(validators[i]);
             activeValidators[priorValidator] = false;
             powers[priorValidator] = 0;
         }

        // Reset validator count, validators array, powers array, and total power
        numbValidators = newValidators.length;
        validators = new address[](numbValidators);
        totalPower = 0;

        // Iterate over the proposed validators
        for (uint256 i = 0; i < numbValidators; i++) {
            // Validators must have power greater than 0
            if(newPowers[i] > 0) {
                 // Set each new validator and their power
                validators[i] = newValidators[i];
                activeValidators[newValidators[i]] = true;
                powers[newValidators[i]] = newPowers[i];

                // Increment validator count and total power
                numbValidators = numbValidators.add(1);
                totalPower = totalPower.add(newPowers[i]);
            }
        }

        // Increment the sequence counter
        seqCounter = seqCounter.add(1);

        emit LogUpdateValidatorSet(
            validators,
            totalPower,
            seqCounter
        );

        return true;
    }

    // NOTE: _contentHash must include cosmosBridgeNonce to prevent replay attack
    function getPowerOfSignatory(
        bytes32 _contentHash,
        bytes memory _signature
    )
        internal
        view
        returns (uint256)
    {
        // Recover the address which originally signed this message
        address recoveredAddr = ECDSA.recover(
            _contentHash,
            _signature
        );

        // Only return the power of active validators
        if(activeValidators[recoveredAddr]) {
            return powers[recoveredAddr];
        } else {
            return 0;
        }
    }


    // TODO: These getter methods should be available automatically and are likely redundant
    function getValidators()
        public
        view
        returns (address[] memory)
    {
        return validators;
    }

    function isActiveValidator(
        address _validator
    )
        public
        view
        returns(bool)
    {
        return activeValidators[_validator];
    }

    function getValidatorPower(
        address _validator
    )
        public
        view
        returns(uint256)
    {
        return powers[_validator];
    }

    function getTotalPower()
        public
        view
        returns (uint256)
    {
        return totalPower;
    }
}