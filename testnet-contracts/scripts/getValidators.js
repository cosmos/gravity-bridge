module.exports = async () => {
    /*******************************************
     *** Set up
     ******************************************/
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");
    const BigNumber = require("bignumber.js")

    // Contract abstraction
    const truffleContract = require("truffle-contract");
    const contract = truffleContract(
        require("../build/contracts/Valset.json")
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
        // Get current accounts
        const accounts = await web3.eth.getAccounts();

        /*******************************************
         *** Contract interaction
         ******************************************/
        await contract.deployed().then(async function (instance) {
            for (let i = 0; i < accounts.length; i++) {
                console.log("Trying " + accounts[i] + "...")
                const isValidator = await instance.isActiveValidator(accounts[i], {
                    from: accounts[0],
                    value: 0,
                    gas: 300000 // 300,000 Gwei
                });
                if (isValidator) {
                    const power = new BigNumber(await instance.getValidatorPower(accounts[i], {
                        from: accounts[0],
                        value: 0,
                        gas: 300000 // 300,000 Gwei
                    }));
                    console.log("Validator " + accounts[i] + " is active! Power:", power.c[0])
                }
            }
        });
    } catch (error) {
        console.error({ error })
    }
};
