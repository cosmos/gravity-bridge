pragma solidity ^0.5.0;

import "./CosmosBank.sol";
import "./EthereumBank.sol";
import "./Valset.sol";

/**
 * @title BridgeBank
 * @dev Bank contract which coordinates asset-related functionality.
 *      CosmosBank manages the minting and burning of tokens which
 *      represent Cosmos based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum.
 **/

contract BridgeBank is CosmosBank, EthereumBank {

    using SafeMath for uint256;
    
    address public operator;
    Valset public valset;

    /*
    * @dev: Constructor, sets operator
    */
    constructor (
        address _operatorAddress,
        address _valsetAddress
    )
        public
    {
        operator = _operatorAddress;
        valset = Valset(_valsetAddress);
    }

    modifier onlyOperator() {
        require(
            msg.sender == operator,
            'Must be BridgeBank operator.'
        );
        _;
    }

   /*
    * @dev: Fallback function allows operator to send funds to the bank directly
    *       This feature is used for testing and is available at the operator's own risk.
    */
    function() external payable onlyOperator {}

    /*
    * @dev: Creates a new BridgeToken
    *
    * @param _symbol: The new BridgeToken's symbol
    * @return: The new BridgeToken contract's address
    */
    function createNewBridgeToken(
        string memory _symbol
    )
        public
        onlyOperator
        returns(address)
    {
        return deployNewCosmosToken(_symbol);
    }

    // TODO: Restrict to validators
    /*
     * @dev: Mints new BankTokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
\    */
     function mintBankTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _cosmosTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
    {
        return mintNewBankTokens(
            _cosmosSender,
            _intendedRecipient,
            _cosmosTokenAddress,
            _symbol,
            _amount
        );
    }

    /*
     * @dev: Locks received Ethereum funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        payable
        availableNonce
        returns(bytes32 _id)
    {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
          require(
              _token == address(0),
              "Ethereum deposits require the 'token' address to be the null address"
            );
          require(
              msg.value == _amount,
              "The transactions value must be equal the specified amount (in wei)"
            );

          // Set the the symbol to ETH
          symbol = "ETH";
          // ERC20 deposit
        } else {
          require(
              CosmosToken(_token).transferFrom(msg.sender, address(this), _amount),
              "Contract token allowances insufficient to complete this lock request"
          );
          // Set symbol to the ERC20 token's symbol
          symbol = CosmosToken(_token).symbol();
        }

        return newEthereumDeposit(
            msg.sender,
            _recipient,
            _token,
            symbol,
            _amount
        );
    }

    // TODO: Restrict this to operator (for now)
    /*
     * @dev: Unlocks Ethereum deposits.
     *
     *       Replicate _id hash off-chain with sha3(cosmosSender, ethereumRecipient, amount) + nonce
     *
     * @param _id: Unique key of the CosmosDeposit.
     */
    function unlock(
        bytes32 _id
    )
        public
        canDeliver(_id)
        returns (bool)
    {
        // TODO: Refactor this refundant check
        require(isLockedEthereumDeposit(_id), "Must be locked");

        // Unlock the deposit and transfer funds
        return unlockEthereumDeposit(_id);

    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getEthereumDepositStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLockedEthereumDeposit(_id);
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getCosmosDepositStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLockedCosmosDeposit(_id);
    }

    /*
    * @dev: Allows access to an Ethereum deposit's information via its unique identifier.
    *
    * @param _id: The deposit to be viewed.
    * @return: Original sender's Ethereum address.
    * @return: Intended Cosmos recipient's address in bytes.
    * @return: The lock deposit's currency, denoted by a token address.
    * @return: The amount locked in the deposit.
    * @return: The deposit's unique nonce.
    */
    function viewEthereumDeposit(
        bytes32 _id
    )
        public
        view
        returns(address, bytes memory, address, uint256, uint256)
    {
        return getEthereumDeposit(_id);
    }

    /*
    * @dev: Allows access to a Cosmos deposit's information via its unique identifier.
    *
    * @param _id: The deposit to be viewed.
    * @return: Original sender's Ethereum address.
    * @return: Intended Cosmos recipient's address in bytes.
    * @return: The lock deposit's currency, denoted by a token address.
    * @return: The amount locked in the deposit.
    * @return: The deposit's unique nonce.
    */
    function viewCosmosDeposit(
        bytes32 _id
    )
        public
        view
        returns(bytes memory, address payable, address, uint256)
    {
        return getCosmosDeposit(_id);
    }

}