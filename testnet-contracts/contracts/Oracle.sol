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
        uint256 _weightedSignedPower,
        uint256 _weightedTotalPower,
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
    *       Pubically available method which attempts to process a prophecy claim
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

        // Process the claim
        (bool valid,
            uint256 weightedSignedPower,
            uint256 weightedTotalPower
        ) = processClaim(_bridgeClaimID);

        require(
            valid,
            "The cumulative power of signatory validators does not meet the threshold"
        );

        // Update the BridgeClaim's status
        completeCosmosBridgeClaim(
            _bridgeClaimID
        );

        emit LogProphecyProcessed(
            _bridgeClaimID,
            weightedSignedPower,
            weightedTotalPower,
            msg.sender
        );
    }

    /*
    * @dev: processProphecyClaim
    *       Operator accessor method which checks if a prophecy claim has passed
    *       the validity threshold without actually completing the claim
    */
    function checkProphecyClaim(
        uint256 _bridgeClaimID
    )
        public
        view
        onlyOperator
        returns(bool, uint256, uint256)
    {
        require(
            cosmosBridge.isBridgeClaimActive(
                _bridgeClaimID
            ) == true,
            "Can only check active bridge claims"
        );
        return processClaim(
            _bridgeClaimID
        );
    }

    /*
    * @dev: processClaim
    *       Attempts to process a prophecy claim. The claim is considered valid if
    *       all active signatory validator powers pass the validation threshold
    */
    function processClaim(
        uint256 _bridgeClaimID
    )
        internal
        view
        returns(bool, uint256, uint256)
    {
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

        // Calculate if weighted signed power has reached threshold of weighted total power
        uint256 weightedSignedPower = signedPower.mul(3);
        uint256 weightedTotalPower = totalPower.mul(2);
        bool hasReachedThreshold = weightedSignedPower >= weightedTotalPower;

        return(
            hasReachedThreshold,
            weightedSignedPower,
            weightedTotalPower
        );
    }

    /*
    * @dev: updateBridgeClaimStatus
    *       Completes a BridgeClaim on the CosmosBridge
    */
    function completeCosmosBridgeClaim(
        uint256 _bridgeClaimID
    )
        internal
    {
        cosmosBridge.completeBridgeClaim(
            _bridgeClaimID
        );
    }
}