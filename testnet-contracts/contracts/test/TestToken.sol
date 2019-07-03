pragma solidity ^0.5.0;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract TestToken is ERC20Mintable { 

  using SafeMath for uint256;

  string public constant name = "Test Token";
  string public constant symbol = "TEST";
  uint8 public constant decimals = 18;
  
}