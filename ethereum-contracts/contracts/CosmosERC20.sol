pragma solidity ^0.4.17;

import "./ERC20.sol";

contract CosmosERC20 is ERC20 {
    bytes public name;
    uint private _totalSupply;
    mapping (address => uint) balances;
    mapping (address => mapping (address => uint)) allowed;

    address public controller;

    event Mint(address to, uint tokens);
    event Burn(address from, uint tokens);

    function totalSupply() public constant returns (uint) {
        return _totalSupply;
    }

    function balanceOf(address tokenOwner) public constant returns (uint balance) {
        return balances[tokenOwner];
    }

    function allowance(address tokenOwner, address spender) public constant returns (uint remaining) {
        return allowed[tokenOwner][spender];
    }

    function transfer(address to, uint tokens) public returns (bool success) {
        if (!(balances[msg.sender] >= tokens)) return false;
        balances[msg.sender] -= tokens;
        balances[to] += tokens;
        Transfer(msg.sender, to, tokens);
        return true;
    }

    function approve(address spender, uint tokens) public returns (bool success) {
        allowed[msg.sender][spender] = tokens;
        Approval(msg.sender, spender, tokens);
        return true;
    }

    function transferFrom(address from, address to, uint tokens) public returns (bool success) {
        if (!(balances[from] >= tokens)) return false;
        if (!(allowed[from][msg.sender] >= tokens)) return false;
        balances[from] -= tokens;
        allowed[from][msg.sender] -= tokens;
        balances[to] += tokens;
        Transfer(from, to, tokens);
        return true;
    }

    function mint(address to, uint tokens) public returns (bool success) {
        if (msg.sender != controller) return false;
        balances[to] += tokens;
        _totalSupply += tokens;
        Mint(to, tokens);
        return true;
    }

    function burn(address from, uint tokens) public returns (bool success) {
        if (msg.sender != controller) return false;
        if (!(balances[from] >= tokens)) return false;
        balances[from] -= tokens;
        _totalSupply -= tokens;
        Burn(from, tokens);
        return true;
    }

    function CosmosERC20(address _controller, bytes _name) public {
        _totalSupply = 0;
        controller = _controller;
        name = _name;
    }
}
