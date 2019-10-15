pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";

contract CosmosBridge {

    using SafeMath for uint256;
    Valset public valset;

    /*
    * @dev: Public variable declarations
    */
    uint256 public bridgeClaimCount;
    mapping(uint256 => BridgeClaim) public bridgeClaims;

    enum Status {
        Inactive,
        Active
    }

    struct BridgeClaim {
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
    event LogNewBridgeClaim(
        uint256 _bridgeClaimCount,
        uint256 _nonce,
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        address _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    /*
    * @dev: Modifier to restrict access to completed BridgeClaims
    */
    modifier isProcessing(
        uint256 _bridgeClaimID
    )
    {
        require(
            isBridgeClaimActive(_bridgeClaimID),
            "Bridge claim is not active"
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
        address _valset
    )
        public
    {
        bridgeClaimCount = 0;
        valset = Valset(_valset);
    }

    /*
    * @dev: newBridgeClaim
    *       Creates a new bridge claim, adding it to the bridgeClaims mapping
    */
    function newBridgeClaim(
        uint256 _nonce,
        bytes memory _cosmosSender,
        address payable _ethereumReceiver,
        address _tokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyValidator
    {
        // Increment the bridge claim count
        bridgeClaimCount = bridgeClaimCount.add(1);

        address originalValidator = msg.sender;

        // Create the new BridgeClaim
        BridgeClaim memory bridgeClaim = BridgeClaim(
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount,
            Status.Active
        );

        // Add the new BridgeClaim to the mapping
        bridgeClaims[bridgeClaimCount] = bridgeClaim;

        emit LogNewBridgeClaim(
            bridgeClaimCount,
            _nonce,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount
        );
    }

    // TODO: add an internal function which mints new BankTokens via BridgeBank

    /*
    * @dev: isBridgeClaimActive
    *       Returns boolean indicating if the BridgeClaim is active
    */
    function isBridgeClaimActive(
        uint256 _bridgeClaimID
    )
        public
        view
        returns(bool)
    {
        return bridgeClaims[_bridgeClaimID].status == Status.Active;
    }

    /*
    * @dev: isBridgeClaimValidatorActive
    *       Returns boolean indicating if the validator that originally
    *       submitted the BridgeClaim is still an active validator
    */
    function isBridgeClaimValidatorActive(
        uint256 _bridgeClaimID
    )
        public
        view
        returns(bool)
    {
        return valset.isActiveValidator(
            bridgeClaims[_bridgeClaimID].originalValidator
        );
    }
}
