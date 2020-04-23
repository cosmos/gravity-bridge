pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";


/**
 * @title CosmosBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Cosmos blockchain.
 **/

contract CosmosBank {
    using SafeMath for uint256;

    string COSMOS_NATIVE_ASSET_PREFIX = "PEGGY";

    uint256 public bridgeTokenCount;
    mapping(string => address) controlledBridgeTokens;
    mapping(bytes32 => CosmosDeposit) cosmosDeposits;

    struct CosmosDeposit {
        bytes cosmosSender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool locked;
    }

    /*
     * @dev: Event declarations
     */
    event LogNewBridgeToken(address _token, string _symbol);

    event LogBridgeTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    /*
     * @dev: Constructor, sets bridgeTokenCount
     */
    constructor() public {
        bridgeTokenCount = 0;
    }

    /*
     * @dev: Get a token symbol's corresponding bridge token address.
     *
     * @param _symbol: The token's symbol/denom without 'PEGGY' prefix.
     * @return: Address associated with the given symbol. Returns address(0) if none is found.
     */
    function getUnprefixedBridgeToken(string memory _symbol)
        public
        view
        returns (address)
    {
        return (
            controlledBridgeTokens[concat(COSMOS_NATIVE_ASSET_PREFIX, _symbol)]
        );
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
    ) internal returns (bytes32) {
        bytes32 depositID = keccak256(
            abi.encodePacked(_cosmosSender, _ethereumRecipient, _token, _amount)
        );

        cosmosDeposits[depositID] = CosmosDeposit(
            _cosmosSender,
            _ethereumRecipient,
            _token,
            _amount,
            true
        );

        return depositID;
    }

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     */
    function deployNewBridgeToken(string memory _symbol)
        internal
        returns (string memory, address)
    {
        bridgeTokenCount = bridgeTokenCount.add(1);

        // Deploy new bridge token contract
        string memory prefixedSymbol = concat(
            COSMOS_NATIVE_ASSET_PREFIX,
            _symbol
        );
        BridgeToken newBridgeToken = (new BridgeToken)(prefixedSymbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);
        controlledBridgeTokens[prefixedSymbol] = newBridgeTokenAddress;

        emit LogNewBridgeToken(newBridgeTokenAddress, prefixedSymbol);

        return (prefixedSymbol, newBridgeTokenAddress);
    }

    /*
     * @dev: Mints new cosmos tokens
     *
     * @param _cosmosSender: The sender's Cosmos address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _cosmosTokenAddress: The currency type
     * @param _symbol: comsos token symbol
     * @param _amount: number of comsos tokens to be minted
     */
    function mintNewBridgeTokens(
        bytes memory _cosmosSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    ) internal {
        require(
            controlledBridgeTokens[_symbol] == _bridgeTokenAddress,
            "Token must be a controlled bridge token"
        );

        // Mint bridge tokens
        require(
            BridgeToken(_bridgeTokenAddress).mint(_intendedRecipient, _amount),
            "Attempted mint of bridge tokens failed"
        );

        newCosmosDeposit(
            _cosmosSender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _amount
        );

        emit LogBridgeTokenMint(
            _bridgeTokenAddress,
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
    function isLockedCosmosDeposit(bytes32 _id) internal view returns (bool) {
        return (cosmosDeposits[_id].locked);
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
    function getCosmosDeposit(bytes32 _id)
        internal
        view
        returns (bytes memory, address payable, address, uint256)
    {
        CosmosDeposit memory deposit = cosmosDeposits[_id];

        return (
            deposit.cosmosSender,
            deposit.ethereumRecipient,
            deposit.bridgeTokenAddress,
            deposit.amount
        );
    }

    /*
     * @dev: Performs low gas-comsuption string concatenation
     *
     * @param _prefix: start of the string
     * @param _suffix: end of the string
     */
    function concat(string memory _prefix, string memory _suffix)
        internal
        pure
        returns (string memory)
    {
        return string(abi.encodePacked(_prefix, _suffix));
    }
}
