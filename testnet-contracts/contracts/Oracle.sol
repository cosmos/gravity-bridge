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
    uint256 public consensusThreshold; // e.g. 75 = 75%

    // Tracks the number of OracleClaims made on an individual BridgeClaim
    mapping(uint256 => address[]) public oracleClaimValidators;
    mapping(uint256 => mapping(address => bool)) public hasMadeClaim;

    /*
     * @dev: Event declarations
     */
    event LogNewOracleClaim(
        uint256 _prophecyID,
        bytes32 _message,
        address _validatorAddress,
        bytes _signature
    );

    event LogProphecyProcessed(
        uint256 _prophecyID,
        uint256 _prophecyPowerCurrent,
        uint256 _prophecyPowerThreshold,
        address _submitter
    );

    /*
     * @dev: Modifier to restrict access to the operator.
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "Must be the operator.");
        _;
    }

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier onlyValidator() {
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to current ValSet validators
     */
    modifier isPending(uint256 _prophecyID) {
        require(
            cosmosBridge.isProphecyClaimActive(_prophecyID) == true,
            "The prophecy must be pending for this operation"
        );
        _;
    }

    /*
     * @dev: Constructor
     */
    constructor(
        address _operator,
        address _valset,
        address _cosmosBridge,
        uint256 _consensusThreshold
    ) public {
        require(
            _consensusThreshold > 0,
            "Consensus threshold must be positive."
        );
        operator = _operator;
        cosmosBridge = CosmosBridge(_cosmosBridge);
        valset = Valset(_valset);
        consensusThreshold = _consensusThreshold;
    }

    /*
     * @dev: newOracleClaim
     *       Allows validators to make new OracleClaims on an existing Prophecy
     */
    function newOracleClaim(
        uint256 _prophecyID,
        bytes32 _message,
        bytes memory _signature
    ) public onlyValidator isPending(_prophecyID) {
        address validatorAddress = msg.sender;

        // Validate the msg.sender's signature
        require(
            validatorAddress == valset.recover(_message, _signature),
            "Invalid message signature."
        );

        // Confirm that this address has not already made an oracle claim on this prophecy
        require(
            !hasMadeClaim[_prophecyID][validatorAddress],
            "Cannot make duplicate oracle claims from the same address."
        );

        hasMadeClaim[_prophecyID][validatorAddress] = true;
        oracleClaimValidators[_prophecyID].push(validatorAddress);

        emit LogNewOracleClaim(
            _prophecyID,
            _message,
            validatorAddress,
            _signature
        );

        // Process the prophecy
        (
            bool valid,
            uint256 prophecyPowerCurrent,
            uint256 prophecyPowerThreshold
        ) = getProphecyThreshold(_prophecyID);

        if (valid) {
            completeProphecy(_prophecyID);

            emit LogProphecyProcessed(
                _prophecyID,
                prophecyPowerCurrent,
                prophecyPowerThreshold,
                msg.sender
            );
        }
    }

    /*
     * @dev: processBridgeProphecy
     *       Pubically available method which attempts to process a bridge prophecy
     */
    function processBridgeProphecy(uint256 _prophecyID)
        public
        isPending(_prophecyID)
    {
        // Process the prophecy
        (
            bool valid,
            uint256 prophecyPowerCurrent,
            uint256 prophecyPowerThreshold
        ) = getProphecyThreshold(_prophecyID);

        require(
            valid,
            "The cumulative power of signatory validators does not meet the threshold"
        );

        // Update the BridgeClaim's status
        completeProphecy(_prophecyID);

        emit LogProphecyProcessed(
            _prophecyID,
            prophecyPowerCurrent,
            prophecyPowerThreshold,
            msg.sender
        );
    }

    /*
     * @dev: checkBridgeProphecy
     *       Operator accessor method which checks if a prophecy has passed
     *       the validity threshold, without actually completing the prophecy.
     */
    function checkBridgeProphecy(uint256 _prophecyID)
        public
        view
        onlyOperator
        isPending(_prophecyID)
        returns (bool, uint256, uint256)
    {
        require(
            cosmosBridge.isProphecyClaimActive(_prophecyID) == true,
            "Can only check active prophecies"
        );
        return getProphecyThreshold(_prophecyID);
    }

    /*
     * @dev: processProphecy
     *       Calculates the status of a prophecy. The claim is considered valid if the
     *       combined active signatory validator powers pass the consensus threshold.
     *       The threshold is x% of Total power, where x is the consensusThreshold param.
     */
    function getProphecyThreshold(uint256 _prophecyID)
        internal
        view
        returns (bool, uint256, uint256)
    {
        uint256 signedPower = 0;
        uint256 totalPower = valset.totalPower();

        // Iterate over the signatory addresses
        for (
            uint256 i = 0;
            i < oracleClaimValidators[_prophecyID].length;
            i = i.add(1)
        ) {
            address signer = oracleClaimValidators[_prophecyID][i];

            // Only add the power of active validators
            if (valset.isActiveValidator(signer)) {
                signedPower = signedPower.add(valset.getValidatorPower(signer));
            }
        }

        // Prophecy must reach total signed power % threshold in order to pass consensus
        uint256 prophecyPowerThreshold = totalPower.mul(consensusThreshold);
        // consensusThreshold is a decimal multiplied by 100, so signedPower must also be multiplied by 100
        uint256 prophecyPowerCurrent = signedPower.mul(100);
        bool hasReachedThreshold = prophecyPowerCurrent >=
            prophecyPowerThreshold;

        return (
            hasReachedThreshold,
            prophecyPowerCurrent,
            prophecyPowerThreshold
        );
    }

    /*
     * @dev: completeProphecy
     *       Completes a prophecy by completing the corresponding BridgeClaim
     *       on the CosmosBridge.
     */
    function completeProphecy(uint256 _prophecyID) internal {
        cosmosBridge.completeProphecyClaim(_prophecyID);
    }
}
