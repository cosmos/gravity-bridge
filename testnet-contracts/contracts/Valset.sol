pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract Valset {

    using SafeMath for uint256;

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

    function recover(
        string memory _message,
        bytes memory _signature
    )
        public
        pure
        returns (address)
    {
        bytes32 message = ethMessageHash(_message);
        return verify(message, _signature);
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

   /*
    * @dev: ECDSA methods for accessibliity
    *
    */
    function verify(
        bytes32 h,
        bytes memory signature
    )
        internal
        pure
        returns (address)
    {
        bytes32 r;
        bytes32 s;
        uint8 v;

        // Check the signature length
        if (signature.length != 65) {
            return (address(0));
        }

        // Divide the signature in r, s and v variables
        // ecrecover takes the signature parameters, and the only way to get them
        // currently is to use assembly.
        // solium-disable-next-line security/no-inline-assembly
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }

        // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
        if (v < 27) {
            v += 27;
        }

        // If the version is correct return the signer address
        if (v != 27 && v != 28) {
            return (address(0));
        } else {
            // solium-disable-next-line arg-overflow
            return ecrecover(h, v, r, s);
        }
    }

    /**
    * @dev prefix a bytes32 value with "\x19Ethereum Signed Message:" and hash the result
    */
    function ethMessageHash(string memory message) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(
            "\x19Ethereum Signed Message:\n32", keccak256(abi.encodePacked(message)))
        );
    }
}