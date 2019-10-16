pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./Valset.sol";
import "./BridgeBank/BridgeBank.sol";

contract CosmosBridge {

    using SafeMath for uint256;

    /*
    * @dev: Public variable declarations
    */
    address public operator;
    Valset public valset;
    address public oracle;
    bool public hasOracle;
    BridgeBank public bridgeBank;
    bool public hasBridgeBank;

    uint256 public bridgeClaimCount;
    mapping(uint256 => BridgeClaim) public bridgeClaims;

    enum Status {
        Empty,
        Pending,
        Completed
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
    event LogOracleSet(
        address _oracle
    );

    event LogBridgeBankSet(
        address _bridgeBank
    );

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

    event LogBridgeClaimCompleted(
        uint256 _bridgeClaimID
    );

    /*
    * @dev: Modifier to restrict access to completed BridgeClaims
    */
    modifier isPending(
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
    * @dev: Modifier to restrict access to the operator.
    */
    modifier onlyOperator()
    {
        require(
            msg.sender == operator,
            'Must be the operator.'
        );
        _;
    }

      /*
    * @dev: The bridge is not active until oracle and bridge bank are set
    */
    modifier isActive()
    {
        require(
            hasOracle == true && hasBridgeBank == true,
            "The Operator must set the oracle and bridge bank for bridge activation"
        );
        _;
    }

    /*
    * @dev: Constructor
    */
    constructor(
        address _operator,
        address _valset
    )
        public
    {
        bridgeClaimCount = 0;
        operator = _operator;
        valset = Valset(_valset);
        hasOracle = false;
        hasBridgeBank = false;
    }

    /*
    * @dev: setOracle
    */
    function setOracle(
        address _oracle
    )
        public
        onlyOperator
    {
        require(
            !hasOracle,
            "The Oracle cannot be updated once it has been set"
        );

        hasOracle = true;
        oracle = _oracle;

        emit LogOracleSet(
            oracle
        );
    }

    /*
    * @dev: setBridgeBank
    */
    function setBridgeBank(
        address payable _bridgeBank
    )
        public
        onlyOperator
    {
        require(
            !hasBridgeBank,
            "The Bridge Bank cannot be updated once it has been set"
        );

        hasBridgeBank = true;
        bridgeBank = BridgeBank(_bridgeBank);

        emit LogBridgeBankSet(
            address(bridgeBank)
        );
    }

    // TODO: BridgeClaims can only be created for BridgeTokens on BridgeBank's whitelist.
    //       If the operator is responsible for adding them, then the automatic relay will
    //       of BridgeClaims will fail until operator has called BrideBank.createNewBridgeToken()
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
        isActive
    {
        require(
            valset.isActiveValidator(msg.sender),
            "Must be an active validator"
        );

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
            Status.Pending
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

    /*
    * @dev: completeBridgeClaim
    *       Allows for the completion of bridge claims once processed by the Oracle
    */
    function completeBridgeClaim(
        uint256 _bridgeClaimID
    )
        public
        isPending(_bridgeClaimID)
    {
        require(
            msg.sender == oracle,
            "Only the Oracle may complete bridge claims"
        );

        bridgeClaims[_bridgeClaimID].status = Status.Completed;

        issueBridgeTokens(_bridgeClaimID);

        emit LogBridgeClaimCompleted(
            _bridgeClaimID
        );
    }

    /*
    * @dev: issueBridgeTokens
    *       Issues a request for the BridgeBank to mint new BridgeTokens
    */
    function issueBridgeTokens(
        uint256 _bridgeClaimID
    )
        internal
    {
        BridgeClaim memory bridgeClaim = bridgeClaims[_bridgeClaimID];

        bridgeBank.mintBridgeTokens(
            bridgeClaim.cosmosSender,
            bridgeClaim.ethereumReceiver,
            bridgeClaim.tokenAddress,
            bridgeClaim.symbol,
            bridgeClaim.amount
        );
    }

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
        return bridgeClaims[_bridgeClaimID].status == Status.Pending;
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
