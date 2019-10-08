const CosmosBridge = artifacts.require("CosmosBridge");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("CosmosBridge", function(accounts) {
  const provider = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("CosmosBridge smart contract deployment", function() {
    beforeEach(async function() {
      this.cosmosBridge = await CosmosBridge.new();
    });

    it("should deploy the CosmosBridge", async function() {
      this.cosmosBridge.should.exist;
    });
  });
});
