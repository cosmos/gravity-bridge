pragma solidity ^0.5.0;

import "./CosmosBank.sol";
import "./EthereumBank.sol";
import "../Oracle.sol";
import "../CosmosBridge.sol";


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
    Oracle public oracle;
    CosmosBridge public cosmosBridge;

    /*
     * @dev: Constructor, sets operator
     */
    constructor(
        address _operatorAddress,
        address _oracleAddress,
        address _cosmosBridgeAddress
    ) public {
        operator = _operatorAddress;
        oracle = Oracle(_oracleAddress);
        cosmosBridge = CosmosBridge(_cosmosBridgeAddress);
    }

    /*
     * @dev: Modifier to restrict access to operator
     */
    modifier onlyOperator() {
        require(msg.sender == operator, "Must be BridgeBank operator.");
        _;
    }

    /*
     * @dev: Modifier to restrict access to the oracle
     */
    modifier onlyOracle() {
        require(
            msg.sender == address(oracle),
            "Access restricted to the oracle"
        );
        _;
    }

    /*
     * @dev: Modifier to restrict access to the cosmos bridge
     */
    modifier onlyCosmosBridge() {
        require(
            msg.sender == address(cosmosBridge),
            "Access restricted to the cosmos bridge"
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
    function createNewBridgeToken(string memory _symbol)
        public
        onlyCosmosBridge
        returns (address)
    {
        return deployNewBridgeToken(_symbol);
    }

    /*
     * @dev: Mints new BankTokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
     */
    function mintBridgeTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    ) public onlyCosmosBridge {
        return
            mintNewBridgeTokens(
                _cosmosSender,
                _intendedRecipient,
                _bridgeTokenAddress,
                _symbol,
                _amount
            );
    }

    /*
     * @dev: Burns BridgeTokens representing native Cosmos assets.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function burn(bytes memory _recipient, address _token, uint256 _amount)
        public
        availableNonce()
    {
        BridgeToken(_token).burnFrom(msg.sender, _amount);
        string memory symbol = BridgeToken(_token).symbol();
        burnFunds(msg.sender, _recipient, _token, symbol, _amount);
    }

    /*
     * @dev: Locks received Ethereum/ERC20 funds.
     *
     * @param _recipient: bytes representation of destination address.
     * @param _token: token address in origin chain (0x0 if ethereum)
     * @param _amount: value of deposit
     */
    function lock(bytes memory _recipient, address _token, uint256 _amount)
        public
        payable
        availableNonce()
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
            symbol = "ETH";
        // ERC20 deposit
        } else {
            require(
                BridgeToken(_token).transferFrom(
                    msg.sender,
                    address(this),
                    _amount
                ),
                "Contract token allowances insufficient to complete this lock request"
            );
            // Set symbol to the ERC20 token's symbol
            symbol = BridgeToken(_token).symbol();
        }

        lockFunds(msg.sender, _recipient, _token, symbol, _amount);
    }

    /*
     * @dev: Unlocks Ethereum and ERC20 tokens held on the contract.
     *
     * @param _recipient: recipient's Ethereum address
     * @param _token: token contract address
     * @param _symbol: token symbol
     * @param _amount: wei amount or ERC20 token count
     */
    function unlock(
        address payable _recipient,
        string memory _symbol,
        uint256 _amount
    ) public onlyCosmosBridge {
        // Confirm that the bank has sufficient locked balances of this token type
        require(
            getLockedFunds(_symbol) >= _amount,
            "The Bank does not hold enough locked tokens to fulfill this request."
        );

        // Confirm that the bank holds sufficient balances to complete the unlock
        address tokenAddress = lockedTokenList[_symbol];
        if (tokenAddress == address(0)) {
            require(
                address(this).balance >= _amount,
                "Insufficient ethereum balance for delivery."
            );
        } else {
            require(
                BridgeToken(tokenAddress).balanceOf(address(this)) >= _amount,
                "Insufficient ERC20 token balance for delivery."
            );
        }
        unlockFunds(_recipient, tokenAddress, _symbol, _amount);
    }

    /*
     * @dev: Exposes an item's current status.
     *
     * @param _id: The item in question.
     * @return: Boolean indicating the lock status.
     */
    function getCosmosDepositStatus(bytes32 _id) public view returns (bool) {
        return isLockedCosmosDeposit(_id);
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
    function viewCosmosDeposit(bytes32 _id)
        public
        view
        returns (bytes memory, address payable, address, uint256)
    {
        return getCosmosDeposit(_id);
    }
}
