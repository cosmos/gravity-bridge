module.exports = async () => {
    /*******************************************
     *** Set up
     ******************************************/
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    const NETWORK_ROPSTEN =
        process.argv[4] === "--network" && process.argv[5] === "ropsten";

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

    // Map containing named events associated with a known topic hash
    var eventTopics = new Map()
    eventTopics.set("0x50e466de4726c2437aa7498d554322f5599f31f0f69f9ce036ad96db77590491", "LogNewOracleClaim")
    eventTopics.set("0x802cd873de701272ec903860b690986bd460b5bcd57e30ac1fdfdeece10528ac", "LogUnlock")
    eventTopics.set("0x79e7c1c0bd54f11809c3bf6023c242783602d61ceff272c6bba6f8559c24ad0d", "LogProphecyCompleted")
    eventTopics.set("0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc", "LogProphecyProcessed")
    eventTopics.set("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "Transfer")

    try {
        // TODO: add as argument
        const txHash = "0x455a31543fca6aad846f0bb6920559881e1e9b924a47907148d5a0033a0bd56e"
        const receipt = await web3.eth.getTransactionReceipt(txHash)

        if (receipt) {
            for (let i = 0; i < receipt.logs.length; i++) {
                console.log("Log #" + i)
                const log = receipt.logs[i]
                if (log.topics) {
                    const knownEvent = eventTopics.has(log.topics[0])
                    if (knownEvent) {
                        const eventName = eventTopics.get(log.topics[0])
                        console.log("Event: ", eventName)
                    } else {
                        console.log("Topic: ", log.topics[0])
                    }
                }
                if (log.data) {
                    console.log("Data: " + log.data)
                }
                console.log()
            }
        }
        return
    } catch (error) {
        console.error({ error })
    }
};
