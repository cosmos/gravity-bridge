pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../Oracle.sol";

contract TestOracle is Oracle {

    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
        Oracle(
            initValidatorAddresses,
            initValidatorPowers
        )
    {
        // Intentionally left blank
    }

    //Wrapper function to test internal method
    function callNewOracleClaim(
        uint256 _cosmosBridgeNonce,
        address _validatorAddress,
        bytes32 _contentHash,
        bytes memory _signature
    )
        public
    {
        return newOracleClaim(
            _cosmosBridgeNonce,
            _validatorAddress,
            _contentHash,
            _signature
        );
    }

   //Wrapper function to test internal method
    function callProcessProphecyClaim(
        uint256 _cosmosBridgeClaimId,
        bytes32 _hash,
        address[] memory _signers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    )
        public
    {
        return processProphecyClaim(
            _cosmosBridgeClaimId,
            msg.sender,
            _hash,
            _signers,
            _v,
            _r,
            _s
        );
    }
}
