pragma solidity ^0.6.6;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./Peggy.sol";

pragma experimental ABIEncoderV2;

// Reentrant evil erc20
contract ReentrantERC20 {
    address state_peggyAddress;

    constructor(address _peggyAddress) public {
        state_peggyAddress = _peggyAddress;
    }

    function transfer(address recipient, uint256 amount) public returns (bool) {
        // _currentValidators, _currentPowers, _currentValsetNonce, _v, _r, _s, _args);(
        address[] memory addresses = new address[](0);
        bytes32[] memory bytes32s = new bytes32[](0);
        uint256[] memory uint256s = new uint256[](0);
        bytes memory bytess = new bytes(0);
        uint256 zero = 0;
        LogicCallArgs memory args;

        {
            args = LogicCallArgs(
                uint256s,
                addresses,
                uint256s,
                addresses,
                address(0),
                bytess,
                zero,
                bytes32(0),
                zero
            );
        }
        
        Peggy(state_peggyAddress).submitLogicCall(
            addresses, 
            uint256s, 
            zero, 
            new uint8[](0), 
            bytes32s, 
            bytes32s,
            args
        );
    }
}
