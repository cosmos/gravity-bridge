pragma solidity ^0.4.11;

import './Valset.sol';

contract ETGate is Valset {
    event Lock();

    function Lock() external { // deposit ether/ERC20s

    }

    event Unlock();

    function Unlock() external { // withdraw ether/ERC20s

    }

    event Burn();

    function Burn() external { // withdraw atom/photons

    }

    event Mint();

    function Mint() external { // deposit atom/photons

    }
}
