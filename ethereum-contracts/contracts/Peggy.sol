pragma solidity ^0.4.17;

import "./CosmosERC20.sol";
import "./Valset.sol";

contract Peggy is Valset {
    mapping (bytes => CosmosERC20) cosmosTokenAddress;

    function getCosmosTokenAddress(bytes name) internal constant returns (address) {
        return cosmosTokenAddress[name];
    }

    event Lock(bytes to, uint64 value, address token, bytes indexed chain);


    /// Locks received funds to the consensus of the peg zone
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transference
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     * @param chain       bytes respresentation of the destination chain
     */
    function lock(bytes to, uint64 value, address token, bytes chain) external payable returns (bool) {
        if (token == address(0)) {
            assert(msg.value == value);
        } else {
            assert(ERC20(token).transferFrom(msg.sender, this, value));
        }
        Lock(to, value, token, chain);
        return true;
    }

    /// Unlocks Ethereum tokens according to the information from the pegzone. Called by the relayers.
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transference
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     * @param chain       bytes respresentation of the destination chain (not used in MVP, for incentivization of relayers)
     * @param idxs        indexes of each validator
     * @param v           recovery id
     * @param r           output of ECDSA signature
     * @param s           output of ECDSA signature
     */
    event Unlock(address to, uint64 value, address token/*, bytes indexed chain*/);
    
    function unlock(
        address[2] addressArg, 
        uint64 value, 
        uint16[] idxs, 
        uint8[] v, 
        bytes32[] r,
        bytes32[] s 
    ) external returns (bool) {
        bytes32 hash = keccak256(byte(1), addressArg[0], value/*, chain.length, chain*/);
        require(Valset.verify(hash, idxs, v, r, s));
        if (addressArg[1] == address(0)) {
            assert(addressArg[0].send(value));
        } else {
            assert(ERC20(addressArg[1]).transfer(addressArg[1], value));
        }
        Unlock(addressArg[0], value, addressArg[1]);
        return true;
    }

    function Peggy(address[] initAddress, uint64[] initPowers) public Valset(initAddress, initPowers) {
    }
}
