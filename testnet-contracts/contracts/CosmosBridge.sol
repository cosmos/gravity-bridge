pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";

contract CosmosBridge {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
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

    /*
    * @dev: Event declarations
    */
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

    /*
    * @dev: Modifier to restrict access to completed CosmosBridgeClaims
    */
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

    /*
    * @dev: Constructor
    */
    constructor()
        public
    {
        cosmosBridgeNonce = 0;
    }

    /*
    * @dev: newCosmosBridgeClaim
    *       Creates a new cosmos bridge claim, adding it to the cosmosBridgeClaims mapping
    */
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

    /*
    * @dev: getCosmosBridgeClaimStatus
    *       Returns the current status of a CosmosBridgeClaim
    */
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
