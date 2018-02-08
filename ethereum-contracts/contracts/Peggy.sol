pragma solidity ^0.4.17;

import "./CosmosERC20.sol";
import "./Valset.sol";

contract Peggy is Valset {

    mapping (bytes => CosmosERC20) cosmosTokenAddress;

    /* Events  */

    event Unlock(address to, uint64 value, address token);
    event Lock(bytes to, uint64 value, address token);


    /* Functions */

    function getCosmosTokenAddress(bytes name) internal constant returns (address) {
      return cosmosTokenAddress[name];
    }

    /* Adapted from https://ethereum.stackexchange.com/questions/15350/how-to-convert-an-bytes-to-address-in-solidity */
    function bytesToAddress (bytes b) internal pure returns (address) {
      uint result = 0;
      uint i = 0;
      if (b[0] == 48 && b[1] == 120) {
        i = 2; // if address starts with 'Ox' begin in index 2
      }
      for (i; i < b.length; i++) {
          uint c = uint(b[i]);
          if (c >= 48 && c <= 57) {
              result = result * 16 + (c - 48);
          }
          if(c >= 65 && c<= 90) {
              result = result * 16 + (c - 55);
          }
          if(c >= 97 && c<= 122) {
              result = result * 16 + (c - 87);
          }
      }
      return address(result);
    }

    /// Locks received funds to the consensus of the peg zone
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transference
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     */
    function lock(bytes to, uint64 value, address token) external payable returns (bool) {
        if (token == address(0)) {
            require(msg.value == value);
            assert(bytesToAddress(to).send(value));

        } else {

            assert(ERC20(token).transferFrom(msg.sender, this, value)); // 'this' is the Peggy contract address
        }
        Lock(to, value, token);
        return true;
    }

    /// Unlocks Ethereum tokens according to the information from the pegzone. Called by the relayers.
    /*
     * @param to          bytes representation of destination address
     * @param value       value of transference
     * @param token       token address in origin chain (0x0 if Ethereum, Cosmos for other values)
     * @param chain       bytes respresentation of the destination chain (not used in MVP, for incentivization of relayers)
     * @param signers     indexes of each validator
     * @param v           array of recoverys id
     * @param r           array of outputs of ECDSA signature
     * @param s           array of outputs of ECDSA signature
     */
    function unlock(
        /* address[2] addressArg, // 0: token, 1: to */
        address token,
        address to,
        uint64 value,
        uint16[] signers,
        uint8[] v,
        bytes32[] r,
        bytes32[] s
    ) external returns (bool) {
        bytes32 hashData = keccak256(byte(1), token, value); /*, chain.length, chain*/
        require(Valset.verifyValidators(hashData, signers, v, r, s));
        if (token == address(0)) {
            assert(to.send(value));
        } else {
            assert(ERC20(token).transfer(to, value));
        }
        Unlock(to, value, token);
        return true;
    }

    function Peggy(address[] initAddress, uint64[] initPowers) public Valset(initAddress, initPowers) {

    }
}
