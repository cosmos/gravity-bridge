pragma solidity ^0.4.11;

import './CosmosERC20.sol';
import './Valset.sol';

contract Peggy is Valset {
    mapping (bytes => CosmosERC20) cosmosTokenAddress;

    function getCosmosTokenAddress(bytes name) constant returns (address) {
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
    function lock(bytes to, uint64 value, address token, bytes chain) external payable {
        if (token == address(0)) {
            assert(msg.value == value);
        } else {
            assert(ERC20(token).transferFrom(msg.sender, this, value));
        }

        Lock(to, value, token, chain);
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
    event Unlock(address to, uint64 value, address token, bytes indexed chain);

    function unlock(address to, uint64 value, address token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external {
        bytes32 hash = keccak256(byte(1), to, value, chain.length, chain);
        assert(Valset.verify(hash, idxs, v, r, s));

        if (token == address(0)) {
            assert(to.send(value));
        } else {
            assert(ERC20(token).transfer(to, value));
        }

        Unlock(to, value, token, chain);
    }

    event Mint(address to, uint64 value, bytes token, bytes indexed chain);

    /// Mints 1:1 backed credit for atoms/photons. Sends the transaction to the smart contract
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transference
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     * @param chain       bytes respresentation of the destination chain
     * @param idxs        indexes of each validator
     * @param v           recovery id.
     * @param r           output of ECDSA signature.
     * @param s           output of ECDSA signature.
     */
    function mint(address to, uint64 value, bytes token, bytes chain, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external {
        require(getCosmosTokenAddress(token) != 0);

        bytes32 hash = keccak256(byte(2), to, value, token.length, token, chain.length, chain/*, signatureNonce++*/);
        assert(Valset.verify(hash, idxs, v, r, s));

        assert(CosmosERC20(getCosmosTokenAddress(token)).mint(to, value));

        Mint(to, value, token, chain);
    }

    event Burn(bytes to, uint64 value, bytes token, bytes indexed chain);

    /// burns receiving funds, unlocking funds in the Cosmos Hub
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transfer
     * @param token       bytes representation of Cosmos token
     * @param chain       bytes respresentation of the Cosmos chain
     */
    function burn(bytes to, uint64 value, bytes token, bytes chain) external {
        assert(CosmosERC20(getCosmosTokenAddress(token)).burn(msg.sender, value));

        Burn(to, value, token, chain/*, witnessNonce++*/);
    }

    event Register(bytes name, address token);

    /// Deploys new CosmosERC20 contract and stores it in a mapping.
    /// Registers new Cosmos token name with its CosmosERC20 address. Called by the relayers.
    /*
     * @param string       bytes representation of destination address
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     * @param idxs        indexes of each validator
     * @param v           recovery id.
     * @param r           output of ECDSA signature.
     * @param s           output of ECDSA signature.
     */
    function register(bytes name, address token, uint16[] idxs, uint8[] v, bytes32[] r, bytes32[] s) external {
        bytes32 hash = keccak256(byte(3), name.length, name, token);
        assert(Valset.verify(hash, idxs, v, r, s));

        cosmosTokenAddress[name] = new CosmosERC20(this, name);
    }

    function Peggy(address[] initAddress, uint64[] initPower) Valset(initAddress, initPower) {

    }
}
