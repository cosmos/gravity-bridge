pragma solidity ^0.6.6;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";


// This is the coin we test with- the REAL Bitcoin
contract TestERC20 is ERC20 {
	constructor() public ERC20("Bitcoin MAX", "MAX") {
		_mint(0xc783df8a850f42e7F7e57013759C285caa701eB6, 10000);
		_mint(0xeAD9C93b79Ae7C1591b1FB5323BD777E86e150d4, 10000);
		_mint(0xE5904695748fe4A84b40b3fc79De2277660BD1D3, 10000);
		_mint(0x92561F28Ec438Ee9831D00D1D59fbDC981b762b2, 10000);
		_mint(0x2fFd013AaA7B5a7DA93336C2251075202b33FB2B, 10000);
	}
}
