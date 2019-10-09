pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import "./Processor.sol";
import "./CosmosBridge.sol";
import "./Oracle.sol";
import "./Bank.sol";

  /*
   *  @title: Peggy
   *  @dev: Peg zone contract two-way asset transfers between Ethereum and
   *        to Cosmos, facilitated by a set of validators. This contract is
   *        NOT intended to be used in production (yet).
   */
contract Peggy is CosmosBridge, Oracle, Bank, Processor {

    bool public active;
    mapping(bytes32 => bool) public ids;

    event LogLock(
        bytes32 _id,
        address _from,
        bytes _to,
        address _token,
        string _symbol,
        uint256 _value,
        uint256 _nonce
    );

    event LogUnlock(
        bytes32 _id,
        address _to,
        address _token,
        uint256 _value,
        uint256 _nonce
    );

    event LogWithdraw(
        bytes32 _id,
        address _to,
        address _token,
        uint256 _value,
        uint256 _nonce
    );

    event LogLockingPaused(
        uint256 _time
    );

    event LogLockingActivated(
        uint256 _time
    );

    /*
    * @dev: Modifier which restricts lock functionality when paused.
    *
    */
    modifier whileActive()
    {
        require(
            active == true,
            'Lock functionality is currently paused.'
        );
        _;
    }
    /*
    * @dev: Constructor, initalizes provider and active status.
    *
    */
    constructor(
        address[] memory initValidatorAddresses,
        uint256[] memory initValidatorPowers
    )
        public
        Bank()
        CosmosBridge()
        Oracle(
            initValidatorAddresses,
            initValidatorPowers
        )
    {
        provider = msg.sender;
        active = true;
        emit LogLockingActivated(now);
    }

    /*
     * @dev: Locks received funds and creates new items.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of item
     */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        payable
        availableNonce()
        whileActive()
        returns(bytes32 _id)
    {
        string memory symbol;

        //Actions based on token address type
        if (msg.value != 0) {
          require(_token == address(0));
          require(msg.value == _amount);
          symbol = "ETH";
        } else {
          require(TokenERC20(_token).transferFrom(msg.sender, address(this), _amount));
          symbol = TokenERC20(_token).symbol();
        }

        //Create an item with a unique key.
        bytes32 id = create(
            msg.sender,
            _recipient,
            _token,
            _amount
        );

        emit LogLock(
            id,
            msg.sender,
            _recipient,
            _token,
            symbol,
            _amount,
            getNonce()
        );

        return id;
    }

    /*
    * @dev: makeOracleClaimOnCosmosBridgeClaim
    *       Processes a validator's claim on an existing CosmosBridgeClaim
    */
    function makeOracleClaimOnCosmosBridgeClaim(
        uint256 _cosmosBridgeNonce,
        bytes32 _contentHash,
        bytes memory _signature
    )
        public
        onlyValidator()
        isProcessing(
            _cosmosBridgeNonce
        )
        returns(bool)
    {

        // Create a new oracle claim
        newOracleClaim(
            _cosmosBridgeNonce,
            msg.sender,
            _contentHash,
            _signature
        );
    }

    /*
    * @dev: processProphecyOnOracleClaims
    *       Processes an attempted prophecy on a CosmosBridgeClaim's OracleClaims
    */
   function processProphecyOnOracleClaims(
        uint256 _cosmosBridgeNonce,
        bytes32 _hash,
        address[] memory _signers,
        uint8[] memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s
    )
        public
    {
        require(
            cosmosBridgeClaims[_cosmosBridgeNonce].status == Status.Completed,
            "Cannot process an prophecy on an already completed CosmosBridgeClaim"
        );

        // Pull the CosmosBridgeClaim from storage
        CosmosBridgeClaim memory cosmosBridgeClaim = cosmosBridgeClaims[_cosmosBridgeNonce];

        // Attempt to process the prophecy claim (throws if unsuccessful)
        processProphecyClaim(
            _cosmosBridgeNonce,
            msg.sender,
            _hash,
            _signers,
            _v,
            _r,
            _s
        );

        // Update the CosmosBridgeClaim's status to completed
        cosmosBridgeClaim.status = Status.Completed;
        
        deliver(
            cosmosBridgeClaim.tokenAddress,
            cosmosBridgeClaim.symbol,
            cosmosBridgeClaim.amount,
            cosmosBridgeClaim.ethereumReceiver
        );
    }

    /*
     * @dev: Unlocks ethereum/erc20 tokens, called by provider.
     *
     *       This is a shortcut utility method for testing purposes.
     *       In the future bidirectional system, unlocking functionality
     *       will be guarded by validator signatures.
     *
     * @param _id: Unique key of the item.
     */
    // TODO: Rework Processor and unlocking system to be compatible with prophecy processing
    function unlock(
        bytes32 _id
    )
        external
        onlyProvider
        canDeliver(_id)
        returns (bool)
    {
        require(isLocked(_id));

        // Transfer item's funds and unlock it
        (address payable sender,
            address token,
            uint256 amount,
            uint256 uniqueNonce) = complete(_id);

        //Emit unlock event
        emit LogUnlock(
            _id,
            sender,
            token,
            amount,
            uniqueNonce
        );
        return true;
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLocked(_id);
    }

    /*
    * @dev: Allows access to an item's information via its unique identifier.
    *
    * @param _id: The item to be viewed.
    * @return: Original sender's address.
    * @return: Intended receiver's address in bytes.
    * @return: The token's address.
    * @return: The amount locked in the item.
    * @return: The item's unique nonce.
    */
    function viewItem(
        bytes32 _id
    )
        public 
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getItem(_id);
    }

    /*
    * @dev: Provider can pause fund locking without impacting other functionality.
    */
    function pauseLocking()
        public
        onlyProvider
    {
        require(active);
        active = false;
        emit LogLockingPaused(now);
    }

    /*
    * @dev: Provider can activate fund locking without impacting other functionality.
    */
    function activateLocking()
        public
        onlyProvider
    {
        require(!active);
        active = true;
        emit LogLockingActivated(now);
    }
}
