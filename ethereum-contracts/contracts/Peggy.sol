pragma solidity ^0.5.0;

import "./EscrowProcessor.sol";

  /*
   *  @title: Peggy
   *  @dev: Peg zone contract for testing one-way transfers from Ethereum
   *        to Cosmos, facilitated by a trusted relayer. This contract is
   *        NOT intended to be used in production and users are empowered
   *        to withdraw their locked funds at any time.
   */
contract Peggy is EscrowProcessor {

    bool public active;
    address public relayer;
    mapping(bytes32 => bool) public ids;

    event LogLock(
        bytes32 _id,
        address _from,
        bytes _to,
        address _token,
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
    * @dev: Modifier to restrict access to the relayer.
    *
    */
    modifier onlyRelayer()
    {
        require(
            msg.sender == relayer,
            'Must be the specified relayer.'
        );
        _;
    }

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
    * @dev: Constructor, initalizes relayer and active status.
    *
    */
    constructor()
        public
    {
        relayer = msg.sender;
        active = true;
        emit LogLockingActivated(now);
    }

    /* 
     * @dev: Locks received funds and creates new escrows.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of escrow
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
         //Actions based on token address type
        if (msg.value != 0) {
          require(_token == address(0));
          require(msg.value == _amount);
        } else {
          require(ERC20(_token).transferFrom(msg.sender, address(this), _amount));
        }

        //Create an escrow with a unique key.
        bytes32 escrowKey = createEscrow(
            msg.sender,
            _recipient,
            _token,
            _amount
        );

        emit LogLock(
            escrowKey,
            msg.sender,
            _recipient,
            _token,
            _amount,
            getNonce()
        );

        return escrowKey;
    }

    /*
     * @dev: Unlocks ethereum/erc20 tokens, called by relayer.
     *
     *       This is a shortcut utility method for testing purposes.
     *       In the future bidirectional system, unlocking functionality
     *       will be guarded by validator signatures.
     *
     * @param _escrowId: Unique key of the escrow.
     */
    function unlock(
        bytes32 _escrowId
    )
        onlyRelayer
        canDeliver(_escrowId)
        external
        returns (bool)
    {
        require(isEscrow(_escrowId));

        // Transfer escrow's funds and delete it from memory
        (address payable sender,
            address token,
            uint256 amount,
            uint256 escrowNonce) = completeEscrow(_escrowId);

        //Emit unlock event
        emit LogUnlock(
            _escrowId,
            sender,
            token,
            amount,
            escrowNonce
        );
        return true;
    }

    /*
     * @dev: Withdraws ethereum/erc20 tokens, called original sender.
     *
     *       This is a backdoor utility method included for testing,
     *       purposes, allowing users to withdraw their funds. This
     *       functionality will be removed in production.
     *
     * @param _escrowId: Unique key of the escrow.
     */
    function withdraw(
        bytes32 _escrowId
    )
        onlySender(_escrowId)
        canDeliver(_escrowId)
        external
        returns (bool)
    {
        require(isEscrow(_escrowId));

        // Transfer escrow's funds and delete it from memory
        (address payable sender,
            address token,
            uint256 amount,
            uint256 escrowNonce) = completeEscrow(_escrowId);

        //Emit withdraw event
        emit LogWithdraw(
            _escrowId,
            sender,
            token,
            amount,
            escrowNonce
        );

        return true;
    }

    /*
    * @dev: Exposes an escrow's current status.
    *
    * @param _escrowId: The escrow id in question.
    * @return: Boolean indicating the lock status.
    */
    function isLocked(
        bytes32 _escrowId
    )
        public 
        view
        returns(bool)
    {
        return isEscrow(_escrowId);
    }

    /*
    * @dev: Allows access to an escrow's information via its unique identifier.
    *
    * @param _escrowId: The escrow to be viewed.
    * @return: Original sender's address.
    * @return: Intended receiver's address in bytes.
    * @return: The token's address.
    * @return: The amount locked in the escrow.
    * @return: The escrow's unique nonce.
    */
    function viewEscrow(
        bytes32 _escrowId
    )
        public 
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getEscrow(_escrowId);
    }

    /*
    * @dev: Relayer can pause fund locking without impacting other functionality.
    */
    function pauseLocking()
        public
        onlyRelayer
    {
        require(active);
        active = false;
        emit LogLockingPaused(now);
    }

    /*
    * @dev: Relayer can activate fund locking without impacting other functionality.
    */
    function activateLocking()
        public
        onlyRelayer
    {
        require(!active);
        active = true;
        emit LogLockingActivated(now);
    }
}
