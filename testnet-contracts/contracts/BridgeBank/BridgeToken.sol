pragma solidity ^0.5.0;

import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";


/**
 * @title BridgeToken
 * @dev Mintable, ERC20Burnable, ERC20 compatible BankToken for use by BridgeBank
 **/

contract BridgeToken is ERC20Mintable, ERC20Burnable, ERC20Detailed {
    constructor(string memory _symbol)
        public
        ERC20Detailed(_symbol, _symbol, 18)
    {
        // Intentionally left blank
    }
}
