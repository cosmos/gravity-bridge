module.exports = async () => {
    /*******************************************
     *** Set up
     ******************************************/
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    // Contract abstraction
    const truffleContract = require("truffle-contract");
    const contract = truffleContract(
        require("../build/contracts/BridgeBank.json")
    );

    /*******************************************
     *** Constants
     ******************************************/
    const NETWORK_ROPSTEN =
        process.argv[4] === "--network" && process.argv[5] === "ropsten";

    /*******************************************
     *** Web3 provider
     *** Set contract provider based on --network flag
     ******************************************/
    let provider;
    if (NETWORK_ROPSTEN) {
        provider = new HDWalletProvider(
            process.env.MNEMONIC,
            "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
        );
    } else {
        provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
    }

    const web3 = new Web3(provider);
    contract.setProvider(web3.currentProvider);
    try {
        // TODO: move to arguments
        const tokenSymbol = "TEST"

        /*******************************************
         *** Contract interaction
         ******************************************/
        await contract.deployed().then(async function (instance) {
            const tokenAddress = await instance.hasLockedFunds(tokenSymbol, 1)
            console.log("Symbol:", tokenAddress)
            console.log("Token address:", tokenSymbol)
        })
    } catch (error) {
        console.error({ error })
    }
};
