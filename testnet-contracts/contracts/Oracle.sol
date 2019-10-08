pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "./Valset.sol";

contract Oracle is Valset{

    struct Claim {
        bytes32 cosmosClaimID;
        address validatorAddress;
        address payable ethereumRecipient;
        uint256 amount;
        bool isClaim;
    }

    // Maps CosmosBridgeClaim id to OracleClaims made on it by validators
    mapping(uint256 => Claim[]) public oracleClaims;
    // Maps validator's address to each CosmosBridgeClaim they've made
    mapping(address => uint256[]) public validatorClaims;

    event LogNewOracleClaim(
        bytes32 _cosmosClaimID,
        address _validator,
        address _ethereumRecipient,
        uint256 _amount
    );

    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
        Valset(
            initValidatorAddresses,
            initValidatorPowers
        )
    {
        // Intentionally left blank
    }

    function newOracleClaim(
        bytes32 _cosmosClaimID,
        address payable _ethereumRecipient,
        uint256 _amount,
        address _validatorAddress
    )
        internal
        returns(Claim memory)
    {
        // Create a new claim
        Claim memory claim = Claim(
            _cosmosClaimID,
            _validatorAddress,
            _ethereumRecipient,
            _amount,
            true
        );

        emit LogNewOracleClaim(
            _cosmosClaimID,
            _validatorAddress,
            _ethereumRecipient,
            _amount
        );

        return (claim);
    }

    // Adds an oracle claim to the CosmosBridgeClaim's claims mapping,
    // as well as this validator's claims mapping
    function addOracleClaim(
        Claim memory _claim,
        uint256 _cosmosBridgeNonce
    )
        internal
        returns(bool)
    {
        address validator = msg.sender;

        require(
            activeValidators[validator],
            "Must be a validator to make a claim"
        );

        // Add the oracle claim to this transaction's claims
        oracleClaims[_cosmosBridgeNonce].push(_claim);
        // Add the oracle claim to this validator's claims
        validatorClaims[validator].push(_cosmosBridgeNonce);

        return true;
    }
}