pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Oracle.sol";

contract CosmosBridge is Oracle {

    using SafeMath for uint256;

    uint256 public cosmosBridgeNonce;
    mapping(uint256 => CosmosBridgeClaim) public cosmosBridgeClaims;

    struct CosmosBridgeClaim {
        uint256 nonce;
        bytes cosmosSender;
        address payable ethereumReceiver;
        address validatorAddress;
        bytes tokenAddress;
        string symbol;
        uint256 amount;
        bool isClaim;
    }

    event LogNewCosmosBridgeClaim(
        uint256 _cosmosBridgeNonce,
        uint256 _nonce,
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        bytes _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    constructor()
        public
    {
        cosmosBridgeNonce = 0;
    }

    // Creates a new cosmos bridge claim, adding it to the cosmosBridgeClaims mapping
    function newCosmosBridgeClaim(
        uint256 _nonce,
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        bytes memory _tokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        returns(bool)
    {
        // Increment the CosmosBridge nonce
        cosmosBridgeNonce = cosmosBridgeNonce.add(1);

        // Parse validator's address
        address validatorAddress = msg.sender;

        // Create the new CosmosBridgeClaim
        CosmosBridgeClaim memory cosmosBridgeClaim = CosmosBridgeClaim(
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            validatorAddress,
            _tokenAddress,
            _symbol,
            _amount,
            true
        );

        // Add the new CosmosBridgeClaim to the mapping
        cosmosBridgeClaims[cosmosBridgeNonce] = cosmosBridgeClaim;

        emit LogNewCosmosBridgeClaim(
            cosmosBridgeNonce,
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            validatorAddress,
            _tokenAddress,
            _symbol,
            _amount
        );

        return true;
    }

    // Processes a validator's claim on an existing CosmosBridgeClaim
    function processOracleClaimOnCosmosBridgeClaim(
        uint256 _cosmosBridgeNonce,
        CosmosBridgeClaim memory cosmosClaim
    )
        internal
        returns(bool)
    {
        require(
            cosmosBridgeClaims[_cosmosBridgeNonce].isClaim,
            "Cannot make an Oracle Claim on an empty Cosmos Bridge Claim"
        );

        // Create unique id by hashing sender and nonce
        bytes32 oracleID = keccak256(
            abi.encodePacked(
                cosmosClaim.cosmosSender,
                cosmosClaim.nonce
            )
        );

        // Create a new claim
        Claim memory claim = newClaim(
            oracleID,
            cosmosClaim.validatorAddress,
            cosmosClaim.ethereumReceiver,
            cosmosClaim.amount
        );

        return addOracleClaim(claim);
    }

}
