pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Mintable.sol";
// import "./ProxyData.sol";

/**
 * @title BridgeNFT
 * @dev Mintable, ERC721 compatible BankNFT for use by BridgeBank
 **/

contract BridgeNFT is ERC721Mintable, ERC721Full {

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