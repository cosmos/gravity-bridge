pragma solidity ^0.5.0;


contract BridgeRegistry {
    address public cosmosBridge;
    address public bridgeBank;
    address public oracle;
    address public valset;

    event LogContractsRegistered(
        address _cosmosBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    );

    constructor(
        address _cosmosBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    ) public {
        cosmosBridge = _cosmosBridge;
        bridgeBank = _bridgeBank;
        oracle = _oracle;
        valset = _valset;

        emit LogContractsRegistered(cosmosBridge, bridgeBank, oracle, valset);
    }
}
