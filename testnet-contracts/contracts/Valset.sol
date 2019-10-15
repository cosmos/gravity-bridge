pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/cryptography/ECDSA.sol";

contract Valset {

    using SafeMath for uint256;
     using ECDSA for bytes32;

    /*
    * @dev: Variable declarations
    */
    address public operator;
    uint256 public validatorCount;
    uint256 public totalPower;
    uint256 public seqCounter = 0;

    address[] public validators;
    mapping(address => bool) public activeValidators;
    mapping(address => uint256) public powers;

    /*
    * @dev: Event declarations
    */
    event LogUpdateValidatorSet(
        address[] _newValidators,
        uint256 _totalPower,
        uint256 _seqCounter
    );

    /*
    * @dev: Modifier to restrict access to the operator.
    */
    modifier onlyOperator()
    {
        require(
            msg.sender == operator,
            'Must be the operator.'
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address[] memory _initValidatorAddresses,
        uint256[] memory _initValidatorPowers
    )
        public
    {
        operator = _operator;
        totalPower = 0;
        validatorCount = _initValidatorAddresses.length;

        // Iterate over initial validators array
        for(uint256 i = 0; i < validatorCount; i++) {
            // Set each initial validator as active validator
            activeValidators[_initValidatorAddresses[i]] = true;
            // Set each initial validator's power
            powers[_initValidatorAddresses[i]] = _initValidatorPowers[i];
            // Add each validator's power to the total power
            totalPower = totalPower.add(_initValidatorPowers[i]);

            // TODO: This will be implemented for updateValidatorsPower()
            // Set each initial validator in validator array
            // validators[i] = initValidatorAddresses[i];
        }
    }
    /*
    * @dev: ECDSA methods for accessibliity
    *
    */
    function recover(
        bytes32 h,
        bytes memory signature
    )
        public
        pure
        returns (address)
    {
        return h.recover(signature);
    }

    function toEthSignedMessageHash(
        bytes32 h
    )
        public
        pure
        returns (bytes32)
    {
        return h.toEthSignedMessageHash();
    }

    // TODO: Implement individudal validator removal
    // function removeValidator(
    //     address _validator
    // )
    //     public
    //     onlyOperator
    // {
    //     require(
    //         isActiveValidator(_validator),
    //         "Can only remove active valdiators"
    //     );
    //     // ....
    // }

    /*
    * @dev: updateValidatorsPower
    *       Allows the provider to update both the validators and powers
    */
   function updateValidatorsPower(
        address[] memory newValidators,
        uint256[] memory newPowers
    )
        public
        onlyOperator
        returns (bool)
    {
        require(
            newValidators.length == newPowers.length,
            "Each validator must have a corresponding power"
        );

        // Reset active validators mapping and powers mapping
        for (uint256 i = 0; i < validatorCount; i++) {
            address priorValidator = validators[i];
            delete(validators[i]);
            activeValidators[priorValidator] = false;
            powers[priorValidator] = 0;
        }

        // Reset validator count, validators array, powers array, and total power
        validatorCount = newValidators.length;
        validators = new address[](validatorCount);
        totalPower = 0;

        // Iterate over the proposed validators
        for (uint256 i = 0; i < validatorCount; i++) {
            // Validators must have power greater than 0
            if(newPowers[i] > 0) {
                 // Set each new validator and their power
                validators[i] = newValidators[i];
                activeValidators[newValidators[i]] = true;
                powers[newValidators[i]] = newPowers[i];

                // Increment validator count and total power
                validatorCount = validatorCount.add(1);
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

    /*
    * @dev: Getter methods
    *
    */

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