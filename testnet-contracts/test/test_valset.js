const Valset = artifacts.require("Valset");

const { toEthSignedMessageHash, fixSignature } = require("./helpers/helpers");

const EVMRevert = "revert";
const BigNumber = web3.BigNumber;

require("chai")
  .use(require("chai-as-promised"))
  .use(require("chai-bignumber")(BigNumber))
  .should();

contract("Valset", function (accounts) {
  const operator = accounts[0];

  const userOne = accounts[1];
  const userTwo = accounts[2];
  const userThree = accounts[3];

  describe("Valset contract deployment", function () {
    beforeEach(async function () {
      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should deploy the Valset and correctly set the current valset version", async function () {
      this.valset.should.exist;

      const valsetValsetVersion = await this.valset.currentValsetVersion();
      Number(valsetValsetVersion).should.be.bignumber.equal(1);
    });

    it("should correctly set initial validators and initial validator count", async function () {
      const userOneValidator = await this.valset.isActiveValidator.call(
        userOne
      );
      const userTwoValidator = await this.valset.isActiveValidator.call(
        userTwo
      );
      const userThreeValidator = await this.valset.isActiveValidator.call(
        userThree
      );
      const valsetValidatorCount = await this.valset.validatorCount();

      userOneValidator.should.be.equal(true);
      userTwoValidator.should.be.equal(true);
      userThreeValidator.should.be.equal(true);
      Number(valsetValidatorCount).should.be.bignumber.equal(
        this.initialValidators.length
      );
    });

    it("should correctly set initial validator powers ", async function () {
      const userOnePower = await this.valset.getValidatorPower.call(userOne);
      const userTwoPower = await this.valset.getValidatorPower.call(userTwo);
      const userThreePower = await this.valset.getValidatorPower.call(
        userThree
      );

      Number(userOnePower).should.be.bignumber.equal(this.initialPowers[0]);
      Number(userTwoPower).should.be.bignumber.equal(this.initialPowers[1]);
      Number(userThreePower).should.be.bignumber.equal(this.initialPowers[2]);
    });

    it("should correctly set the initial total power", async function () {
      const valsetTotalPower = await this.valset.totalPower();

      Number(valsetTotalPower).should.be.bignumber.equal(
        this.initialPowers[0] + this.initialPowers[1] + this.initialPowers[2]
      );
    });
  });

  describe("Dynamic validator set", function () {
    describe("Adding validators", function () {
      beforeEach(async function () {
        this.initialValidators = [userOne];
        this.initialPowers = [5];

        this.userTwoPower = 11;
        this.userThreePower = 44;

        this.valset = await Valset.new(
          operator,
          this.initialValidators,
          this.initialPowers
        );
      });

      it("should correctly update the valset when the operator adds a new validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await this.valset.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(1);

        // Confirm initial total power
        const priorTotalPower = await this.valset.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          this.initialPowers[0]
        );

        // Operator adds a validator
        await this.valset.addValidator(userTwo, this.userTwoPower, {
          from: operator
        }).should.be.fulfilled;

        // Confirm that userTwo has been set as a validator
        const isUserTwoValidator = await this.valset.isActiveValidator.call(
          userTwo
        );
        isUserTwoValidator.should.be.equal(true);

        // Confirm that userTwo's power has been correctly set
        const userTwoSetPower = await this.valset.getValidatorPower.call(
          userTwo
        );
        Number(userTwoSetPower).should.be.bignumber.equal(this.userTwoPower);

        // Confirm updated validator count
        const postValsetValidatorCount = await this.valset.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(2);

        // Confirm updated total power
        const postTotalPower = await this.valset.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(
          this.initialPowers[0] + this.userTwoPower
        );
      });

      it("should emit a LogValidatorAdded event upon the addition of a new validator", async function () {
        // Get the event logs from the addition of a new validator
        const { logs } = await this.valset.addValidator(
          userTwo,
          this.userTwoPower,
          {
            from: operator
          }
        );
        const event = logs.find(e => e.event === "LogValidatorAdded");

        // Confirm that the event data is correct
        event.args._validator.should.be.equal(userTwo);
        Number(event.args._power).should.be.bignumber.equal(this.userTwoPower);
        Number(event.args._currentValsetVersion).should.be.bignumber.equal(1);
        Number(event.args._validatorCount).should.be.bignumber.equal(2);
        Number(event.args._totalPower).should.be.bignumber.equal(
          this.initialPowers[0] + this.userTwoPower
        );
      });

      it("should allow the operator to add multiple new validators", async function () {
        await this.valset.addValidator(userTwo, this.userTwoPower, {
          from: operator
        }).should.be.fulfilled;
        await this.valset.addValidator(userThree, this.userThreePower, {
          from: operator
        }).should.be.fulfilled;
        await this.valset.addValidator(accounts[4], 77, {
          from: operator
        }).should.be.fulfilled;
        await this.valset.addValidator(accounts[5], 23, {
          from: operator
        }).should.be.fulfilled;

        // Confirm updated validator count
        const postValsetValidatorCount = await this.valset.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(5);

        // Confirm updated total power
        const valsetTotalPower = await this.valset.totalPower();
        Number(valsetTotalPower).should.be.bignumber.equal(
          this.initialPowers[0] + this.userTwoPower + this.userThreePower + 100 // (23 + 77)
        );
      });
    });

    describe("Updating validator's power", function () {
      beforeEach(async function () {
        this.initialValidators = [userOne];
        this.initialPowers = [5];

        this.userTwoPower = 11;
        this.userThreePower = 44;

        this.valset = await Valset.new(
          operator,
          this.initialValidators,
          this.initialPowers
        );
      });

      it("should allow the operator to update a validator's power", async function () {
        const NEW_POWER = 515;

        // Confirm userOne's initial power
        const userOneInitialPower = await this.valset.getValidatorPower.call(
          userOne
        );
        Number(userOneInitialPower).should.be.bignumber.equal(
          this.initialPowers[0]
        );

        // Confirm initial total power
        const priorTotalPower = await this.valset.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          this.initialPowers[0]
        );

        // Operator updates the validator's initial power
        await this.valset.updateValidatorPower(userOne, NEW_POWER, {
          from: operator
        }).should.be.fulfilled;

        // Confirm userOne's power has increased
        const userOnePostPower = await this.valset.getValidatorPower.call(
          userOne
        );
        Number(userOnePostPower).should.be.bignumber.equal(NEW_POWER);

        // Confirm total power has been updated
        const postTotalPower = await this.valset.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(NEW_POWER);
      });

      it("should emit a LogValidatorPowerUpdated event upon the update of a validator's power", async function () {
        const NEW_POWER = 111;

        // Get the event logs from the update of a validator's power
        const { logs } = await this.valset.updateValidatorPower(
          userOne,
          NEW_POWER,
          {
            from: operator
          }
        );
        const event = logs.find(e => e.event === "LogValidatorPowerUpdated");

        // Confirm that the event data is correct
        event.args._validator.should.be.equal(userOne);
        Number(event.args._power).should.be.bignumber.equal(NEW_POWER);
        Number(event.args._currentValsetVersion).should.be.bignumber.equal(1);
        Number(event.args._validatorCount).should.be.bignumber.equal(1);
        Number(event.args._totalPower).should.be.bignumber.equal(NEW_POWER);
      });
    });

    describe("Removing validators", function () {
      beforeEach(async function () {
        this.initialValidators = [userOne, userTwo];
        this.initialPowers = [33, 21];

        this.valset = await Valset.new(
          operator,
          this.initialValidators,
          this.initialPowers
        );
      });

      it("should correctly update the valset when the operator removes a validator", async function () {
        // Confirm initial validator count
        const priorValsetValidatorCount = await this.valset.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(
          this.initialValidators.length
        );

        // Confirm initial total power
        const priorTotalPower = await this.valset.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          this.initialPowers[0] + this.initialPowers[1]
        );

        // Operator removes a validator
        await this.valset.removeValidator(userTwo, {
          from: operator
        }).should.be.fulfilled;

        // Confirm that userTwo is no longer an active validator
        const isUserTwoValidator = await this.valset.isActiveValidator.call(
          userTwo
        );
        isUserTwoValidator.should.be.equal(false);

        // Confirm that userTwo's power has been reset
        const userTwoPower = await this.valset.getValidatorPower.call(userTwo);
        Number(userTwoPower).should.be.bignumber.equal(0);

        // Confirm updated validator count
        const postValsetValidatorCount = await this.valset.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(1);

        // Confirm updated total power
        const postTotalPower = await this.valset.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(this.initialPowers[0]);
      });

      it("should emit a LogValidatorRemoved event upon the removal of a validator", async function () {
        // Get the event logs from the update of a validator's power
        const { logs } = await this.valset.removeValidator(userTwo, {
          from: operator
        });
        const event = logs.find(e => e.event === "LogValidatorRemoved");

        // Confirm that the event data is correct
        event.args._validator.should.be.equal(userTwo);
        Number(event.args._power).should.be.bignumber.equal(0);
        Number(event.args._currentValsetVersion).should.be.bignumber.equal(1);
        Number(event.args._validatorCount).should.be.bignumber.equal(1);
        Number(event.args._totalPower).should.be.bignumber.equal(
          this.initialPowers[0]
        );
      });
    });

    describe("Updating the entire valset", function () {
      beforeEach(async function () {
        this.initialValidators = [userOne, userTwo];
        this.initialPowers = [33, 21];

        this.secondValidators = [userThree, accounts[4], accounts[5]];
        this.secondPowers = [4, 19, 50];

        this.valset = await Valset.new(
          operator,
          this.initialValidators,
          this.initialPowers
        );
      });

      it("should correctly update the valset", async function () {
        // Confirm current valset version number
        const priorValsetVersion = await this.valset.currentValsetVersion();
        Number(priorValsetVersion).should.be.bignumber.equal(1);

        // Confirm initial validator count
        const priorValsetValidatorCount = await this.valset.validatorCount();
        Number(priorValsetValidatorCount).should.be.bignumber.equal(
          this.initialValidators.length
        );

        // Confirm initial total power
        const priorTotalPower = await this.valset.totalPower();
        Number(priorTotalPower).should.be.bignumber.equal(
          this.initialPowers[0] + this.initialPowers[1]
        );

        // Operator resets the valset
        await this.valset.updateValset(
          this.secondValidators,
          this.secondPowers,
          {
            from: operator
          }
        ).should.be.fulfilled;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidator = await this.valset.isActiveValidator.call(
          userOne
        );
        isUserOneValidator.should.be.equal(false);
        const isUserTwoValidator = await this.valset.isActiveValidator.call(
          userTwo
        );
        isUserTwoValidator.should.be.equal(false);

        // Confirm that all three secondary validators are now active validators
        const isUserThreeValidator = await this.valset.isActiveValidator.call(
          userThree
        );
        isUserThreeValidator.should.be.equal(true);
        const isUserFourValidator = await this.valset.isActiveValidator.call(
          accounts[4]
        );
        isUserFourValidator.should.be.equal(true);
        const isUserFiveValidator = await this.valset.isActiveValidator.call(
          accounts[5]
        );
        isUserFiveValidator.should.be.equal(true);

        // Confirm updated valset version number
        const postValsetVersion = await this.valset.currentValsetVersion();
        Number(postValsetVersion).should.be.bignumber.equal(2);

        // Confirm updated validator count
        const postValsetValidatorCount = await this.valset.validatorCount();
        Number(postValsetValidatorCount).should.be.bignumber.equal(
          this.secondValidators.length
        );

        // Confirm updated total power
        const postTotalPower = await this.valset.totalPower();
        Number(postTotalPower).should.be.bignumber.equal(
          this.secondPowers[0] + this.secondPowers[1] + this.secondPowers[2]
        );
      });

      it("should allow active validators to remain active if they are included in the new valset", async function () {
        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorFirstValsetVersion = await this.valset.isActiveValidator.call(
          userOne
        );
        isUserOneValidatorFirstValsetVersion.should.be.equal(true);

        // Operator resets the valset
        await this.valset.updateValset(
          [this.initialValidators[0]],
          [this.initialPowers[0]],
          {
            from: operator
          }
        ).should.be.fulfilled;

        // Confirm that both initial validators are no longer an active validators
        const isUserOneValidatorSecondValsetVersion = await this.valset.isActiveValidator.call(
          userOne
        );
        isUserOneValidatorSecondValsetVersion.should.be.equal(true);
      });

      it("should emit LogValsetReset and LogValsetUpdated events upon the update of the valset", async function () {
        // Get the event logs from the valset update
        const { logs } = await this.valset.updateValset(
          this.secondValidators,
          this.secondPowers,
          {
            from: operator
          }
        ).should.be.fulfilled;

        // Get the LogValsetReset event
        const eventLogValsetReset = logs.find(
          e => e.event === "LogValsetReset"
        );

        // Confirm that the LogValsetReset event data is correct
        Number(
          eventLogValsetReset.args._newValsetVersion
        ).should.be.bignumber.equal(2);
        Number(
          eventLogValsetReset.args._validatorCount
        ).should.be.bignumber.equal(0);
        Number(eventLogValsetReset.args._totalPower).should.be.bignumber.equal(
          0
        );

        // Get the LogValsetUpdated event
        const eventLogValasetUpdated = logs.find(
          e => e.event === "LogValsetUpdated"
        );

        // Confirm that the LogValsetUpdated event data is correct
        Number(
          eventLogValasetUpdated.args._newValsetVersion
        ).should.be.bignumber.equal(2);
        Number(
          eventLogValasetUpdated.args._validatorCount
        ).should.be.bignumber.equal(this.secondValidators.length);
        Number(
          eventLogValasetUpdated.args._totalPower
        ).should.be.bignumber.equal(
          this.secondPowers[0] + this.secondPowers[1] + this.secondPowers[2]
        );
      });
    });
  });

  describe("Gas recovery", function () {
    beforeEach(async function () {
      this.initialValidators = [userOne, userTwo];
      this.initialPowers = [50, 60];

      this.secondValidators = [userThree];
      this.secondPowers = [5];

      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should not allow the gas recovery of storage in use by active validators", async function () {
      // Operator attempts to recover gas from userOne's storage slot
      await this.valset
        .recoverGas(1, userOne, {
          from: operator
        })
        .should.be.rejectedWith(EVMRevert);
    });

    it("should allow the gas recovery of inactive validator storage", async function () {
      // Confirm that both initial validators are active validators
      const isUserOneValidatorPrior = await this.valset.isActiveValidator.call(
        userOne
      );
      isUserOneValidatorPrior.should.be.equal(true);
      const isUserTwoValidatorPrior = await this.valset.isActiveValidator.call(
        userTwo
      );
      isUserTwoValidatorPrior.should.be.equal(true);

      // Operator updates the valset, making userOne and userTwo inactive validators
      await this.valset.updateValset(this.secondValidators, this.secondPowers, {
        from: operator
      }).should.be.fulfilled;

      // Confirm that both initial validators are no longer an active validators
      const isUserOneValidatorPost = await this.valset.isActiveValidator.call(
        userOne
      );
      isUserOneValidatorPost.should.be.equal(false);
      const isUserTwoValidatorPost = await this.valset.isActiveValidator.call(
        userTwo
      );
      isUserTwoValidatorPost.should.be.equal(false);

      // Operator recovers gas from inactive validator userOne
      await this.valset.recoverGas(1, userOne, {
        from: operator
      }).should.be.fulfilled;

      // Operator recovers gas from inactive validator userTwo
      await this.valset.recoverGas(1, userTwo, {
        from: operator
      }).should.be.fulfilled;
    });
  });

  describe("Signature verification", function () {
    beforeEach(async function () {
      // Create hash using Solidity's Sha3 hashing function
      this.cosmosBridgeNonce = 3;
      this.cosmosSender = web3.utils.utf8ToHex(
        "985cfkop78sru7gfud4wce83kuc9rmw89rqtzmy"
      );
      this.nonce = 17;
      this.message = web3.utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: this.nonce }
      );

      this.initialValidators = [userOne, userTwo, userThree];
      this.initialPowers = [5, 8, 12];

      this.valset = await Valset.new(
        operator,
        this.initialValidators,
        this.initialPowers
      );
    });

    it("should correctly validate signatures", async function () {
      // Create the signature
      const signature = await web3.eth.sign(this.message, userOne);

      // Recover the signer address from the generated message and signature.
      const signer = await this.valset.recover(this.message, signature);
      signer.should.be.equal(userOne);
    });

    it("should not validate signatures on a different hashed message", async function () {
      // Create the signature
      const signature = await web3.eth.sign(this.message, userOne);

      // Set up a different message (has increased nonce)
      const differentMessage = web3.utils.soliditySha3(
        { t: "uint256", v: this.cosmosBridgeNonce },
        { t: "bytes", v: this.cosmosSender },
        { t: "uint256", v: Number(this.nonce + 50) }
      );

      // Recover the signer address from a different message
      const signer = await this.valset.recover(differentMessage, signature);
      signer.should.not.be.equal(userOne);
    });

    it("should not validate signatures from a different address", async function () {
      // Create the signature
      const signature = await web3.eth.sign(this.message, userTwo);

      // Recover the signer address from the generated message and signature.
      const signer = await this.valset.recover(this.message, signature);
      signer.should.not.be.equal(userOne);
    });
  });
});