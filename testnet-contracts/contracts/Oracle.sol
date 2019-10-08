pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "./Valset.sol";

contract Oracle is Valset{

    struct OracleClaim {
        uint256 cosmosBridgeNonce;
        address validatorAddress;
        // TODO: Replace with validator's signed hash of the (ethereum recipient, amount, and cosmos bridge nonce)
        bytes contentHash;
        bool isClaim;
    }

    // Maps CosmosBridgeClaim id to OracleClaims made on it by validators
    mapping(uint256 => OracleClaim[]) public oracleClaims;
    // Maps validator's address to each CosmosBridgeClaim they've made
    mapping(address => uint256[]) public validatorOracleClaims;

    event LogNewOracleClaim(
        uint256 _cosmosBridgeNonce,
        address _validator,
        bytes _contentHash
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
        uint256 _bridgeContractNonce,
        address _validatorAddress,
        bytes memory _contentHash
    )
        internal
    {
        // Create a new claim
        OracleClaim memory oracleClaim = OracleClaim(
            _bridgeContractNonce,
            _validatorAddress,
            _contentHash,
            true
        );

        // Add the oracle claim to this CosmosBridgeClaim's OracleClaims
        oracleClaims[_cosmosBridgeNonce].push(_claim);
        // Add the oracle claim to this validator's OracleClaims
        validatorClaims[validator].push(_cosmosBridgeNonce);

        emit LogNewOracleClaim(
            _bridgeContractNonce,
            _validatorAddress,
            _contentHash
        );
    }

    function processProphecyClaim(
        bytes memory _contentHash,
        address[] memory _signers,
        bytes[] memory _signatures
    )
        internal
        returns(uint256)
    {
        uint256 signedPower = 0;
        mapping(address => bool) storage usedAddresses;

        // Iterate over the signatory addresses
        for (uint256 i = 0; i < _signers.length; i = i.add(1)) {
            // Only consider this signer if it's an unused address
            if(!usedAddresses[_signers[i]]) {
                // Get the power of the address that signed the hash
                uint256 signerPower = getPowerOfSignatory(
                    _contentHash,
                    _signatures[i]
                );

                // Mark this signer's address as used
                usedAddresses[_signers[i]] = true;

                // Add this signer's power to the total signed power
                signedPower = signedPower.add(signerPower);
            }
        }

        require(
            signedPower.mul(3) > totalPower.mul(2),
            "The cumulative power of signatory validators does not meet the threshold"
        );

        return signedPower;
    }
}