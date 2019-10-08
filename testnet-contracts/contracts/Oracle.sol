pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Recoverer.sol";

contract Oracle is Recoverer {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
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

    /*
    * @dev: Temporary buffer which records address use in a processProphecyClaim attempt,
    *       since the Solidity compiler does not yet implement memory (local) mappings
    */
    mapping(address => bool) internal usedAddresses;

    struct OracleClaim {
        address validatorAddress;
        bytes32 contentHash;
        bytes signature;
        bool isClaim;
    }

    /*
    * @dev: Event declarations
    */
    event LogNewOracleClaim(
        uint256 _cosmosBridgeClaimId,
        address _validatorAddress,
        bytes32 _contentHash,
        bytes _signature
    );

    event LogProphecyProcessed(
        uint256 _cosmosBridgeClaimId,
        uint256 _signedPower,
        uint256 _totalPower,
        address _submitter
    );
    
    event LogUpdateValidatorSet(
        address[] _newValidators,
        uint256 _totalPower,
        uint256 _seqCounter
    );

    /*
    * @dev: Modifier to restrict access to current validators.
    */
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
    */
    modifier onlyProvider()
    {
        require(
            msg.sender == provider,
            'Must be the specified provider.'
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
    {
        totalPower = 0;
        numbValidators = initValidatorAddresses.length;

        // Iterate over initial validators array
        for(uint256 i = 0; i < numbValidators; i++) {
            // Set each initial validator as active validator
            activeValidators[initValidatorAddresses[i]] = true;
            // Set each initial validator's power
            powers[initValidatorAddresses[i]] = initValidatorPowers[i];
            // Add each validator's power to the total power
            totalPower = totalPower.add(initValidatorPowers[i]);

            // TODO: This will be implemented for updateValidatorsPower()
            // Set each initial validator in validator array
            // validators[i] = initValidatorAddresses[i];
        }
    }

    /*
    * @dev: newOracleClaim
    *       Make a new OracleClaim on an existing CosmosBridgeClaim
    */
    function newOracleClaim(
        uint256 _cosmosBridgeNonce,
        address _validatorAddress,
        bytes32 _contentHash,
        bytes memory _signature
    )
        internal
    {
        // Create a new claim
        OracleClaim memory oracleClaim = OracleClaim(
            _validatorAddress,
            _contentHash,
            _signature,
            true
        );

        // Add the oracle claim to this CosmosBridgeClaim's OracleClaims
        oracleClaims[_cosmosBridgeNonce].push(oracleClaim);
        // Add the oracle claim to this validator's OracleClaims
        validatorOracleClaims[_validatorAddress].push(_cosmosBridgeNonce);

        emit LogNewOracleClaim(
            _cosmosBridgeNonce,
            _validatorAddress,
            _contentHash,
            _signature
        );
    }

    /*
    * @dev: processProphecyClaim
    *       Attempts to process a prophecy claim submission by reaching
    *       validator consensus
    */
    function processProphecyClaim(
        uint256 _cosmosBridgeClaimId,
        address _submitter,
        bytes32 _hash,
        address[] memory _signers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    )
        internal
    {
        uint256 signedPower = 0;

        // Iterate over the signatory addresses
        for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
            address signer = _signers[i];

            // Only consider this signer if it's an unused address
            if(!usedAddresses[signer]) {
                // Mark this signer's address as used
                usedAddresses[signer] = true;

                // Validate the signature
                bool valid = isValidSignature(
                    signer,
                    _hash,
                    _v[i],
                    _r[i],
                    _s[i]
                );

                // Only add the power of active validators
                if(valid && activeValidators[signer]) {
                    signedPower = signedPower.add(powers[signer]);
                }
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

        emit LogProphecyProcessed(
            _cosmosBridgeClaimId,
            signedPower,
            totalPower,
            _submitter
        );
    }

    /*
    * @dev: updateValidatorsPower
    *       Allows the provider to update both the validators and powers
    */
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

    /*
    * @dev: Getters
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