pragma solidity ^0.4.11;

import './CosmosERC20.sol';
import './Valset.sol';

contract Peggy is Valset {
    mapping (bytes => CosmosERC20) cosmosTokenAddress;

    function getCosmosTokenAddress(bytes name) constant returns (address) {
        return cosmosTokenAddress[name];
    }

    event Lock(bytes to, uint64 value, address token, bytes indexed chain);

    function lock(bytes to, uint64 value, address token, bytes chain) external payable { 
        if (token == 0) {
            assert(msg.value == value);
        } else {
            assert(ERC20(token).transferFrom(msg.sender, this, value));
        }

        Lock(to, value, token, chain);
    }

    event Unlock(address to, uint64 value, address token, bytes indexed chain);

    function unlock(address to, uint64 value, address token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external { 
        bytes32 hash = keccak256(byte(1), to, value, chain.length, chain);
        assert(Valset.verify(hash, idxs, v, r, s));

        if (token == 0) {
            assert(to.send(value));
        } else {
            assert(ERC20(token).transfer(to, value));
        }

        Unlock(to, value, token, chain);
    }

    event Mint(address to, uint64 value, bytes token, bytes indexed chain);

    function mint(address to, uint64 value, bytes token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external { 
        require(getCosmosTokenAddress(token) != 0);

        bytes32 hash = keccak256(byte(2), to, value, token.length, token, chain.length, chain/*, signatureNonce++*/);
        assert(Valset.verify(hash, idxs, v, r, s));

        assert(CosmosERC20(getCosmosTokenAddress(token)).mint(to, value));

        Mint(to, value, token, chain);
    }

    event Burn(bytes to, uint64 value, bytes token, bytes indexed chain);

    function burn(bytes to, uint64 value, bytes token, bytes chain) external { 
        assert(CosmosERC20(getCosmosTokenAddress(token)).burn(msg.sender, value));

        Burn(to, value, token, chain/*, witnessNonce++*/);
    }

    event Register(bytes name, address token);

    function register(bytes name, address token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external {
        bytes32 hash = keccak256(byte(3), name.length, name, token);
        assert(Valset.verify(hash, idxs, v, r, s));

        cosmosTokenAddress[name] = new CosmosERC20(this, name);
    }

    function Peggy(address[] initAddress, uint64[] initPower) Valset(initAddress, initPower) {

    }
}


