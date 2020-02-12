pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC721/ERC721Mintable.sol";
// import "./ProxyData.sol";

/**
 * @title BridgeNFT
 * @dev Mintable, ERC721 compatible BankNFT for use by BridgeBank
 **/

contract BridgeNFT is ERC721Mintable, ERC721Full {
    string public _name;
    string public _symbol;

    constructor(
        string memory _sym
    )
        public
        ERC721Full(
            _sym,
            _sym
        )
    {
        // Intentionally left blank
    }

    function init(string memory _sym) public {
        require(keccak256(abi.encodePacked(_name)) == keccak256(abi.encodePacked("")), "already initialized");
        _name = _sym;
        _symbol = _sym;
    }
    function name() external view returns (string memory) {
        return _name;
    }
    function symbol() external view returns (string memory) {
        return _symbol;
    }
}