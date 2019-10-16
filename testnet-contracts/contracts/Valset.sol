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