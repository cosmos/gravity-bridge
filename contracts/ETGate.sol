pragma solidity ^0.4.11;

import "./ERC20.sol";
import "./TendermintLC.sol";
import "./TendermintUtil.sol";

contract ETGate is TendermintLC, TendermintUtil {
    // entry points

    event Deposit(address to, uint64 value, address token, bytes chain, uint64 seq);

    function deposit(address to, uint64 value, address token, bytes chain) external payable {
        if (token == 0) assert(value == msg.value);
        else            assert(ERC20(token).transferFrom(msg.sender, this, value));
        deposited[chain][token] += value;
        Deposit(to, value, token, chain, depositSeq++);
    }

    event Withdraw(address to, uint64 value, address token, bytes chain, uint seq);

    function withdraw(
        // withdrawal data
        address to,
        uint64  value,
        address token,
        bytes   chain,
        uint64    seq,
        // uint fee,
        // TendermintLC data
        uint      height,
        int8[]    proofInnerHeight,
        int[]     proofInnerSize,
        bytes20[] proofInnerHash,
        bytes     proofInnerDirection
    ) external {
        require(withdrawable(height, value, token, chain, seq));

        TendermintLC.verify(writeUvarint(seq), 
                            bytes20ToBytes(ripemd160(to, value, token, chain)), 
                            height,
                            proofInnerHeight,
                            proofInnerSize,
                            proofInnerHash,
                            proofInnerDirection);

        if (token == 0) assert(to.send(value));
        else            assert(ERC20(token).transfer(to, value));

        Withdraw(to, value, token, chain, seq);
    }

    // helper functions

    function bytes20ToBytes(bytes20 x) internal pure returns (bytes) {
        bytes memory res = new bytes(20);
        for (uint i = 0; i < 20; i++) {
            res[i] = x[i];
        }
        return res;
    }

    function withdrawable(uint height, uint64 value, address token, bytes chain, uint seq) constant returns (bool) {
        return available(height) &&              // is the header submitted by the relayers?
               continuous(seq) &&                // is the sequence continuous?
               deposited[chain][token] >= value; // is the zone holding enough value?
        }

    // structs and state variables

    // deposit()

    uint64 public depositSeq = 0;

    mapping (bytes => mapping (address => uint)) public deposited;
    function getDeposited(bytes k1, address k2) constant returns (uint) { return deposited[k1][k2]; }
}
