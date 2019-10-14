pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

// import "./CosmosBridge.sol";
// import "./Oracle.sol";
// import "./EthereumBank.sol";
// import "./CosmosBank.sol";

  /*
   *  @title: Peggy
   *  @dev: Peg zone contract two-way asset transfers between Ethereum and
   *        to Cosmos, facilitated by a set of validators. This contract is
   *        NOT intended to be used in production (yet).
   */
contract Peggy {

    bool public active;
    mapping(bytes32 => bool) public ids;
    address public operator;
   
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
    constructor()
        public
    {
        operator = msg.sender;
        active = true;
        emit LogLockingActivated(now);
    }

    /*
     * @dev: Locks received Ethereum funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    // function lock(
    //     bytes memory _recipient,
    //     address _token,
    //     uint256 _amount
    // )
    //     public
    //     payable
    //     availableNonce()
    //     whileActive()
    //     returns(bytes32 _id)
    // {
    //     string memory symbol;

    //     // Ethereum deposit
    //     if (msg.value != 0) {
    //       require(
    //           _token == address(0),
    //           "Ethereum deposits require the 'token' address to be the null address"
    //         );
    //       require(
    //           msg.value == _amount,
    //           "The transactions value must be equal the specified amount (in wei)"
    //         );

    //       // Set the the symbol to ETH
    //       symbol = "ETH";
    //       // ERC20 deposit
    //     } else {
    //       require(
    //           TokenERC20(_token).transferFrom(msg.sender, address(this), _amount),
    //           "Contract token allowances insufficient to complete this lock request"
    //       );
    //       // Set symbol to the ERC20 token's symbol
    //       symbol = TokenERC20(_token).symbol();
    //     }

        //Create a deposit with a unique key.
    //     bytes32 id = lockEthereumDeposit(
    //         msg.sender,
    //         _recipient,
    //         _token,
    //         symbol,
    //         _amount
    //     );

    //     // emit LogLock(
    //     //     id,
    //     //     msg.sender,
    //     //     _recipient,
    //     //     _token,
    //     //     symbol,
    //     //     _amount,
    //     //     getNonce()
    //     // );

    //     return id;
    // }

    /*
     * @dev: Unlocks Cosmos deposits.
     *
     *       Replicate _id hash off-chain with sha3(cosmosSender, ethereumRecipient, amount) + nonce
     *
     * @param _id: Unique key of the CosmosDeposit.
     */
    // function unlock(
    //     bytes32 _id
    // )
    //     public
    //     onlyProvider
    //     canDeliver(_id)
    //     returns (bool)
    // {
    //     // TODO: Refactor this refundant check
    //     require(isLockedEthereumDeposit(_id), "Must be locked");

    //     // Unlock the deposit and transfer funds
    //     (address payable sender,
    //         address token,
    //         uint256 amount,
    //         uint256 uniqueNonce) = unlockEthereumDeposit(_id);

        //Emit unlock event
        // emit LogUnlock(
        //     _id,
        //     sender,
        //     token,
        //     amount,
        //     uniqueNonce
        // );
    //     return true;
    // }

    /*
    * @dev: makeOracleClaimOnCosmosBridgeClaim
    *       Processes a validator's claim on an existing CosmosBridgeClaim
    */
    // function makeOracleClaimOnCosmosBridgeClaim(
    //     uint256 _cosmosBridgeNonce,
    //     bytes32 _contentHash,
    //     bytes memory _signature
    // )
    //     public
    //     onlyValidator()
    //     isProcessing(
    //         _cosmosBridgeNonce
    //     )
    //     returns(bool)
    // {

    //     // Create a new oracle claim
    //     newOracleClaim(
    //         _cosmosBridgeNonce,
    //         msg.sender,
    //         _contentHash,
    //         _signature
    //     );
    // }

    /*
    * @dev: processProphecyOnOracleClaims
    *       Processes an attempted prophecy on a CosmosBridgeClaim's OracleClaims
    */
//    function processProphecyOnOracleClaims(
//         uint256 _cosmosBridgeNonce,
//         // TODO: Replace _hash with its individual components
//         // bytes32 _hash,
//         // address[] memory _signers,
//         // uint8[] memory _v,
//         // bytes32[] memory _r,
//         // bytes32[] memory _s
//     )
//         public
//     {
//         require(
//             cosmosBridgeClaims[_cosmosBridgeNonce].status == Status.Completed,
//             "Cannot process an prophecy on an already completed CosmosBridgeClaim"
//         );

        // Pull the CosmosBridgeClaim from storage
        // CosmosBridgeClaim memory cosmosBridgeClaim = cosmosBridgeClaims[_cosmosBridgeNonce];

        // Attempt to process the prophecy claim (throws if unsuccessful)
        // processProphecyClaim(
        //     _cosmosBridgeNonce,
        //     msg.sender,
        //     _hash,
        //     _signers,
        //     _v,
        //     _r,
        //     _s
        // );

        // Update the CosmosBridgeClaim's status to completed
        // cosmosBridgeClaim.status = Status.Completed;

        // BridgeBank.mintCosmosToken(
        //     cosmosBridgeClaim.cosmosSender,
        //     cosmosBridgeClaim.ethereumReceiver,
        //     cosmosBridgeClaim.tokenAddress,
        //     cosmosBridgeClaim.symbol,
        //     cosmosBridgeClaim.amount
        // );
    // }

    // /*
    // * @dev: Exposes an item's current status.
    // *
    // * @param _id: The item in question.
    // * @return: Boolean indicating the lock status.
    // */
    // function getEthereumDepositStatus(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(bool)
    // {
    //     return isLockedEthereumDeposit(_id);
    // }

    // /*
    // * @dev: Exposes an item's current status.
    // *
    // * @param _id: The item in question.
    // * @return: Boolean indicating the lock status.
    // */
    // function getCosmosDepositStatus(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(bool)
    // {
    //     return isLockedCosmosDeposit(_id);
    // }

    // /*
    // * @dev: Allows access to an Ethereum deposit's information via its unique identifier.
    // *
    // * @param _id: The deposit to be viewed.
    // * @return: Original sender's Ethereum address.
    // * @return: Intended Cosmos recipient's address in bytes.
    // * @return: The lock deposit's currency, denoted by a token address.
    // * @return: The amount locked in the deposit.
    // * @return: The deposit's unique nonce.
    // */
    // function viewEthereumDeposit(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(address, bytes memory, address, uint256, uint256)
    // {
    //     return getEthereumDeposit(_id);
    // }

    // /*
    // * @dev: Allows access to a Cosmos deposit's information via its unique identifier.
    // *
    // * @param _id: The deposit to be viewed.
    // * @return: Original sender's Ethereum address.
    // * @return: Intended Cosmos recipient's address in bytes.
    // * @return: The lock deposit's currency, denoted by a token address.
    // * @return: The amount locked in the deposit.
    // * @return: The deposit's unique nonce.
    // */
    // function viewCosmosDeposit(
    //     bytes32 _id
    // )
    //     public
    //     view
    //     returns(bytes memory, address payable, address, uint256)
    // {
    //     return getCosmosDeposit(_id);
    // }

    /*
    * @dev: Provider can pause fund locking without impacting other functionality.
    */
    function pauseLocking()
        public
        // onlyOperator
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
        // onlyOperator
    {
        require(!active);
        active = true;
        emit LogLockingActivated(now);
    }
}
