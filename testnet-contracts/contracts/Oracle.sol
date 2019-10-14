pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";

contract Oracle {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
    Valset public valset;
    address public operator;

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
        bytes signature; // TODO: break into components (or just validate when submitted...)
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
    * @dev: Modifier to restrict access to current ValSet validators
    */
    modifier onlyValidator()
    {
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address _valset
    )
        public
    {
        operator = _operator;
        valset = Valset(_valset);
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
    // function processProphecyClaim(
    //     uint256 _cosmosBridgeClaimId,
    //     address _submitter,
    //     bytes32 _hash,
    //     address[] memory _signers,
    //     uint8[] memory _v,
    //     bytes32[] memory _r,
    //     bytes32[] memory _s
    // )
    //     internal
    // {
    //     uint256 signedPower = 0;

    //     // Iterate over the signatory addresses
    //     for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
    //         address signer = _signers[i];

    //         // Only consider this signer if it's an unused address
    //         if(!usedAddresses[signer]) {
    //             // Mark this signer's address as used
    //             usedAddresses[signer] = true;

    //             // Validate the signature
    //             bool valid = isValidSignature(
    //                 signer,
    //                 _hash,
    //                 _v[i],
    //                 _r[i],
    //                 _s[i]
    //             );

    //             // Only add the power of active validators
    //             if(valid && activeValidators[signer]) {
    //                 signedPower = signedPower.add(powers[signer]);
    //             }
    //         }
    //     }

    //     // Reset usedAddresses mapping
    //     for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
    //         if(usedAddresses[_signers[i]]) {
    //             usedAddresses[_signers[i]] = false;
    //         }
    //     }

    //     require(
    //         signedPower.mul(3) > totalPower.mul(2),
    //         "The cumulative power of signatory validators does not meet the threshold"
    //     );

    //     emit LogProphecyProcessed(
    //         _cosmosBridgeClaimId,
    //         signedPower,
    //         totalPower,
    //         _submitter
    //     );
    // }
}