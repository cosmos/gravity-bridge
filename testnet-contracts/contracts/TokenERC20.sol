pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";

contract TokenERC20 is IERC20 {
  string public name;
  string public symbol;
  uint8 public decimals;
}