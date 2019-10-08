pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Oracle.sol";

contract CosmosBridge is Oracle {

    using SafeMath for uint256;

    uint256 public cosmosBridgeNonce;
    mapping(uint256 => CosmosBridgeClaim) public cosmosBridgeClaims;

    enum Status {
        Processing,
        Completed
    }

    struct CosmosBridgeClaim {
        uint256 nonce;
        bytes cosmosSender;
        address payable ethereumReceiver;
        address originalValidator;
        bytes tokenAddress;
        string symbol;
        uint256 amount;
        bool isClaim;
        Status status;
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

    event LogProphecyProcessed(
        uint256 _cosmosBridgeNonce,
        uint256 _signedPower,
        uint256 _totalPower,
        address _submitter
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
        isValidator(msg.sender)
        returns(bool)
    {
        address originalValidator = msg.sender;

        // Increment the CosmosBridge nonce
        cosmosBridgeNonce = cosmosBridgeNonce.add(1);

        // Create the new CosmosBridgeClaim
        CosmosBridgeClaim memory cosmosBridgeClaim = CosmosBridgeClaim(
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount,
            true,
            Status.Processing
        );

        // Add the new CosmosBridgeClaim to the mapping
        cosmosBridgeClaims[cosmosBridgeNonce] = cosmosBridgeClaim;

        emit LogNewCosmosBridgeClaim(
            cosmosBridgeNonce,
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount
        );

        return true;
    }

    // Processes a validator's claim on an existing CosmosBridgeClaim
    function processOracleClaimOnCosmosBridgeClaim(
        uint256 _cosmosBridgeNonce,
        bytes memory contentHash
    )
        public
        isValidator(msg.sender)
        returns(bool)
    {
        require(
            cosmosBridgeClaims[_cosmosBridgeNonce].isClaim,
            "Cannot make an Oracle Claim on an empty Cosmos Bridge Claim"
        );

        // Create a new oracle claim
        newOracleClaim(
            _cosmosBridgeNonce,
            msg.sender,
            contentHash
        );
    }

    function processProphecyOnCosmosBridgeClaim(
        bytes32 _cosmosBridgeNonce,
        address[] memory signers,
        bytes[] memory signatures
    )
        public
    {
        // Pull the CosmosBridgeClaim from storage
        CosmosBridgeClaim memory cosmosBridgeClaim = cosmosBridgeClaims[_cosmosBridgeNonce];

        // Recreate the hash validators have signed
        bytes32 contentHash = keccak256(
            abi.encodePacked(
                _cosmosBridgeNonce,
                cosmosBridgeClaim.cosmosSender,
                cosmosBridgeClaim.nonce
            )
        );

        // Attempt to process the prophecy claim (throws if unsuccessful)
        uint256 signedPower = processProphecyClaim(
            contentHash,
            signers,
            signatures
        );

        // Update the CosmosBridgeClaim's status to completed
        cosmosBridgeClaims[_cosmosBridgeNonce].status = Status.Completed;

        emit LogProphecyProcessed(
            _cosmosBridgeNonce,
            signedPower,
            totalPower,
            msg.sender
        );
    }

    function getCosmosBridgeClaimStatus(
        uint256 _cosmosBridgeNonce
    )
        public
        returns(uint256 status)
    {
        return cosmosBridgeClaims[_cosmosBridgeNonce.status];
    }

}
