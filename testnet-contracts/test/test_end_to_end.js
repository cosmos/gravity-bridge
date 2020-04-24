const Valset = artifacts.require("Valset");
const CosmosBridge = artifacts.require("CosmosBridge");
const Oracle = artifacts.require("Oracle");
const BridgeBank = artifacts.require("BridgeBank");
const BridgeToken = artifacts.require("BridgeToken");

var bigInt = require("big-integer");

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(web3.BigNumber))
  .should();

contract("CosmosBridge", function (accounts) {
  // System operator
  const operator = accounts[0];

  // Initial validator accounts
  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];
  const userFour = accounts[4];

  // User account
  const userSeven = accounts[7];

  // Contract's enum ClaimType can be represented a sequence of integers
  const CLAIM_TYPE_BURN = 1;
  const CLAIM_TYPE_LOCK = 2;

  // Consensus threshold
  const consensusThreshold = 70;

  describe("CosmosBridge smart contract deployment", function () {
    beforeEach(async function () {
      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );
    });

    it("should deploy the CosmosBridge with the correct parameters", async function () {
      this.cosmosBridge.should.exist;

      const claimCount = await this.cosmosBridge.prophecyClaimCount();
      Number(claimCount).should.be.bignumber.equal(0);

      const cosmosBridgeValset = await this.cosmosBridge.valset();
      cosmosBridgeValset.should.be.equal(this.valset.address);
    });
  });

  describe("Claim flow", function () {
    beforeEach(async function () {
      // Set up ProphecyClaim values
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.ethereumReceiver = userSeven;
      this.ethTokenAddress = "0x0000000000000000000000000000000000000000";
      this.symbol = "ETH";
      this.nativeCosmosAssetDenom = "ATOM";
      this.prefixedNativeCosmosAssetDenom = "PEGGYATOM";
      this.amountWei = 100;
      this.amountNativeCosmos = 815;

      // Deploy Valset contract
      this.initialValidators = [userOne, userTwo, userThree, userFour];
      this.initialPowers = [30, 20, 21, 29];
      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );

      // Deploy CosmosBridge contract
      this.cosmosBridge = await CosmosBridge.new(operator, this.valset.address);

      // Deploy Oracle contract
      this.oracle = await Oracle.new(
        operator,
        this.valset.address,
        this.cosmosBridge.address,
        consensusThreshold
      );

      // Deploy BridgeBank contract
      this.bridgeBank = await BridgeBank.new(
        operator,
        this.oracle.address,
        this.cosmosBridge.address
      );

      // Operator sets Oracle
      await this.cosmosBridge.setOracle(this.oracle.address, {
        from: operator
      });

      // Operator sets Bridge Bank
      await this.cosmosBridge.setBridgeBank(this.bridgeBank.address, {
        from: operator
      });
    });

    it("Burn prophecy claim flow", async function () {
      console.log("\t[Attempt burn -> unlock]");

      // --------------------------------------------------------
      //  Lock ethereum on contract in advance of burn
      // --------------------------------------------------------
      await this.bridgeBank.lock(
        this.ethereumReceiver,
        this.ethTokenAddress,
        this.amountWei,
        {
          from: userOne,
          value: this.amountWei
        }
      ).should.be.fulfilled;

      const contractBalanceWei = await web3.eth.getBalance(
        this.bridgeBank.address
      );

      // Confirm that the contract has been loaded with funds
      contractBalanceWei.should.be.bignumber.equal(this.amountWei);

      // --------------------------------------------------------
      //  Check receiver's account balance prior to the claims
      // --------------------------------------------------------
      const priorRecipientBalance = await web3.eth.getBalance(userSeven);

      // --------------------------------------------------------
      //  Create a new burn prophecy claim on cosmos bridge
      // --------------------------------------------------------
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_BURN,
        this.cosmosSender,
        this.ethereumReceiver,
        this.symbol,
        this.amountWei,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId = Number(event.args._prophecyID);
      const claimCosmosSender = event.args._cosmosSender;
      const claimEthereumReceiver = event.args._ethereumReceiver;
      const claimTokenAddress = event.args._tokenAddress;
      const claimAmountWei = Number(event.args._amount);

      // --------------------------------------------------------
      //  Generate validator signatures and submit oracle claims
      // --------------------------------------------------------

      // Create hash using Solidity's Sha3 hashing function
      const message = web3.utils.soliditySha3(
        {
          t: "uint256",
          v: claimProphecyId
        },
        {
          t: "bytes",
          v: claimCosmosSender
        },
        {
          t: "address payable",
          v: claimEthereumReceiver
        },
        {
          t: "address",
          v: claimTokenAddress
        },
        {
          t: "uint256",
          v: claimAmountWei
        }
      );

      let signature1 = await web3.eth.sign(message, userOne);
      let signature2 = await web3.eth.sign(message, userTwo);
      let signature4 = await web3.eth.sign(message, userFour);

      // Validator userOne submits an oracle claim
      await this.oracle.newOracleClaim(claimProphecyId, message, signature1, {
        from: userOne
      }).should.be.fulfilled;

      // Validator userTwo submits an oracle claim
      await this.oracle.newOracleClaim(claimProphecyId, message, signature2, {
        from: userTwo
      }).should.be.fulfilled;

      // Validator userThree submits an oracle claim
      const { logs: pLogs } = await this.oracle.newOracleClaim(
        claimProphecyId,
        message,
        signature4,
        {
          from: userFour
        }
      ).should.be.fulfilled;

      // Validator userFour's oracle claim surpasses threshold and prophecy claim is processed
      const pEvent = pLogs.find(e => e.event === "LogProphecyProcessed");
      const processedProphecyID = Number(pEvent.args._prophecyID);
      const processedPowerCurrent = Number(pEvent.args._prophecyPowerCurrent);
      const processedPowerThreshold = Number(
        pEvent.args._prophecyPowerThreshold
      );

      processedProphecyID.should.be.bignumber.equal(claimProphecyId);
      console.log(
        "\tPower Threshold:",
        processedPowerThreshold + ".",
        "Processed at",
        processedPowerCurrent + "."
      );

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------
      const postRecipientBalance = bigInt(
        String(await web3.eth.getBalance(userSeven))
      );

      var expectedBalance = bigInt(String(priorRecipientBalance)).plus(
        String(this.amountWei)
      );

      const receivedFunds = expectedBalance.equals(postRecipientBalance);
      receivedFunds.should.be.equal(true);
    });

    it("Lock prophecy claim flow", async function () {
      console.log("\t[Attempt lock -> mint] (new)");
      const priorRecipientBalance = 0;

      // --------------------------------------------------------
      //  Create a new lock prophecy claim on cosmos bridge
      // --------------------------------------------------------
      const { logs } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.nativeCosmosAssetDenom,
        this.amountNativeCosmos,
        {
          from: userOne
        }
      ).should.be.fulfilled;

      const event = logs.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId = Number(event.args._prophecyID);
      const claimCosmosSender = event.args._cosmosSender;
      const claimEthereumReceiver = event.args._ethereumReceiver;
      const claimTokenAddress = event.args._tokenAddress;
      const claimAmount = Number(event.args._amount);
      // Check that the bridge token is a controlled bridge token
      const bridgeTokenAddr = await this.bridgeBank.getBridgeToken(
        this.prefixedNativeCosmosAssetDenom
      );
      claimTokenAddress.should.be.equal(bridgeTokenAddr);

      // --------------------------------------------------------
      //  Generate validator signatures and submit oracle claims
      // --------------------------------------------------------

      // Create hash using Solidity's Sha3 hashing function
      const message = web3.utils.soliditySha3(
        {
          t: "uint256",
          v: claimProphecyId
        },
        {
          t: "bytes",
          v: claimCosmosSender
        },
        {
          t: "address payable",
          v: claimEthereumReceiver
        },
        {
          t: "address",
          v: claimTokenAddress
        },
        {
          t: "uint256",
          v: claimAmount
        }
      );

      let signature1 = await web3.eth.sign(message, userOne);
      let signature2 = await web3.eth.sign(message, userTwo);
      let signature3 = await web3.eth.sign(message, userThree);

      // Validator userOne submits an oracle claim
      await this.oracle.newOracleClaim(claimProphecyId, message, signature1, {
        from: userOne
      }).should.be.fulfilled;

      // Validator userTwo submits an oracle claim
      await this.oracle.newOracleClaim(claimProphecyId, message, signature2, {
        from: userTwo
      }).should.be.fulfilled;

      // Validator userThree submits an oracle claim
      const { logs: pLogs } = await this.oracle.newOracleClaim(
        claimProphecyId,
        message,
        signature3,
        {
          from: userThree,
          gas: 3000000
        }
      ).should.be.fulfilled;

      // Validator userThree's oracle claim surpasses threshold and prophecy claim is processed
      const pEvent = pLogs.find(e => e.event === "LogProphecyProcessed");
      const processedProphecyID = Number(pEvent.args._prophecyID);
      const processedPowerCurrent = Number(pEvent.args._prophecyPowerCurrent);
      const processedPowerThreshold = Number(
        pEvent.args._prophecyPowerThreshold
      );

      processedProphecyID.should.be.bignumber.equal(claimProphecyId);

      console.log(
        "\tPower Threshold:",
        processedPowerThreshold + ".",
        "Processed at",
        processedPowerCurrent + "."
      );

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------

      this.bridgeToken = await BridgeToken.at(bridgeTokenAddr);

      const postRecipientBalance = bigInt(
        String(await this.bridgeToken.balanceOf(claimEthereumReceiver))
      );

      var expectedBalance = bigInt(String(priorRecipientBalance)).plus(
        String(this.amountNativeCosmos)
      );

      const receivedFunds = expectedBalance.equals(postRecipientBalance);
      receivedFunds.should.be.equal(true);

      // --------------------------------------------------------
      //  Now we'll do a 2nd lock prophecy claim of the native cosmos asset
      // --------------------------------------------------------
      console.log("\t[Attempt lock -> mint] (existing)");
      const { logs: logs2 } = await this.cosmosBridge.newProphecyClaim(
        CLAIM_TYPE_LOCK,
        this.cosmosSender,
        this.ethereumReceiver,
        this.nativeCosmosAssetDenom,
        this.amountNativeCosmos,
        {
          from: userOne
        }
      ).should.be.fulfilled;


      const event2 = logs2.find(e => e.event === "LogNewProphecyClaim");
      const claimProphecyId2 = Number(event2.args._prophecyID);
      const claimCosmosSender2 = event2.args._cosmosSender;
      const claimEthereumReceiver2 = event2.args._ethereumReceiver;
      const claimTokenAddress2 = event2.args._tokenAddress;
      const claimAmount2 = Number(event2.args._amount);

      // Check that the token contract representing the asset is the same
      bridgeTokenAddr.should.be.equal(claimTokenAddress2);
      // --------------------------------------------------------
      //  Generate validator signatures and submit oracle claims
      // --------------------------------------------------------

      // Create hash using Solidity's Sha3 hashing function
      const message2 = web3.utils.soliditySha3(
        {
          t: "uint256",
          v: claimProphecyId2
        },
        {
          t: "bytes",
          v: claimCosmosSender2
        },
        {
          t: "address payable",
          v: claimEthereumReceiver2
        },
        {
          t: "address",
          v: claimTokenAddress2
        },
        {
          t: "uint256",
          v: claimAmount2
        }
      );

      let signature2_2 = await web3.eth.sign(message2, userTwo);
      let signature2_3 = await web3.eth.sign(message2, userThree);
      let signature2_4 = await web3.eth.sign(message2, userFour);

      // Validator userOne submits an oracle claim
      await this.oracle.newOracleClaim(
        claimProphecyId2,
        message2,
        signature2_2,
        {
          from: userTwo
        }
      ).should.be.fulfilled;

      // Validator userTwo submits an oracle claim
      await this.oracle.newOracleClaim(
        claimProphecyId2,
        message2,
        signature2_3,
        {
          from: userThree
        }
      ).should.be.fulfilled;

      // Validator userFour submits an oracle claim
      const { logs: pLogs2 } = await this.oracle.newOracleClaim(
        claimProphecyId2,
        message2,
        signature2_4,
        {
          from: userFour,
          gas: 3000000
        }
      ).should.be.fulfilled;

      // Validator userFour's oracle claim surpasses threshold and prophecy claim is processed
      const pEvent2 = pLogs2.find(e => e.event === "LogProphecyProcessed");
      const processedProphecyID2 = Number(pEvent2.args._prophecyID);
      const processedPowerCurrent2 = Number(pEvent2.args._prophecyPowerCurrent);
      const processedPowerThreshold2 = Number(
        pEvent2.args._prophecyPowerThreshold
      );

      processedProphecyID2.should.be.bignumber.equal(claimProphecyId2);

      console.log(
        "\tPower Threshold:",
        processedPowerThreshold2 + ".",
        "Processed at",
        processedPowerCurrent2 + "."
      );

      // --------------------------------------------------------
      //  Check receiver's account balance after the claim is processed
      // --------------------------------------------------------

      const postRecipientBalance2 = bigInt(
        String(await this.bridgeToken.balanceOf(claimEthereumReceiver2))
      );

      var expectedBalance2 = bigInt(String(postRecipientBalance)).plus(
        String(this.amountNativeCosmos)
      );

      const receivedFunds2 = expectedBalance2.equals(postRecipientBalance2);
      receivedFunds2.should.be.equal(true);
    });
  });
});
