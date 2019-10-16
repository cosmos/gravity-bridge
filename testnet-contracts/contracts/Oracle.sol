pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./CosmosBridge.sol";

contract Oracle {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
    CosmosBridge public cosmosBridge;
    Valset public valset;
    address public operator;

    // Tracks the number of OracleClaims made on an individual BridgeClaim
    mapping(uint256 => address[]) public oracleClaimValidators;
    mapping(uint256 => mapping(address => bool)) public hasMadeClaim;

    /*
    * @dev: Event declarations
    */
    event LogNewOracleClaim(
        uint256 _bridgeClaimID,
        address _validatorAddress,
        bytes32 _message,
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
        address _valset,
        address _cosmosBridge
    )
        public
    {
        operator = _operator;
        cosmosBridge = CosmosBridge(_cosmosBridge);
        valset = Valset(_valset);
    }

    /*
    * @dev: newOracleClaim
    *       Allows validators to make new OracleClaims on an existing BridgeClaim
    */
    function newOracleClaim(
        uint256 _bridgeClaimID,
        bytes32 _message,
        bytes memory _signature
    )
        public
        onlyValidator
    {
        address validatorAddress = msg.sender;

        require(
            cosmosBridge.isBridgeClaimActive(
                _bridgeClaimID
            ) == true,
            "Can only make oracle claims upon active bridge claims"
        );

        // Validate the msg.sender's signature
        require(
            validatorAddress == valset.recover(
                _message,
                _signature
            ),
            "Invalid message signature."
        );

        // Confirm that this address has not already made a claim
        require(
            !hasMadeClaim[_bridgeClaimID][validatorAddress],
            "Cannot make duplicate oracle claims from the same address."
        );

        hasMadeClaim[_bridgeClaimID][validatorAddress] = true;
        oracleClaimValidators[_bridgeClaimID].push(validatorAddress);

        emit LogNewOracleClaim(
            _bridgeClaimID,
            validatorAddress,
            _message,
            _signature
        );
    }

    /*
    * @dev: processProphecyClaim
    *       Attempts to process a prophecy claim using validated validator powers
    */
    function processProphecyClaim(
        uint256 _bridgeClaimID
    )
        public
    {
        require(
            cosmosBridge.isBridgeClaimActive(
                _bridgeClaimID
            ) == true,
            "Can only attempt to process active bridge claims"
        );

        uint256 signedPower = 0;
        uint256 totalPower = valset.totalPower();

        // Iterate over the signatory addresses
        for (uint256 i = 0; i < oracleClaimValidators[_bridgeClaimID].length; i = i.add(1)) {
            address signer = oracleClaimValidators[_bridgeClaimID][i];

                // Only add the power of active validators
                if(valset.isActiveValidator(signer)) {
                    signedPower = signedPower.add(
                        valset.getValidatorPower(
                            signer
                        )
                    );
                }
        }

        require(
            signedPower.mul(3) > totalPower.mul(2),
            "The cumulative power of signatory validators does not meet the threshold"
        );

        // Update the BridgeClaim's status
        updateBridgeClaimStatus(
            _bridgeClaimID
        );

        emit LogProphecyProcessed(
            _bridgeClaimID,
            signedPower.mul(3),
            totalPower.mul(2),
            msg.sender
        );
    }

    /*
    * @dev: updateBridgeClaimStatus
    *       Completes a BridgeClaim on the CosmosBridge
    */
    function updateBridgeClaimStatus(
        uint256 _bridgeClaimID
    )
        internal
    {
        cosmosBridge.completeBridgeClaim(
            _bridgeClaimID
        );
    }
}