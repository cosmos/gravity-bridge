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
    uint256 public totalPower;
    uint256 public currentValsetVersion;
    uint256 public validatorCount;
    mapping (bytes32 => bool) public validators;
    mapping(bytes32 => uint256) public powers;

    /*
    * @dev: Event declarations
    */
    event LogValidatorAdded(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorPowerUpdated(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValidatorRemoved(
        address _validator,
        uint256 _power,
        uint256 _currentValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetReset(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    event LogValsetUpdated(
        uint256 _newValsetVersion,
        uint256 _validatorCount,
        uint256 _totalPower
    );

    /*
    * @dev: Modifier which restricts access to the operator.
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
        address[] memory _initValidators,
        uint256[] memory _initPowers
    )
        public
    {
        operator = _operator;
        currentValsetVersion = 0;

        updateValset(
            _initValidators,
            _initPowers
        );
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
    * @dev: addValidator
    */
    function addValidator(
        address _validatorAddress,
        uint256 _validatorPower
    )
        public
        onlyOperator
    {
        addValidatorInternal(
            _validatorAddress,
            _validatorPower
        );
    }

    /*
    * @dev: updateValidatorPower
    */
    function updateValidatorPower(
        address _validatorAddress,
        uint256 _newValidatorPower
    )
        public
        onlyOperator
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        require(
            validators[key],
            "Can only update the power of active valdiators"
        );


        // Adjust total power by new validator power
        uint256 priorPower = powers[key];
        totalPower = totalPower.sub(priorPower);
        totalPower = totalPower.add(_newValidatorPower);

        // Set validator's new power
        powers[key] = _newValidatorPower;

        emit LogValidatorPowerUpdated(
            _validatorAddress,
            _newValidatorPower,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: removeValidator
    */
    function removeValidator(
        address _validatorAddress
    )
        public
        onlyOperator
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        require(
            validators[key],
            "Can only remove active valdiators"
        );

        // Update validator count and total power
        validatorCount = validatorCount.sub(1);
        totalPower = totalPower.sub(powers[key]);

        // Delete validator and power
        delete validators[key];
        delete powers[key];

        emit LogValidatorRemoved(
            _validatorAddress,
            0,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: updateValset
    */
    function updateValset(
        address[] memory _validators,
        uint256[] memory _powers
    )
        public
        onlyOperator
    {
       require(
           _validators.length == _powers.length,
           "Every validator must have a corresponding power"
       );

       resetValset();

       for(uint256 i = 0; i < _validators.length; i = i.add(1)) {
           addValidatorInternal(_validators[i], _powers[i]);
       }

        emit LogValsetUpdated(
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: isActiveValidator
    */
    function isActiveValidator(
        address _validatorAddress
    )
        public
        view
        returns(bool)
    {
        // Recreate the unique key for this address given the current mapping version
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        // Return bool indicating if this address is an active validator
        return validators[key];
    }

    /*
    * @dev: getValidatorPower
    */
    function getValidatorPower(
        address _validatorAddress
    )
        external
        view
        returns(uint256)
    {
        // Recreate the unique key for this address given the current mapping version
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        return powers[key];
    }

    /*
    * @dev: recoverGas
    */
    function recoverGas(
        uint256 _valsetVersion,
        address _validatorAddress
    )
        external
        onlyOperator
    {
        require(
            _valsetVersion < currentValsetVersion,
            "Gas recovery only allowed for previous validator sets"
        );

        // Recreate the unique key used to identify this validator in the given version
        bytes32 key = keccak256(
            abi.encodePacked(
                _valsetVersion,
                _validatorAddress
            )
        );

        // Delete from mappings and recover gas
        delete(validators[key]);
        delete(powers[key]);
    }

    /*
    * @dev: addValidatorInternal
    */
    function addValidatorInternal(
        address _validatorAddress,
        uint256 _validatorPower
    )
        internal
    {
        // Create a unique key which for this validator's position in the current version of the mapping
        bytes32 key = keccak256(
            abi.encodePacked(
                currentValsetVersion,
                _validatorAddress
            )
        );

        validatorCount = validatorCount.add(1);
        totalPower = totalPower.add(_validatorPower);

        // Set validator as active and set their power
        validators[key] = true;
        powers[key] = _validatorPower;

        emit LogValidatorAdded(
            _validatorAddress,
            _validatorPower,
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }

    /*
    * @dev: resetValset
    */
    function resetValset()
        internal
    {
        currentValsetVersion = currentValsetVersion.add(1);
        validatorCount = 0;
        totalPower = 0;

        emit LogValsetReset(
            currentValsetVersion,
            validatorCount,
            totalPower
        );
    }
}