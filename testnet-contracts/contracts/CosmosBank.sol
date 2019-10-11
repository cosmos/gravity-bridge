pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./CosmosToken.sol";

/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible CosmosTokens
 **/

contract CosmosBank {

    using SafeMath for uint256;

    mapping(address => bool) public cosmosTokens;
    uint256 public cosmosTokenCount;
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
    event LogCosmosTokenDeploy(
        address _token,
        string _symbol
    );

    event LogCosmosTokenMint(
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
     * @dev: Mints new cosmos tokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
\    d*/
     function mintCosmosToken(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _cosmosTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // If no comsos token address, deploy a new comsos token
        address cosmosTokenAddress = _cosmosTokenAddress;
        if(!cosmosTokens[cosmosTokenAddress]) {
            cosmosTokenAddress = deployNewCosmosToken(_symbol);
        } else {
            cosmosTokenAddress = _cosmosTokenAddress;
        }

        // Must be cosmos token controlled by the CosmosBank
        require(
            cosmosTokens[cosmosTokenAddress],
            "Invalid cosmos token address"
        );

        // Mint bank tokens
        require(
            CosmosToken(cosmosTokenAddress).mint(address(this), _amount),
            "Attempted mint of cosmos tokens failed"
        );

        newCosmosDeposit(
            _cosmosSender,
            _intendedRecipient,
            cosmosTokenAddress,
            _amount
        );

        emit LogCosmosTokenMint(cosmosTokenAddress, _symbol, _amount, _intendedRecipient);
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

        // TODO: cosmosToken contract deployment puts Peggy over gas limit, causing deployment to fail
        // Deploy new cosmos token contract
        // CosmosToken newCosmosToken = (new CosmosToken)(_symbol);

        // Set address in tokens mapping
        // address newCosmosTokenAddress = address(newCosmosToken);
        address newCosmosTokenAddress = address(0);
        cosmosTokens[newCosmosTokenAddress] = true;

        emit LogCosmosTokenDeploy(
            newCosmosTokenAddress,
            _symbol
        );

        return newCosmosTokenAddress;
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
        returns(bytes memory, address payable, address, uint256) //, uint256)
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