pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract CosmosBridge {

    using SafeMath for uint256;

    uint256 public cosmosBridgeNonce;
    mapping(uint256 => CosmosBridgeClaim) public cosmosBridgeClaims;

    enum Status {
        Active,
        Completed
    }

    struct CosmosBridgeClaim {
        uint256 nonce;
        bytes cosmosSender;
        address payable ethereumReceiver;
        address originalValidator;
        address tokenAddress;
        string symbol;
        uint256 amount;
        Status status;
    }

    event LogNewCosmosBridgeClaim(
        uint256 _cosmosBridgeNonce,
        uint256 _nonce,
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        address _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    event LogProphecyProcessed(
        uint256 _cosmosBridgeNonce,
        uint256 _signedPower,
        uint256 _totalPower,
        address _submitter
    );

    modifier isProcessing(
        uint256 _cosmosBridgeNonce
    )
    {
        require(
            cosmosBridgeClaims[_cosmosBridgeNonce].status == Status.Active,
            "Cannot make an OracleClaim on an already completed CosmosBridgeClaim"
        );
        _;
    }

    constructor()
        public
    {
        cosmosBridgeNonce = 0;
    }

    // TODO: Peggy public function protected by onlyValidator()
    // Creates a new cosmos bridge claim, adding it to the cosmosBridgeClaims mapping
    function newCosmosBridgeClaim(
        uint256 _nonce,
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        address _tokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        internal
        returns(bool)
    {
        // Increment the CosmosBridge nonce
        cosmosBridgeNonce = cosmosBridgeNonce.add(1);

        address originalValidator = msg.sender;

        // Create the new CosmosBridgeClaim
        CosmosBridgeClaim memory cosmosBridgeClaim = CosmosBridgeClaim(
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount,
            Status.Active
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


    function getCosmosBridgeClaimStatus(
        uint256 _cosmosBridgeNonce
    )
        public
        view
        returns(Status status)
    {
        return cosmosBridgeClaims[_cosmosBridgeNonce].status;
    }

}
