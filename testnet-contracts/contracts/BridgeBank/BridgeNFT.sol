pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Mintable.sol";

/**
 * @title BridgeToken
 * @dev Mintable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC721Mintable, ERC721Full {

    constructor(
        string memory _symbol
    )
        public
        ERC721Full(
            _symbol,
            _symbol
        )
    {
        // Intentionally left blank
    }
}