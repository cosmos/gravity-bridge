pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";

/**
 * @title CosmosToken
 * @dev Mintable ERC20 token controlled by CosmosBank
 **/

contract CosmosToken is ERC20Mintable {

    using SafeMath for uint256;

    uint8 public constant decimals = 18;

    string public symbol;

    constructor(
        string memory _symbol
    )
        public
    {
        symbol = _symbol;
    }
}