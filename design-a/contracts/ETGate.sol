pragma solidity ^0.4.11;

import './CosmosERC20.sol';
import './Valset.sol';

contract ETGate is Valset {
    mapping (string => CosmosERC20) cosmosTokenAddress;

    function getCosmosTokenAddress(string name) constant returns (address) {
        return cosmosTokenAddress[name];
    }

    uint64 private _nonce = 0;

    event Lock(bytes to, uint64 value, address token, bytes indexed chain, uint64 indexed nonce);

    function lock(bytes to, uint64 value, address token, bytes chain) external payable { 

    }

    event Unlock(address to, uint64 value, address token);

    function unlock(address to, uint64 value, address token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external { 

    }

    event Mint(address to, uint64 value, bytes token);

    function mint(address to, uint64 value, bytes token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external { 
    }

    event Burn(bytes to, uint64 value, bytes token, bytes chain);

    function burn(bytes to, uint64 value, bytes token, bytes chain) external { 

    }

    event Register(string name, address token);

    function register(string name, address token, uint16[] idxs, uint8[] v, bytes32 r, bytes32[] s) external {
        
    }

    function ETGate(address[] initAddress, uint64[] initPower) Valset(initAddress, initPower) {

    }
}
