pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../../node_modules/openzeppelin-solidity/contracts/cryptography/ECDSA.sol";

contract Oracle {

    using SafeMath for uint256;
    using ECDSA for address;

    address public provider;
    uint256 public numbValidators;
    uint256 public totalPower;
    uint256 public seqCounter = 0;

    address[] public validators;
    mapping(address => bool) public activeValidators;
    mapping(address => uint256) public powers;
    // Maps CosmosBridgeClaim id to OracleClaims made on it by validators
    mapping(uint256 => OracleClaim[]) public oracleClaims;
    // Maps validator's address to each CosmosBridgeClaim they've made
    mapping(address => uint256[]) public validatorOracleClaims;

    // Temporary buffer which records address use in a processProphecyClaim attempt
    // since the Solidity compiler does not yet implement memory (local) mappings
    mapping(address => bool) internal usedAddresses;

    struct OracleClaim {
        uint256 cosmosBridgeNonce;
        address validatorAddress;
        bytes32 contentHash;
        bool isClaim;
    }

    event LogNewOracleClaim(
        uint256 _cosmosBridgeNonce,
        address _validator,
        bytes32 _contentHash
    );
    event LogUpdateValidatorSet(
        address[] _newValidators,
        uint256 _totalPower,
        uint256 _seqCounter
    );

    modifier onlyValidator()
    {
        require(
            activeValidators[msg.sender],
            "Must be an active validator"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the provider.
    *
    */
    modifier onlyProvider()
    {
        require(
            msg.sender == provider,
            'Must be the specified provider.'
        );
        _;
    }

    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
    {
        // Set validator count and validators array
        numbValidators = initValidatorAddresses.length;
        validators = initValidatorAddresses;

        // Iterate over validators array and set each validator's power
        for(uint256 i = 0; i < numbValidators; i++) {
            powers[validators[i]] = initValidatorPowers[i];
        }
    }

    function newOracleClaim(
        uint256 _cosmosBridgeNonce,
        address _validatorAddress,
        bytes32 _contentHash
    )
        internal
    {
        // Create a new claim
        OracleClaim memory oracleClaim = OracleClaim(
            _cosmosBridgeNonce,
            _validatorAddress,
            _contentHash,
            true
        );

        // Add the oracle claim to this CosmosBridgeClaim's OracleClaims
        oracleClaims[_cosmosBridgeNonce].push(oracleClaim);
        // Add the oracle claim to this validator's OracleClaims
        validatorOracleClaims[_validatorAddress].push(_cosmosBridgeNonce);

        emit LogNewOracleClaim(
            _cosmosBridgeNonce,
            _validatorAddress,
            _contentHash
        );
    }

    function updateValidatorsPower(
        address[] memory newValidators,
        uint256[] memory newPowers
    )
        public
        onlyProvider()
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

    function processProphecyClaim(
        bytes32 _contentHash,
        address[] memory _signers,
        bytes[] memory _signatures
    )
        internal
        returns(uint256)
    {
        uint256 signedPower = 0;

        // Iterate over the signatory addresses
        for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
            address signer = _signers[i];

            // Only consider this signer if it's an unused address
            if(!usedAddresses[signer]) {
                // Mark this signer's address as used
                usedAddresses[signer] = true;

                // Get the power of the address that signed the hash
                uint256 signerPower = getPowerOfSignatory(
                    _contentHash,
                    _signatures[i]
                );


                // Add this signer's power to the total signed power
                signedPower = signedPower.add(signerPower);
            }
        }

        // Reset usedAddresses mapping
        for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
            if(usedAddresses[_signers[i]]) {
                usedAddresses[_signers[i]] = false;
            }
        }

        require(
            signedPower.mul(3) > totalPower.mul(2),
            "The cumulative power of signatory validators does not meet the threshold"
        );

        return signedPower;
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