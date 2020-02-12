const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeToken = artifacts.require("BridgeToken");
const BridgeNFT = artifacts.require("BridgeNFT");
const NFTFactory = artifacts.require("NFTFactory");
const BridgeBank = artifacts.require("BridgeBank");

const { toEthSignedMessageHash, fixSignature } = require("./helpers/helpers");
// const ethers = require('ethers')
// const utils = ethers.utils

// const inBytes = utils.formatBytes32String("test");
// var web3Abi = require('web3-eth-abi');
const Web3Utils = require("web3-utils");
const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("NFTFactory", function(accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];


  describe("NFTFactory deployment and basics", function() {
    beforeEach(async function() {
      this.bridgeNFT = await BridgeNFT.new("TEST", {
        // gas: 4612388,
        from: operator
      });

      this.nftFactory = await NFTFactory.new(this.bridgeNFT.address, {
        // gas: 4612388,
        from: operator
      });

    });


    it("should have the correct BridgeNFT master and NFTFactory addresses", async function() {
        proxyMaster = await this.nftFactory.target()
        proxyMaster.should.be.equal(this.bridgeNFT.address);
      })


    it("should produce a new NFT proxy", async function() {
        data = web3.eth.abi.encodeFunctionCall({
            name: 'init',
            type: 'function',
            inputs: [{
                type: 'string',
                name: '_sym'
            }]
        }, ['FOOBAR']);

        expectedAddress = await this.nftFactory.createProxy.call(data)
        const { logs } = await this.nftFactory.createProxy(data)

        // Get the event logs and compare to expected bridge nft address and symbol
        const event = logs.find(e => e.event === "ProxyDeployed");
        event.args.proxyAddress.should.be.equal(expectedAddress);

        var newNFT = await BridgeNFT.at(expectedAddress)
        newName = await newNFT.name()
        newName.should.be.equal('FOOBAR');

      })
  });

});
