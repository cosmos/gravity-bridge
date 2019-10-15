pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BankToken.sol";

/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible tokens which
 *      represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank {

    using SafeMath for uint256;

    uint256 public cosmosTokenCount;
    mapping(address => bool) public bankTokenWhitelist;
    mapping(bytes32 => CosmosDeposit) cosmosDeposits;

    struct CosmosDeposit {
        bytes cosmosSender;
        address payable ethereumRecipient;
        address cosmosTokenAddress;
        uint256 amount;
        // uint256 nonce;
        bool locked;
    }

    /*
    * @dev: Event declarations
    */
    event LogNewBankToken(
        address _token,
        string _symbol
    );

    event LogBankTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    /*
    * @dev: Constructor, sets cosmosTokenCount
    */
    constructor () public {
        cosmosTokenCount = 0;
    }

    /*
    * @dev: Creates a new CosmosDeposit with a unique ID
    *
    * @param _cosmosSender: The sender's Cosmos address in bytes.
    * @param _ethereumRecipient: The intended recipient's Ethereum address.
    * @param _token: The currency type
    * @param _amount: The amount in the deposit.
    * @return: The newly created CosmosDeposit's unique id.
    */
    // TODO: Only called by validators
    function newCosmosDeposit(
        bytes memory _cosmosSender,
        address payable _ethereumRecipient,
        address _token,
        uint256 _amount
        // uint256 _nonce
    )
        internal
        returns(bytes32)
    {
        // cosmosDepositNonce++;

        bytes32 depositID = keccak256(
            abi.encodePacked(
                _cosmosSender,
                _ethereumRecipient,
                _token,
                _amount
                // _nonce
            )
        );

        cosmosDeposits[depositID] = CosmosDeposit(
            _cosmosSender,
            _ethereumRecipient,
            _token,
            _amount,
            // _nonce,
            true
        );

        return depositID;
    }

    /*
     * @dev: Deploys a new cosmos token contract
     *
     * @param _symbol: cosmos token symbol
     */
    function deployNewCosmosToken(
        string memory _symbol
    )
        internal
        returns(address)
    {
        cosmosTokenCount = cosmosTokenCount.add(1);

        // Deploy new cosmos token contract
        BankToken newCosmosToken = (new BankToken)(_symbol);

        // Set address in tokens mapping
        address newCosmosTokenAddress = address(newCosmosToken);
        bankTokenWhitelist[newCosmosTokenAddress] = true;

        emit LogNewBankToken(
            newCosmosTokenAddress,
            _symbol
        );

        return newCosmosTokenAddress;
    }

    // TODO: Only called by validators
    /*
     * @dev: Mints new cosmos tokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
\    */
     function mintNewBankTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _cosmosTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // Must be whitelisted token
        require(
            bankTokenWhitelist[_cosmosTokenAddress],
            "Token must be on CosmosBank's whitelist"
        );

        // Mint bank tokens
        require(
            BankToken(_cosmosTokenAddress).mint(
                address(this),
                _amount
            ),
            "Attempted mint of cosmos tokens failed"
        );

        newCosmosDeposit(
            _cosmosSender,
            _intendedRecipient,
            _cosmosTokenAddress,
            _amount
        );

        emit LogBankTokenMint(
            _cosmosTokenAddress,
            _symbol,
            _amount,
            _intendedRecipient
        );
    }

    /*
    * @dev: Checks if an individual CosmosDeposit exists.
    *
    * @param _id: The unique CosmosDeposit's id.
    * @return: Boolean indicating if the CosmosDeposit exists in memory.
    */
    function isLockedCosmosDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bool)
    {
        return(cosmosDeposits[_id].locked);
    }

  /*
    * @dev: Gets an item's information
    *
    * @param _Id: The item containing the desired information.
    * @return: Sender's address.
    * @return: Recipient's address in bytes.
    * @return: Token address.
    * @return: Amount of ethereum/erc20 in the item.
    * @return: Unique nonce of the item.
    */
    function getCosmosDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bytes memory, address payable, address, uint256)
    {
        CosmosDeposit memory deposit = cosmosDeposits[_id];

        return(
            deposit.cosmosSender,
            deposit.ethereumRecipient,
            deposit.cosmosTokenAddress,
            deposit.amount
        ); // deposit.nonce
    }
}