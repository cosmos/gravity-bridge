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

    uint256 public prophecyClaimCount;
    mapping(uint256 => ProphecyClaim) public prophecyClaims;

    enum Status {
        Null,
        Pending,
        Success,
        Failed
    }

    struct ProphecyClaim {
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

    event LogNewProphecyClaim(
        uint256 _prophecyID,
        bytes _cosmosSender,
        address payable _ethereumReceiver,
        address _validatorAddress,
        address _tokenAddress,
        string _symbol,
        uint256 _amount
    );

    event LogProphecyCompleted(
        uint256 _prophecyID
    );

    /*
    * @dev: Modifier which only allows access to currently pending prophecies
    */
    modifier isPending(
        uint256 _prophecyID
    )
    {
        require(
            isProphecyClaimActive(_prophecyID),
            "Prophecy claim is not active"
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
        prophecyClaimCount = 0;
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

    /*
    * @dev: newProphecyClaim
    *       Creates a new prophecy claim, adding it to the prophecyClaims mapping.
    *       ProphecyClaims can only be created for BridgeTokens on BridgeBank's whitelist.
    *        If the operator is responsible for adding them, then the automatic relay will
    *       of ProphecyClaim will fail until operator has called BridgeBank.createNewBridgeToken().
    */
    function newProphecyClaim(
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

        // Increment the prophecy claim count
        prophecyClaimCount = prophecyClaimCount.add(1);

        address originalValidator = msg.sender;

        // Create the new ProphecyClaim
        ProphecyClaim memory prophecyClaim = ProphecyClaim(
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount,
            Status.Pending
        );

        // Add the new ProphecyClaim to the mapping
        prophecyClaims[prophecyClaimCount] = prophecyClaim;

        emit LogNewProphecyClaim(
            prophecyClaimCount,
            _cosmosSender,
            _ethereumReceiver,
            originalValidator,
            _tokenAddress,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: completeProphecyClaim
    *       Allows for the completion of ProphecyClaims once processed by the Oracle
    */
    function completeProphecyClaim(
        uint256 _prophecyID
    )
        public
        isPending(_prophecyID)
    {
        require(
            msg.sender == oracle,
            "Only the Oracle may complete prophecies"
        );

        prophecyClaims[_prophecyID].status = Status.Success;

        issueBridgeTokens(_prophecyID);

        emit LogProphecyCompleted(
            _prophecyID
        );
    }

    /*
    * @dev: issueBridgeTokens
    *       Issues a request for the BridgeBank to mint new BridgeTokens
    */
    function issueBridgeTokens(
        uint256 _prophecyID
    )
        internal
    {
        ProphecyClaim memory prophecyClaim = prophecyClaims[_prophecyID];

        bridgeBank.mintBridgeTokens(
            prophecyClaim.cosmosSender,
            prophecyClaim.ethereumReceiver,
            prophecyClaim.tokenAddress,
            prophecyClaim.symbol,
            prophecyClaim.amount
        );
    }

    /*
    * @dev: isProphecyClaimActive
    *       Returns boolean indicating if the ProphecyClaim is active
    */
    function isProphecyClaimActive(
        uint256 _prophecyID
    )
        public
        view
        returns(bool)
    {
        return prophecyClaims[_prophecyID].status == Status.Pending;
    }

    /*
    * @dev: isProphecyValidatorActive
    *       Returns boolean indicating if the validator that originally
    *       submitted the ProphecyClaim is still an active validator
    */
    function isProphecyClaimValidatorActive(
        uint256 _prophecyID
    )
        public
        view
        returns(bool)
    {
        return valset.isActiveValidator(
            prophecyClaims[_prophecyID].originalValidator
        );
    }
}
