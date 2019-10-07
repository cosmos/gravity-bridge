pragma solidity ^0.5.0;

contract Oracle {

    struct Claim {
        bytes32 cosmosClaimID;
        address validatorAddress;
        address payable ethereumRecipient;
        uint256 amount;
        bool isClaim;
    }

    mapping(bytes32 => Claim[]) public cosmosBridgeClaims;
    mapping(address => Claim[]) public validatorClaims;

    event LogNewOracleClaim(
        bytes32 _cosmosClaimID,
        address _validator,
        address _ethereumRecipient,
        uint256 _amount
    );

    function newClaim(
        bytes32 _cosmosClaimID,
        address _validatorAddress,
        address payable _ethereumRecipient,
        uint256 _amount
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
        Claim memory _claim
    )
        internal
        returns(bool)
    {
        bytes32 cosmosClaimID = _claim.cosmosClaimID;
        address validator = _claim.validatorAddress;

        // Add the oracle claim to this transaction's claims
        cosmosBridgeClaims[cosmosClaimID].push(_claim);
        // Add the oracle claim to this validator's claims
        validatorClaims[validator].push(_claim);

        return true;
    }
}