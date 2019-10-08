pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";

/**
 * @title BankToken
 * @dev Mintable ERC20 token controlled by bank contract
 **/

contract BankToken is ERC20Mintable {

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